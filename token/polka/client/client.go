package client

import (
	"errors"
	gsrpc "mytoken/core/lib/sublib/substrate-rpc-client"
	"mytoken/core/lib/sublib/substrate-rpc-client/types"
)

type Client struct {
	Api      *gsrpc.SubstrateAPI
	Metadata *types.Metadata
	rpcUrl   string
}

// Dial 调取
func Dial(rpcUrl string) (*Client, error) {
	if rpcUrl == "" {
		return nil, errors.New("rpcUrl is empty")
	}
	api, err := gsrpc.NewSubstrateAPI(rpcUrl)
	if err != nil {
		return nil, err
	}
	c := &Client{
		rpcUrl: rpcUrl,
		Api:    api,
	}
	err = c.LoadMetadata(true)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// LoadMetadata
//
//	@Description: 更新metadata
//	@receiver c
//	@param needUpdate 是否需要更新metadata
//	@return error
func (c *Client) LoadMetadata(needUpdate bool) error {
	if needUpdate || c.Metadata == nil {
		meta, err := c.Api.RPC.State.GetMetadataLatest()
		if err != nil {
			return err
		}
		c.Metadata = meta
	}
	return nil
}
