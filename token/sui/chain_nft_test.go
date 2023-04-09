package sui

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/types"
	"testing"
	"time"
)

func TestMintNFT(t *testing.T) {
	owner := M1Account(t)
	nftService := NewNFT(defaultChain)

	var (
		timeNow = time.Now().Format("06-01-02 15:04")
		nftName = "ComingChat NFT at " + timeNow
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "http://www.wjblog.top/images/my_head-touch-icon-next.png"
	)
	//生成未签名
	unsignTx, err := nftService.BuildUnSignMintNFTTx(owner.Address, nftName, nftDesc, nftUrl, types.SuiCoinType, decimal.NewFromInt(config.MaxGasBudget))
	require.Nil(t, err)
	//签名并提交交易
	transferService := NewTransaction(defaultChain)
	result, err := transferService.SignAndSendTx(owner, unsignTx)
	require.Nil(t, err)
	t.Log(result)
}

/*func TestFetchNfts(t *testing.T) {
	// owner := "0xd059ab4c5c7d2be6537101f76c41f25cdf5cc26e"
	owner := M1Account(t).Address
	nftService := NewNFT(defaultChain)
	nfts, err := nftService.BuildUnSignMintNFTTx(owner)
	require.Nil(t, err)
	for name, group := range nfts {
		t.Log("=======================================")
		t.Logf("group: %v, count: %v", name, len(group))
		for idx, nft := range group {
			t.Logf("%4v: %v", idx, nft)
		}
	}
}*/

/*


func TestTransferNFT(t *testing.T) {
	account := M1Account(t)
	receiver := M2Account(t).Address()

	defaultChain := TestnetChain()

	nfts, err := defaultChain.FetchNFTs(account.Address())
	require.Nil(t, err)
	var nft *base.NFT
out:
	for _, group := range nfts {
		for _, n := range group {
			nft = n
			break out
		}
	}
	require.NotNil(t, nft)

	txn, err := defaultChain.TransferNFT(account.Address(), receiver, nft.Id, "", MaxGasBudget)
	require.Nil(t, err)
	signedTxn, err := txn.SignWithAccount(account)
	require.Nil(t, err)
	hash, err := defaultChain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)
	t.Log("transfer nft success, hash = ", hash)
}
*/
