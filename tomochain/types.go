package tomochain

import (
	"encoding/json"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tomochaincommon "github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/params"
	"math/big"
)

var (

	// Blockchain is TomoChain.
	Blockchain string = "TomoChain"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "88"


	// TestnetNetwork is the value of the network
	TestnetNetwork string = "89"

	// DevnetNetwork is the value of the network
	DevnetNetwork string = "1992"

	// MainnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	MainnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  params.TomoMainnetGenesisHash.Hex(),
		Index: GenesisBlockIndex,
	}

	// TestnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the testnet genesis block.
	TestnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  params.TestnetGenesisHash.Hex(),
		Index: GenesisBlockIndex,
	}

	// DevnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the testnet genesis block.
	DevnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  "",
		Index: GenesisBlockIndex,
	}

	// MainnetTomoArguments are the arguments to start a mainnet tomochain instance.
	MainnetTomoArguments = `--config=/app/tomochain/tomochain.toml --gcmode=archive  --store-reward --tomox.dbengine "leveldb" `

	// TestnetTomoArguments are the arguments to start a mainnet tomochain instance.
	TestnetTomoArguments = ` -- testnet --config=/app/tomochain/tomochain.toml --gcmode=archive --store-reward --tomox.dbengine "leveldb" `

	// DevnetTomoArguments are the arguments to start a mainnet tomochain instance.
	DevnetTomoArguments = `--config=/app/tomochain/tomochain.toml --gcmode=archive  --store-reward --tomox.dbengine "leveldb" --ethstats "rosetta:nguyennguyen@localhost:3333"  `
)

var CallMethods = []string{
	common.RPC_METHOD_GET_TRANSACTION_RECEIPT,
}
type rpcBlock struct {
	Hash         tomochaincommon.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
	UncleHashes  []tomochaincommon.Hash    `json:"uncles"`
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *tomochaincommon.Hash    `json:"blockHash,omitempty"`
	From        *tomochaincommon.Address `json:"from,omitempty"`
}

type rpcTransaction struct {
	tx *tomochaintypes.Transaction
	txExtraInfo
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func (tx *rpcTransaction) LoadedTransaction() *loadedTransaction {
	ethTx := &loadedTransaction{
		Transaction: tx.tx,
		From:        tx.txExtraInfo.From,
		BlockNumber: tx.txExtraInfo.BlockNumber,
		BlockHash:   tx.txExtraInfo.BlockHash,
	}
	return ethTx
}

type loadedTransaction struct {
	Transaction *tomochaintypes.Transaction
	From        *tomochaincommon.Address
	BlockNumber *string
	BlockHash   *tomochaincommon.Hash
	FeeAmount   *big.Int
	Miner       string
	Status      bool

	Trace    *Call
	RawTrace json.RawMessage
	Receipt  *tomochaintypes.Receipt
}

type rpcCall struct {
	Result *Call `json:"result"`
}

type rpcRawCall struct {
	Result json.RawMessage `json:"result"`
}

// Call is an TomoChain debug trace.
type Call struct {
	Type         string                  `json:"type"`
	From         tomochaincommon.Address `json:"from"`
	To           tomochaincommon.Address `json:"to"`
	Value        *big.Int                `json:"value"`
	GasUsed      *big.Int                `json:"gasUsed"`
	Revert       bool
	ErrorMessage string  `json:"error"`
	Calls        []*Call `json:"calls"`
}

type flatCall struct {
	Type         string                  `json:"type"`
	From         tomochaincommon.Address `json:"from"`
	To           tomochaincommon.Address `json:"to"`
	Value        *big.Int                `json:"value"`
	GasUsed      *big.Int                `json:"gasUsed"`
	Revert       bool
	ErrorMessage string `json:"error"`
}

func (t *Call) flatten() *flatCall {
	return &flatCall{
		Type:         t.Type,
		From:         t.From,
		To:           t.To,
		Value:        t.Value,
		GasUsed:      t.GasUsed,
		Revert:       t.Revert,
		ErrorMessage: t.ErrorMessage,
	}
}

// UnmarshalJSON is a custom unmarshaler for Call.
func (t *Call) UnmarshalJSON(input []byte) error {
	type CustomTrace struct {
		Type         string         `json:"type"`
		From         tomochaincommon.Address `json:"from"`
		To           tomochaincommon.Address `json:"to"`
		Value        *hexutil.Big   `json:"value"`
		GasUsed      *hexutil.Big   `json:"gasUsed"`
		Revert       bool
		ErrorMessage string  `json:"error"`
		Calls        []*Call `json:"calls"`
	}
	var dec CustomTrace
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	t.Type = dec.Type
	t.From = dec.From
	t.To = dec.To
	if dec.Value != nil {
		t.Value = (*big.Int)(dec.Value)
	} else {
		t.Value = new(big.Int)
	}
	if dec.GasUsed != nil {
		t.GasUsed = (*big.Int)(dec.Value)
	} else {
		t.GasUsed = new(big.Int)
	}
	if dec.ErrorMessage != "" {
		// Any error surfaced by the decoder means that the transaction
		// has reverted.
		t.Revert = true
	}
	t.ErrorMessage = dec.ErrorMessage
	t.Calls = dec.Calls
	return nil
}
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

type rpcProgress struct {
	StartingBlock hexutil.Uint64
	CurrentBlock  hexutil.Uint64
	HighestBlock  hexutil.Uint64
	PulledStates  hexutil.Uint64
	KnownStates   hexutil.Uint64
}


// GetTransactionReceiptInput is the input to the call
// method "eth_getTransactionReceipt".
type GetTransactionReceiptInput struct {
	TxHash string `json:"tx_hash"`
}


// CallType returns a boolean indicating
// if the provided trace type is a call type.
func CallType(t string) bool {
	callTypes := []string{
		common.CallOpType,
		common.CallCodeOpType,
		common.DelegateCallOpType,
		common.StaticCallOpType,
	}

	for _, callType := range callTypes {
		if callType == t {
			return true
		}
	}

	return false
}

// CreateType returns a boolean indicating
// if the provided trace type is a create type.
func CreateType(t string) bool {
	createTypes := []string{
		common.CreateOpType,
		common.Create2OpType,
	}

	for _, createType := range createTypes {
		if createType == t {
			return true
		}
	}

	return false
}