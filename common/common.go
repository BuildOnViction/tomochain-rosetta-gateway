// Copyright (c) 2020 TomoChain

package common

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
)

type TransactionType uint64

const (
	// TomoChain network information
	TomoChainBlockchain       = "TomoChain"
	TomoChainMainnetNetWorkId = 88
	TomoChainTestnetNetWorkId = 89
	TomoChainDevnetNetWorkId  = 99
	ExtraSeal                 = 65
	Epoch                     = 900
	DefaultGasLimit           = 10000000
	HardForkUpdateTxFee       = 13523400 // tx fee transfer to masternode owner

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	// This is very important for tracking account balance from genesis to current state
	// This makes sure that: currentBalance = balanceAtGenesis + all balance-changing events
	HistoricalBalanceSupported = true

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
	RPC_METHOD_DEBUG_TRACE_TRANSACTION  = "debug_traceTransaction"
	RPC_METHOD_GET_TRANSACTION_RECEIPT  = "eth_getTransactionReceipt"
	RPC_METHOD_GET_BALANCE              = "eth_getBalance"
	RPC_METHOD_GET_REWARD_BY_HASH       = "eth_getRewardByHash"
	RPC_METHOD_GET_CHAIN_ID             = "eth_chainId"
	RPC_METHOD_GET_OWNER_BY_COINBASE    = "eth_getOwnerByCoinbase"

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

var (
	SpecialRewardAddrMap = map[string]string{
		"0x5248bfb72fd4f234e062d3e9bb76f08643004fcd": "29410",
		"0x5ac26105b35ea8935be382863a70281ec7a985e9": "23551",
		"0x09c4f991a41e7ca0645d7dfbfee160b55e562ea4": "25821",
		"0xb3157bbc5b401a45d6f60b106728bb82ebaa585b": "20051",
		"0x741277a8952128d5c2ffe0550f5001e4c8247674": "23937",
		"0x10ba49c1caa97d74b22b3e74493032b180cebe01": "27320",
		"0x07048d51d9e6179578a6e3b9ee28cdc183b865e4": "29758",
		"0x4b899001d73c7b4ec404a771d37d9be13b8983de": "26148",
		"0x85cb320a9007f26b7652c19a2a65db1da2d0016f": "27216",
		"0x06869dbd0e3a2ea37ddef832e20fa005c6f0ca39": "29449",
		"0x82e48bc7e2c93d89125428578fb405947764ad7c": "28084",
		"0x1f9a78534d61732367cbb43fc6c89266af67c989": "29287",
		"0x7c3b1fa91df55ff7af0cad9e0399384dc5c6641b": "21574",
		"0x5888dc1ceb0ff632713486b9418e59743af0fd20": "28836",
		"0xa512fa1c735fc3cc635624d591dd9ea1ce339ca5": "25515",
		"0x0832517654c7b7e36b1ef45d76de70326b09e2c7": "22873",
		"0xca14e3c4c78bafb60819a78ff6e6f0f709d2aea7": "24968",
		"0x652ce195a23035114849f7642b0e06647d13e57a": "24091",
		"0x29a79f00f16900999d61b6e171e44596af4fb5ae": "20790",
		"0xf9fd1c2b0af0d91b0b6754e55639e3f8478dd04a": "23331",
		"0xb835710c9901d5fe940ef1b99ed918902e293e35": "28273",
		"0x04dd29ce5c253377a9a3796103ea0d9a9e514153": "29956",
		"0x2b4b56846eaf05c1fd762b5e1ac802efd0ab871c": "24911",
		"0x1d1f909f6600b23ce05004f5500ab98564717996": "25637",
		"0x0dfdcebf80006dc9ab7aae8c216b51c6b6759e86": "26378",
		"0x2b373890a28e5e46197fbc04f303bbfdd344056f": "21109",
		"0xa8a3ef3dc5d8e36aee76f3671ec501ec31e28254": "22072",
		"0x4f3d18136fe2b5665c29bdaf74591fc6625ef427": "21650",
		"0x175d728b0e0f1facb5822a2e0c03bde93596e324": "21588",
		"0xd575c2611984fcd79513b80ab94f59dc5bab4916": "28971",
		"0x0579337873c97c4ba051310236ea847f5be41bc0": "28344",
		"0xed12a519cc15b286920fc15fd86106b3e6a16218": "24443",
		"0x492d26d852a0a0a2982bb40ec86fe394488c419e": "26623",
		"0xce5c7635d02dc4e1d6b46c256cae6323be294a32": "28459",
		"0x8b94db158b5e78a6c032c7e7c9423dec62c8b11c": "21803",
		"0x0e7c48c085b6b0aa7ca6e4cbcc8b9a92dc270eb4": "21739",
		"0x206e6508462033ef8425edc6c10789d241d49acb": "21883",
		"0x7710e7b7682f26cb5a1202e1cff094fbf7777758": "28907",
		"0xcb06f949313b46bbf53b8e6b2868a0c260ff9385": "28932",
		"0xf884e43533f61dc2997c0e19a6eff33481920c00": "27780",
		"0x8b635ef2e4c8fe21fc2bda027eb5f371d6aa2fc1": "23115",
		"0x10f01a27cf9b29d02ce53497312b96037357a361": "22716",
		"0x693dd49b0ed70f162d733cf20b6c43dc2a2b4d95": "20020",
		"0xe0bec72d1c2a7a7fb0532cdfac44ebab9f6f41ee": "23071",
		"0xc8793633a537938cb49cdbbffd45428f10e45b64": "24652",
		"0x0d07a6cbbe9fa5c4f154e5623bfe47fb4d857d8e": "21907",
		"0xd4080b289da95f70a586610c38268d8d4cf1e4c4": "22719",
		"0x8bcfb0caf41f0aa1b548cae76dcdd02e33866a1b": "29062",
		"0xabfef22b92366d3074676e77ea911ccaabfb64c1": "23110",
		"0xcc4df7a32faf3efba32c9688def5ccf9fefe443d": "21397",
		"0x7ec1e48a582475f5f2b7448a86c4ea7a26ea36f8": "23105",
		"0xe3de67289080f63b0c2612844256a25bb99ac0ad": "29721",
		"0x3ba623300cf9e48729039b3c9e0dee9b785d636e": "25917",
		"0x402f2cfc9c8942f5e7a12c70c625d07a5d52fe29": "24712",
		"0xd62358d42afbde095a4ca868581d85f9adcc3d61": "24449",
		"0x3969f86acb733526cd61e3c6e3b4660589f32bc6": "29579",
		"0x67615413d7cdadb2c435a946aec713a9a9794d39": "26333",
		"0xfe685f43acc62f92ab01a8da80d76455d39d3cb3": "29825",
		"0x3538a544021c07869c16b764424c5987409cba48": "22746",
		"0xe187cf86c2274b1f16e8225a7da9a75aba4f1f5f": "23734",
	}

	SpecialRewardBlockMap = map[uint64]string{
		9073579: "0x5248bfb72fd4f234e062d3e9bb76f08643004fcd",
		9147130: "0x5ac26105b35ea8935be382863a70281ec7a985e9",
		9147195: "0x09c4f991a41e7ca0645d7dfbfee160b55e562ea4",
		9147200: "0xb3157bbc5b401a45d6f60b106728bb82ebaa585b",
		9147206: "0x741277a8952128d5c2ffe0550f5001e4c8247674",
		9147212: "0x10ba49c1caa97d74b22b3e74493032b180cebe01",
		9147217: "0x07048d51d9e6179578a6e3b9ee28cdc183b865e4",
		9147223: "0x4b899001d73c7b4ec404a771d37d9be13b8983de",
		9147229: "0x85cb320a9007f26b7652c19a2a65db1da2d0016f",
		9147234: "0x06869dbd0e3a2ea37ddef832e20fa005c6f0ca39",
		9147240: "0x82e48bc7e2c93d89125428578fb405947764ad7c",
		9147246: "0x1f9a78534d61732367cbb43fc6c89266af67c989",
		9147251: "0x7c3b1fa91df55ff7af0cad9e0399384dc5c6641b",
		9147257: "0x5888dc1ceb0ff632713486b9418e59743af0fd20",
		9147263: "0xa512fa1c735fc3cc635624d591dd9ea1ce339ca5",
		9147268: "0x0832517654c7b7e36b1ef45d76de70326b09e2c7",
		9147274: "0xca14e3c4c78bafb60819a78ff6e6f0f709d2aea7",
		9147279: "0x652ce195a23035114849f7642b0e06647d13e57a",
		9147285: "0x29a79f00f16900999d61b6e171e44596af4fb5ae",
		9147291: "0xf9fd1c2b0af0d91b0b6754e55639e3f8478dd04a",
		9147296: "0xb835710c9901d5fe940ef1b99ed918902e293e35",
		9147302: "0x04dd29ce5c253377a9a3796103ea0d9a9e514153",
		9147308: "0x2b4b56846eaf05c1fd762b5e1ac802efd0ab871c",
		9147314: "0x1d1f909f6600b23ce05004f5500ab98564717996",
		9147319: "0x0dfdcebf80006dc9ab7aae8c216b51c6b6759e86",
		9147325: "0x2b373890a28e5e46197fbc04f303bbfdd344056f",
		9147330: "0xa8a3ef3dc5d8e36aee76f3671ec501ec31e28254",
		9147336: "0x4f3d18136fe2b5665c29bdaf74591fc6625ef427",
		9147342: "0x175d728b0e0f1facb5822a2e0c03bde93596e324",
		9145281: "0xd575c2611984fcd79513b80ab94f59dc5bab4916",
		9145315: "0x0579337873c97c4ba051310236ea847f5be41bc0",
		9145341: "0xed12a519cc15b286920fc15fd86106b3e6a16218",
		9145367: "0x492d26d852a0a0a2982bb40ec86fe394488c419e",
		9145386: "0xce5c7635d02dc4e1d6b46c256cae6323be294a32",
		9145414: "0x8b94db158b5e78a6c032c7e7c9423dec62c8b11c",
		9145436: "0x0e7c48c085b6b0aa7ca6e4cbcc8b9a92dc270eb4",
		9145463: "0x206e6508462033ef8425edc6c10789d241d49acb",
		9145493: "0x7710e7b7682f26cb5a1202e1cff094fbf7777758",
		9145519: "0xcb06f949313b46bbf53b8e6b2868a0c260ff9385",
		9145549: "0xf884e43533f61dc2997c0e19a6eff33481920c00",
		9147352: "0x8b635ef2e4c8fe21fc2bda027eb5f371d6aa2fc1",
		9147357: "0x10f01a27cf9b29d02ce53497312b96037357a361",
		9147363: "0x693dd49b0ed70f162d733cf20b6c43dc2a2b4d95",
		9147369: "0xe0bec72d1c2a7a7fb0532cdfac44ebab9f6f41ee",
		9147375: "0xc8793633a537938cb49cdbbffd45428f10e45b64",
		9147380: "0x0d07a6cbbe9fa5c4f154e5623bfe47fb4d857d8e",
		9147386: "0xd4080b289da95f70a586610c38268d8d4cf1e4c4",
		9147392: "0x8bcfb0caf41f0aa1b548cae76dcdd02e33866a1b",
		9147397: "0xabfef22b92366d3074676e77ea911ccaabfb64c1",
		9147403: "0xcc4df7a32faf3efba32c9688def5ccf9fefe443d",
		9147408: "0x7ec1e48a582475f5f2b7448a86c4ea7a26ea36f8",
		9147414: "0xe3de67289080f63b0c2612844256a25bb99ac0ad",
		9147420: "0x3ba623300cf9e48729039b3c9e0dee9b785d636e",
		9147425: "0x402f2cfc9c8942f5e7a12c70c625d07a5d52fe29",
		9147431: "0xd62358d42afbde095a4ca868581d85f9adcc3d61",
		9147437: "0x3969f86acb733526cd61e3c6e3b4660589f32bc6",
		9147442: "0x67615413d7cdadb2c435a946aec713a9a9794d39",
		9147448: "0xfe685f43acc62f92ab01a8da80d76455d39d3cb3",
		9147453: "0x3538a544021c07869c16b764424c5987409cba48",
		9147459: "0xe187cf86c2274b1f16e8225a7da9a75aba4f1f5f",
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

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    hexutil.Big     `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
}

func SupportedOperationTypes() []string {
	return OperationTypes
}
