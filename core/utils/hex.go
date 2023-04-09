package utils

import (
	"encoding/hex"
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

// HexEncode encodes bytes to a hex string. Contrary to hex.EncodeToString, this function prefixes the hex string
// with "0x"
func HexEncodeToString(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
