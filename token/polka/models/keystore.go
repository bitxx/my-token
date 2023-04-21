package models

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
	"log"
	"mytoken/core/lib/sublib/merlin"
	"mytoken/core/lib/sublib/ristretto255"
	"mytoken/core/lib/sublib/schnorrkel"
	"mytoken/core/lib/sublib/subkey"
	"mytoken/core/utils/hexutil"
	"mytoken/core/utils/u8util"
)

const (
	saltLength = 32
	pubLength  = 32

	secLength = 64

	seedLength = 32

	scryptLength = 32 + (3 * 4)
	nonceLength  = 24

	defaultN int64 = 1 << 15
	defaultP int64 = 1
	defaultR int64 = 8
)

var (
	pkcs8Divider = []byte{161, 35, 3, 33, 0}
	pkcs8Header  = []byte{48, 83, 2, 1, 1, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32}
	seedOffset   = len(pkcs8Header)
	divOffset    = seedOffset + secLength
)

type Keystore struct {
	Encoded  string    `json:"encoded"`
	Encoding *encoding `json:"encoding"`
	Address  string    `json:"address"`
	Meta     *meta     `json:"meta"`
}

type encoding struct {
	Content []string `json:"content"`
	Type    []string `json:"type"`
	Version string   `json:"version"`
}

type meta struct {
	GenesisHash string `json:"genesisHash"`
	Name        string `json:"name"`
	WhenCreated int64  `json:"whenCreated"`
}

type keyring struct {
	privateKey [64]byte
	publicKey  [32]byte
}

func (k *Keystore) CheckPassword(password string) (publicKey string, err error) {
	keyPair, err := k.decodeKeystore(k, password)
	if err != nil {
		return "", err
	}
	publicKey = k.publicKeyHex(keyPair.publicKey)
	return
}

