// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	tomochaincommon "github.com/tomochain/tomochain/common"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"math/big"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cast"

	"github.com/tomochain/tomochain/crypto"
)

type constructionAPIService struct {
	client tc.TomoChainClient
}

// NewConstructionAPIService creates a new instance of an ConstructionAPIService.
func NewConstructionAPIService(client tc.TomoChainClient) server.ConstructionAPIServicer {
	return &constructionAPIService{
		client: client,
	}
}

// ConstructionCombine implements the /construction/combine endpoint.
func (s *constructionAPIService) ConstructionCombine(
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	b, err := hex.DecodeString(request.UnsignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}
	var unsignTx *tomochaintypes.Transaction
	err = json.Unmarshal(b, &unsignTx)
	if err != nil {
		terr := common.ErrUnableToParseTx
		terr.Message += err.Error()
		return nil, terr
	}
	if len(request.Signatures) != 1 {
		terr := common.ErrInvalidInputParam
		terr.Message += "need exact 1 signature"
		return nil, terr
	}

	rawSig := request.Signatures[0].Bytes
	if len(rawSig) != 65 {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid signature length"
		return nil, terr
	}

	chainId, err := strconv.Atoi(request.NetworkIdentifier.Network)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid network"
		return nil, terr
	}
	signedTx, err := unsignTx.WithSignature(tomochaintypes.NewEIP155Signer(big.NewInt(cast.ToInt64(chainId))), rawSig)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += "cannot sign transaction"
		return nil, terr
	}
	signedTxData, err := json.Marshal(signedTx)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += "cannot marshal signed transaction"
		return nil, terr
	}
	return &types.ConstructionCombineResponse{
		SignedTransaction: hex.EncodeToString(signedTxData),
	}, nil
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *constructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	if len(request.PublicKey.Bytes) == 0 || request.PublicKey.CurveType != types.Secp256k1 {
		terr := common.ErrInvalidInputParam
		terr.Message += "unsupported public key type"
		return nil, terr
	}

	rawPub := request.PublicKey.Bytes
	addr := tomochaincommon.BytesToAddress(crypto.Keccak256(rawPub[1:])[12:])

	return &types.ConstructionDeriveResponse{
		Address: addr.String(),
	}, nil
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *constructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid signed transaction format: " + err.Error()
		return nil, terr
	}
	var tx *tomochaintypes.Transaction
	err = json.Unmarshal(tran, &tx)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += "unable to unmarshal signed transaction " + err.Error()
		return nil, terr
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: tx.Hash().String(),
		},
	}, nil
}

// FIXME: required options
// sender (string): address of sender
// to (string): destination address
// gas_limit (uint64) : gas limit of the transaction
// gas_price (uint64): gas price in wei
// value (uint64)
// data ([]bytes) : data include method name, argument if this tx call a contract

func parseMetaDataToCallMsg(options map[string]interface{}) (tomochain.CallMsg, *types.Error) {
	sender, ok := options[common.METADATA_SENDER]
	if !ok {
		terr := common.ErrInvalidInputParam
		terr.Message += "empty sender address"
		return tomochain.CallMsg{}, terr
	}

	to, ok := options[common.METADATA_RECIPIENT]
	if !ok {
		terr := common.ErrInvalidInputParam
		terr.Message += "empty sender address"
		return tomochain.CallMsg{}, terr
	}
	destinationAddress := tomochaincommon.HexToAddress(to.(string))

	gasLimit, ok := options[common.METADATA_GAS_LIMIT]
	if !ok {
		terr := common.ErrInvalidInputParam
		terr.Message += "empty gasLimit"
		return tomochain.CallMsg{}, terr
	}

	gasPrice, ok := options[common.METADATA_GAS_PRICE]
	if !ok {
		gasPrice = big.NewInt(tomochaincommon.DefaultMinGasPrice)
	}

	v, ok := options[common.METADATA_TRANSACTION_VALUE]
	if !ok {
		v = big.NewInt(0)
	}

	d, ok := options[common.METADATA_TRANSACTION_DATA]
	if !ok {
		d = []byte{}
	}

	callMsg := tomochain.CallMsg{
		From:            tomochaincommon.HexToAddress(sender.(string)),
		To:              &destinationAddress,
		Gas:             gasLimit.(uint64),
		GasPrice:        new(big.Int).SetUint64(gasPrice.(uint64)),
		Value:           new(big.Int).SetUint64(v.(uint64)),
		Data:            d.([]byte),
		BalanceTokenFee: nil,
	}
	return callMsg, nil
}

