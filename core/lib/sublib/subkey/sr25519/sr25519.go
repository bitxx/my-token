package sr25519

import (
	"errors"
	"fmt"
	"mytoken/core/lib/sublib/merlin"
	"mytoken/core/lib/sublib/schnorrkel"
	"mytoken/core/lib/sublib/subkey"
)

const (
	miniSecretKeyLength = 32

	secretKeyLength = 64

	signatureLength = 64
)

type keyRing struct {
	seed   []byte
	secret *schnorrkel.SecretKey
	pub    *schnorrkel.PublicKey
}

func (kr keyRing) Sign(msg []byte) (signature []byte, err error) {
	sig, err := kr.secret.Sign(signingContext(msg))
	if err != nil {
		return signature, err
	}

	s := sig.Encode()
	return s[:], nil
}

func (kr keyRing) Verify(msg []byte, signature []byte) bool {
	var sigs [signatureLength]byte
	copy(sigs[:], signature)
	sig := new(schnorrkel.Signature)
	if err := sig.Decode(sigs); err != nil {
		return false
	}
	ok, err := kr.pub.Verify(sig, signingContext(msg))
	if err != nil || !ok {
		return false
	}

	return true
}

func signingContext(msg []byte) *merlin.Transcript {
	return schnorrkel.NewSigningContext([]byte("substrate"), msg)
}

// Public returns the public key in bytes
func (kr keyRing) Public() []byte {
	bytes := kr.pub.Encode()
	return bytes[:]
}

func (kr keyRing) Seed() []byte {
	return kr.seed
}

func (kr keyRing) AccountID() []byte {
	return kr.Public()
}

func (kr keyRing) SS58Address(network uint16) string {
	return subkey.SS58Encode(kr.AccountID(), network)
}

func deriveKeySoft(secret *schnorrkel.SecretKey, cc [32]byte) (*schnorrkel.SecretKey, error) {
	t := merlin.NewTranscript("SchnorrRistrettoHDKD")
	t.AppendMessage([]byte("sign-bytes"), nil)
	ek, err := secret.DeriveKey(t, cc)
	if err != nil {
		return nil, err
	}
	return ek.Secret()
}

func deriveKeyHard(secret *schnorrkel.SecretKey, cc [32]byte) (*schnorrkel.MiniSecretKey, error) {
	t := merlin.NewTranscript("SchnorrRistrettoHDKD")
	t.AppendMessage([]byte("sign-bytes"), nil)
	t.AppendMessage([]byte("chain-code"), cc[:])
	s := secret.Encode()
	t.AppendMessage([]byte("secret-key"), s[:])
	mskb := t.ExtractBytes([]byte("HDKD-hard"), miniSecretKeyLength)
	msk := [miniSecretKeyLength]byte{}
	copy(msk[:], mskb)
	return schnorrkel.NewMiniSecretKeyFromRaw(msk)
}

type Scheme struct{}

func (s Scheme) String() string {
	return "Sr25519"
}

func (s Scheme) Generate() (subkey.KeyPair, error) {
	ms, err := schnorrkel.GenerateMiniSecretKey()
	if err != nil {
		return nil, err
	}

	secret := ms.ExpandEd25519()
	pub, err := secret.Public()
	if err != nil {
		return nil, err
	}

	seed := ms.Encode()
	return keyRing{
		seed:   seed[:],
		secret: secret,
		pub:    pub,
	}, nil
}

func (s Scheme) FromSeed(seed []byte) (subkey.KeyPair, error) {
	switch len(seed) {
	case miniSecretKeyLength:
		var mss [32]byte
		copy(mss[:], seed)
		ms, err := schnorrkel.NewMiniSecretKeyFromRaw(mss)
		if err != nil {
			return nil, err
		}

		return keyRing{
			seed:   seed,
			secret: ms.ExpandEd25519(),
			pub:    ms.Public(),
		}, nil

	case secretKeyLength:
		var key, nonce [32]byte
		copy(key[:], seed[0:32])
		copy(nonce[:], seed[32:64])
		secret := schnorrkel.NewSecretKey(key, nonce)
		pub, err := secret.Public()
		if err != nil {
			return nil, err
		}

		return keyRing{
			seed:   seed,
			secret: secret,
			pub:    pub,
		}, nil
	}

	return nil, errors.New("invalid seed length")
}

func (s Scheme) FromPhrase(phrase, pwd string) (subkey.KeyPair, error) {
	ms, err := schnorrkel.MiniSecretKeyFromMnemonic(phrase, pwd)
	if err != nil {
		return nil, err
	}

	secret := ms.ExpandEd25519()
	pub, err := secret.Public()
	if err != nil {
		return nil, err
	}

	seed := ms.Encode()
	return keyRing{
		seed:   seed[:],
		secret: secret,
		pub:    pub,
	}, nil
}

func (s Scheme) Derive(pair subkey.KeyPair, djs []subkey.DeriveJunction) (subkey.KeyPair, error) {
	kr := pair.(keyRing)
	secret := kr.secret
	seed := kr.seed
	var err error
	for _, dj := range djs {
		if dj.IsHard {
			ms, err := deriveKeyHard(secret, dj.ChainCode)
			if err != nil {
				return nil, err
			}

			secret = ms.ExpandEd25519()
			if seed != nil {
				es := ms.Encode()
				seed = es[:]
			}
			continue
		}

		secret, err = deriveKeySoft(secret, dj.ChainCode)
		if err != nil {
			return nil, err
		}
		seed = nil
	}

	pub, err := secret.Public()
	if err != nil {
		return nil, err
	}

	return &keyRing{seed: seed, secret: secret, pub: pub}, nil
}

func (s Scheme) FromPublicKey(bytes []byte) (subkey.PublicKey, error) {
	if len(bytes) != 32 {
		return nil, fmt.Errorf("expected 32 bytes")
	}
	arr := [32]byte{}
	copy(arr[:], bytes[:32])
	key, err := schnorrkel.NewPublicKey(arr)
	if err != nil {
		return nil, err
	}

	return &keyRing{pub: key}, nil
}
