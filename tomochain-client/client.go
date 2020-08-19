// Copyright (c) 2020 TomoChain

package tomochain_client

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	"github.com/tomochain/tomochain/common"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/ethclient"
	"math/big"
	"strconv"
	"sync"
)

type (
	// TomoChainClient is the TomoChain blockchain client interface.
	TomoChainClient interface {
		// GetChainID returns the network chain context, derived from the
		// genesis document.
		GetChainID(ctx context.Context) (*big.Int, error)

		// GetBlockByNumber returns the TomoChain block at given height.
		GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)

		// GetBlockByHash returns the TomoChain block at given hash.
		GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)

		// GetLatestBlock returns latest TomoChain block.
		GetLatestBlock(ctx context.Context) (*types.Block, error)

		// GetGenesisBlock returns the TomoChain genesis block.
		GetGenesisBlock(ctx context.Context) (*types.Block, error)

		// GetAccount returns the TomoChain staking account for given owner address
		// at given height.
		GetAccount(ctx context.Context, blockHash common.Hash, owner string) (*types.AccountBalanceResponse, error)

		// SubmitTx submits the given encoded transaction to the node.
		SubmitTx(ctx context.Context, tx tomochaintypes.Transaction) (txid string, err error)

		// GetTransactions returns transactions of the block.
		GetTransactions(ctx context.Context, hash common.Hash) ([]*types.Transaction, error)

		// GetConfig returns the config.
		GetConfig() *config.Config

		SuggestGasPrice(ctx context.Context) (uint64, error)

		EstimateGas(ctx context.Context, msg tomochain.CallMsg) (uint64, error)
	}
)

type (
	// TomoChainRpcClient is an implementation of TomoChain client using local rpc/ipc connection.
	TomoChainRpcClient struct {
		sync.RWMutex
		rpcConn *ethclient.Client
		cfg     *config.Config
	}
)

// NewTomoChainClient returns an implementation of TomoChainClient
func NewTomoChainClient(cfg *config.Config) (cli *TomoChainRpcClient, err error) {
	rpcClient, err := ethclient.Dial(cfg.Server.Endpoint)
	if err != nil {
		return nil, err
	}
	return &TomoChainRpcClient{
		rpcConn: rpcClient,
		cfg:     cfg,
	}, nil
}

func (c *TomoChainRpcClient) GetChainID(ctx context.Context) (*big.Int, error) {
	id, err := strconv.Atoi(c.cfg.NetworkIdentifier.Network)
	if err != nil {
		return common.Big0, err
	}
	return big.NewInt(int64(id)), nil
}

func (c *TomoChainRpcClient) GetBlockByNumber(ctx context.Context, number *big.Int) (ret *types.Block, err error) {
	block, err := c.rpcConn.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) GetBlockByHash(ctx context.Context, hash common.Hash) (ret *types.Block, err error) {
	block, err := c.rpcConn.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	return c.GetBlockByNumber(ctx, nil)
}

func (c *TomoChainRpcClient) GetGenesisBlock(ctx context.Context) (*types.Block, error) {
	block, err := c.rpcConn.BlockByNumber(ctx, common.Big0)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) SuggestGasPrice(ctx context.Context) (uint64, error) {
	suggestedGasPrice, err := c.rpcConn.SuggestGasPrice(ctx)
	if err != nil {
		return uint64(0), err
	}
	return suggestedGasPrice.Uint64(), nil
}

func (c *TomoChainRpcClient) EstimateGas(ctx context.Context, msg tomochain.CallMsg) (uint64, error) {
	gas, err := c.rpcConn.EstimateGas(ctx, msg)
	if err != nil {
		return uint64(0), err
	}
	return gas, err
}

func (c *TomoChainRpcClient) GetAccount(ctx context.Context, blockHash common.Hash, owner string) (ret *types.AccountBalanceResponse, err error) {
	//TODO: get account metadata
	return ret, nil
}

func (c *TomoChainRpcClient) GetTransactions(ctx context.Context, hash common.Hash) (ret []*types.Transaction, err error) {
	block, err := c.rpcConn.BlockByHash(ctx, hash)
	if err != nil {
		return []*types.Transaction{}, err
	}
	return c.PackTransaction(block.Transactions()), nil
}

func (c *TomoChainRpcClient) SubmitTx(ctx context.Context, tx tomochaintypes.Transaction) (txid string, err error) {
	// TODO
	return txid, nil
}

func (c *TomoChainRpcClient) GetConfig() *config.Config {
	return c.cfg
}

func (c *TomoChainRpcClient) PackBlockData(block *tomochaintypes.Block) (ret *types.Block) {
	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: block.Number().Int64(),
			Hash:  block.Hash().String(),
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: block.Number().Int64() - 1,
			Hash:  block.ParentHash().String(),
		},
		Timestamp:    new(big.Int).Mul(block.Time(), big.NewInt(1000)).Int64(),
		Transactions: c.PackTransaction(block.Transactions()),
		Metadata:     map[string]interface{}{},
	}
}

func (c *TomoChainRpcClient) PackTransaction(transactions tomochaintypes.Transactions) (ret []*types.Transaction) {
	//TODO
	return ret
}
