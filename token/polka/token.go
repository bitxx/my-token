package polka

import (
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	subTypes "mytoken/core/lib/sublib/substrate-rpc-client/types"
)

type Token struct {
	Chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{chain}
}

// TokenInfo
// Warning: polka chain is not currently supported
func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("substrate chain is not currently supported")
}

func (t *Token) BalanceOf(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	cli, err := t.Chain.Client()
	if err != nil {
		return
	}
	account := Account{}
	publicKey, err := account.DecodeAddressToPublicKey(address)
	if err != nil {
		return
	}
	data, err := cli.GetLatestStorage(publicKey)
	if err != nil {
		return
	}
	totalInt := decimal.NewFromBigInt(data.Data.Free.Int, 0).Add(decimal.NewFromBigInt(data.Data.Reserved.Int, 0))
	locked := decimal.Max(decimal.NewFromBigInt(data.Data.MiscFrozen.Int, 0), decimal.NewFromBigInt(data.Data.FreeFrozen.Int, 0))
	usableInt := decimal.NewFromBigInt(data.Data.Free.Int, 0).Sub(locked)

	return &base.Balance{
		Total:  totalInt,
		Usable: usableInt,
	}, nil
}

// BuildUnSignTokenTransferTx
//
//	@Description: 未签名交易
//	@receiver t
//	@param receiverAddress
//	@param amount
//	@return extrinsic
//	@return err
func (t *Token) BuildUnSignTokenTransferTx(receiverAddress string, amount decimal.Decimal) (extrinsic subTypes.Extrinsic, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	mulAddr, err := DefaultAccount.DecodeAddressToMultiAddress(receiverAddress)
	if err != nil {
		return
	}
	cli, err := t.Chain.Client()
	if err != nil {
		return
	}
	return cli.BuildBalancesTransfer(&mulAddr, subTypes.NewUCompact(amount.BigInt()))
}

/*func (t *Token) EstimateTransferFees(account *Account, receiverAddress string, amount decimal.Decimal) (f *decimal.Decimal, err error) {
	txn, err := t.BuildUnSignTokenTransferTx(account, receiverAddress, amount, decimal.NewFromInt(config.MaxGasForTransfer))
	if err != nil {
		return
	}
	transaction := &models.Transaction{}
	err = json.Unmarshal(txn, &transaction)
	if err != nil {
		return
	}
	return NewTransaction(t.chain).EstimateGasFee(transaction.TxBytes)
}*/
