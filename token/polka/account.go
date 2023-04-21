package polka

import (
	"encoding/json"
	"errors"
	"mytoken/core/base"
	"mytoken/core/lib/sublib/subkey"
	"mytoken/core/lib/sublib/subkey/sr25519"
	"mytoken/core/lib/sublib/substrate-rpc-client/signature"
	"mytoken/core/lib/sublib/substrate-rpc-client/types"
	"mytoken/core/utils/hexutil"
	"mytoken/token/polka/models"
)

var DefaultAccount = Account{}

//sr25519
// 加密涉及三部分：seed key nonce
// seed为密钥或者助记词
// key 可以理解成是一个密码，签名时候会用到，该key为seed使用512hash生成，因此不可逆转为seed
// nonce 512hash生成key时产生的
// 生成的keystore，核心信息在encoded中，是标准的sr25519格式，可以获取到公钥、nonce、key，在keystore中，无法获取到私钥seed
// seed通过hash512可转换为key，来验证私钥的正确性，key可以用来签名文件等
// 密钥助记词以及keystore均可获取到key，但keystore无法获取到seed

type Account struct {
	keypair  *signature.KeyringPair //助记词和私钥导入会用到
	keystore *models.Keystore       // keystore使用，该方式无法导出正常私钥
	Address  string
	Network  int
}

// NewAccountWithMnemonic 只有助记词可以导出私钥
func NewAccountWithMnemonic(mnemonic string, network int) (*Account, error) {
	if len(mnemonic) == 0 {
		return nil, errors.New("mnemonic is empty")
	}
	keyringPair, err := signature.KeyringPairFromSecret(mnemonic, uint16(network))
	if err != nil {
		return nil, err
	}

	return &Account{
		keypair: &keyringPair,
		Address: keyringPair.Address,
		Network: network,
	}, nil
}

func NewAccountWithPrivateKey(prikey string, network int) (*Account, error) {
	//旧版
	/*seed, err := utils.HexDecodeString(prikey)
	if err != nil {
		return nil, err
	}
	kyr, err := sr25519.Scheme{}.FromSeed(seed)
	if err != nil {
		return nil, err
	}

	ss58Address := kyr.SS58Address(uint16(network))
	var pk = kyr.Public()
	keypair := signature.KeyringPair{
		URI:       prikey,
		Address:   ss58Address,
		publicKey: pk,
	}*/

	//简化版
	if len(prikey) == 0 {
		return nil, errors.New("prikey is empty")
	}
	keyringPair, err := signature.KeyringPairFromSecret(prikey, uint16(network))
	if err != nil {
		return nil, err
	}
	return &Account{
		keypair: &keyringPair,
		Address: keyringPair.Address,
		Network: network,
	}, nil
}

// NewAccountWithKeystore
// keystore解析出的私钥，不能用于恢复助记词导出的私钥等信息
// 解析出的keystore中能够获取到公钥，该公钥转换出的地址和keystore中明文标记的地址一致
func NewAccountWithKeystore(keystoreString, password string, network int) (*Account, error) {
	var keyStore models.Keystore
	err := json.Unmarshal([]byte(keystoreString), &keyStore)
	if err != nil {
		return nil, err
	}
	pub, err := keyStore.CheckPassword(password)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	address, err := account.EncodePublicKeyToAddress(pub, network)
	if err != nil {
		return nil, err
	}
	account.Address = address
	account.Network = network
	account.keystore = &keyStore
	return account, nil
}

// Keystore 暂未实现
func (a *Account) Keystore() (jsonKeystore string, err error) {
	if a.keypair != nil {
		return "", err
	}
	if a.keystore != nil {
		jsonKeystoreByte, e := json.Marshal(a.keystore)
		if e != nil {
			err = e
			return
		}
		return string(jsonKeystoreByte), nil
	}
	return "", errors.New("no wallet")
}

// PrivateKey @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	if a.keypair == nil {
		return nil, errors.New("not support export private key")
	}

	scheme := sr25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, a.keypair.URI)
	if err != nil {
		return nil, err
	}
	return kyr.Seed(), nil
}

// PrivateKeyHex @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	data, err := a.PrivateKey()
	if err != nil {
		return "", err
	}
	return hexutil.HexEncodeToString(data), nil
}

// PublicKey @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.keypair.PublicKey
}

// PublicKeyHex @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return hexutil.HexEncodeToString(a.keypair.PublicKey)
}

func (a *Account) EncodePublicKeyToAddress(publicKey string, network int) (string, error) {
	pubByte := hexutil.HexToBytes(publicKey)
	address := subkey.SS58Encode(pubByte, uint16(network))
	if len(address) == 0 {
		return "", errors.New("address is empty")
	}
	return address, nil
}

func (a *Account) DecodeAddressToPublicKeyHex(address string) (string, error) {
	_, pkByte, err := subkey.SS58Decode(address)
	if err != nil {
		return "", err
	}
	publicKey := hexutil.HexEncodeToString(pkByte)

	return publicKey, nil
}

func (a *Account) DecodeAddressToMultiAddress(address string) (types.MultiAddress, error) {
	publicKey, err := a.DecodeAddressToPublicKeyHex(address)
	if err != nil {
		return types.MultiAddress{}, err
	}
	return types.NewMultiAddressFromHexAccountID(publicKey)
}

func (a *Account) DecodeAddressToPublicKey(address string) ([]byte, error) {
	_, pkByte, err := subkey.SS58Decode(address)
	if err != nil {
		return nil, err
	}
	return pkByte, err
}

func (a *Account) IsValidAddress(address string) bool {
	_, err := a.DecodeAddressToPublicKeyHex(address)
	return err == nil
}

func (a *Account) Sign(message []byte, password string) (data []byte, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if a.keypair != nil {
		return signature.Sign(message, a.keypair.URI)
	} else if a.keystore != nil {
		return a.keystore.Sign(message, password)
	}
	return nil, errors.New("no wallet")
}

func (a *Account) SignHex(messageHex string, password string) (string, error) {
	bytes, err := hexutil.HexDecodeString(messageHex)
	if err != nil {
		return "", err
	}
	signed, err := a.Sign(bytes, password)
	if err != nil {
		return "", err
	}
	return hexutil.HexEncodeToString(signed), nil
}

/*
// 内置账号，主要用来给用户未签名的交易签一下名
// 然后给用户去链上查询手续费，保护用户资产安全
func mockAccount() *Account {
	mnemonic := "infant carbon above canyon corn collect finger drip area feature mule autumn"
	a, _ := NewAccountWithMnemonic(mnemonic, 44)
	return a
}
*/
