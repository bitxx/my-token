package sui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mytoken/core/base"
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/config"
	"mytoken/token/sui/models"
	"mytoken/token/sui/types"
	"sort"
)

const (
	pickSmaller = iota // pick smaller coins to match amount
	pickBigger         // pick bigger coins to match amount
	pickByOrder        // pick coins by coins order to match amount
)

type Token struct {
	chain *Chain

	rType types.ResourceType
}

func NewTokenMain(chain *Chain) *Token {
	token, _ := NewToken(chain, types.SuiCoinType)
	return token
}

func NewToken(chain *Chain, tag string) (*Token, error) {
	token, err := types.NewResourceType(tag)
	if err != nil {
		return nil, err
	}
	return &Token{chain, *token}, nil
}

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) CoinType() string {
	return fmt.Sprintf("0x2::coin::Coin<%v>", t.rType.ShortString())
}

// TokenInfo token信息获取
func (t *Token) TokenInfo() (info json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	cli, err := t.chain.Client()
	if err != nil {
		return
	}
	info, err = cli.GetCoinMetadata(context.TODO(), t.rType.ShortString())
	return info, nil
}

// BalanceOf 根据地址获取地址余额
func (t *Token) BalanceOf(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	coins, err := t.GetCoins(address, t.rType.ShortString())
	if err != nil {
		return
	}
	total := decimal.Zero
	for _, coin := range coins {
		total = total.Add(coin.Balance)
	}
	return &base.Balance{
		Total:  total,
		Usable: total,
	}, nil
}

func (t *Token) GetCoins(address string, coinType string) (coins []types.Coin, err error) {
	cli, err := t.chain.Client()
	if err != nil {
		return
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	coinObjects, err := cli.GetCoins(context.TODO(), *addr, &coinType, nil, 200)
	if err != nil {
		return nil, err
	}
	page := models.Coin{}
	err = json.Unmarshal(coinObjects, &page)
	if err != nil {
		return nil, err
	}
	for _, coin := range page.Data {
		coins = append(coins, types.Coin{
			Balance:      coin.Balance,
			CoinType:     coin.CoinType,
			Digest:       coin.Digest,
			Version:      coin.Version,
			CoinObjectId: coin.CoinObjectID,
		})
	}
	// sort by balance descend
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Balance.Cmp(coins[j].Balance) > 0
	})
	return coins, nil
}

// BuildUnSignTokenTransferTx 生成从账户1转账指定余额到账户2的交易信息，未签名
func (t *Token) BuildUnSignTokenTransferTx(account *Account, receiverAddress string, amount, gasBudget decimal.Decimal) (txn json.RawMessage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if err != nil {
		return
	}
	coinIds, err := t.PickCoinIds(account.Address, t.rType.ShortString(), amount)
	if err != nil {
		return
	}
	if len(coinIds) <= 0 {
		err = errors.New("get coins error")
		return
	}

	cli, err := t.chain.Client()
	if err != nil {
		return
	}
	signer, _ := types.NewAddressFromHex(account.Address)
	recipient, err := types.NewAddressFromHex(receiverAddress)
	if len(coinIds) >= 2 {
		txn, err = cli.PayAllSui(context.Background(), *signer, *recipient, coinIds, decimal.NewFromInt(int64(len(coinIds)-1)).Mul(gasBudget))
	} else {
		txn, err = cli.TransferSui(context.Background(), *signer, *recipient, coinIds[0], amount, gasBudget)
	}
	if err != nil {
		return
	}
	return
}

// EstimateTransferFees
//
//	@Description: 估算账户1转账到账户2的手续费
//	@receiver t
//	@param account
//	@param receiverAddress
//	@param amount
//	@return f
//	@return err
func (t *Token) EstimateTransferFees(account *Account, receiverAddress string, amount decimal.Decimal) (f *decimal.Decimal, err error) {
	txn, err := t.BuildUnSignTokenTransferTx(account, receiverAddress, amount, decimal.NewFromInt(config.MaxGasForTransfer))
	if err != nil {
		return
	}
	transaction := &models.Transaction{}
	err = json.Unmarshal(txn, &transaction)
	if err != nil {
		return
	}
	return NewTransaction(t.chain).EstimateGasFee(transaction.TxBytes)
}

