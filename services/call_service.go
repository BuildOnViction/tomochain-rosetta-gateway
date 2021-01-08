package services


import (
	"context"
	"errors"
	"github.com/tomochain/tomochain-rosetta-gateway/common"

	"github.com/tomochain/tomochain-rosetta-gateway/configuration"
	"github.com/tomochain/tomochain-rosetta-gateway/tomochain"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// CallAPIService implements the server.CallAPIServicer interface.
type CallAPIService struct {
	config *configuration.Configuration
	client Client
}

// NewCallAPIService creates a new instance of a CallAPIService.
func NewCallAPIService(cfg *configuration.Configuration, client Client) *CallAPIService {
	return &CallAPIService{
		config: cfg,
		client: client,
	}
}

// Call implements the /call endpoint.
func (s *CallAPIService) Call(
	ctx context.Context,
	request *types.CallRequest,
) (*types.CallResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, common.ErrUnavailableOffline
	}

	response, err := s.client.Call(ctx, request)
	if errors.Is(err, tomochain.ErrCallParametersInvalid) {
		return nil, common.ErrCallParametersInvalid
	}
	if errors.Is(err, tomochain.ErrCallOutputMarshal) {
		return nil, common.ErrCallOutputMarshal
	}
	if errors.Is(err, tomochain.ErrCallMethodInvalid) {
		return nil, common.ErrCallMethodInvalid
	}
	if err != nil {
		return nil, common.ErrTomo
	}

	return response, nil
}
