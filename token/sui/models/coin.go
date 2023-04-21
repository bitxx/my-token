package models

import (
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
)

type Coin struct {
	Data []struct {
		CoinType            string          `json:"coinType"`
		CoinObjectID        types.ObjectId  `json:"coinObjectId"`
		Version             decimal.Decimal `json:"version"`
		Digest              string          `json:"digest"`
		Balance             decimal.Decimal `json:"balance"`
		LockedUntilEpoch    interface{}     `json:"lockedUntilEpoch"`
		PreviousTransaction string          `json:"previousTransaction"`
	} `json:"data"`
	NextCursor  string `json:"nextCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}
