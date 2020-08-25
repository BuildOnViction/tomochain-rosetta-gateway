// Copyright (c) 2020 TomoChain

package tomochain_client

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	tomochaincommon "github.com/tomochain/tomochain/common"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/ethclient"
	"github.com/tomochain/tomochain/rpc"
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
		GetBlockByHash(ctx context.Context, hash tomochaincommon.Hash) (*types.Block, error)

		// GetLatestBlock returns latest TomoChain block.
		GetLatestBlock(ctx context.Context) (*types.Block, error)

		// GetGenesisBlock returns the TomoChain genesis block.
		GetGenesisBlock(ctx context.Context) (*types.Block, error)

		// GetAccount returns the TomoChain staking account for given owner address
		// at given height.
		GetAccount(ctx context.Context, blockHash tomochaincommon.Hash, owner string) (*types.AccountBalanceResponse, error)

		// SubmitTx submits the given encoded transaction to the node.
		SubmitTx(ctx context.Context, tx tomochaintypes.Transaction) (txid string, err error)

		// GetBlockTransactions returns transactions of the block.
		GetBlockTransactions(ctx context.Context, hash tomochaincommon.Hash) ([]*types.Transaction, error)

		// GetMempool returns all transactions in mempool
		GetMempool(ctx context.Context) ([]tomochaincommon.Hash, error)

		// GetMempoolTransactions returns the specified transaction according to the hash in the mempool
		GetMempoolTransaction(ctx context.Context, hash tomochaincommon.Hash) (*types.Transaction, error)

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
		ethClient *ethclient.Client
		cfg       *config.Config
	}
)

// NewTomoChainClient returns an implementation of TomoChainClient
func NewTomoChainClient(cfg *config.Config) (cli *TomoChainRpcClient, err error) {
	ethClient, err := ethclient.Dial(cfg.Server.Endpoint)
	if err != nil {
		return nil, err
	}
	return &TomoChainRpcClient{
		ethClient: ethClient,
		cfg:       cfg,
	}, nil
}

func (c *TomoChainRpcClient) ConnectRpc() (*rpc.Client, error) {
	return rpc.Dial(c.cfg.Server.Endpoint)
}

func (c *TomoChainRpcClient) GetChainID(ctx context.Context) (*big.Int, error) {
	id, err := strconv.Atoi(c.cfg.NetworkIdentifier.Network)
	if err != nil {
		return tomochaincommon.Big0, err
	}
	return big.NewInt(int64(id)), nil
}

func (c *TomoChainRpcClient) GetBlockByNumber(ctx context.Context, number *big.Int) (res *types.Block, err error) {
	block, err := c.ethClient.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) GetBlockByHash(ctx context.Context, hash tomochaincommon.Hash) (res *types.Block, err error) {
	block, err := c.ethClient.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	return c.GetBlockByNumber(ctx, nil)
}

func (c *TomoChainRpcClient) GetGenesisBlock(ctx context.Context) (*types.Block, error) {
	block, err := c.ethClient.BlockByNumber(ctx, tomochaincommon.Big0)
	if err != nil {
		return nil, err
	}
	return c.PackBlockData(block), nil
}

func (c *TomoChainRpcClient) SuggestGasPrice(ctx context.Context) (uint64, error) {
	suggestedGasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return uint64(0), err
	}
	return suggestedGasPrice.Uint64(), nil
}

func (c *TomoChainRpcClient) EstimateGas(ctx context.Context, msg tomochain.CallMsg) (uint64, error) {
	gas, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return uint64(0), err
	}
	return gas, err
}

//TODO: internal transfer via smart contract must be done
// https://www.rosetta-api.org/docs/all_balance_changing.html

func (c *TomoChainRpcClient) GetAccount(ctx context.Context, blockHash tomochaincommon.Hash, owner string) (res *types.AccountBalanceResponse, err error) {
	block, err := c.GetBlockByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	res = &types.AccountBalanceResponse{}
	res.BlockIdentifier = block.BlockIdentifier

	balance, err := c.ethClient.BalanceAt(ctx, tomochaincommon.HexToAddress(owner), big.NewInt(block.BlockIdentifier.Index))
	if err != nil {
		return nil, err
	}
	// TODO: support native coin TOMO only, tokens are not available yet
	res.Balances = []*types.Amount{
		{
			Value:    balance.String(),
			Currency: common.TomoNativeCoin,
		},
	}
	res.Coins = nil

	// attach nonce
	nonce, err := c.ethClient.NonceAt(ctx, tomochaincommon.HexToAddress(owner), big.NewInt(block.BlockIdentifier.Index))
	if err != nil {
		return nil, err
	}
	res.Metadata = map[string]interface{}{
		"sequence_number": nonce,
	}
	return res, nil
}

