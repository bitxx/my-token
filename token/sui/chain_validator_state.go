package sui

import (
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/models"
	"mytoken/token/sui/types"
	"regexp"
)

var cachedValidatorState *ValidatorState

type Validator struct {
	Address         string          `json:"address"`
	Name            string          `json:"name"`
	Desc            string          `json:"desc"`
	ImageUrl        string          `json:"imageUrl"`
	ProjectUrl      string          `json:"projectUrl"`
	APY             float64         `json:"apy"`
	Commission      decimal.Decimal `json:"commission"`
	TotalStaked     decimal.Decimal `json:"totalStaked"`
	DelegatedStaked string          `json:"delegatedStaked"`
	SelfStaked      string          `json:"selfStaked"`
	TotalRewards    decimal.Decimal `json:"totalRewards"`
	GasPrice        decimal.Decimal `json:"gasPrice"`
}

type ValidatorState struct {
	chain *Chain
	// The current epoch in Sui. An epoch takes approximately 24 hours and runs in checkpoints.
	Epoch decimal.Decimal `json:"epoch"`
	// Array of `Validator` elements
	Validators []Validator `json:"validators"`

	// The amount of all tokens staked in the Sui Network.
	TotalStaked decimal.Decimal `json:"totalStaked"`
	// The amount of rewards won by all Sui validators in the last epoch.
	TotalRewards          decimal.Decimal `json:"lastEpochReward"`
	EpochStartTimestampMs decimal.Decimal `json:"epochStartTimestampMs"`
	EpochDurationMs       decimal.Decimal `json:"epochDurationMs"`
}

func NewValidatorState(chain *Chain) *ValidatorState {
	return &ValidatorState{chain: chain}
}

// IsMyValidator 判断是否为自己指定的validator
func (vs *ValidatorState) IsMyValidator(validator *Validator) bool {
	if vs == nil || vs.chain == nil {
		return false
	}
	reg := regexp.MustCompile(config.MyValidatorNameRex)
	nameMatched := reg.MatchString(validator.Name)
	return nameMatched || validator.Address == config.MyValidatorAddress
}

func (vs *ValidatorState) TotalActiveValidatorState(useCache bool) (validatorState *ValidatorState, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if vs == nil || vs.chain == nil {
		return nil, errors.New("param is error")
	}
	if useCache && cachedValidatorState != nil {
		return cachedValidatorState, nil
	}

	sysStateJson, err := vs.chain.GetLatestSysState()
	if err != nil {
		return nil, err
	}

	sysState := models.SysState{}
	err = json.Unmarshal(sysStateJson, &sysState)
	if err != nil {
		return nil, err
	}

	totalRewards := decimal.Decimal{}
	var validators []Validator
	for _, activeValidator := range sysState.ActiveValidators {
		validator := Validator{
			Address:         activeValidator.SuiAddress,
			Name:            activeValidator.Name,
			Desc:            activeValidator.Description,
			ImageUrl:        activeValidator.ImageURL,
			ProjectUrl:      activeValidator.ProjectURL,
			APY:             activeValidator.CalculateAPY(sysState.Epoch),
			Commission:      activeValidator.CommissionRate,
			SelfStaked:      "--",
			DelegatedStaked: "--",
			TotalStaked:     activeValidator.StakingPoolSuiBalance,
			TotalRewards:    activeValidator.RewardsPool,
			GasPrice:        activeValidator.GasPrice,
		}
		if vs.IsMyValidator(&validator) {
			//自己过滤需要优先显示的池子排队在前面
			validators = append([]Validator{validator}, validators...)
		} else {
			validators = append(validators, validator)
		}
		totalRewards = totalRewards.Add(validator.TotalRewards)
	}

	res := &ValidatorState{
		Epoch:                 sysState.Epoch,
		Validators:            validators,
		TotalStaked:           sysState.TotalStake,
		TotalRewards:          totalRewards,
		EpochDurationMs:       sysState.EpochDurationMs,
		EpochStartTimestampMs: sysState.EpochStartTimestampMs,
	}

	cachedValidatorState = res
	return res, nil
}

func (vs *ValidatorState) GetValidator(address string, useCache bool) (v *Validator, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	validatorState, err := vs.TotalActiveValidatorState(useCache)
	if err != nil {
		return nil, err
	}

	for _, v := range validatorState.Validators {
		if types.IsSameStringAddress(address, v.Address) {
			return &v, nil
		}
	}
	return nil, errors.New("not found")
}

// EarningAmountTimeAfterTimestampMs 从当前周期的当前时间，到下一个周期结束，还剩多少时间
// 可以使用：time.Now().UnixMilli()
func (vs *ValidatorState) EarningAmountTimeAfterTimestampMs(timestamp int64) decimal.Decimal {
	ranTime := decimal.NewFromInt(timestamp).Sub(vs.EpochStartTimestampMs)
	leftTime := vs.EpochDurationMs.Mul(decimal.NewFromInt(2)).Sub(ranTime)
	return leftTime
}