// ConstructionMetadata implements the /construction/metadata endpoint.
func (s *constructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	callMsg, terr := parseMetaDataToCallMsg(request.Options)
	if terr != nil {
		return nil, terr
	}
	estimateGas, err := s.client.EstimateGas(ctx, callMsg)
	if err != nil {
		return nil, common.ErrUnableToEstimateGas
	}
	account, err := s.client.GetAccount(ctx, nil, callMsg.From.String())
	if err != nil {
		terr := common.ErrUnableToGetAccount
		terr.Message += err.Error()
		return nil, terr
	}
	meta := account.Metadata

	meta[common.METADATA_GAS_LIMIT] = callMsg.Gas
	meta[common.METADATA_GAS_PRICE] = callMsg.GasPrice
	suggestedFee := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), callMsg.GasPrice)

	return &types.ConstructionMetadataResponse{
		Metadata: meta,
		SuggestedFee: []*types.Amount{
			{
				Value:    suggestedFee.String(),
				Currency: common.TomoNativeCoin,
			},
		},
	}, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *constructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	var tx *tomochaintypes.Transaction
	b, err := hex.DecodeString(request.Transaction)
	if err != nil {
		return nil, common.ErrUnableToParseTx
	}
	err = json.Unmarshal(b, &tx)
	if err != nil {
		return nil, common.ErrUnableToParseTx
	}

	ops, err := parseOpsFromTx(tx)
	if err != nil {
		return nil, common.ErrUnableToParseTx
	}
	meta, err := parseMetaDataFromTx(tx)
	if err != nil {
		return nil, common.ErrUnableToParseTx
	}

	resp := &types.ConstructionParseResponse{
		Operations: ops,
		Metadata:   meta,
	}
	sender := meta[common.METADATA_SENDER].(string)
	if request.Signed {
		resp.Signers = []string{sender}
	}
	return resp, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *constructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	if err := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}
	if len(request.Operations) != 2 {
		terr := common.ErrInvalidInputParam
		terr.Message += "ConstructionPayloadsRequest require 2 operations"
		return nil, terr
	}
	addr := request.Operations[0].Account.Address
	account, err := s.client.GetAccount(ctx, nil, addr)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	nonce, ok := account.Metadata[common.METADATA_SEQUENCE_NUMBER]
	if !ok || nonce == nil {
		return nil, common.ErrUnableToGetNextNonce
	}

	tx := tomochaintypes.NewTransaction(cast.ToUint64(nonce),
		tomochaincommon.HexToAddress(request.Operations[1].Account.Address),
		new(big.Int).SetUint64(cast.ToUint64(request.Operations[0].Amount.Value)),
		cast.ToUint64(request.Metadata[common.METADATA_GAS_LIMIT]),
		new(big.Int).SetUint64(cast.ToUint64(request.Operations[0].Amount.Value)),
		request.Metadata[common.METADATA_TRANSACTION_DATA].([]byte))

	d, err := json.Marshal(tx)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	unsignedTx := hex.EncodeToString(d)
	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedTx,
		Payloads: []*types.SigningPayload{
			{
				Address:       request.Operations[0].Account.Address,
				Bytes:         tx.Hash().Bytes(),
				SignatureType: types.Ecdsa,
			},
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *constructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if err := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}

	options := make(map[string]interface{})
	if len(request.Operations) != 2 {
		return nil, common.ErrConstructionCheck
	}
	// sender
	options[common.METADATA_SENDER] = request.Operations[0].Account.Address
	options[common.METADATA_TRANSACTION_TYPE] = request.Operations[0].Type
	options[common.METADATA_AMOUNT] = request.Operations[0].Amount.Value
	options[common.METADATA_SYMBOL] = request.Operations[0].Amount.Currency.Symbol
	options[common.METADATA_DECIMALS] = request.Operations[0].Amount.Currency.Decimals

	// recipient
	options[common.METADATA_RECIPIENT] = request.Operations[1].Account.Address

	// XXX it is unclear where these meta data should be
	if request.Metadata[common.METADATA_GAS_LIMIT] != nil {
		options[common.METADATA_GAS_LIMIT] = request.Metadata[common.METADATA_GAS_LIMIT]
	}
	if request.Metadata[common.METADATA_GAS_PRICE] != nil {
		options[common.METADATA_GAS_PRICE] = request.Metadata[common.METADATA_GAS_PRICE]
	}

	return &types.ConstructionPreprocessResponse{
		Options: options,
	}, nil
}

// ConstructionSubmit implements the /construction/submit endpoint.
func (s *constructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}

	txID, err := s.client.SubmitTx(ctx, tran)
	if err != nil {
		terr := common.ErrUnableToSubmitTx
		terr.Message += err.Error()
		return nil, terr
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txID,
		},
	}, nil
}

func parseOpsFromTx(transaction *tomochaintypes.Transaction) ([]*types.Operation, error) {
	// TODO:
	return nil, nil
}

func parseMetaDataFromTx(transaction *tomochaintypes.Transaction) (map[string]interface{}, error) {
	// TODO:
	return nil, nil
}
