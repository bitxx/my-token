package models

import (
	"mytoken/token/sui/types"
)

const (
	SignatureSchemeSerializedEd25519   byte = 0
	SignatureSchemeSerializedSecp256k1 byte = 1
)

type Transaction struct {
	TxBytes types.Base64Data `json:"txBytes"`
	Gas     []struct {
		ObjectID types.ObjectId `json:"objectId"`
		Version  int            `json:"version"`
		Digest   string         `json:"digest"`
	} `json:"gas"`
	InputObjects []struct {
		ImmOrOwnedMoveObject struct {
			ObjectID types.ObjectId `json:"objectId"`
			Version  int            `json:"version"`
			Digest   string         `json:"digest"`
		} `json:"ImmOrOwnedMoveObject"`
	} `json:"inputObjects"`
}

/*func (txn *Transaction) SignSerializedSigWith(privateKey ed25519.PrivateKey) (txBytes *types.Base64Data, signatures []types.Base64Data) {
	// IntentBytes See: sui/crates/shared-crypto/src/intent.rs
	// This is currently hardcoded with [IntentScope::TransactionData = 0, Version::V0 = 0, AppId::Sui = 0]
	signTx := bytes.NewBuffer([]byte{0, 0, 0})
	signTx.Write(txn.TxBytes.Data())
	message := blake2b.Sum256(signTx.Bytes())
	signatureData := bytes.NewBuffer([]byte{byte(SignatureSchemeSerializedEd25519)})
	signatureData.Write(ed25519.Sign(privateKey, message[:]))
	signatureData.Write(privateKey.Public().(ed25519.PublicKey))
	return &txn.TxBytes, []types.Base64Data{signatureData.Bytes()}
}*/