func (k *Keystore) Sign(msg []byte, password string) ([]byte, error) {
	kr, err := k.decodeKeystore(k, password)
	if err != nil {
		return nil, err
	}
	if len(msg) > 256 {
		h := blake2b.Sum256(msg)
		msg = h[:]
	}
	signature, err := kr.sign(k.signingContext(msg))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (k *Keystore) signingContext(msg []byte) *merlin.Transcript {
	tml := merlin.NewTranscript("SigningContext")
	tml.AppendMessage([]byte(""), []byte("substrate"))
	tml.AppendMessage([]byte("sign-bytes"), msg)
	return tml
}

func (k *Keystore) decodeKeystore(ks *Keystore, password string) (*keyring, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	var (
		privateKey [64]byte
		publicKey  [32]byte
	)

	if ks.Encoding == nil || ks.Encoding.Version != "3" || len(ks.Encoding.Content) < 2 || ks.Encoding.Content[0] != "pkcs8" || ks.Encoding.Content[1] != "sr25519" {
		return nil, errors.New("unable to decode non-pkcs8 type")
	}

	encrypted, err := base64.RawStdEncoding.DecodeString(ks.Encoded)
	if err != nil {
		return nil, err
	}
	pubKey, secretKey, err := k.decodePolkaKeystoreEncoded(&password, encrypted, ks.Encoding)
	if err != nil {
		return nil, err
	}

	copy(publicKey[:], pubKey[:])
	copy(privateKey[:], secretKey[:])

	//生成pk
	_, pkByte, err := subkey.SS58Decode(ks.Address)
	if err != nil {
		return nil, err
	}
	addrPubKey := hexutil.HexEncodeToString(pkByte)

	if err != nil {
		return nil, err
	}
	if addrPubKey != "0x"+hex.EncodeToString(publicKey[:]) {
		return nil, errors.New("decoded public keys are not equal")
	}

	return &keyring{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (k *Keystore) publicKeyHex(publicKey [32]byte) string {
	return "0x" + hex.EncodeToString(publicKey[:])
}

func (k *Keystore) decodePolkaKeystoreEncoded(passphrase *string, encrypted []byte, encodeType *encoding) ([]byte, []byte, error) {
	var (
		tmpSecret [32]byte
		tmpNonce  [24]byte
		err       error
		password  []byte
	)

	if len(encodeType.Type) < 2 || encodeType.Type[1] != "xsalsa20-poly1305" {
		return nil, nil, errors.New("no encrypted data available to decode")
	}

	if passphrase == nil {
		return nil, nil, errors.New("password required to decode encrypted data")
	}

	if encrypted == nil || len(encrypted) == 0 {
		return nil, nil, errors.New("no encrypted data to decode")
	}

	encoded := encrypted

	if len(encrypted) < 24 {
		return nil, nil, errors.New("encrypted length is less than 24")
	}

	if encodeType.Type[0] == "scrypt" {
		salt := encrypted[:saltLength]

		N := u8util.ToBN(encrypted[32+0:32+4], true).Int64()
		p := u8util.ToBN(encrypted[32+4:32+8], true).Int64()
		r := u8util.ToBN(encrypted[32+8:32+12], true).Int64()

		if N != defaultN || p != defaultP || r != defaultR {
			return nil, nil, errors.New("invalid injected scrypt params found")
		}

		password, err = scrypt.Key([]byte(*passphrase), salt, int(N), int(r), int(p), 64)
		if err != nil {
			return nil, nil, err
		}

		encrypted = encrypted[scryptLength:]

	} else {
		password = []uint8(*passphrase)
	}

	secret := u8util.FixLength(password, 256, true)
	if len(secret) != 32 {
		return nil, nil, errors.New("secret length is not 32")
	}

	copy(tmpSecret[:], secret)
	copy(tmpNonce[:], encrypted[0:nonceLength])

	encoded, err = naclDecrypt(encrypted[nonceLength:], tmpNonce, tmpSecret)

	if err != nil {
		return nil, nil, err
	}
	if encoded == nil || len(encoded) == 0 {
		return nil, nil, errors.New("Encoded is nil")
	}
	header := encoded[:seedOffset]
	if string(header) != string(pkcs8Header) {
		return nil, nil, errors.New("invalid Pkcs8 header found in body")
	}
	// note: check Encoded lengths?
	secretKey := encoded[seedOffset : seedOffset+secLength]
	divider := encoded[divOffset : divOffset+len(pkcs8Divider)]
	if !bytes.Equal(divider, pkcs8Divider) {
		divOffset = seedOffset + seedLength
		secretKey = encoded[seedOffset:divOffset]
		divider = encoded[divOffset : divOffset+len(pkcs8Divider)]
		if !bytes.Equal(divider, pkcs8Divider) {
			return nil, nil, errors.New("invalid Pkcs8 divider found in body")
		}
	}
	pubOffset := divOffset + len(pkcs8Divider)
	publicKey := encoded[pubOffset : pubOffset+pubLength]
	return publicKey, secretKey, nil
}

func (kg *keyring) sign(t *merlin.Transcript) ([]byte, error) {
	var (
		pubKey [32]byte
		rByte  [64]byte
		sck    [32]byte
	)

	copy(pubKey[:], kg.publicKey[:])
	pub, err := schnorrkel.NewPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pubByte := pub.Encode()

	t.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
	t.AppendMessage([]byte("sign:pk"), pubByte[:])

	tRng, err := t.BuildRNG().ReKeyWithWitnessBytes([]byte("signing"), kg.privateKey[32:]).Finalize(rand.Reader)
	if err != nil {
		return nil, err
	}
	_, err = tRng.Read(rByte[:])
	if err != nil {
		return nil, err
	}

	r := ristretto255.NewScalar().FromUniformBytes(rByte[:])
	R := ristretto255.NewElement().ScalarBaseMult(r)
	t.AppendMessage([]byte("sign:R"), R.Encode([]byte{}))

	// form k
	k := ristretto255.NewScalar().FromUniformBytes(t.ExtractBytes([]byte("sign:c"), 64))
	//k.

	// form scalar from secret key x
	key := divideScalarByCofactor(kg.privateKey[:32])
	ms, _ := schnorrkel.NewMiniSecretKeyFromRaw([32]byte(kg.privateKey[:32]))
	ms.ExpandEd25519()
	ss := ms.Encode()
	fmt.Println(hexutil.HexEncodeToString(ss[:]))
	copy(sck[:], key[:])
	x, err := schnorrkel.ScalarFromBytes(sck)
	if err != nil {
		return nil, err
	}

	// s = kx + r
	s := x.Multiply(x, k).Add(x, r)
	r.Zero()

	//encode
	out := [64]byte{}
	renc := R.Encode([]byte{})
	copy(out[:32], renc)
	senc := s.Encode([]byte{})
	copy(out[32:], senc)
	out[63] |= 128

	return out[:], nil
}

func divideScalarByCofactor(s []byte) []byte {
	l := len(s) - 1
	low := byte(0)
	for i := range s {
		r := s[l-i] & 0x07 // remainder
		s[l-i] >>= 3
		s[l-i] += low
		low = r << 5
	}

	return s
}

func naclDecrypt(box []byte, nonce [24]byte, secret [32]byte) ([]byte, error) {
	if box == nil || len(box) == 0 {
		return nil, errors.New("cannot decrypt a nil message")
	}

	var (
		out []byte
		ok  bool
	)
	out, ok = secretbox.Open(out, box, &nonce, &secret)
	if !ok {
		return nil, errors.New("could not decrypt message")
	}

	return out, nil
}
