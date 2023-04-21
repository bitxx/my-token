package hexutil

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

// HexDecodeString decodes bytes from a hex string. Contrary to hex.DecodeString, this function does not error if "0x"
// is prefixed, and adds an extra 0 if the hex string has an odd length.
func HexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 != 0 {
		s = "0" + s
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return string(c)

}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// HexEncode encodes bytes to a hex string. Contrary to hex.EncodeToString, this function prefixes the hex string
// with "0x"
func HexEncodeToString(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// Decode decodes a hex string with 0x prefix.
func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("empty hex string")
	}
	if !has0xPrefix(input) {
		return nil, errors.New("hex string without 0x prefix")
	}
	return hex.DecodeString(input[2:])
}

// DecodeUint64 decodes a hex string with 0x prefix as a quantity.
func DecodeUint64(input string) (uint64, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(raw, 16, 64)
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", errors.New("empty hex string")
	}
	if !has0xPrefix(input) {
		return "", errors.New("hex string without 0x prefix")
	}
	input = input[2:]
	if len(input) == 0 {
		return "", errors.New("hex string \"0x\"")
	}
	if len(input) > 1 && input[0] == '0' {
		return "", errors.New("hex number with leading zero digits")
	}
	return input, nil
}
