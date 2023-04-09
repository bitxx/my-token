package sui

import (
	"mytoken/token/sui/config"
	"mytoken/token/testcase"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFaucet 水龙头领取，根据调试发现，仅支持开发环境领取，测试环境是不支持领取的
func TestFaucet(t *testing.T) {
	res, err := FaucetFundAccount(testcase.Accounts1.Sui.Address, config.DevNetFaucetUrl)
	require.Nil(t, err)
	t.Log("res = ", res)
}
