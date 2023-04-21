package sui

import (
	"context"
	"encoding/json"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
)

type Object struct {
	chain *Chain
}

func NewObject(chain *Chain) *Object {
	return &Object{chain}
}

// BuildUnSignObjectTransferTx gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Object) BuildUnSignObjectTransferTx(sender, receiver, objectId string, gasBudget decimal.Decimal) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderAddress, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	receiverAddress, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}

	object, err := types.NewHexData(objectId)
	if err != nil {
		return nil, err
	}

	//获取gasid
	tokenService := NewTokenMain(c.chain)
	gasId, err := tokenService.PickMaxCoinId(sender, gasBudget)
	if err != nil {
		return nil, err
	}

	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	return cli.TransferObject(context.Background(), *senderAddress, *receiverAddress, *object, &gasId, gasBudget)
}

// FetchObjects 返回当前用户所有的Object
// limit传入0，则返回所有object
func (c *Object) FetchObjects(owner string, query *types.SuiObjectResponseQuery, limit uint64) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	senderAddress, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}

	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	return cli.GetOwnedObjects(context.TODO(), senderAddress, query, nil, limit)
}

// MultiGetObjects 指定id，获取对应的Object
func (c *Object) MultiGetObjects(objIds []string, options *types.SuiObjectDataOptions) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	var objIdsByte []types.ObjectId
	for _, id := range objIds {
		oid, err := types.NewHexData(id)
		if err != nil {
			return nil, err
		}
		objIdsByte = append(objIdsByte, *oid)
	}

	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	return cli.MultiGetObjects(context.TODO(), objIdsByte, options)
}
