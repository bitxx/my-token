package bcs

import (
	"math/big"
	"testing"
)

func TestNewUint128FromUint64(t *testing.T) {
	expected := &big.Int{}
	expected.Add(big.NewInt(0).Lsh(big.NewInt(96), 64), big.NewInt(50))
	result := NewUint128FromUint64(50, 96).Big()
	if result.Cmp(expected) != 0 {
		t.Fatalf("want: %s, got: %s", expected.String(), result.String())
	}
}

func TestNewBigIntFromUint64(t *testing.T) {
	var expected uint64 = (1 << 63) + 12345
	result := NewBigIntFromUint64(expected)

	if result.Uint64() != expected {
		t.Fatalf("want: %d, got: %s", expected, result.String())
	}
}

func TestNewUint128(t *testing.T) {
	s := "1770887431076116955186"
	expected := &big.Int{}
	expected.Add(big.NewInt(0).Lsh(big.NewInt(96), 64), big.NewInt(50))

	result, err := NewUint128(s)
	if err != nil {
		t.Fatal(err)
	}

	if result.Big().Cmp(expected) != 0 {
		t.Fatalf("want: %s, got: %s", expected.String(), result.String())
	}
}
