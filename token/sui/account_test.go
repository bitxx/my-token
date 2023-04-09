package sui

import (
	"fmt"
	"mytoken/token/sui/types"
	"mytoken/token/testcase"
	"testing"

	"github.com/stretchr/testify/require"
)

// Account of os environment M1
func M1Account(t *testing.T) *Account {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	account, err := NewAccountWithMnemonic(scheme, testcase.M1)
	require.Nil(t, err)
	return account
}

func M2Account(t *testing.T) *Account {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	account, err := NewAccountWithMnemonic(scheme, testcase.M2)
	require.Nil(t, err)
	return account
}

func TestAccountInfo(t *testing.T) {
	account := M2Account(t)
	privateKey, _ := account.PrivateKeyHex()
	t.Log("privateKey:", privateKey)
	t.Log("publicKey:", account.PublicKeyHex())
	t.Log("address:", account.Address)
}

func TestAccount(t *testing.T) {
	testAccount := testcase.Accounts1
	account := M1Account(t)
	require.Equal(t, account.Address, testAccount.Sui.Address)

	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address)
}

func TestNewAccountWithKeystore(t *testing.T) {
	testAccount := testcase.Accounts1
	account, err := NewAccountWithKeystore(testAccount.Sui.KeyStore)
	require.Nil(t, err)
	require.Equal(t, account.Address, testAccount.Sui.Address)
}

func TestKeystore(t *testing.T) {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	testAccount := testcase.Accounts2
	account, err := NewAccountWithKeystore(testAccount.Sui.KeyStore)
	pk, err := account.PrivateKeyHex()
	require.Nil(t, err)

	keystore, err := account.Keystore(scheme, pk)
	require.Nil(t, err)
	require.Equal(t, keystore, testAccount.Sui.KeyStore)
}

func TestPublicKeyToAddress(t *testing.T) {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	testAccount := testcase.Accounts2
	account := M1Account(t)
	addr, err := account.EncodePublicKeyToAddress(scheme, testAccount.Sui.PublicKey)
	require.Nil(t, err)
	require.Equal(t, addr, testAccount.Sui.Address)
}

func TestAccountWithPrivatekey(t *testing.T) {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(scheme, mnemonic)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := NewAccountWithPrivateKey(scheme, privateKey)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address, accountFromPrikey.Address)
}

func TestEncodePublicKeyToAddress(t *testing.T) {
	scheme, err := types.NewSignatureScheme(types.SIGNATURE_SCHEME_FLAG_ED25519)
	require.Nil(t, err)

	account := M1Account(t)
	fmt.Println(account.EncodePublicKeyToAddress(scheme, account.PublicKeyHex()))
}

func TestIsValidAddress(t *testing.T) {
	account := M1Account(t)
	fmt.Println(account.IsValidAddress(account.Address))
}
