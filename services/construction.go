// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cast"
	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	tomochaincommon "github.com/tomochain/tomochain/common"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/rlp"
	"math/big"
	"strconv"
)

type transaction struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Value    *big.Int `json:"value"`
	Input    []byte   `json:"input"`
	Nonce    uint64   `json:"nonce"`
	GasPrice *big.Int `json:"gas_price"`
	GasLimit uint64   `json:"gas"`
	ChainID  *big.Int `json:"chain_id"`
}

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
	fmt.Println("\n\n\n\nENDPOINT: ConstructionCombine")
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	b, err := hex.DecodeString(request.UnsignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}
	unsignTx := &transaction{}
	err = rlp.DecodeBytes(b, unsignTx)
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

	tomochainTransaction := tomochaintypes.NewTransaction(
		unsignTx.Nonce,
		tomochaincommon.HexToAddress(unsignTx.To),
		unsignTx.Value,
		unsignTx.GasLimit,
		unsignTx.GasPrice,
		unsignTx.Input,
	)
	signedTx, err := tomochainTransaction.WithSignature(tomochaintypes.NewEIP155Signer(big.NewInt(cast.ToInt64(chainId))), rawSig)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += "cannot sign transaction"
		return nil, terr
	}
	signedTxData, err := rlp.EncodeToBytes(signedTx)
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
	fmt.Println("\n\n\n\nENDPOINT: ConstructionDerive")
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	if len(request.PublicKey.Bytes) == 0 || request.PublicKey.CurveType != types.Secp256k1 {
		terr := common.ErrInvalidInputParam
		terr.Message += "unsupported public key type"
		return nil, terr
	}
	pubkey, err := crypto.DecompressPubkey(request.PublicKey.Bytes)
	if err != nil {
		return nil, common.ErrUnableToDecompressPubkey
	}
	pubBytes := crypto.FromECDSAPub(pubkey)
	addr := tc.PubToAddress(pubBytes)

	return &types.ConstructionDeriveResponse{
		Address: addr.String(),
	}, nil
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *constructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Println("\n\n\n\nENDPOINT: ConstructionHash")
	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid signed transaction format: " + err.Error()
		return nil, terr
	}
	tx := &tomochaintypes.Transaction{}
	err = rlp.DecodeBytes(tran, tx)
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
	destinationAddress := tomochaincommon.HexToAddress(cast.ToString(to))

	gasLimit, ok := options[common.METADATA_GAS_LIMIT]
	if !ok {
		// set default gaslimit
		gasLimit = common.DefaultGasLimit
	}

	gp, ok := options[common.METADATA_GAS_PRICE]
	if !ok {
		gp = tomochaincommon.DefaultMinGasPrice
	}
	gasPrice, _ := new(big.Int).SetString(cast.ToString(gp), 10)

	v, ok := options[common.METADATA_TRANSACTION_AMOUNT]
	if !ok {
		v = 0
	}
	value, _ := new(big.Int).SetString(cast.ToString(v), 10)

	d, ok := options[common.METADATA_TRANSACTION_DATA]
	if !ok || d == nil {
		d = []byte{}
	}

	callMsg := tomochain.CallMsg{
		From:            tomochaincommon.HexToAddress(cast.ToString(sender)),
		To:              &destinationAddress,
		Gas:             cast.ToUint64(gasLimit),
		GasPrice:        gasPrice,
		Value:           new(big.Int).Abs(value),
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
	fmt.Println("\n\n\n\nENDPOINT: ConstructionMetadata")
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
	account, err := s.client.GetAccount(ctx, callMsg.From.String())
	if err != nil {
		terr := common.ErrUnableToGetAccount
		terr.Message += err.Error()
		return nil, terr
	}
	meta := account.Metadata

	meta[common.METADATA_GAS_LIMIT] = callMsg.Gas
	meta[common.METADATA_GAS_PRICE] = callMsg.GasPrice
	meta[common.METADATA_SENDER] = callMsg.From

	v, ok := request.Options[common.METADATA_TRANSACTION_AMOUNT]
	if !ok {
		v = 0
	}
	value, _ := new(big.Int).SetString(cast.ToString(v), 10)
	meta[common.METADATA_AMOUNT] = value

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
	fmt.Println("\n\n\n\nENDPOINT: ConstructionParse")

	if terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tx := &transaction{}

	if !request.Signed {
		// decode unsigned transaction
		b, err := hex.DecodeString(request.Transaction)
		if err != nil {
			terr := common.ErrUnableToParseTx
			terr.Message += err.Error()
			return nil, terr
		}
		err = rlp.DecodeBytes(b, tx)
		if err != nil {
			terr := common.ErrUnableToParseTx
			terr.Message += err.Error()
			return nil, terr
		}
	} else {
		// decode signed transaction
		t := new(tomochaintypes.Transaction)
		b, err := hex.DecodeString(request.Transaction)
		if err != nil {
			terr := common.ErrUnableToParseTx
			terr.Message += err.Error()
			return nil, terr
		}
		err = rlp.DecodeBytes(b, t)
		if err != nil {
			terr := common.ErrUnableToParseTx
			terr.Message += err.Error()
			return nil, terr
		}

		tx.To = t.To().String()
		tx.Value = t.Value()
		tx.Input = t.Data()
		tx.Nonce = t.Nonce()
		tx.GasPrice = t.GasPrice()
		tx.GasLimit = t.Gas()
		tx.ChainID = t.ChainId()

		msg, err := t.AsMessage(tomochaintypes.NewEIP155Signer(t.ChainId()), nil, nil)
		if err != nil {
			terr := common.ErrUnableToGetAccount
			terr.Message += err.Error()
			return nil, terr
		}
		tx.From = msg.From().String()
	}

	// Ensure valid from address
	ok := tomochaincommon.IsHexAddress(tx.From)
	if !ok {
		terr := common.ErrUnableToGetAccount
		terr.Message += fmt.Errorf("%s is not a valid address", tx.From).Error()
		return nil, terr
	}

	// Ensure valid to address
	ok = tomochaincommon.IsHexAddress(tx.From)
	if !ok {
		terr := common.ErrUnableToGetAccount
		terr.Message += fmt.Errorf("%s is not a valid address", tx.To).Error()
		return nil, terr
	}

	ops := []*types.Operation{
		{
			Type: common.TRANSACTION_TYPE_NATIVE_TRANSFER.String(),
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Account: &types.AccountIdentifier{
				Address: tx.From,
			},
			Amount: &types.Amount{
				Value:    new(big.Int).Neg(tx.Value).String(),
				Currency: common.TomoNativeCoin,
			},
		},
		{
			Type: common.TRANSACTION_TYPE_NATIVE_TRANSFER.String(),
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Account: &types.AccountIdentifier{
				Address: tx.To,
			},
			Amount: &types.Amount{
				Value:    tx.Value.String(),
				Currency: common.TomoNativeCoin,
			},
		},
	}

	metaMap := map[string]interface{}{
		common.METADATA_NONCE:     tx.Nonce,
		common.METADATA_GAS_PRICE: tx.GasPrice,
		common.METADATA_CHAIN_ID:  tx.ChainID,
	}
	resp := &types.ConstructionParseResponse{}

	if request.Signed {
		resp = &types.ConstructionParseResponse{
			Operations: ops,
			Signers:    []string{tx.From},
			Metadata:   metaMap,
		}
	} else {
		resp = &types.ConstructionParseResponse{
			Operations: ops,
			Signers:    []string{},
			Metadata:   metaMap,
		}
	}
	return resp, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *constructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	fmt.Println("\n\n\n\nENDPOINT: ConstructionPayloads metadata", request.Operations[0].Account.Address)

	if err := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}
	if len(request.Operations) != 2 {
		terr := common.ErrInvalidInputParam
		terr.Message += "ConstructionPayloadsRequest require 2 operations"
		return nil, terr
	}
	addr := request.Operations[0].Account.Address
	account, err := s.client.GetAccount(ctx, addr)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	nonce, ok := account.Metadata[common.METADATA_SEQUENCE_NUMBER]
	if !ok || nonce == nil {
		return nil, common.ErrUnableToGetNextNonce
	}
	txValue, _ := new(big.Int).SetString(cast.ToString(request.Operations[1].Amount.Value), 10)
	gasPrice, _ := new(big.Int).SetString(cast.ToString(request.Metadata[common.METADATA_GAS_PRICE]), 10)

	var txdata []byte
	if request.Metadata[common.METADATA_TRANSACTION_DATA] != nil {
		txdata = request.Metadata[common.METADATA_TRANSACTION_DATA].([]byte)
	}
	tx := tomochaintypes.NewTransaction(cast.ToUint64(nonce),
		tomochaincommon.HexToAddress(request.Operations[1].Account.Address),
		txValue,
		cast.ToUint64(request.Metadata[common.METADATA_GAS_LIMIT]),
		gasPrice,
		txdata)
	checkFrom := request.Operations[0].Account.Address

	chainId, err := s.client.GetChainID(ctx)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}

	unsignedTx := &transaction{
		From:     checkFrom,
		To:       request.Operations[1].Account.Address,
		Value:    tx.Value(),
		Input:    tx.Data(),
		Nonce:    tx.Nonce(),
		GasPrice: gasPrice,
		GasLimit: tx.Gas(),
		ChainID:  chainId,
	}

	d, err := rlp.EncodeToBytes(unsignedTx)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	unsignedTxEncode := hex.EncodeToString(d)

	signer := tomochaintypes.NewEIP155Signer(chainId)

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedTxEncode,
		Payloads: []*types.SigningPayload{
			{
				Address:       checkFrom,
				Bytes:         signer.Hash(tx).Bytes(),
				SignatureType: types.EcdsaRecovery,
			},
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *constructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	fmt.Println("\n\n\n\nENDPOINT: ConstructionPreprocess", request.Operations[0].Amount)
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
	options[common.METADATA_SYMBOL] = request.Operations[0].Amount.Currency.Symbol
	options[common.METADATA_DECIMALS] = request.Operations[0].Amount.Currency.Decimals

	// recipient
	options[common.METADATA_RECIPIENT] = request.Operations[1].Account.Address
	options[common.METADATA_AMOUNT] = request.Operations[1].Amount.Value

	if request.Metadata[common.METADATA_GAS_LIMIT] != nil {
		options[common.METADATA_GAS_LIMIT] = request.Metadata[common.METADATA_GAS_LIMIT]
	}
	if request.Metadata[common.METADATA_GAS_PRICE] != nil {
		options[common.METADATA_GAS_PRICE] = request.Metadata[common.METADATA_GAS_PRICE]
	}
	fmt.Println("Preprocess Result", options)
	return &types.ConstructionPreprocessResponse{
		Options: options,
	}, nil
}

