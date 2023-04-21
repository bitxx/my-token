package polka

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"testing"
)

func TestTransactionDetail(t *testing.T) {
	txService := NewTransaction(TestnetChain())
	detail, err := txService.TransactionDetail("0xae206d9cda14319830daf956dc596f93b1ef5931b7d5e5b2565be05ec0518dd8")
	require.Nil(t, err)
	t.Log(detail)
}

func TestSignTx(t *testing.T) {
	testAccount := M1Account(t)
	test2Account := M2Account(t)
	tokenService := NewToken(defaultChain)
	extrinsic, err := tokenService.BuildUnSignTokenTransferTx(test2Account.Address, decimal.NewFromInt(199999838881992))
	require.Nil(t, err)
	txService := NewTransaction(tokenService.Chain)
	signedTx, err := txService.SignTx(&extrinsic, testAccount, "")
	t.Log(signedTx)
}

func TestSendSignTx(t *testing.T) {
	testAccount := M1Account(t)
	test2Account := M2Account(t)
	tokenService := NewToken(defaultChain)
	extrinsic, err := tokenService.BuildUnSignTokenTransferTx(test2Account.Address, decimal.NewFromInt(199999838881992))
	require.Nil(t, err)
	txService := NewTransaction(tokenService.Chain)
	signedTx, err := txService.SignTx(&extrinsic, testAccount, "")
	hash, err := txService.SendSignTx(signedTx)
	require.Nil(t, err)
	t.Log(hash)
	//hash：0xd8c303825e524c641dcba563ae4751e48edbcba701ddbf9e4ec05c9c9052ddfe
}

func TestEstimateGasFee(t *testing.T) {
	testAccount := M1Account(t)
	test2Account := M2Account(t)
	tokenService := NewToken(defaultChain)
	extrinsic, err := tokenService.BuildUnSignTokenTransferTx(test2Account.Address, decimal.NewFromInt(99))
	require.Nil(t, err)
	txService := NewTransaction(tokenService.Chain)
	signedTx, err := txService.SignTx(&extrinsic, testAccount, "")
	estimate, err := txService.EstimateGasFee(signedTx)
	require.Nil(t, err)
	t.Log(estimate)
}
