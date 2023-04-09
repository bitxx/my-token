package models

import "mytoken/core/lib/decimal"

type StakePool struct {
	ValidatorAddress string  `json:"validatorAddress"`
	StakingPool      string  `json:"stakingPool"`
	Stakes           []Stake `json:"stakes"`
}

type Stake struct {
	StakedSuiID       string          `json:"stakedSuiId"`
	StakeRequestEpoch decimal.Decimal `json:"stakeRequestEpoch"`
	StakeActiveEpoch  decimal.Decimal `json:"stakeActiveEpoch"`
	Principal         decimal.Decimal `json:"principal"`
	Status            string          `json:"status"`
	EstimatedReward   decimal.Decimal `json:"estimatedReward"`
}
