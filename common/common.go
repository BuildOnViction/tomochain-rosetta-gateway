// Copyright (c) 2020 TomoChain

package common

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
)

type TransactionType uint64

const (
	// tomochain network information
	TomoChainBlockchain       = "TomoChain"
	TomoChainMainnetNetWorkId = 88
	TomoChainTestnetNetWorkId = 89
	TomoChainDevnetNetWorkId  = 99
	ExtraSeal                 = 65
	Epoch                     = 900
	DefaultGasLimit           = 10000000

	// text
	SUCCESS                     = "SUCCESS"
	FAIL                        = "FAIL"
	PENDING                     = "PENDING"
	METADATA_NEW_BALANCE        = "new_balance"
	METADATA_ACCOUNT_SEQUENCE   = "account_sequence"
	METADATA_RECENT_BLOCK_HASH  = "recent_block_hash"
	METADATA_GAS_LIMIT          = "gas_limit"
	METADATA_GAS_PRICE          = "gas_price"
	METADATA_RECEIPT            = "receipt"
	METADATA_TRACE              = "trace"
	METADATA_RECIPIENT          = "recipient"
	METADATA_SENDER             = "sender"
	METADATA_TRANSACTION_TYPE   = "type"
	METADATA_TRANSACTION_AMOUNT = "amount"
	METADATA_TRANSACTION_DATA   = "data"
	METADATA_AMOUNT             = "amount"
	METADATA_SYMBOL             = "symbol"
	METADATA_DECIMALS           = "decimals"
	METADATA_NONCE              = "nonce"
	METADATA_CHAIN_ID           = "chain_id"

	// rpc method name
	RPC_METHOD_SEND_SIGNED_TRANSACTION  = "eth_sendRawTransaction"
	RPC_METHOD_GET_PENDING_TRANSACTIONS = "eth_pendingTransactions"
	RPC_METHOD_GET_TRANSACTION_COUNT    = "eth_getTransactionCount"
	RPC_METHOD_GAS_PRICE                = "eth_gasPrice"
	RPC_METHOD_ESTIMATE_GAS             = "eth_estimateGas"
	RPC_METHOD_GET_BLOCK_BY_NUMBER      = "eth_getBlockByNumber"
	RPC_METHOD_GET_BLOCK_BY_HASH        = "eth_getBlockByHash"
	RPC_METHOD_DEBUG_TRACE_BLOCK        = "debug_traceBlockByHash"
	RPC_METHOD_GET_TRANSACTION_RECEIPT  = "eth_getTransactionReceipt"
	RPC_METHOD_GET_BALANCE              = "eth_getBalance"
	RPC_METHOD_GET_REWARD_BY_HASH       = "eth_getRewardByHash"
	RPC_METHOD_GET_CHAIN_ID             = "eth_chainId"
	// transaction type code
	TRANSACTION_TYPE_NATIVE_TRANSFER           TransactionType = 0
	TRANSACTION_TYPE_IN_CONTRACT_TRANSFER      TransactionType = 1
	TRANSACTION_TYPE_GAS_FEE                   TransactionType = 2
	TRANSACTION_TYPE_CLAIM_FROM_REWARDING_FUND TransactionType = 3
)

// Enum value maps for TransactionType.
var (
	TRANSACTION_TYPE_NAME = map[int32]string{
		0: "transfer",
		1: "in_contract_transfer",
		2: "gas_fee",
		3: "reward",
	}
	TRANSACTION_TYPE_CODE = map[string]int32{
		"transfer":             0,
		"in_contract_transfer": 1,
		"gas_fee":              2,
		"reward":               3,
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
	return TRANSACTION_TYPE_NAME[int32(t)]
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
