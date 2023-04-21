package polka

import (
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	subTypes "mytoken/core/lib/sublib/substrate-rpc-client/types"
	"mytoken/core/lib/sublib/substrate-rpc-client/types/codec"
)

type Transaction struct {
	chain *Chain
}

func NewTransaction(chain *Chain) *Transaction {
	return &Transaction{chain}
}

// SignTx 签名交易
//
//	@Description:
//	@receiver c
//	@param tx
//	@param account
//	@param passwd 若account为助记词生成，则passwd无需输入
//	@return signTx
//	@return err
func (c *Transaction) SignTx(tx *subTypes.Extrinsic, account *Account, passwd string) (signTx string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if tx == nil {
		return "", errors.New("signedTx is empty")
	}
	if account.keystore != nil && passwd == "" {
		return "", errors.New("keystore is used,passwd can not empty")
	}

	cli, err := c.chain.Client()
	if err != nil {
		return
	}

	nonce, err := cli.GetLatestNonce(account.Address)
	if err != nil {
		return
	}

	genesisHash, err := cli.GetBlockHash(0)
	if err != nil {
		return
	}

	runtimeVersion, err := cli.GetLatestRuntimeVersion()
	if err != nil {
		return
	}
	tx.Signature = subTypes.ExtrinsicSignatureV4{
		Nonce: subTypes.NewUCompactFromUInt(uint64(nonce)),
		Era:   subTypes.ExtrinsicEra{IsImmortalEra: true},
		Tip:   subTypes.NewUCompactFromUInt(0),
	}
	methodBytes, err := codec.Encode(tx.Method)
	if err != nil {
		return "", err
	}
	//组装未签名数据
	unsignTx, err := codec.Encode(subTypes.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: subTypes.ExtrinsicPayloadV3{
			Method:      methodBytes,
			Era:         tx.Signature.Era,
			Nonce:       tx.Signature.Nonce,
			Tip:         tx.Signature.Tip,
			SpecVersion: runtimeVersion.SpecVersion,
			GenesisHash: genesisHash,
			BlockHash:   genesisHash,
		},
		TransactionVersion: runtimeVersion.TransactionVersion,
	})
	if err != nil {
		return "", err
	}
	//签名数据
	signed, err := account.Sign(unsignTx, passwd)
	if err != nil {
		return "", err
	}

	//签名数据包装公钥
	multiAddr, err := account.DecodeAddressToMultiAddress(account.Address)

	tx.Signature.Signer = multiAddr
	tx.Signature.Signature = subTypes.MultiSignature{IsSr25519: true, AsSr25519: subTypes.NewSignature(signed)}
	tx.Version |= subTypes.ExtrinsicBitSigned
	return codec.EncodeToHex(tx)

}

func (c *Transaction) SendSignTx(signedTx string) (hashString string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if signedTx == "" {
		return "", errors.New("signedTx is empty")
	}
	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	return cli.AuthorSubmitExtrinsic(signedTx)
}

func (c *Transaction) EstimateGasFee(signedTx string) (fee *decimal.Decimal, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if signedTx == "" {
		return nil, errors.New("signedTx is empty")
	}
	cli, err := c.chain.Client()
	if err != nil {
		return
	}
	info, err := cli.PaymentQueryInfo(signedTx)
	if err != nil {
		return
	}

	estimateFeeStr, ok := info["partialFee"].(string)
	if !ok {
		err = errors.New("get estimated fee result nil")
		return
	}
	estimate, err := decimal.NewFromString(estimateFeeStr)
	if err != nil {
		return
	}
	return &estimate, nil
}

func (c *Transaction) TransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	// need show from subscan website
	return nil, nil
}
