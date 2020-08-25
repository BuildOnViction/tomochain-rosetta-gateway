// Copyright (c) 2020 TomoChain

package services

import (
	"context"
	"encoding/hex"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/protobuf/proto"
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
	if terr := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	tran, err := hex.DecodeString(request.UnsignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}
	if err := proto.Unmarshal(tran, act); err != nil {
		return nil, common.ErrUnmarshal
	}

	if len(request.Signatures) != 1 {
		terr := common.ErrInvalidInputParam
		terr.Message += "need exact 1 signature"
		return nil, terr
	}

	rawPub := request.Signatures[0].PublicKey.Bytes
	if btcec.IsCompressedPubKey(rawPub) {
		pubk, err := btcec.ParsePubKey(rawPub, btcec.S256())
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "invalid pubkey: " + err.Error()
			return nil, terr
		}
		rawPub = pubk.SerializeUncompressed()
	}
	// NOTE set right sender pubkey here
	act.SenderPubKey = rawPub

	rawSig := request.Signatures[0].Bytes
	if len(rawSig) != 65 {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid signature length"
		return nil, terr
	}
	act.Signature = rawSig

	msg, err := proto.Marshal(act)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	return &types.ConstructionCombineResponse{
		SignedTransaction: hex.EncodeToString(msg),
	}, nil
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *constructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	if terr := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	if len(request.PublicKey.Bytes) == 0 || request.PublicKey.CurveType != CurveType {
		terr := common.ErrInvalidInputParam
		terr.Message += "unsupported public key type"
		return nil, terr
	}

	rawPub := request.PublicKey.Bytes
	if btcec.IsCompressedPubKey(rawPub) {
		pubk, err := btcec.ParsePubKey(rawPub, btcec.S256())
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "invalid public key: " + err.Error()
			return nil, terr
		}
		rawPub = pubk.SerializeUncompressed()
	}

	pub, err := crypto.BytesToPublicKey(rawPub)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid public key: " + err.Error()
		return nil, terr
	}
	addr, err := address.FromBytes(pub.Hash())
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid public key: " + err.Error()
		return nil, terr
	}
	return &types.ConstructionDeriveResponse{
		Address: addr.String(),
	}, nil
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *constructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	if terr := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "invalid signed transaction format: " + err.Error()
		return nil, terr
	}
	h := hash.Hash256b(tran)

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			hex.EncodeToString(h[:]),
		},
	}, nil
}

type metadataInputOptions struct {
	senderAddress string
	gasLimit      *uint64
	gasPrice      *uint64
	maxFee        *big.Int
	feeMultiplier *float64
	typ           common.TransactionLogType
}

func parseMetadataInputOptions(options map[string]interface{}) (*metadataInputOptions, *types.Error) {
	opts := &metadataInputOptions{}
	idRaw, ok := options["sender"]
	if !ok {
		terr := common.ErrInvalidInputParam
		terr.Message += "empty sender address"
		return nil, terr
	}

	var err error
	opts.senderAddress, err = cast.ToStringE(idRaw)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}

	if _, ok := options["type"]; !ok {
		terr := common.ErrInvalidInputParam
		terr.Message += "empty operation type"
		return nil, terr
	}
	typ, err := cast.ToStringE(options["type"])
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += "failed to parse type: " + err.Error()
		return nil, terr
	}
	if !common.IsSupportedConstructionType(typ) {
		terr := common.ErrInvalidInputParam
		terr.Message += "unsupported type"
		return nil, terr
	}
	opts.typ = common.TransactionLogType(common.TransactionLogType_value[typ])

	if rawgl, ok := options["gasLimit"]; ok {
		gasLimit, err := cast.ToUint64E(rawgl)
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "failed to parse gasLimit: " + err.Error()
			return nil, terr
		}
		opts.gasLimit = &gasLimit
	}

	if rawgp, ok := options["gasPrice"]; ok {
		gasPrice, err := cast.ToUint64E(rawgp)
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "failed to parse gasPrice: " + err.Error()
			return nil, terr
		}
		opts.gasPrice = &gasPrice
	}

	if rawmp, ok := options["feeMultiplier"]; ok {
		feeMultiplier, err := cast.ToFloat64E(rawmp)
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "failed to parse fee multiplier: " + err.Error()
			return nil, terr
		}
		opts.feeMultiplier = &feeMultiplier
	}

	if rawmf, ok := options["maxFee"]; ok {
		maxFeeStr, err := cast.ToStringE(rawmf)
		if err != nil {
			terr := common.ErrInvalidInputParam
			terr.Message += "failed to parse max fee: " + err.Error()
			return nil, terr
		}
		maxFee, ok := new(big.Int).SetString(maxFeeStr, 10)
		if !ok {
			terr := common.ErrInvalidInputParam
			terr.Message += "failed to parse max fee"
			return nil, terr
		}
		opts.maxFee = maxFee
	}

	return opts, nil
}

