package sui

import (
	"mytoken/token/sui/config"
)

var defaultChain = TestnetChain()

// DevnetChain
//
//	@Description: 开发网
//	@return *Chain
func DevnetChain() *Chain {
	return NewChain(config.DevNetRpcUrl)
}

// TestnetChain
//
//	@Description: 测试网
//	@return *Chain
func TestnetChain() *Chain {
	return NewChain(config.TestnetRpcUrl)
}
