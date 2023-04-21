package models

import (
	"fmt"
	"mytoken/core/lib/decimal"
)

type Balance struct {
	CoinType        string                                    `json:"coinType"`
	CoinObjectCount uint64                                    `json:"coinObjectCount"`
	TotalBalance    decimal.Decimal                           `json:"totalBalance"`
	LockedBalance   map[SafeSuiBigInt[uint64]]decimal.Decimal `json:"lockedBalance"`
}

type Supply struct {
	Value SafeSuiBigInt[uint64] `json:"value"`
}

type SafeBigInt interface {
	~int64 | ~uint64
}

func NewSafeSuiBigInt[T SafeBigInt](num T) SafeSuiBigInt[T] {
	return SafeSuiBigInt[T]{
		data: num,
	}
}

type SafeSuiBigInt[T SafeBigInt] struct {
	data T
}

func (s *SafeSuiBigInt[T]) UnmarshalText(data []byte) error {
	return s.UnmarshalJSON(data)
}

func (s *SafeSuiBigInt[T]) UnmarshalJSON(data []byte) error {
	num := decimal.NewFromInt(0)
	err := num.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	if num.BigInt().IsInt64() {
		s.data = T(num.BigInt().Int64())
		return nil
	}

	if num.BigInt().IsUint64() {
		s.data = T(num.BigInt().Uint64())
		return nil
	}
	return fmt.Errorf("json data [%s] is not T", string(data))
}

func (s SafeSuiBigInt[T]) MarshalJSON() ([]byte, error) {
	return decimal.NewFromInt(int64(s.data)).MarshalJSON()
}

func (s SafeSuiBigInt[T]) Int64() int64 {
	return int64(s.data)
}

func (s SafeSuiBigInt[T]) Uint64() uint64 {
	return uint64(s.data)
}

func (s *SafeSuiBigInt[T]) Decimal() decimal.Decimal {
	return decimal.NewFromInt(s.Int64())
}
