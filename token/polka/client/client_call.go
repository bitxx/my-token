package client

import (
	"errors"
	"mytoken/core/lib/sublib/substrate-rpc-client/client"
	subTypes "mytoken/core/lib/sublib/substrate-rpc-client/types"
)

// AuthorSubmitExtrinsic 提交签名交易
func (c *Client) AuthorSubmitExtrinsic(signedTx string) (hashString string, err error) {
	err = c.Api.Client.Call(&hashString, "author_submitExtrinsic", signedTx)
	return
}

// GetLatestStorage
//
//	@Description: 获取最新metadata
//	@receiver c
//	@param publicKey
//	@return accountInfo
//	@return err
func (c *Client) GetLatestStorage(publicKey []byte) (accountInfo *subTypes.AccountInfo, err error) {
	call, err := subTypes.CreateStorageKey(c.Metadata, "System", "Account", publicKey)
	if err != nil {
		return
	}
	accountInfo = &subTypes.AccountInfo{}
	ok, err := c.Api.RPC.State.GetStorageLatest(call, accountInfo)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("latest storate query error")
		return
	}
	return
}

// BuildBalancesTransfer
//
//	@Description: 转账，生成转账交易信息
//	@receiver c
//	@param mulAddr
//	@param amount
//	@return extrinsic
//	@return err
func (c *Client) BuildBalancesTransfer(mulAddr *subTypes.MultiAddress, amount subTypes.UCompact) (extrinsic subTypes.Extrinsic, err error) {
	callType, err := subTypes.NewCall(c.Metadata, "Balances.transfer", mulAddr, &amount)
	if err != nil {
		return subTypes.Extrinsic{}, err
	}
	return subTypes.NewExtrinsic(callType), nil
}

// GetLatestNonce
//
//	@Description: 获取当前账户最新nonce
//	@receiver c
//	@param walletAddress
//	@return nonce
//	@return err
func (c *Client) GetLatestNonce(walletAddress string) (nonce int64, err error) {
	err = client.CallWithBlockHash(c.Api.Client, &nonce, "system_accountNextIndex", nil, walletAddress)
	return
}

func (c *Client) PaymentQueryInfo(signTx string) (data map[string]interface{}, err error) {
	err = client.CallWithBlockHash(c.Api.Client, &data, "payment_queryInfo", nil, signTx)
	return
}

// GetBlockHash
//
//	@Description: 获取指定块hash
//	@receiver c
//	@param blockNum
//	@return hash
//	@return err
func (c *Client) GetBlockHash(blockNum uint64) (hash subTypes.Hash, err error) {
	return c.Api.RPC.Chain.GetBlockHash(blockNum)
}

func (c *Client) GetLatestRuntimeVersion() (runtimeVersion *subTypes.RuntimeVersion, err error) {
	return c.Api.RPC.State.GetRuntimeVersionLatest()
}
