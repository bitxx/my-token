package sui

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
	"testing"
	"time"
)

const (
	stakeOwnerTestAddr = "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"
	validatorTestAddr  = "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
	stakeId            = "0x5cdb23dacf54329660467b900a2598bb796353fa"
)

func TestGetStakePools(t *testing.T) {
	stakePoolsService := NewStakePool(TestnetChain())

	stakePools, err := stakePoolsService.GetStakePools(stakeOwnerTestAddr, true)
	require.Nil(t, err)
	for _, ss := range stakePools {
		t.Log("stakingPool:", ss.StakingPool, " validatorAddr:", ss.Validator.Address)
		for i, s := range ss.Stakes {
			t.Log("   --- stakeId[", i, "]：", s.StakeId)
		}
	}
}

func TestStakeEarningTimems(t *testing.T) {
	validatorState := ValidatorState{
		Epoch:                 9,
		EpochDurationMs:       86400000,
		EpochStartTimestampMs: 1680760906723,
	}

	ti := validatorState.EarningAmountTimeAfterTimestampMs(time.Now().UnixMilli())
	t.Log(ti)

	stake := Stake{
		RequestEpoch: 8,
	}
	stakePoolsService := NewStakePool(TestnetChain())
	ti2 := stakePoolsService.EarningAmountTimeAfterTimestampMs(time.Now().UnixMilli(), stake, validatorState)
	t.Log(ti2)
}

func TestTotalStakeAtValidator(t *testing.T) {
	stakePoolsService := NewStakePool(TestnetChain())
	total, err := stakePoolsService.TotalStakeAtValidator(validatorTestAddr, stakeOwnerTestAddr, true)
	require.Nil(t, err)
	t.Log(total)
}

func TestAverageApyStake(t *testing.T) {
	stakePoolsService := NewStakePool(TestnetChain())
	stakePools, err := stakePoolsService.GetStakePools(stakeOwnerTestAddr, true)
	require.Nil(t, err)
	apy := stakePoolsService.AverageApy(stakePools)
	t.Log("APY：", apy)
}

func TestAddDelegation(t *testing.T) {
	account := M1Account(t)
	stakePoolsService := NewStakePool(TestnetChain())

	//未签名交易
	unsignTx, err := stakePoolsService.BuildUnSignAddDelegationTx(account.Address, decimal.NewFromInt(1000), types.SuiCoinType, stakeOwnerTestAddr)
	require.Nil(t, err)

	//签名并提交交易
	transferService := NewTransaction(defaultChain)
	result, err := transferService.SignAndSendTx(account, unsignTx)
	require.Nil(t, err)
	t.Log(string(result))
}

func TestWithdrawDelegation(t *testing.T) {
	account := M1Account(t)
	stakePoolsService := NewStakePool(TestnetChain())
	//未签名交易
	unsignTx, err := stakePoolsService.BuildUnSignWithdrawDelegationTx(account.Address, stakeId, types.SuiCoinType)
	require.Nil(t, err)

	//签名并提交交易
	transferService := NewTransaction(defaultChain)
	result, err := transferService.SignAndSendTx(account, unsignTx)
	require.Nil(t, err)
	t.Log(string(result))
}
