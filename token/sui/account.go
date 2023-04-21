package sui

import (
	"encoding/base64"
	"errors"
	"mytoken/core/lib/bip39"
	"mytoken/core/utils"
	"mytoken/core/utils/hexutil"
	"mytoken/token/sui/types"
	"regexp"
)

const ()

type Account struct {
	keyPair types.SuiKeyPair
	Address string
}

func newAccount(scheme types.SignatureScheme, seed []byte) *Account {
	suiKeyPair := types.NewSuiKeyPair(scheme, seed)
	return &Account{
		keyPair: suiKeyPair,
		Address: types.GenAddress(scheme, suiKeyPair.PublicKey()),
	}
}

func NewAccountWithKeystore(keystore string) (*Account, error) {
	ksByte, err := base64.StdEncoding.DecodeString(keystore)
	if err != nil {
		return nil, err
	}
	scheme, err := types.NewSignatureScheme(ksByte[0])
	if err != nil {
		return nil, err
	}
	return newAccount(scheme, ksByte[1:]), nil
}

func NewAccountWithMnemonic(scheme types.SignatureScheme, mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := utils.DeriveForPath("m/44'/784'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return newAccount(scheme, key.Key), nil
}

// NewAccountWithPrivateKey rename for support android.
// Android cannot support both NewAccountWithMnemonic(string) and NewAccountWithPrivateKey(string)
func NewAccountWithPrivateKey(scheme types.SignatureScheme, prikey string) (*Account, error) {
	seed, err := hexutil.HexDecodeString(prikey)
	if err != nil {
		return nil, err
	}
	return newAccount(scheme, seed), nil
}

func (a *Account) Keystore(scheme types.SignatureScheme, privateKey string) (string, error) {
	account, err := NewAccountWithPrivateKey(scheme, privateKey)
	if err != nil {
		return "", err
	}
	pk := []byte{scheme.Flag()}
	pk = append(pk, account.PrivateKey()...)
	return base64.StdEncoding.EncodeToString(pk), nil
}

// PrivateKey @return privateKey data
func (a *Account) PrivateKey() []byte {
	return a.keyPair.PrivateKey()[:32]
}

// PrivateKeyHex @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return hexutil.HexEncodeToString(a.keyPair.PrivateKey()[:32]), nil
}

// PublicKey @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.keyPair.PublicKey()
}

// PublicKeyHex @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return hexutil.HexEncodeToString(a.keyPair.PublicKey())
}

func (a *Account) Sign(msg types.Base64Data, intent types.Intent) (*types.Signature, error) {
	signature, err := types.NewSignatureSecure[types.Base64Data](
		types.NewIntentMessage(intent, msg), &a.keyPair,
	)
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (a *Account) SignHex(msgHex string, intent types.Intent) (*types.Signature, error) {
	msg, err := types.NewBase64Data(msgHex)
	if err != nil {
		return nil, err
	}
	signature, err := types.NewSignatureSecure[types.Base64Data](
		types.NewIntentMessage(intent, *msg), &a.keyPair,
	)
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (a *Account) EncodePublicKeyToAddress(scheme types.SignatureScheme, publicKey string) (string, error) {
	publicBytes, err := hexutil.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	return types.GenAddress(scheme, publicBytes), nil
}

func (a *Account) DecodeAddressToPublicKeyHex(address string) (string, error) {
	return "", errors.New("sui cannot support decode address to public key")
}

func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("sui cannot support decode address to public key")
}

func (a *Account) IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,64}$`)
	return reg.MatchString(address)
}
