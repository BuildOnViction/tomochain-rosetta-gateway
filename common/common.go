// Copyright (c) 2020 TomoChain

package common

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"strconv"
)

type TransactionType uint64

const (
	// tomochain network information
	TomoChainBlockchain       = "tomochain"
	TomoChainMainnetNetWorkId = 88

	// text
	SUCCESS                    = "SUCCESS"
	FAIL                       = "FAIL"
	PENDING                    = "PENDING"
	METADATA_NEW_BALANCE       = "new_balance"
	METADATA_SEQUENCE_NUMBER   = "sequence_number"
	METADATA_RECENT_BLOCK_HASH = "recent_block_hash"
	METADATA_GAS_LIMIT         = "gas_limit"
	METADATA_GAS_PRICE         = "gas_price"
	METADATA_RECIPIENT         = "recipient"
	METADATA_SENDER            = "sender"
	METADATA_TRANSACTION_TYPE  = "type"
	METADATA_TRANSACTION_VALUE = "value"
	METADATA_TRANSACTION_DATA  = "data"
	METADATA_AMOUNT            = "amount"
	METADATA_SYMBOL            = "symbol"
	METADATA_DECIMALS          = "decimals"

	// rpc method name
	RPC_METHOD_SEND_SIGNED_TRANSACTION  = "eth_sendRawTransaction"
	RPC_METHOD_GET_PENDING_TRANSACTIONS = "eth_pendingTransactions"

	// transaction type code
	TRANSACTION_TYPE_IN_CONTRACT_TRANSFER       TransactionType = 0
	TRANSACTION_TYPE_WITHDRAW_BUCKET            TransactionType = 1
	TRANSACTION_TYPE_CREATE_BUCKET              TransactionType = 2
	TRANSACTION_TYPE_DEPOSIT_TO_BUCKET          TransactionType = 3
	TRANSACTION_TYPE_CANDIDATE_SELF_STAKE       TransactionType = 4
	TRANSACTION_TYPE_CANDIDATE_REGISTRATION_FEE TransactionType = 5
	TRANSACTION_TYPE_GAS_FEE                    TransactionType = 6
	TRANSACTION_TYPE_NATIVE_TRANSFER            TransactionType = 7
	TRANSACTION_TYPE_DEPOSIT_TO_REWARDING_FUND  TransactionType = 8
	TRANSACTION_TYPE_CLAIM_FROM_REWARDING_FUND  TransactionType = 9
)

// Enum value maps for TransactionType.
var (
	TRANSACTION_TYPE_NAME = map[int32]string{
		0: "IN_CONTRACT_TRANSFER",
		1: "WITHDRAW_BUCKET",
		2: "CREATE_BUCKET",
		3: "DEPOSIT_TO_BUCKET",
		4: "CANDIDATE_SELF_STAKE",
		5: "CANDIDATE_REGISTRATION_FEE",
		6: "GAS_FEE",
		7: "NATIVE_TRANSFER",
		8: "DEPOSIT_TO_REWARDING_FUND",
		9: "CLAIM_FROM_REWARDING_FUND",
	}
	TRANSACTION_TYPE_CODE = map[string]int32{
		"IN_CONTRACT_TRANSFER":       0,
		"WITHDRAW_BUCKET":            1,
		"CREATE_BUCKET":              2,
		"DEPOSIT_TO_BUCKET":          3,
		"CANDIDATE_SELF_STAKE":       4,
		"CANDIDATE_REGISTRATION_FEE": 5,
		"GAS_FEE":                    6,
		"NATIVE_TRANSFER":            7,
		"DEPOSIT_TO_REWARDING_FUND":  8,
		"CLAIM_FROM_REWARDING_FUND":  9,
	}
)

type RPCTransaction struct {
	BlockHash        common.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex hexutil.Uint    `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

var (
	TomoNativeCoin = &types.Currency{
		Symbol:   "TOMO",
		Decimals: 18,
	}
)

func (t TransactionType) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func SupportedOperationTypes() []string {
	opTyps := make([]string, 0, len(TRANSACTION_TYPE_NAME))
	for _, name := range TRANSACTION_TYPE_NAME {
		opTyps = append(opTyps, name)
	}
	return opTyps
}

func SupportedConstructionTypes() []string {
	return []string{
		TRANSACTION_TYPE_NATIVE_TRANSFER.String(),
	}
}
