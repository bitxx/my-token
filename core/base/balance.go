package base

import "mytoken/core/lib/decimal"

type Balance struct {
	Total  decimal.Decimal
	Usable decimal.Decimal
}
