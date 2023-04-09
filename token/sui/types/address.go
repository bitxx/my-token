package types

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"strings"
)

const (
	SIGNATURE_SCHEME_FLAG_ED25519 = 0x0
	ADDRESS_LENGTH                = 64
)

type Address = HexData

func GenAddress(scheme SignatureScheme, publicKey []byte) string {
	tmp := []byte{scheme.Flag()}
	tmp = append(tmp, publicKey...)
	addrBytes := blake2b.Sum256(tmp)
	return "0x" + hex.EncodeToString(addrBytes[:])[:ADDRESS_LENGTH]
}

// NewAddressFromHex
/**
 * Creates Address from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 * e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*Address, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	data, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = ADDRESS_LENGTH / 2
	if len(data) > addressLength {
		return nil, fmt.Errorf("hex string is too long. Address's length is %v data", addressLength)
	}

	res := [addressLength]byte{}
	copy(res[addressLength-len(data):], data[:])
	return &Address{
		data: res[:],
	}, nil
}

// ShortString Returns the address with leading zeros trimmed, e.g. 0x2
func (a Address) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a.data), "0")
}

func IsSameStringAddress(addr1, addr2 string) bool {
	if strings.HasPrefix(addr1, "0x") {
		addr1 = addr1[2:]
	}
	if strings.HasPrefix(addr2, "0x") {
		addr2 = addr2[2:]
	}
	addr1 = strings.TrimLeft(addr1, "0")
	return strings.TrimLeft(addr1, "0") == strings.TrimLeft(addr2, "0")
}