func estimateGasAction(opts *metadataInputOptions) (*iotextypes.Action, *types.Error) {
	// need a valid pubkey to estimate, just use one
	rawPrivKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	rawPubKey := rawPrivKey.PubKey()
	act := &iotextypes.Action{
		SenderPubKey: rawPubKey.SerializeUncompressed(),
	}

	switch opts.typ {
	// XXX once support send out payload, need to pass payload here to get right gaslimit
	case iotextypes.TransactionLogType_NATIVE_TRANSFER:
		act.Core = &iotextypes.ActionCore{
			Action: &iotextypes.ActionCore_Transfer{
				Transfer: &iotextypes.Transfer{},
			},
		}
	}
	return act, nil
}

// ConstructionMetadata implements the /construction/metadata endpoint.
func (s *constructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	if terr := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}

	opts, terr := parseMetadataInputOptions(request.Options)
	if terr != nil {
		return nil, terr
	}
	account, err := s.client.GetAccount(ctx, 0, opts.senderAddress)
	if err != nil {
		terr := common.ErrUnableToGetAccount
		terr.Message += err.Error()
		return nil, terr
	}
	meta := account.Metadata

	var gasLimit, gasPrice uint64
	if opts.gasLimit == nil {
		estAct, terr := estimateGasAction(opts)
		if terr != nil {
			return nil, terr
		}
		gasLimit, err = s.client.EstimateGas(ctx, estAct)
		if err != nil {
			terr := common.ErrUnableToEstimateGas
			terr.Message += err.Error()
			return nil, terr
		}
	} else {
		gasLimit = *opts.gasLimit
	}

	if opts.gasPrice == nil {
		gasPrice, err = s.client.SuggestGasPrice(ctx)
		if err != nil {
			terr := common.ErrUnableToGetSuggestGas
			terr.Message += err.Error()
			return nil, terr
		}
	} else {
		gasPrice = *opts.gasPrice
	}

	// apply fee multiplier
	if opts.feeMultiplier != nil {
		base := new(big.Float).SetUint64(gasPrice)
		multiplier := new(big.Float).SetFloat64(*opts.feeMultiplier)
		gasPrice, _ = new(big.Float).Mul(base, multiplier).Uint64()
	}

	meta["gasLimit"] = gasLimit
	meta["gasPrice"] = gasPrice
	suggestedFee := new(big.Int).Mul(
		new(big.Int).SetUint64(gasPrice),
		new(big.Int).SetUint64(gasLimit))

	// check if maxFee >= fee
	if opts.maxFee != nil {
		if opts.maxFee.Cmp(suggestedFee) < 0 {
			return nil, common.ErrExceededFee
		}
	}

	return &types.ConstructionMetadataResponse{
		Metadata: meta,
		SuggestedFee: []*types.Amount{
			&types.Amount{
				Value: suggestedFee.String(),
				Currency: &types.Currency{
					Symbol:   s.client.GetConfig().Currency.Symbol,
					Decimals: s.client.GetConfig().Currency.Decimals,
				},
			},
		},
	}, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *constructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	if terr := common.ValidateNetworkIdentifier(ctx, request.NetworkIdentifier); terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.Transaction)
	if err != nil {
		return nil, common.ErrUnableToParseTx
	}

	act := &iotextypes.Action{}
	if err := proto.Unmarshal(tran, act); err != nil {
		return nil, common.ErrUnableToParseTx
	}

	sender, terr := s.checkIoAction(act, request.Signed)
	if terr != nil {
		return nil, terr
	}
	ops, meta := s.ioActionToOps(sender, act)

	resp := &types.ConstructionParseResponse{
		Operations: ops,
		Metadata:   meta,
	}
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
	if err := common.ValidateNetworkIdentifier(ctx, request.NetworkIdentifier); err != nil {
		return nil, err
	}
	if err := s.checkOperationAndMeta(request.Operations, request.Metadata, true); err != nil {
		return nil, err
	}

	act := s.opsToIoAction(request.Operations, request.Metadata)
	msg, err := proto.Marshal(act)
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	unsignedTx := hex.EncodeToString(msg)

	core, err := proto.Marshal(act.GetCore())
	if err != nil {
		terr := common.ErrServiceInternal
		terr.Message += err.Error()
		return nil, terr
	}
	h := hash.Hash256b(core)
	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedTx,
		Payloads: []*types.SigningPayload{
			&types.SigningPayload{
				Address:       request.Operations[0].Account.Address,
				Bytes:         h[:],
				SignatureType: SignatureType,
			},
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *constructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if err := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier); err != nil {
		return nil, err
	}

	if err := s.checkOperationAndMeta(request.Operations, request.Metadata, false); err != nil {
		return nil, err
	}

	options := make(map[string]interface{})
	options["sender"] = request.Operations[0].Account.Address
	options["type"] = request.Operations[0].Type
	options["amount"] = request.Operations[1].Amount.Value
	options["symbol"] = request.Operations[1].Amount.Currency.Symbol
	options["decimals"] = request.Operations[1].Amount.Currency.Decimals
	options["recipient"] = request.Operations[1].Account.Address

	// XXX it is unclear where these meta data should be
	if request.Metadata["gasLimit"] != nil {
		options["gasLimit"] = request.Metadata["gasLimit"]
	}
	if request.Metadata["gasPrice"] != nil {
		options["gasPrice"] = request.Metadata["gasPrice"]
	}

	// check and set max fee and fee multiplier
	if len(request.MaxFee) != 0 {
		maxFee := request.MaxFee[0]
		if maxFee.Currency.Symbol != s.client.GetConfig().Currency.Symbol ||
			maxFee.Currency.Decimals != s.client.GetConfig().Currency.Decimals {
			terr := common.ErrConstructionCheck
			terr.Message += "invalid currency"
			return nil, terr
		}
		options["maxFee"] = maxFee.Value
	}
	if request.SuggestedFeeMultiplier != nil {
		options["feeMultiplier"] = *request.SuggestedFeeMultiplier
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
	terr := common.ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	tran, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}

	act := &iotextypes.Action{}
	if err := proto.Unmarshal(tran, act); err != nil {
		terr := common.ErrInvalidInputParam
		terr.Message += err.Error()
		return nil, terr
	}

	txID, err := s.client.SubmitTx(ctx, act)
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

