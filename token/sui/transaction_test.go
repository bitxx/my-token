package sui

import (
	"testing"
)

const testHash = "8RXvmEj8oQB4S1zLA5wJ5V6DGMfm6w2UQD9tYWCXQyCu"

func TestTransactionDetail(t *testing.T) {
	transaction := NewTransaction(DevnetChain())
	resp, err := transaction.TransactionDetail(testHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp))
}
