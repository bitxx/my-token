package polka

import (
	"github.com/stretchr/testify/require"
	"mytoken/token/testcase"
	"testing"
)

const (
	defaultNetworkId = 30 //phala
)

// Account of os environment M1
func M1Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M1, defaultNetworkId)
	require.Nil(t, err)
	return account
}

// Account of os environment M2
func M2Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M2, defaultNetworkId)
	require.Nil(t, err)
	return account
}

func TestAccountInfo(t *testing.T) {
	testAccount := testcase.Accounts1
	account := M1Account(t)
	require.Equal(t, account.Address, testAccount.Phala.Address)
	privateKey, _ := account.PrivateKeyHex()
	t.Log("privateKey:", privateKey)
	t.Log("publicKey:", account.PublicKeyHex())
	t.Log("address:", account.Address)
}

func TestNewAccountWithPrivatekey(t *testing.T) {
	account := M1Account(t)
	privateKey, _ := account.PrivateKeyHex()
	accountFromPrikey, err := NewAccountWithPrivateKey(privateKey, defaultNetworkId)
	require.Nil(t, err)
	require.Equal(t, accountFromPrikey.Address, testcase.Accounts1.Phala.Address)
}

func TestNewAccountWithKeystore(t *testing.T) {
	testAccount := testcase.Accounts1
	account, err := NewAccountWithKeystore(testAccount.Phala.KeyStore, "testtest", defaultNetworkId)
	require.Nil(t, err)
	require.Equal(t, account.Address, testAccount.Phala.Address)
}
