package sui

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/types"
	"mytoken/token/testcase"
	"testing"
)

func TestCoinType(t *testing.T) {
	token := NewTokenMain(defaultChain)
	t.Log(token.CoinType())
}

func TestTokenInfo(t *testing.T) {
	token, err := NewToken(defaultChain, types.SuiCoinType)
	require.Nil(t, err)

	info, err := token.TokenInfo()
	require.Nil(t, err)

	t.Log(string(info))

	mainToken := NewTokenMain(defaultChain)
	tokenInfo, err := mainToken.TokenInfo()
	require.Nil(t, err)

	t.Log(string(tokenInfo))
}

func TestBalance(t *testing.T) {
	token := NewTokenMain(defaultChain)
	b, err := token.BalanceOf(testcase.Accounts2.Sui.Address)
	require.Nil(t, err)

	t.Log(b.Total)
	t.Log(b.Usable)
}

func TestSendTokenTransferSign(t *testing.T) {
	fromAccount := M1Account(t)
	toAddress := testcase.Accounts2.Sui.Address

	token := NewTokenMain(defaultChain)
	//未签名交易
	unsignTx, err := token.BuildUnSignTokenTransferTx(fromAccount, toAddress, decimal.NewFromInt(10000000), decimal.NewFromInt(config.MaxGasForTransfer))
	require.Nil(t, err)

	//签名并提交交易
	transferService := NewTransaction(defaultChain)
	result, err := transferService.SignAndSendTx(fromAccount, unsignTx)
	require.Nil(t, err)

	t.Log(string(result))
}

func TestEstimateFees(t *testing.T) {
	fromAccount := M1Account(t)
	toAddress := testcase.Accounts2.Sui.Address

	token := NewTokenMain(defaultChain)
	resp, err := token.EstimateTransferFees(fromAccount, toAddress, decimal.NewFromInt(10000))
	require.Nil(t, err)

	t.Log(resp.String())
}
