package sui

import (
	"context"
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/models"
	"mytoken/token/sui/types"
	"sync"
)

var cachedStakePoolsMap sync.Map

type Stake struct {
	StakeId      string          `json:"stakeId"`
	Principal    decimal.Decimal `json:"principal"`
	RequestEpoch int64           `json:"requestEpoch"`
	ActiveEpoch  int64           `json:"activeEpoch"`
	Status       string          `json:"status"`
	EarnedAmount decimal.Decimal `json:"earnedAmount"`
}

type StakePool struct {
	chain            *Chain
	ValidatorAddress string     `json:"validatorAddress"`
	StakingPool      string     `json:"stakingPool"`
	Validator        *Validator `json:"validator"`
	Stakes           []Stake    `json:"stakes"`
}

func NewStakePool(chain *Chain) *StakePool {
	return &StakePool{chain: chain}
}

// GetStakePools 获取某地址所有质押池
func (sp *StakePool) GetStakePools(owner string, useCache bool) (stakePools []StakePool, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if sp == nil || sp.chain == nil || owner == "" {
		return nil, errors.New("param is error")
	}

	if useCache {
		if cachedPools, ok := cachedStakePoolsMap.Load(owner); ok {
			if v, ok := cachedPools.([]StakePool); ok {
				return v, nil
			}
		}
	}

	addr, err := types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}

	cli, err := sp.chain.Client()
	if err != nil {
		return
	}

	stakePoosJson, err := cli.GetStakePools(context.TODO(), *addr)
	if err != nil {
		return nil, err
	}

	var stakePoolModes []models.StakePool
	err = json.Unmarshal(stakePoosJson, &stakePoolModes)
	if err != nil {
		return nil, err
	}

	validatorStateService := NewValidatorState(sp.chain)
	validatorState, err := validatorStateService.TotalActiveValidatorState(useCache)
	if err != nil {
		return nil, err
	}
	for _, stakePoolModel := range stakePoolModes {
		myValidator := &Validator{}
		for _, v := range validatorState.Validators {
			if types.IsSameStringAddress(stakePoolModel.ValidatorAddress, v.Address) {
				myValidator = &v
				break
			}
		}
		if myValidator == nil {
			continue
		}

		var stakes []Stake
		for _, s := range stakePoolModel.Stakes {
			stake := Stake{
				StakeId:      s.StakedSuiID,
				Principal:    s.Principal,
				RequestEpoch: s.StakeRequestEpoch.IntPart(),
				ActiveEpoch:  s.StakeActiveEpoch.IntPart(),
				Status:       s.Status,
				EarnedAmount: s.EstimatedReward,
			}
			stakes = append(stakes, stake)
		}

		resSS := StakePool{
			ValidatorAddress: stakePoolModel.ValidatorAddress,
			StakingPool:      stakePoolModel.StakingPool,
			Validator:        myValidator,
			Stakes:           stakes,
		}
		stakePools = append(stakePools, resSS)
	}

	cachedStakePoolsMap.Store(owner, stakePools)
	return stakePools, nil
}

// EarningAmountTimeAfterTimestampMs 某个stake产生收益的时间
func (sp *StakePool) EarningAmountTimeAfterTimestampMs(timestamp int64, stake Stake, validatorStateInfo ValidatorState) int64 {
	rewardEpoch := stake.RequestEpoch + 2
	leftEpoch := rewardEpoch - validatorStateInfo.Epoch

	ranTime := timestamp - validatorStateInfo.EpochStartTimestampMs
	leftTime := validatorStateInfo.EpochDurationMs*leftEpoch - ranTime
	return leftTime
}

// TotalStakeAtValidator 某账户在某验证节点中的所有质押
func (sp *StakePool) TotalStakeAtValidator(validator, owner string, useCache bool) (sui *decimal.Decimal, err error) {
	var stakePools []StakePool
	if useCache {
		if cachedStakes, ok := cachedStakePoolsMap.Load(owner); ok {
			if v, ok := cachedStakes.([]StakePool); ok {
				stakePools = v
			}
		}
	}
	if stakePools == nil {
		stakePools, err = sp.GetStakePools(owner, useCache)
		if err != nil {
			return nil, err
		}
	}

	total := decimal.Zero
	for _, stakePool := range stakePools {
		if types.IsSameStringAddress(stakePool.ValidatorAddress, validator) {
			for _, v := range stakePool.Stakes {
				total = total.Add(v.Principal)
			}
			break
		}
	}
	return &total, nil
}

// AverageApy 指定质押池集合的平均apy
func (sp *StakePool) AverageApy(stakePools []StakePool) float64 {
	if stakePools == nil {
		return 0
	}
	totalApy := float64(0)
	added := map[string]bool{}
	for _, pool := range stakePools {
		if added[pool.ValidatorAddress] != true && pool.Validator != nil {
			totalApy = totalApy + pool.Validator.APY
			added[pool.ValidatorAddress] = true
			continue
		}
	}

	if count := len(added); count > 0 {
		return totalApy / float64(count)
	} else {
		return 0
	}
}

// BuildUnSignAddDelegationTx 生成未签名交易-新增质押
func (sp *StakePool) BuildUnSignAddDelegationTx(owner string, amount decimal.Decimal, coinType string, validatorAddress string) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}

	validator, err := types.NewAddressFromHex(validatorAddress)
	if err != nil {
		return
	}

	tokenServcie := NewTokenMain(sp.chain)
	coinIds, gasId, err := tokenServcie.PickCoinIdsAndGasId(owner, coinType, amount, decimal.NewFromInt(config.MaxGasBudget))
	if err != nil {
		return
	}

	cli, err := sp.chain.Client()
	if err != nil {
		return
	}
	return cli.RequestAddStake(context.Background(), *signer, coinIds, amount.BigInt().Uint64(), *validator, gasId, config.MaxGasBudget)
}

// BuildUnSignWithdrawDelegationTx 生成未签名交易-提现交易
func (sp *StakePool) BuildUnSignWithdrawDelegationTx(owner, stakeId, coinType string) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ownerByte, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	stakeIdByte, err := types.NewHexData(stakeId)
	if err != nil {
		return
	}
	tokenServcie := NewTokenMain(sp.chain)
	gasId, err := tokenServcie.PickMaxCoinId(owner, coinType, decimal.NewFromInt(config.MaxGasBudget))
	if err != nil {
		return
	}

	cli, err := sp.chain.Client()
	if err != nil {
		return
	}
	return cli.RequestWithdrawStake(context.Background(), *ownerByte, *stakeIdByte, &gasId, config.MaxGasBudget)
}
