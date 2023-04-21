package sui

import (
	"testing"
)

const testHash = "CfnWfh8BeZwFyvUBhXjaN4QmBVmKzMGxRu4xABhrKasc"

func TestTransactionDetail(t *testing.T) {
	transaction := NewTransaction(defaultChain)
	resp, err := transaction.TransactionDetail(testHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp))
}