func (c *TomoChainRpcClient) GetBlockTransactions(ctx context.Context, hash tomochaincommon.Hash) (res []*types.Transaction, err error) {
	block, err := c.ethClient.BlockByHash(ctx, hash)
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

func (c *TomoChainRpcClient) PackBlockData(block *tomochaintypes.Block) (res *types.Block) {
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
	}
}

func (c *TomoChainRpcClient) PackTransaction(transactions tomochaintypes.Transactions) []*types.Transaction {
	result := []*types.Transaction{}
	for _, tx := range transactions {
		result = append(result, &types.Transaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: tx.Hash().String(),
			},
			Operations: []*types.Operation{
				// sender
				{
					OperationIdentifier: nil,
					RelatedOperations:   nil,
					Type:                common.TransactionLogType_name[int32(common.TransactionLogType_NATIVE_TRANSFER)],
					Status:              common.StatusSuccess,
					Account: &types.AccountIdentifier{
						Address: tx.From().String(),
					},
					Amount: &types.Amount{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Value:    tx.Value().Text(10),
						Currency: common.TomoNativeCoin,
					},
				},
				// recipient
				{
					OperationIdentifier: nil,
					RelatedOperations:   nil,
					Type:                common.TransactionLogType_name[int32(common.TransactionLogType_NATIVE_TRANSFER)],
					Status:              common.StatusSuccess,
					Account: &types.AccountIdentifier{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Address: (*(tx.To())).String(),
					},
					Amount: &types.Amount{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Value:    tx.Value().String(),
						Currency: common.TomoNativeCoin,
					},
				},
			},
		})
	}
	return result
}

// GetMempool returns all transactions in mempool
func (c *TomoChainRpcClient) GetMempool(ctx context.Context) ([]tomochaincommon.Hash, error) {
	rpcClient, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	pendingTxs := []*common.RPCTransaction{}
	err = rpcClient.CallContext(ctx, &pendingTxs, "eth_pendingTransactions")
	if err != nil {
		return nil, err
	}
	pendingTxHash := []tomochaincommon.Hash{}
	for _, tx := range pendingTxs {
		pendingTxHash = append(pendingTxHash, tx.Hash)
	}
	return pendingTxHash, nil
}

// GetMempoolTransactions returns the specified transaction according to the hash in the mempool
func (c *TomoChainRpcClient) GetMempoolTransaction(ctx context.Context, hash tomochaincommon.Hash) (*types.Transaction, error) {
	rpcClient, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	pendingTxs := []*common.RPCTransaction{}
	err = rpcClient.CallContext(ctx, &pendingTxs, "eth_pendingTransactions")
	if err != nil {
		return nil, err
	}

	for _, tx := range pendingTxs {
		if tx.Hash.String() == hash.String() {
			//TODO: not support internal transaction yet
			return &types.Transaction{
				TransactionIdentifier: &types.TransactionIdentifier{
					Hash: tx.Hash.String(),
				},
				Operations: []*types.Operation{
					// sender
					{
						OperationIdentifier: nil,
						RelatedOperations:   nil,
						Type:                common.TransactionLogType_name[int32(common.TransactionLogType_NATIVE_TRANSFER)],
						Status:              common.StatusSuccess,
						Account: &types.AccountIdentifier{
							Address: tx.From.String(),
						},
						Amount: &types.Amount{
							//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
							Value:    tx.Value.ToInt().String(),
							Currency: common.TomoNativeCoin,
						},
					},
					// recipient
					{
						OperationIdentifier: nil,
						RelatedOperations:   nil,
						Type:                common.TransactionLogType_name[int32(common.TransactionLogType_NATIVE_TRANSFER)],
						Status:              common.StatusSuccess,
						Account: &types.AccountIdentifier{
							//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
							Address: (*(tx.To)).String(),
						},
						Amount: &types.Amount{
							//TODO: right for native transfer only, wrong for internal transaction with other tokens or contract transfer
							Value:    tx.Value.ToInt().String(),
							Currency: common.TomoNativeCoin,
						},
					},
				},
			}, nil
		}
	}
	return nil, nil
}
