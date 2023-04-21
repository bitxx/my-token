package sui

import (
	"context"
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/models"
	"mytoken/token/sui/types"
)

type Transaction struct {
	chain *Chain
}

func NewTransaction(chain *Chain) *Transaction {
	return &Transaction{chain}
}

// SignAndSendTx 签名并提交交易
func (c *Transaction) SignAndSendTx(account *Account, unsignTxn json.RawMessage) (result json.RawMessage, err error) {
	signTx, txBytes, err := c.SignTx(account, unsignTxn)
	if err != nil {
		return nil, err
	}
	return c.SendSignTx(txBytes, []types.Signature{*signTx})
}

// SignTx 签名交易
func (c *Transaction) SignTx(account *Account, unsignTxn json.RawMessage) (signTx *types.Signature, txBytes *types.Base64Data, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	model := &models.Transaction{}
	err = json.Unmarshal(unsignTxn, &model)
	if err != nil {
		return nil, nil, err
	}
	signTx, err = account.Sign(model.TxBytes, types.DefaultIntent())
	return signTx, &model.TxBytes, err
}

// SendSignTx 提交已签名交易到链上
//
//	@Description: 提交已签名交易到链上
//	@receiver c
//	@param signedTx
//	@return hash
//	@return err
func (c *Transaction) SendSignTx(txBytes *types.Base64Data, signatures []types.Signature) (result json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if txBytes == nil || signatures == nil {
		return nil, errors.New("signedTx is empty")
	}
	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	options := &types.TransactionBlockOptions{ShowInput: true, ShowEffects: true, ShowEvents: true, ShowObjectChanges: true, ShowBalanceChanges: true}
	return cli.ExecuteTransactionBlock(context.TODO(), *txBytes, signatures, options, types.TxnRequestTypeWaitForLocalExecution)
}

func (c *Transaction) EstimateGasFee(txBytes types.Base64Data) (fee *decimal.Decimal, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	resp, err := cli.DryRunTransaction(context.TODO(), txBytes)
	if err != nil {
		return
	}
	dry := &models.DryRunTransaction{}
	err = json.Unmarshal(resp, dry)
	if err != nil {
		return
	}
	return dry.GasFee()
}

// TransactionDetail details through transaction hash
func (c *Transaction) TransactionDetail(hash string) (result json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	options := types.TransactionBlockOptions{ShowInput: true, ShowEffects: true, ShowEvents: true, ShowObjectChanges: true, ShowBalanceChanges: true}
	return cli.GetTransactionBlock(context.Background(), hash, options)
}