func (s *constructionAPIService) opsToIoAction(ops []*types.Operation, meta map[string]interface{}) *iotextypes.Action {
	act := &iotextypes.Action{
		Core: &iotextypes.ActionCore{
			GasLimit: cast.ToUint64(meta["gasLimit"]),
			GasPrice: new(big.Int).SetUint64(cast.ToUint64(meta["gasPrice"])).String(),
			Nonce:    cast.ToUint64(meta["nonce"]),
		},
		// NOTE use SenderPubKey field to temporary pass sender address,
		// since rosetta-cli need to verify sender address in operations.
		SenderPubKey: []byte(ops[0].Account.Address),
	}

	switch iotextypes.TransactionLogType(iotextypes.TransactionLogType_value[ops[0].Type]) {
	case iotextypes.TransactionLogType_NATIVE_TRANSFER:
		act.Core.Action = &iotextypes.ActionCore_Transfer{Transfer: opsToIoTransfer(ops)}
	}
	return act
}

func (s *constructionAPIService) ioActionToOps(sender string, act *iotextypes.Action) ([]*types.Operation, map[string]interface{}) {
	meta := make(map[string]interface{})
	meta["nonce"] = act.GetCore().GetNonce()
	meta["gasLimit"] = act.GetCore().GetGasLimit()
	gasPrice, _ := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	meta["gasPrice"] = gasPrice.Uint64()

	actCore := act.GetCore()
	var ops []*types.Operation
	switch {
	case actCore.GetTransfer() != nil:
		ops = ioTransferToOps(sender, actCore.GetTransfer(), &types.Currency{
			Symbol:   s.client.GetConfig().Currency.Symbol,
			Decimals: s.client.GetConfig().Currency.Decimals,
		})
	}
	return ops, meta
}

