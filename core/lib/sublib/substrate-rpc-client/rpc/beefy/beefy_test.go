package beefy

import (
	"os"
	"testing"

	"mytoken/core/lib/sublib/substrate-rpc-client/client"
	"mytoken/core/lib/sublib/substrate-rpc-client/config"
)

var testBeefy Beefy

func TestMain(m *testing.M) {
	cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	testBeefy = NewBeefy(cl)
	os.Exit(m.Run())
}
