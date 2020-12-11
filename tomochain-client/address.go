package tomochain_client

import (
	"github.com/tomochain/tomochain/common"
	"log"
)

// ChecksumAddress ensures an Ethereum hex address
// is in Checksum Format. If the address cannot be converted,
// it returns !ok.
func ChecksumAddress(address string) (string, bool) {
	ok := common.IsHexAddress(address)
	if !ok {
		return "", false
	}

	return common.HexToAddress(address).Hex(), true
}

// MustChecksum ensures an address can be converted
// into a valid checksum. If it does not, the program
// will exit.
func MustChecksum(address string) string {
	addr, ok := ChecksumAddress(address)
	if !ok {
		log.Fatalf("invalid address %s", address)
	}

	return addr
}
