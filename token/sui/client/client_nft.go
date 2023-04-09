package client

import (
	"context"
	"encoding/json"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
)

// MintNFT
// Create an unsigned transaction to mint a nft at testnet
func (c *Client) MintNFT(ctx context.Context, signer types.Address, nftName, nftDescription, nftUri string, gasId *types.ObjectId, gasBudget decimal.Decimal) (json.RawMessage, error) {
	packageId, _ := types.NewAddressFromHex("0x2")
	args := []any{
		nftName, nftDescription, nftUri,
	}
	return c.MoveCall(ctx, signer, *packageId, "testnet_nft", "mint", []string{}, args, gasId, gasBudget)
}

//
//func (c *Client) GetNFTsOwnedByAddress(ctx context.Context, address types.Address) ([]types.SuiObjectResponse, error) {
//	return c.BatchGetObjectsOwnedByAddress(ctx, address, types.SuiObjectDataOptions{
//		ShowType:    true,
//		ShowContent: true,
//		ShowOwner:   true,
//	}, "0x2::devnet_nft::DevNetNFT")
//}
