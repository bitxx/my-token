package sui

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/types"
	"mytoken/token/testcase"
	"testing"
)

func TestTransferObjectWithSign(t *testing.T) {
	testObjectId := "0x16213dc30181ef4fd43fb2705f8bb3c1569499687a45f77880c7859544a06629"
	fromAccount := M1Account(t)
	toAddress := testcase.Accounts2.Sui.Address
	transaction := NewObject(DevnetChain())
	//生成未签名交易
	unsignTx, err := transaction.BuildUnSignObjectTransferTx(fromAccount.Address, toAddress, testObjectId, types.SuiCoinType, decimal.NewFromInt(config.PerObjectMaxGasForPay))
	if err != nil {
		t.Fatal(err)
	}

	//签名并提交交易
	transferService := NewTransaction(defaultChain)
	result, err := transferService.SignAndSendTx(fromAccount, unsignTx)
	require.Nil(t, err)

	t.Log(string(result))
}

func TestGetObjects(t *testing.T) {
	query := &types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType:                true,
			ShowContent:             true,
			ShowBcs:                 true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
			ShowDisplay:             true,
		},
	}
	objectService := NewObject(defaultChain)
	objects, err := objectService.FetchObjects(testcase.Accounts1.Sui.Address, query, 0)
	require.Nil(t, err)
	t.Log(string(objects))
}

func TestMultiGetObjects(t *testing.T) {
	objectService := NewObject(defaultChain)
	objIds := []string{
		"0x0a2b788d9df9a2f57d913b033f25fe3a27709280fd4be9b79937e77b29422981",
	}
	opt := types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	}
	objects, err := objectService.MultiGetObjects(objIds, &opt)
	require.Nil(t, err)
	t.Log(string(objects))
}
