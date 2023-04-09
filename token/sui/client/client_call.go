package client

import (
	"context"
	"encoding/json"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
)

// MARK - Getter Function

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (c *Client) GetBalance(ctx context.Context, owner types.Address, coinType string) (json.RawMessage, error) {
	if coinType == "" {
		return c.CallContext(ctx, getBalance, owner)
	} else {
		return c.CallContext(ctx, getBalance, owner, coinType)
	}
}

func (c *Client) GetAllBalances(ctx context.Context, owner types.Address) (json.RawMessage, error) {
	return c.CallContext(ctx, getAllBalances, owner)
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(ctx context.Context, owner types.Address, coinType *string, cursor *types.ObjectId, limit uint) (json.RawMessage, error) {
	return c.CallContext(ctx, getCoins, owner, coinType, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetAllCoins(ctx context.Context, owner types.Address, cursor *types.ObjectId, limit uint) (json.RawMessage, error) {
	return c.CallContext(ctx, getAllCoins, owner, cursor, limit)
}

func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (json.RawMessage, error) {
	return c.CallContext(ctx, getCoinMetadata, coinType)
}

func (c *Client) GetObject(ctx context.Context, objID types.ObjectId, options *types.SuiObjectDataOptions) (json.RawMessage, error) {
	return c.CallContext(ctx, getObject, objID, options)
}

func (c *Client) MultiGetObjects(ctx context.Context, objIDs []types.ObjectId, options *types.SuiObjectDataOptions) (json.RawMessage, error) {
	return c.CallContext(ctx, multiGetObjects, objIDs, options)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (c *Client) GetOwnedObjects(ctx context.Context, address *types.Address, query *types.SuiObjectResponseQuery, cursor *types.CheckpointedObjectId, limit uint64) (json.RawMessage, error) {
	if limit > 0 {
		return c.CallContext(ctx, getOwnedObjects, address, query, cursor, limit)
	} else {
		return c.CallContext(ctx, getOwnedObjects, address, query, cursor)
	}
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (json.RawMessage, error) {
	return c.CallContext(ctx, getTotalSupply, coinType)
}

func (c *Client) GetTotalTransactionBlocks(ctx context.Context) (json.RawMessage, error) {
	return c.CallContext(ctx, getTotalTransactionBlocks)
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
/*func (c *Client) BatchGetObjectsOwnedByAddress(ctx context.Context, address types.Address, options types.SuiObjectDataOptions, filterType string) ([]types.SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return c.BatchGetFilteredObjectsOwnedByAddress(
		ctx, address, options, func(sod *types.SuiObjectData) bool {
			return filterType == "" || filterType == *sod.Type
		},
	)
}*/

/*func (c *Client) BatchGetFilteredObjectsOwnedByAddress(
	ctx context.Context,
	address types.Address,
	options types.SuiObjectDataOptions,
	filter func(*types.SuiObjectData) bool,
) ([]types.SuiObjectResponse, error) {
	query := types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	filteringObjs, err := c.GetOwnedObjects(ctx, address, &query, nil, 0)
	if err != nil {
		return nil, err
	}
	objIds := make([]types.ObjectId, 0)
	for _, obj := range filteringObjs.Data {
		if obj.Data == nil {
			continue // error obj
		}
		if filter != nil && filter(obj.Data) == false {
			continue // ignore objects if non-specified type
		}
		objIds = append(objIds, obj.Data.ObjectId)
	}

	return c.MultiGetObjects(ctx, objIds, &options)
}*/

func (c *Client) GetTransactionBlock(ctx context.Context, digest string, options types.TransactionBlockOptions) (json.RawMessage, error) {
	return c.CallContext(ctx, getTransactionBlock, digest, options)
}

func (c *Client) GetReferenceGasPrice(ctx context.Context) (json.RawMessage, error) {
	return c.CallContext(ctx, getReferenceGasPrice)
}

func (c *Client) GetEvents(ctx context.Context, digest string) (json.RawMessage, error) {
	return c.CallContext(ctx, getEvents, digest)
}

func (c *Client) TryGetPastObject(ctx context.Context, objectId types.ObjectId, version uint64, options *types.SuiObjectDataOptions) (json.RawMessage, error) {
	return c.CallContext(ctx, tryGetPastObject, objectId, version, options)
}

func (c *Client) DevInspectTransactionBlock(ctx context.Context, senderAddress types.Address, txByte types.Base64Data, gasPrice decimal.Decimal, epoch uint64) (json.RawMessage, error) {
	return c.CallContext(ctx, devInspectTransactionBlock, senderAddress, txByte, gasPrice.BigInt().Uint64(), epoch)
}

func (c *Client) DryRunTransaction(ctx context.Context, txBytes types.Base64Data) (json.RawMessage, error) {
	return c.CallContext(ctx, dryRunTransactionBlock, txBytes)
}

// MARK - Write Function

// TODO
func (c *Client) ExecuteTransactionBlock(ctx context.Context, txBytes types.Base64Data, signatures []types.Signature, options *types.TransactionBlockOptions, requestType string) (json.RawMessage, error) {
	return c.CallContext(ctx, executeTransactionBlock, txBytes, signatures, options, requestType)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(ctx context.Context, signer, recipient types.Address, objID types.ObjectId, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, transferObject, signer, objID, gasId, gasBudget.BigInt().Uint64(), recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(ctx context.Context, signer, recipient types.Address, suiObjID types.ObjectId, amount, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, transferSui, signer, suiObjID, gasBudget.BigInt().Uint64(), recipient, amount.BigInt().Uint64())
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (c *Client) PayAllSui(ctx context.Context, signer, recipient types.Address, inputCoins []types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, payAllSui, signer, inputCoins, recipient, gasBudget.BigInt().Uint64())
}

func (c *Client) Pay(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amounts []decimal.Decimal, gas *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	var uAmounts []uint64
	for _, a := range amounts {
		uAmounts = append(uAmounts, a.BigInt().Uint64())
	}
	return c.CallContext(ctx, pay, signer, inputCoins, recipients, uAmounts, gas, gasBudget.BigInt().Uint64())
}

func (c *Client) PaySui(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amounts []decimal.Decimal, gasBudget decimal.Decimal) (json.RawMessage, error) {
	var uAmounts []uint64
	for _, a := range amounts {
		uAmounts = append(uAmounts, a.BigInt().Uint64())
	}
	return c.CallContext(ctx, paySui, signer, inputCoins, recipients, uAmounts, gasBudget.BigInt().Uint64())
}

// SplitCoin Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, signer types.Address, Coin types.ObjectId, splitAmounts []decimal.Decimal, gas *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	var uAmounts []uint64
	for _, a := range splitAmounts {
		uAmounts = append(uAmounts, a.BigInt().Uint64())
	}
	return c.CallContext(ctx, splitCoin, signer, Coin, uAmounts, gas, gasBudget.BigInt().Uint64())
}

// SplitCoinEqual Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, signer types.Address, Coin types.ObjectId, splitCount uint64, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, splitCoinEqual, signer, Coin, splitCount, gasId, gasBudget.BigInt().Uint64())
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, signer types.Address, primaryCoin, coinToMerge types.ObjectId, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, mergeCoins, signer, primaryCoin, coinToMerge, gasId, gasBudget.BigInt().Uint64())
}

func (c *Client) Publish(ctx context.Context, sender types.Address, compiledModules []*types.Base64Data, dependencies []types.ObjectId, gas types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, publish, sender, compiledModules, dependencies, gas, gasBudget.BigInt().Uint64())
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: not support param `typeArguments` yet.
// So now only methods with `typeArguments` are supported
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (c *Client) MoveCall(ctx context.Context, signer types.Address, packageId types.ObjectId, module, function string, typeArgs []string, arguments []any, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, moveCall, signer, packageId, module, function, typeArgs, arguments, gasId, gasBudget.BigInt().Uint64())
}

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (c *Client) BatchTransaction(ctx context.Context, signer types.Address, txnParams []map[string]interface{}, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	return c.CallContext(ctx, batchTransaction, signer, txnParams, gasId, gasBudget.BigInt().Uint64())
}

func (c *Client) QueryTransactionBlocks(ctx context.Context, query types.SuiTransactionBlockResponseQuery, cursor string, limit *uint, descendingOrder bool) (json.RawMessage, error) {
	return c.CallContext(ctx, queryTransactionBlocks, query, cursor, limit, descendingOrder)
}

func (c *Client) QueryEvents(ctx context.Context, query types.EventFilter, cursor *types.EventId, limit *uint, descendingOrder bool) (json.RawMessage, error) {
	return c.CallContext(ctx, queryEvents, query, cursor, limit, descendingOrder)
}