// PickCoinIdsAndGasId
// @Description: 支付amount时候，选出合适的coins和gas的id
// @receiver t
// @param owner
// @param coinType
// @param amount
// @param maxGasBudgetForStake
// @return coinIds
// @return gasId
func (t *Token) PickCoinIdsAndGasId(owner, coinType string, amount, gasAmount decimal.Decimal) (coinIds []types.ObjectId, gasId *types.ObjectId, err error) {
	if gasAmount.Cmp(decimal.Zero) <= 0 {
		return nil, nil, errors.New("gasAmount need >=0 ")
	}
	if amount.Cmp(decimal.Zero) <= 0 {
		return nil, nil, errors.New("amount need >=0 ")
	}

	coins, err := t.GetCoins(owner, coinType)
	if err != nil {
		return
	}
	needCoins, gasCoin, err := t.pickCoinsAndGas(coins, amount, gasAmount, pickBigger)
	for _, coin := range needCoins {
		coinIds = append(coinIds, coin.CoinObjectId)
	}
	gasId = &gasCoin.CoinObjectId
	return
}

// PickCoinIds 根据amount获取所有需要的 id
func (t *Token) PickCoinIds(owner, coinType string, amount decimal.Decimal) (coinIds []types.ObjectId, err error) {
	coins, err := t.GetCoins(owner, coinType)
	if err != nil {
		return
	}
	needCoins, err := t.pickCoins(coins, amount, pickBigger)
	for _, coin := range needCoins {
		coinIds = append(coinIds, coin.CoinObjectId)
	}
	return
}

// PickMaxCoinId 获取最大的一个coin id
func (t *Token) PickMaxCoinId(owner, coinType string, amount decimal.Decimal) (gasId types.ObjectId, err error) {
	coins, err := t.GetCoins(owner, coinType)
	if err != nil {
		return
	}

	gasCoin, err := t.pickOneCoinForAmount(coins, amount)
	if err != nil {
		return
	}
	gasId = gasCoin.CoinObjectId

	return
}

// pickCoinsAndGas 根据 amount、gas挑选出所需要的coins和gas
func (t *Token) pickCoinsAndGas(allCoins []types.Coin, amount decimal.Decimal, gasAmount decimal.Decimal, pickMethod int) ([]types.Coin, *types.Coin, error) {
	if gasAmount.Cmp(decimal.Zero) <= 0 {
		return nil, nil, errors.New("gasAmount need >=0 ")
	}
	if amount.Cmp(decimal.Zero) <= 0 {
		return nil, nil, errors.New("amount need >=0 ")
	}

	// find smallest to match gasAmount
	var gasCoin *types.Coin
	var selectIndex int
	for i := range allCoins {
		if allCoins[i].Balance.Cmp(gasAmount) < 0 {
			continue
		}

		if nil == gasCoin || gasCoin.Balance.Cmp(allCoins[i].Balance) > 0 {
			gasCoin = &allCoins[i]
			selectIndex = i
		}
	}
	if nil == gasCoin {
		return nil, nil, errors.New("coins not match request")
	}

	lastCoins := make([]types.Coin, 0, len(allCoins)-1)
	lastCoins = append(lastCoins, allCoins[0:selectIndex]...)
	lastCoins = append(lastCoins, allCoins[selectIndex+1:]...)
	pickCoins, err := t.pickCoins(lastCoins, amount, pickMethod)
	return pickCoins, gasCoin, err
}

func (t *Token) pickOneCoinForAmount(allCoins []types.Coin, amount decimal.Decimal) (*types.Coin, error) {
	for _, coin := range allCoins {
		if coin.Balance.Cmp(amount) >= 0 {
			return &coin, nil
		}
	}
	return nil, errors.New("no coin is enough to cover the gas")
}

// pickCoins 根据 amount挑选出所需要的coins
func (t *Token) pickCoins(allCoins []types.Coin, amount decimal.Decimal, pickMethod int) ([]types.Coin, error) {
	var sortedCoins []types.Coin
	if pickMethod == pickByOrder {
		sortedCoins = allCoins
	} else {
		sortedCoins = make([]types.Coin, len(allCoins))
		copy(sortedCoins, allCoins)
		sort.Slice(sortedCoins, func(i, j int) bool {
			if pickMethod == pickSmaller {
				return sortedCoins[i].Balance.Cmp(sortedCoins[j].Balance) < 0
			} else {
				return sortedCoins[i].Balance.Cmp(sortedCoins[j].Balance) >= 0
			}
		})
	}

	result := make([]types.Coin, 0)
	total := decimal.Zero
	for _, coin := range sortedCoins {
		result = append(result, coin)
		total = total.Add(coin.Balance)
		if total.Cmp(amount) >= 0 {
			return result, nil
		}
	}

	return nil, errors.New("coins not match request")
}
