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
	Phala       AccountCase
	Polka2      AccountCase
	Polka44     AccountCase
	Solana      AccountCase
	Aptos       AccountCase
	Sui         AccountCase
	Starcoin    AccountCase
}

var M1 = "current cage icon tattoo suffer bleak adapt orange fix remind million educate"
var M2 = "regular spin recall duck calm ring awesome raven letter model faint favorite"
var Mterra = "fine garment lab pigeon card side skirt box pink profit tell tourist"

var Accounts1 = AccountGroup{
	BtcMainnet:  AccountCase{M1, "", "", "", ""},
	BtcSignet:   AccountCase{M1, "", "", "", ""},
	Cosmos:      AccountCase{M1, "", "", "", ""},
	Terra:       AccountCase{Mterra, "", "", "", ""},
	DogeMainnet: AccountCase{M1, "", "", "", ""},
	DogeTestnet: AccountCase{M1, "", "", "", ""},
	Ethereum:    AccountCase{M1, "", "", "", ""},
	Phala:       AccountCase{M1, "44LYsBdxwSrmMRSYd2oLTVJsZ2kPVGW8p4TyQjnzhhzLkdge", "0xa48e40754f104a3e10422142d492fe73c5f9394c3e2103adf49f7e7a51dcf602", "0xcc23d6ecda194611c5eaaf7639f826e85a626521a53994689bdd1e8d0110c066", "{\"encoded\":\"osafYPwTSQXej8P/e82vW+ulBy5IYMqbJSMDJCzJjRYAgAAAAQAAAAgAAAAjTJ4GOMT3ubIwevG6ebGvNdA0dRjYuR9fmermW9uFWUiFks1PqnB6tgEBJUguWn91NYxOgL2pWu+3dWmTEuhGSHAtYCMpKyWSQp9ZG/ckOxld907XGtWO5V8Vs/Cne8ytczBJl0bjOsy/QCNxlErwwtFFeb3CvYpIfnJxKo3rW54YlyucIOOpSJtavd2+ncd6IhcCFXI754weOvaO\",\"encoding\":{\"content\":[\"pkcs8\",\"sr25519\"],\"type\":[\"scrypt\",\"xsalsa20-poly1305\"],\"version\":\"3\"},\"address\":\"5FnTy68oYwMVE8J7VRMczJVyJrTv9U4P5M7YnwWcFy8C8Xck\",\"meta\":{\"genesisHash\":\"\",\"name\":\"testM1\",\"whenCreated\":1681128297037}}"},
	Polka2:      AccountCase{M1, "", "", "", ""},
	Polka44:     AccountCase{M1, "", "", "", ""},
	Solana:      AccountCase{M1, "", "", "", ""},
	Aptos:       AccountCase{M1, "", "", "", ""},
	Sui:         AccountCase{M1, "0x9f63286a92e97d1410558236b55139d0fbe728764e30f562fd0e400405bb452a", "0x3fa31aa81777aadfdbdea05b5046859bb06ebd3c9622c1c2209a1cb332e93109", "0x736e58a1ecdc034b8a989e8bde2d8370a9ccfcec312c2a06c6129f8ad89860a1", "AHNuWKHs3ANLipiei94tg3CpzPzsMSwqBsYSn4rYmGCh"},
	Starcoin:    AccountCase{M1, "", "", "", ""},
}

var Accounts2 = AccountGroup{
	Sui:   AccountCase{M2, "0x1eb82d5b4ea515ba25061b4f682e1fa30fecd07ffb76a83f5a3067ae74ce8c45", "0xf6b56fe07ddb7bf595a2e1d3a87e3cea6e092eca13728940225e527fc557a3fa", "0x2b96925dafd96caf848d30e1ee83cc45160b99e45660eed06a24ee406dc5e110", "ACuWkl2v2WyvhI0w4e6DzEUWC5nkVmDu0Gok7kBtxeEQ"},
	Phala: AccountCase{M2, "44aCzJJpcAc7RYRociX2YGo1ktfwSGXhQUrcS1mJLbLdF49i", "0xaef8c63f2fa493bc7de9614bd9fe8055c52c4f6ca3d9d4ed639d102eb1c26a42", "0x00cf504756173bba8b57d51ed89a7108863f97f4e04f7386c9a262c0f7d1c45b", ""},
}
