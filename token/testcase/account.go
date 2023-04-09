package testcase

type AccountCase struct {
	Mnemonic   string
	Address    string
	PublicKey  string
	PrivateKey string
	KeyStore   string
}

type AccountGroup struct {
	BtcMainnet  AccountCase
	BtcSignet   AccountCase
	Cosmos      AccountCase
	Terra       AccountCase
	DogeMainnet AccountCase
	DogeTestnet AccountCase
	Ethereum    AccountCase
	Polka0      AccountCase
	Polka2      AccountCase
	Polka44     AccountCase
	Solana      AccountCase
	Aptos       AccountCase
	Sui         AccountCase
	Starcoin    AccountCase
}

var M1 = "助记词"
var M2 = "助记词"
var Mterra = "助记词"

var Accounts1 = AccountGroup{
	BtcMainnet:  AccountCase{M1, "", "", "", ""},
	BtcSignet:   AccountCase{M1, "", "", "", ""},
	Cosmos:      AccountCase{M1, "", "", "", ""},
	Terra:       AccountCase{Mterra, "", "", "", ""},
	DogeMainnet: AccountCase{M1, "", "", "", ""},
	DogeTestnet: AccountCase{M1, "", "", "", ""},
	Ethereum:    AccountCase{M1, "", "", "", ""},
	Polka0:      AccountCase{M1, "", "", "", ""},
	Polka2:      AccountCase{M1, "", "", "", ""},
	Polka44:     AccountCase{M1, "", "", "", ""},
	Solana:      AccountCase{M1, "", "", "", ""},
	Aptos:       AccountCase{M1, "", "", "", ""},
	Sui:         AccountCase{M1, "地址", "公钥", "私钥", "keystore"},
	Starcoin:    AccountCase{M1, "", "", "", ""},
}

var Accounts2 = AccountGroup{
	Sui: AccountCase{M2, "地址", "公钥", "私钥", "keystore"},
}

var EmptyMnemonic = AccountCase{Mnemonic: ""}
var ErrorMnemonic = AccountCase{Mnemonic: ""}
