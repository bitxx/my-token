package base

type Chain interface {
	//FetchTransactionDetail(hash string) (*TransactionDetail, error)

	// FetchTransactionStatus Fetch transaction status through transaction hash
	//FetchTransactionStatus(hash string) TransactionStatus

	// BatchFetchTransactionStatus Batch fetch the transaction status, the hash list and the return value,
	// which can only be passed as strings separated by ","
	// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
	// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
	//BatchFetchTransactionStatus(hashListString string) string

	// -----------------------------
	// polka
	// GetSignDataFromChain(t *Transaction, walletAddress string) ([]byte, error)

	// EstimateFeeForTransaction(transaction *Transaction) (s string, err error)

	// FetchScriptHashForMiniX(transferTo, amount string) (*MiniXScriptHash, error)
}
