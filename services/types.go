package services

import (
	"context"
	"encoding/json"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"math/big"

	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	tomochaincommon "github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
)

// Client is used by the servicers to get block
// data and to submit transactions.
type Client interface {
	// GetChainID returns the network chain context, derived from the
	// genesis document.
	GetChainID(ctx context.Context) (*big.Int, error)

	Status(context.Context) (
		*RosettaTypes.BlockIdentifier,
		int64,
		*RosettaTypes.SyncStatus,
		[]*RosettaTypes.Peer,
		error,
	)

	Block(
		context.Context,
		*RosettaTypes.PartialBlockIdentifier,
	) (*RosettaTypes.Block, error)

	Balance(
		context.Context,
		*RosettaTypes.AccountIdentifier,
		*RosettaTypes.PartialBlockIdentifier,
	) (*RosettaTypes.AccountBalanceResponse, error)

	PendingNonceAt(context.Context, tomochaincommon.Address) (uint64, error)

	NonceAt(ctx context.Context, account tomochaincommon.Address, blockNumber string) (uint64, error)

	SuggestGasPrice(ctx context.Context) (*big.Int, error)

	EstimateGas(ctx context.Context, msg common.CallArgs) (uint64, error)

	// SubmitTx submits the given encoded transaction to the node.
	SubmitTx(ctx context.Context, signedTx hexutil.Bytes) (txid string, err error)

	Call(
		ctx context.Context,
		request *RosettaTypes.CallRequest,
	) (*RosettaTypes.CallResponse, error)
}

type options struct {
	From string `json:"from"`
}

type metadata struct {
	Nonce    uint64   `json:"nonce"`
	GasPrice *big.Int `json:"gas_price"`
}

type metadataWire struct {
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gas_price"`
}

func (m *metadata) MarshalJSON() ([]byte, error) {
	mw := &metadataWire{
		Nonce:    hexutil.Uint64(m.Nonce).String(),
		GasPrice: hexutil.EncodeBig(m.GasPrice),
	}

	return json.Marshal(mw)
}

func (m *metadata) UnmarshalJSON(data []byte) error {
	var mw metadataWire
	if err := json.Unmarshal(data, &mw); err != nil {
		return err
	}

	nonce, err := hexutil.DecodeUint64(mw.Nonce)
	if err != nil {
		return err
	}

	gasPrice, err := hexutil.DecodeBig(mw.GasPrice)
	if err != nil {
		return err
	}

	m.GasPrice = gasPrice
	m.Nonce = nonce
	return nil
}

type parseMetadata struct {
	Nonce    uint64   `json:"nonce"`
	GasPrice *big.Int `json:"gas_price"`
	ChainID  *big.Int `json:"chain_id"`
}

type parseMetadataWire struct {
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gas_price"`
	ChainID  string `json:"chain_id"`
}

func (p *parseMetadata) MarshalJSON() ([]byte, error) {
	pmw := &parseMetadataWire{
		Nonce:    hexutil.Uint64(p.Nonce).String(),
		GasPrice: hexutil.EncodeBig(p.GasPrice),
		ChainID:  hexutil.EncodeBig(p.ChainID),
	}

	return json.Marshal(pmw)
}


type transactionWire struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Input    string `json:"input"`
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gas_price"`
	GasLimit string `json:"gas"`
	ChainID  string `json:"chain_id"`
}

func (t *transaction) MarshalJSON() ([]byte, error) {
	tw := &transactionWire{
		From:     t.From,
		To:       t.To,
		Value:    hexutil.EncodeBig(t.Value),
		Input:    hexutil.Encode(t.Input),
		Nonce:    hexutil.EncodeUint64(t.Nonce),
		GasPrice: hexutil.EncodeBig(t.GasPrice),
		GasLimit: hexutil.EncodeUint64(t.GasLimit),
		ChainID:  hexutil.EncodeBig(t.ChainID),
	}

	return json.Marshal(tw)
}

func (t *transaction) UnmarshalJSON(data []byte) error {
	var tw transactionWire
	if err := json.Unmarshal(data, &tw); err != nil {
		return err
	}

	value, err := hexutil.DecodeBig(tw.Value)
	if err != nil {
		return err
	}

	input, err := hexutil.Decode(tw.Input)
	if err != nil {
		return err
	}

	nonce, err := hexutil.DecodeUint64(tw.Nonce)
	if err != nil {
		return err
	}

	gasPrice, err := hexutil.DecodeBig(tw.GasPrice)
	if err != nil {
		return err
	}

	gasLimit, err := hexutil.DecodeUint64(tw.GasLimit)
	if err != nil {
		return err
	}

	chainID, err := hexutil.DecodeBig(tw.ChainID)
	if err != nil {
		return err
	}

	t.From = tw.From
	t.To = tw.To
	t.Value = value
	t.Input = input
	t.Nonce = nonce
	t.GasPrice = gasPrice
	t.GasLimit = gasLimit
	t.ChainID = chainID
	t.GasPrice = gasPrice
	return nil
}
