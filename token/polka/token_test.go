package polka

import (
	"github.com/stretchr/testify/require"
	"mytoken/core/lib/decimal"
	"mytoken/token/testcase"
	"testing"
)

func TestBalance(t *testing.T) {
	token := NewToken(defaultChain)
	b, err := token.BalanceOf(testcase.Accounts2.Phala.Address)
	require.Nil(t, err)

	t.Log(b.Total)
	t.Log(b.Usable)
}

func TestBuildUnSignTokenTransferTx(t *testing.T) {
	tok := NewToken(defaultChain)
	extrinsic, err := tok.BuildUnSignTokenTransferTx("5TE1T7Znw5eaDqzpxCm8MBocbZLbbeeZ4GcnsRN37s2amd5s", decimal.NewFromInt(99))
	require.Nil(t, err)
	t.Log("version：", extrinsic.Version)
}
