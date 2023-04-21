package polka

import (
	"mytoken/token/polka/client"
)

type Chain struct {
	RpcUrl    string
	ScanUrl   string
	RpcClient *client.Client
}

// NewChain
// @param rpcUrl will be used to get metadata, query balance, estimate fee, send signed tx.
// @param scanUrl will be used to query transaction details
func NewChain(rpcUrl string) *Chain {
	return &Chain{
		RpcUrl: rpcUrl,
	}
}

func (c *Chain) Client() (*client.Client, error) {
	if c.RpcClient != nil {
		return c.RpcClient, nil
	}
	var err error
	c.RpcClient, err = client.Dial(c.RpcUrl)
	return c.RpcClient, err
}
