package sui

import (
	"context"
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/models"
	"mytoken/token/sui/types"
	"sort"
)

type NFT struct {
	chain *Chain
}

func NewNFT(chain *Chain) *NFT {
	return &NFT{chain: chain}
}

// BuildUnSignNFTTransferTx 生成nft交易信息
// @Description: 生成nft交易信息
// @receiver c
// @param sender
// @param receiver
// @param nftId
// @param gasId
// @param gasBudget
// @return txn
// @return err
func (n *NFT) BuildUnSignNFTTransferTx(sender, receiver, nftId string, gasBudget decimal.Decimal) (txn json.RawMessage, err error) {
	return NewObject(n.chain).BuildUnSignObjectTransferTx(sender, receiver, nftId, gasBudget)
}

func (n *NFT) BuildUnSignMintNFTTx(creator, name, description, uri, coinType string, gasBudget decimal.Decimal) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(creator)
	if err != nil {
		return nil, errors.New("Invalid creator address")
	}

	tokenServcie := NewTokenMain(n.chain)
	gasId, err := tokenServcie.PickMaxCoinId(creator, gasBudget)
	if err != nil {
		return
	}

	cli, err := n.chain.Client()
	if err != nil {
		return
	}
	return cli.MintNFT(context.Background(), *signer, name, description, uri, &gasId, gasBudget)
}

// FetchNFTs 获取nft所有
func (n *NFT) FetchNFTs(owner string, limit uint64) (res map[string][]*base.NFT, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	objectService := NewObject(n.chain)
	query := &types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	txnJson, err := objectService.FetchObjects(owner, query, limit)
	if err != nil {
		return nil, err
	}
	objectIndex := &models.ObjectIndex{}
	err = json.Unmarshal(txnJson, objectIndex)
	if err != nil {
		return nil, err
	}
	nftObjs := objectIndex.Data

	nfts := []*base.NFT{}
	for _, obj := range nftObjs {
		nft := &base.NFT{}
		nft.Id = obj.ObjectID.ShortString()
		nfts = append(nfts, nft)
	}

	sort.Slice(nfts, func(i, j int) bool {
		return nfts[i].Timestamp > nfts[j].Timestamp
	})
	group := make(map[string][]*base.NFT)
	group["Other"] = nfts
	return group, nil
}

/*
func transformNFT(nft *types.SuiObjectResponse) *base.NFT {
	if nft == nil || nft.Data == nil || nft.Data.Content == nil || nft.Data.Content.Data.MoveObject == nil {
		return nil
	}
	fields := struct {
		Id struct {
			Id string `json:"id"`
		} `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Url         string `json:"url"`
	}{}
	metaBytes, err := json.Marshal(nft.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(metaBytes, &fields)
	if err != nil {
		return nil
	}
	if fields.Name == "" && fields.Url == "" {
		return nil
	}

	return &base.NFT{
		HashString: *nft.Data.PreviousTransaction,

		Id:          fields.Id.Id,
		Name:        fields.Name,
		Description: fields.Description,
		Image:       strings.Replace(fields.Url, "ipfs://", "https://ipfs.io/ipfs/", 1),
	}
}*/

/**/

/*

func (c *Chain) FetchNFTs(owner string) (res map[string][]*base.NFT, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	address, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	nftObjects, err := client.BatchGetFilteredObjectsOwnedByAddress(context.Background(), *address, types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowPreviousTransaction: true,
	}, func(sod *types.SuiObjectData) bool {
		if strings.HasPrefix(*sod.Type, "0x2::coin::Coin<") {
			return false
		}
		return true
	})
	if err != nil {
		return
	}
	nfts := []*base.NFT{}
	for _, obj := range nftObjects {
		nft := transformNFT(&obj)
		if nft != nil {
			nfts = append(nfts, nft)
		}
	}

	sort.Slice(nfts, func(i, j int) bool {
		return nfts[i].Timestamp > nfts[j].Timestamp
	})
	group := make(map[string][]*base.NFT)
	group["Other"] = nfts
	return group, nil
}


*/
