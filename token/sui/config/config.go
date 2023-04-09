package config

const (
	PerObjectMaxGasForPay = 1000000
	MaxGasForTransfer     = 2000000

	MaxGasBudget = 12000000

	DevNetFaucetUrl  = "https://faucet.devnet.sui.io/gas"
	TestNetFaucetUrl = "https://faucet.testnet.sui.io/gas"
	DevNetRpcUrl     = "https://fullnode.devnet.sui.io"

	//https://fullnode.testnet.sui.io
	//https://sui-rpc-pt.testnet-pride.com
	//https://rpc-sui-testnet.cosmostation.io
	//https://testnet.artifact.systems/sui
	//https://rpc-testnet.suiscan.xyz
	//https://sui-testnet.brightlystake.com
	//https://sui-rpc.testnet.lgns.net
	//https://sui-testnet-rpc.bartestnet.com
	//https://sui-testnet-endpoint.blockvision.org
	//https://sui-testnet-rpc.allthatnode.com
	//https://sui-testnet-rpc-germany.allthatnode.com
	TestnetRpcUrl = "https://fullnode.testnet.sui.io"

	//我的质押池默认地址
	MyValidatorAddress = "0x"
	MyValidatorNameRex = `(?i)^My[ ._-]*Validator$`
)