// ConstructionSubmit implements the /construction/submit \n\n\n\nENDPOINT.
func (s *constructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Println("\n\n\n\nENDPOINT: ConstructionSubmit")
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
	value := tomochaincommon.Big0
	from := tomochaincommon.Address{}
	if v := transaction.Value(); v != nil {
		value = transaction.Value()
	}
	if f := transaction.From(); f != nil {
		from = *f
	}
	return []*types.Operation{
		// sender
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type:   common.TRANSACTION_TYPE_NATIVE_TRANSFER.String(),
			Status: "",
			Account: &types.AccountIdentifier{
				Address: from.String(),
			},
			//FIXME: native transfer only
			Amount: &types.Amount{
				Value:    new(big.Int).Neg(value).String(),
				Currency: common.TomoNativeCoin,
			},
		},

		// recipient
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Type:   common.TRANSACTION_TYPE_NATIVE_TRANSFER.String(),
			Status: "",
			Account: &types.AccountIdentifier{
				Address: (*(transaction.To())).String(),
			},
			//FIXME: native transfer only
			Amount: &types.Amount{
				Value:    value.String(),
				Currency: common.TomoNativeCoin,
			},
		},
	}, nil
}

func parseMetaDataFromTx(transaction *tomochaintypes.Transaction) (map[string]interface{}, error) {
	meta := map[string]interface{}{}

	value := tomochaincommon.Big0
	if v := transaction.Value(); v != nil {
		value = transaction.Value()
	}
	if f := transaction.From(); f != nil {
		meta[common.METADATA_SENDER] = (*f).String()
	}

	meta[common.METADATA_RECIPIENT] = (*(transaction.To())).String()
	meta[common.METADATA_GAS_LIMIT] = transaction.Gas()
	meta[common.METADATA_GAS_PRICE] = transaction.GasPrice().Uint64()
	meta[common.METADATA_TRANSACTION_AMOUNT] = value.Uint64()
	meta[common.METADATA_TRANSACTION_DATA] = transaction.Data()
	return meta, nil
}
