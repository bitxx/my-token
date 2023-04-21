package schnorrkel

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"mytoken/core/lib/sublib/merlin"
	"mytoken/core/lib/sublib/ristretto255"
	"strings"
)

func challengeScalar(t *merlin.Transcript, msg []byte) *ristretto255.Scalar {
	scb := t.ExtractBytes(msg, 64)
	sc := ristretto255.NewScalar()
	sc.FromUniformBytes(scb)
	return sc
}

// https://github.com/w3f/schnorrkel/blob/718678e51006d84c7d8e4b6cde758906172e74f8/src/scalars.rs#L18
func divideScalarByCofactor(s []byte) []byte {
	l := len(s) - 1
	low := byte(0)
	for i := range s {
		r := s[l-i] & 0x07 // remainder
		s[l-i] >>= 3
		s[l-i] += low
		low = r << 5
	}

	return s
}

// NewRandomElement returns a random ristretto element
func NewRandomElement() (*ristretto255.Element, error) {
	e := ristretto255.NewElement()
	s := [64]byte{}
	_, err := rand.Read(s[:])
	if err != nil {
		return nil, err
	}

	return e.FromUniformBytes(s[:]), nil
}

// NewRandomScalar returns a random ristretto scalar
func NewRandomScalar() (*ristretto255.Scalar, error) {
	s := [64]byte{}
	_, err := rand.Read(s[:])
	if err != nil {
		return nil, err
	}

	ss := ristretto255.NewScalar()
	sc := ss.FromUniformBytes(s[:])
	if sc.Equal(ristretto255.NewScalar()) == 1 {
		return nil, errors.New("scalar generated was zero")
	}

	return sc, nil
}

// ScalarFromBytes returns a ristretto scalar from the input bytes
// performs input mod l where l is the group order
func ScalarFromBytes(b [32]byte) (*ristretto255.Scalar, error) {
	s := ristretto255.NewScalar()
	err := s.Decode(b[:])
	if err != nil {
		return nil, err
	}

	return s, nil
}

// HexToBytes turns a 0x prefixed hex string into a byte slice
func HexToBytes(in string) ([]byte, error) {
	if len(in) < 2 {
		return nil, errors.New("invalid string")
	}

	if strings.Compare(in[:2], "0x") != 0 {
		return nil, errors.New("could not byteify non 0x prefixed string")
	}

	if len(in)%2 != 0 {
		return nil, errors.New("cannot decode a odd length string")
	}

	in = in[2:]
	out, err := hex.DecodeString(in)
	return out, err
}
