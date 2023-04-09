package sui

import (
	"context"
	"encoding/json"
	"mytoken/core/base"
	"mytoken/token/sui/client"
	"mytoken/token/sui/config"
)

var defaultChain = TestnetChain()

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string
}

func NewChainWithRpcUrl(rpcUrl string) *Chain {
	return &Chain{RpcUrl: rpcUrl}
}

func (c *Chain) Client() (*client.Client, error) {
	if c.rpcClient != nil {
		return c.rpcClient, nil
	}
	var err error
	c.rpcClient, err = client.Dial(c.RpcUrl)
	return c.rpcClient, err
}

// DevnetChain
//
//	@Description: 开发网
//	@return *Chain
func DevnetChain() *Chain {
	return NewChainWithRpcUrl(config.DevNetRpcUrl)
}

// TestnetChain
//
//	@Description: 测试网
//	@return *Chain
func TestnetChain() *Chain {
	return NewChainWithRpcUrl(config.TestnetRpcUrl)
}

func (c *Chain) GetLatestSysState() (chainState json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return nil, err
	}
	return cli.GetLatestSuiSystemState(context.Background())
}

/*func (c *Chain) BaseMoveCall(address, packageId, module, funcName string, typArgs []string, arg []any) (json.RawMessage, error) {
	cli, err := c.Client()
	if err != nil {
		return nil, err
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}
	packageIdHex, err := types.NewHexData(packageId)
	if err != nil {
		return nil, err
	}
	suiToken := NewTokenMain(c)
	coins, err := suiToken.GetCoins(address, suiToken.rType.ShortString())
	if err != nil {
		return nil, err
	}
	coin, err := coins.pickOneCoinForAmount(config.GasBudget)
	if err != nil {
		return nil, err
	}
	tx, err := cli.MoveCall(
		context.Background(),
		*addr,
		*packageIdHex,
		module,
		funcName,
		typArgs,
		arg,
		&coin.Reference.ObjectId,
		config.GasBudget,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}*/
