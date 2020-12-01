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



	// MinerRewardOpType is used to describe
	// a miner block reward.
	MinerRewardOpType = "MINER_REWARD"

	// FeeOpType is used to represent fee operations.
	FeeOpType = "FEE"

	// CallOpType is used to represent CALL trace operations.
	CallOpType = "CALL"

	// CreateOpType is used to represent CREATE trace operations.
	CreateOpType = "CREATE"

	// Create2OpType is used to represent CREATE2 trace operations.
	Create2OpType = "CREATE2"

	// SelfDestructOpType is used to represent SELFDESTRUCT trace operations.
	SelfDestructOpType = "SELFDESTRUCT"

	// CallCodeOpType is used to represent CALLCODE trace operations.
	CallCodeOpType = "CALLCODE"

	// DelegateCallOpType is used to represent DELEGATECALL trace operations.
	DelegateCallOpType = "DELEGATECALL"

	// StaticCallOpType is used to represent STATICCALL trace operations.
	StaticCallOpType = "STATICCALL"

	// DestructOpType is a synthetic operation used to represent the
	// deletion of suicided accounts that still have funds at the end
	// of a transaction.
	DestructOpType = "DESTRUCT"
)

var (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"

	TomoNativeCoin = &types.Currency{
		Symbol:   "TOMO",
		Decimals: 18,
	}
	OperationTypes = []string{
		MinerRewardOpType,
		FeeOpType,
		CallOpType,
		CreateOpType,
		Create2OpType,
		SelfDestructOpType,
		CallCodeOpType,
		DelegateCallOpType,
		StaticCallOpType,
		DestructOpType,
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

func SupportedOperationTypes() []string {
	return OperationTypes
}