func (s *constructionAPIService) checkOperationAndMeta(ops []*types.Operation, meta map[string]interface{}, mustMeta bool) *types.Error {

	terr := common.ErrConstructionCheck
	if len(ops) == 0 {
		terr.Message += "operation numbers are no expected"
		return terr
	}
	typ := ops[0].Type
	if !common.IsSupportedConstructionType(typ) {
		terr.Message += "unsupported construction type"
		return terr
	}
	switch iotextypes.TransactionLogType(iotextypes.TransactionLogType_value[typ]) {
	case iotextypes.TransactionLogType_NATIVE_TRANSFER:
		if terr := checkTransferOps(ops, &types.Currency{
			Symbol:   s.client.GetConfig().Currency.Symbol,
			Decimals: s.client.GetConfig().Currency.Decimals,
		}); terr != nil {
			return terr
		}
	}

	// check metadata exists
	if mustMeta {
		if meta["gasLimit"] == nil || meta["gasPrice"] == nil || meta["nonce"] == nil {
			terr.Message += "metadata not complete"
			return terr
		}
	}

	// check gas
	if meta["gasLimit"] != nil {
		if _, err := cast.ToUint64E(meta["gasLimit"]); err != nil {
			terr.Message += "invalid gas limit"
			return terr
		}
	}
	if meta["gasPrice"] != nil {
		if _, err := cast.ToUint64E(meta["gasPrice"]); err != nil {
			terr.Message += "invalid gas price"
			return terr
		}
	}
	if meta["nonce"] != nil {
		if _, err := cast.ToUint64E(meta["nonce"]); err != nil {
			terr.Message += "invalid nonce"
			return terr
		}
	}
	return nil
}

func (s *constructionAPIService) checkIoAction(act *iotextypes.Action, signed bool) (sender string, terr *types.Error) {
	terr = common.ErrConstructionCheck
	if _, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10); !ok {
		terr.Message += "invalid gas price"
		return "", terr
	}

	if !signed {
		// NOTE use SenderPubKey to pass sender address
		return string(act.GetSenderPubKey()), nil
	}

	// check pubkey and address
	if len(act.GetSenderPubKey()) == 0 {
		terr.Message += "invalid pub key"
		return "", terr
	}

	pub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		terr.Message += "invalid pub key"
		return "", terr
	}
	senderAddr, err := address.FromBytes(pub.Hash())
	if err != nil {
		terr.Message += "invalid io address"
		return "", terr
	}
	sender = senderAddr.String()

	core, err := proto.Marshal(act.GetCore())
	if err != nil {
		terr = common.ErrServiceInternal
		terr.Message += err.Error()
		return "", terr
	}
	h := hash.Hash256b(core)
	if !pub.Verify(h[:], act.GetSignature()) {
		terr.Message += "invalid signature"
		return "", terr
	}
	return sender, nil
}
