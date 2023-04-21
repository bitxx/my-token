package sui

import (
	"context"
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/token/sui/client"
)

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string
	ScanUrl   string
}

func NewChain(rpcUrl, scanUrl string) *Chain {
	return &Chain{RpcUrl: rpcUrl, ScanUrl: scanUrl}
}

func (c *Chain) Client() (*client.Client, error) {
	if c.RpcUrl == "" {
		return nil, errors.New("rpcUrl is empty")
	}
	if c.rpcClient != nil {
		return c.rpcClient, nil
	}
	var err error
	c.rpcClient, err = client.Dial(c.RpcUrl)
	return c.rpcClient, err
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
