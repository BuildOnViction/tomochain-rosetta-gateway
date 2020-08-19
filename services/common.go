// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/types"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
)

const (
	TomoChainBlockchain = "tomochain"
	TomoChainNetWorkId  = 88
	StatusSuccess       = "success"
	StatusFail          = "fail"
)

// ValidateNetworkIdentifier validates the network identifier.
func ValidateNetworkIdentifier(ctx context.Context, client tc.TomoChainClient, ni *types.NetworkIdentifier) *types.Error {
	if ni != nil {
		if ni.Blockchain != TomoChainBlockchain {
			return ErrInvalidBlockchain
		}
		if ni.SubNetworkIdentifier != nil {
			return ErrInvalidSubnetwork
		}

		if chainId, err := client.GetChainID(ctx); err != nil || chainId.Uint64() != TomoChainNetWorkId {
			return ErrInvalidNetwork
		}
	} else {
		return ErrMissingNID
	}
	return nil
}

func SupportedOperationTypes() []string {
	opTyps := make([]string, 0, len(TransactionLogType_name))
	for _, name := range TransactionLogType_name {
		opTyps = append(opTyps, name)
	}
	return opTyps
}

func SupportedConstructionTypes() []string {
	return []string{
		TransactionLogType_NATIVE_TRANSFER.String(),
	}
}

func IsSupportedConstructionType(typ string) bool {
	for _, styp := range SupportedConstructionTypes() {
		if typ == styp {
			return true
		}
	}
	return false
}
