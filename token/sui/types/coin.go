package types

import (
	"mytoken/core/lib/decimal"
)

const (
	SuiCoinType = "0x2::sui::SUI"
)

type CoinStruct struct {
	CoinType            string          `json:"coinType"`
	CoinObjectId        ObjectId        `json:"coinObjectId"`
	Version             decimal.Decimal `json:"version"`
	Digest              string          `json:"digest"`
	Balance             decimal.Decimal `json:"balance"`
	LockedUntilEpoch    *int            `json:"lockedUntilEpoch"`
	PreviousTransaction string          `json:"previousTransaction"`
}

type Coin = CoinStruct
