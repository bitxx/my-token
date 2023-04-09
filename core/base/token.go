package base

type TokenInfo struct {
	Name    string
	Symbol  string
	Decimal int16
}

type Token interface {
	Chain() Chain

	TokenInfo() (*TokenInfo, error)

	Balance(address string) (*Balance, error)
}
