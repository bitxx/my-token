package client

import (
	"context"
	"encoding/json"
	"mytoken/token/sui/types"
)

func (c *Client) GetLatestSuiSystemState(ctx context.Context) (json.RawMessage, error) {
	return c.CallContext(ctx, getLatestSuiSystemState)
}

func (c *Client) GetStakePools(ctx context.Context, owner types.Address) (json.RawMessage, error) {
	return c.CallContext(ctx, getStakes, owner)
}

func (c *Client) GetStakesByIds(ctx context.Context, stakedSuiIds []types.ObjectId) (json.RawMessage, error) {
	return c.CallContext(ctx, getStakesByIds, stakedSuiIds)
}

func (c *Client) RequestAddStake(ctx context.Context, signer types.Address, coins []types.ObjectId, amount uint64, validator types.Address, gas *types.ObjectId, gasBudget uint64) (json.RawMessage, error) {
	return c.CallContext(ctx, requestAddStake, signer, coins, amount, validator, gas, gasBudget)
}

func (c *Client) RequestWithdrawStake(ctx context.Context, signer types.Address, stakedSuiId types.ObjectId, gas *types.ObjectId, gasBudget uint64) (json.RawMessage, error) {
	return c.CallContext(ctx, requestWithdrawStake, signer, stakedSuiId, gas, gasBudget)
}
