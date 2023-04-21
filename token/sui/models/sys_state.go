package models

import (
	"math"
	"mytoken/core/lib/decimal"
)

type SysState struct {
	Epoch                                 decimal.Decimal   `json:"epoch"`
	ProtocolVersion                       decimal.Decimal   `json:"protocolVersion"`
	SystemStateVersion                    decimal.Decimal   `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  decimal.Decimal   `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       decimal.Decimal   `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     decimal.Decimal   `json:"referenceGasPrice"`
	SafeMode                              bool              `json:"safeMode"`
	SafeModeStorageRewards                decimal.Decimal   `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            decimal.Decimal   `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                decimal.Decimal   `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       decimal.Decimal   `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 decimal.Decimal   `json:"epochStartTimestampMs"`
	EpochDurationMs                       decimal.Decimal   `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                decimal.Decimal   `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     decimal.Decimal   `json:"maxValidatorCount"`
	MinValidatorJoiningStake              decimal.Decimal   `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            decimal.Decimal   `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        decimal.Decimal   `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          decimal.Decimal   `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   decimal.Decimal   `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       decimal.Decimal   `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount decimal.Decimal   `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              decimal.Decimal   `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              int64             `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            decimal.Decimal   `json:"totalStake"`
	ActiveValidators                      []ActiveValidator `json:"activeValidators"`
	PendingActiveValidatorsID             string            `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           decimal.Decimal   `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []interface{}     `json:"pendingRemovals"`
	StakingPoolMappingsID                 string            `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               decimal.Decimal   `json:"stakingPoolMappingsSize"`
	InactivePoolsID                       string            `json:"inactivePoolsId"`
	InactivePoolsSize                     decimal.Decimal   `json:"inactivePoolsSize"`
	ValidatorCandidatesID                 string            `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               decimal.Decimal   `json:"validatorCandidatesSize"`
	AtRiskValidators                      [][]interface{}   `json:"atRiskValidators"`
	ValidatorReportRecords                [][]interface{}   `json:"validatorReportRecords"`
}

type ActiveValidator struct {
	SuiAddress                   string          `json:"suiAddress"`
	ProtocolPubkeyBytes          string          `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes           string          `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes            string          `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes       string          `json:"proofOfPossessionBytes"`
	Name                         string          `json:"name"`
	Description                  string          `json:"description"`
	ImageURL                     string          `json:"imageUrl"`
	ProjectURL                   string          `json:"projectUrl"`
	NetAddress                   string          `json:"netAddress"`
	P2PAddress                   string          `json:"p2pAddress"`
	PrimaryAddress               string          `json:"primaryAddress"`
	WorkerAddress                string          `json:"workerAddress"`
	NextEpochProtocolPubkeyBytes interface{}     `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   interface{}     `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  interface{}     `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   interface{}     `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          interface{}     `json:"nextEpochNetAddress"`
	NextEpochP2PAddress          interface{}     `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      interface{}     `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       interface{}     `json:"nextEpochWorkerAddress"`
	VotingPower                  decimal.Decimal `json:"votingPower"`
	OperationCapID               string          `json:"operationCapId"`
	GasPrice                     decimal.Decimal `json:"gasPrice"`
	CommissionRate               decimal.Decimal `json:"commissionRate"`
	NextEpochStake               decimal.Decimal `json:"nextEpochStake"`
	NextEpochGasPrice            decimal.Decimal `json:"nextEpochGasPrice"`
	NextEpochCommissionRate      decimal.Decimal `json:"nextEpochCommissionRate"`
	StakingPoolID                string          `json:"stakingPoolId"`
	StakingPoolActivationEpoch   decimal.Decimal `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch interface{}     `json:"stakingPoolDeactivationEpoch"`
	StakingPoolSuiBalance        decimal.Decimal `json:"stakingPoolSuiBalance"`
	RewardsPool                  decimal.Decimal `json:"rewardsPool"`
	PoolTokenBalance             decimal.Decimal `json:"poolTokenBalance"`
	PendingStake                 decimal.Decimal `json:"pendingStake"`
	PendingTotalSuiWithdraw      decimal.Decimal `json:"pendingTotalSuiWithdraw"`
	PendingPoolTokenWithdraw     decimal.Decimal `json:"pendingPoolTokenWithdraw"`
	ExchangeRatesID              string          `json:"exchangeRatesId"`
	ExchangeRatesSize            decimal.Decimal `json:"exchangeRatesSize"`
}

func (av *ActiveValidator) CalculateAPY(epoch decimal.Decimal) float64 {
	var (
		stakingPoolSuiBalance      = av.StakingPoolSuiBalance
		stakingPoolActivationEpoch = av.StakingPoolActivationEpoch
		poolTokenBalance           = av.PoolTokenBalance
	)

	// If the staking pool is active then we calculate its APY. Or if staking started in epoch 0
	if stakingPoolActivationEpoch.Cmp(decimal.Zero) == 0 {
		numEpochsParticipated := epoch.Sub(stakingPoolActivationEpoch)
		pow1, _ := stakingPoolSuiBalance.Sub(poolTokenBalance).Div(poolTokenBalance).Add(decimal.NewFromInt(1)).Float64()
		pow2, _ := decimal.NewFromInt(365).Div(numEpochsParticipated).Float64()
		apy := (math.Pow(pow1, pow2) - 1) * 100
		if apy > 100000 {
			return 0
		} else {
			return apy
		}
	} else {
		return 0
	}
}
