package polka

import (
	"github.com/stretchr/testify/require"
	"mytoken/token/polka/config"
	"testing"
)

var defaultChain = TestnetChain()

func TestnetChain() *Chain {
	return NewChain(config.ProdnetRpcUrl)
}

func TestNewChain(t *testing.T) {
	_, err := defaultChain.Client()
	require.Nil(t, err)
}
