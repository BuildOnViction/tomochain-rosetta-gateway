// Copyright (c) 2020 TomoChain

package tomochain_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	tomochaincommon "github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/consensus/posv"
	tomochaintypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
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
		GetAccount(ctx context.Context, owner string) (*types.AccountBalanceResponse, error)

		// SubmitTx submits the given encoded transaction to the node.
		SubmitTx(ctx context.Context, signedTx hexutil.Bytes) (txid string, err error)

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

func (c *TomoChainRpcClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	//FIXME: ethClient.BlockByNumber wrapped block data, always return block without validator data
	//block, err := c.ethClient.BlockByNumber(ctx, number)

	client, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	var raw json.RawMessage
	var blockNumber interface{}
	if number != nil {
		blockNumber = hexutil.EncodeBig(number)
	} else {
		blockNumber = "latest"
	}
	err = client.CallContext(ctx, &raw, common.RPC_METHOD_GET_BLOCK_BY_NUMBER, blockNumber, true)
	if err != nil {
		fmt.Println("getBlockByNumber error when calling rpc")
		return nil, err
	}
	var data map[string]interface{}
	if err = json.Unmarshal(raw, &data); err != nil {
		fmt.Println("getBlockByNumber error when unmarshalling raw data")
		return nil, err
	}
	hash := ""
	if data["hash"] != nil {
		hash = (data["hash"]).(string)
	}
	header := &tomochaintypes.Header{}
	body := &tomochaintypes.Body{}
	if err := json.Unmarshal(raw, &header); err != nil {
		fmt.Println("GetBlockByNumber: Unmarshal header error", err)
		return nil, err
	}
	var coinbase tomochaincommon.Address
	if header.Number.Uint64() > 0 {
		coinbase, err = GetCoinbaseFromHeader(header)
		if err != nil {
			fmt.Println("GetBlockByNumber: error get CoinbaseFromHeader", err)
			return nil, err
		}
	}
	header.Coinbase = coinbase
	if err := json.Unmarshal(raw, &body); err != nil {
		fmt.Println("GetBlockByNumber: Unmarshal body error", err)
		return nil, err
	}
	block := tomochaintypes.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles)
	result := &types.Block{}
	if result, err = c.PackBlockData(ctx, block, tomochaincommon.HexToHash(hash)); err != nil {
		fmt.Println("GetBlockByNumber: packing block data error", err)
		return nil, err
	}
	return result, nil
}

func (c *TomoChainRpcClient) GetBlockByHash(ctx context.Context, hash tomochaincommon.Hash) (res *types.Block, err error) {
	//FIXME: ethClient.GetBlockByHash wrapped block data, always return block without validator data
	//block, err := c.ethClient.GetBlockByHash(ctx, hash)

	client, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	var raw json.RawMessage
	err = client.CallContext(ctx, &raw, common.RPC_METHOD_GET_BLOCK_BY_HASH, hash, true)
	if err != nil {
		fmt.Println("GetBlockByHash: Call RPC err", err)
		return nil, err
	}
	var data map[string]interface{}
	if err = json.Unmarshal(raw, &data); err != nil {
		fmt.Println("GetBlockByHash: Unmarshal block data error", err)
		return nil, err
	}
	hashString := ""
	if data["hash"] != nil {
		hashString = (data["hash"]).(string)
	}
	header := &tomochaintypes.Header{}
	body := &tomochaintypes.Body{}
	if err := json.Unmarshal(raw, &header); err != nil {
		fmt.Println("GetBlockByHash: Unmarshal header error", err)
		return nil, err
	}
	coinbase, err := GetCoinbaseFromHeader(header)
	if err != nil {
		return nil, err
	}
	header.Coinbase = coinbase
	if err := json.Unmarshal(raw, &body); err != nil {
		fmt.Println("GetBlockByHash: Unmarshal body error", err)
		return nil, err
	}
	block := tomochaintypes.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles)
	result := &types.Block{}
	if result, err = c.PackBlockData(ctx, block, tomochaincommon.HexToHash(hashString)); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *TomoChainRpcClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	return c.GetBlockByNumber(ctx, nil)
}

func (c *TomoChainRpcClient) GetGenesisBlock(ctx context.Context) (*types.Block, error) {
	return c.GetBlockByNumber(ctx, tomochaincommon.Big0)
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

func (c *TomoChainRpcClient) GetAccount(ctx context.Context, owner string) (res *types.AccountBalanceResponse, err error) {
	block, err := c.GetBlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	res = &types.AccountBalanceResponse{}
	res.BlockIdentifier = block.BlockIdentifier

	balance, err := c.ethClient.BalanceAt(ctx, tomochaincommon.HexToAddress(owner), nil)
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
	nonce, err := c.ethClient.NonceAt(ctx, tomochaincommon.HexToAddress(owner), nil)
	if err != nil {
		return nil, err
	}
	res.Metadata = map[string]interface{}{
		common.METADATA_SEQUENCE_NUMBER: nonce,
	}
	return res, nil
}

func (c *TomoChainRpcClient) GetBlockTransactions(ctx context.Context, hash tomochaincommon.Hash) (res []*types.Transaction, err error) {
	block, err := c.ethClient.BlockByHash(ctx, hash)
	if err != nil {
		return []*types.Transaction{}, err
	}
	return c.PackTransaction(ctx, block, hash)
}

func (c *TomoChainRpcClient) SubmitTx(ctx context.Context, signedTx hexutil.Bytes) (string, error) {
	rpcClient, err := c.ConnectRpc()
	if err != nil {
		return "", err
	}
	defer rpcClient.Close()

	hash := tomochaincommon.Hash{}
	err = rpcClient.CallContext(ctx, &hash, common.RPC_METHOD_SEND_SIGNED_TRANSACTION, signedTx)
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func (c *TomoChainRpcClient) GetConfig() *config.Config {
	return c.cfg
}

func (c *TomoChainRpcClient) PackBlockData(ctx context.Context, block *tomochaintypes.Block, finalBlockHash tomochaincommon.Hash) (*types.Block, error) {
	if block == nil {
		return nil, nil
	}
	var parent *types.BlockIdentifier
	if block.Number().Int64() > 0 {
		parent = &types.BlockIdentifier{
			Index: block.Number().Int64() - 1,
			Hash:  block.ParentHash().String(),
		}
	}
	var (
		transactions []*types.Transaction
		err          error
	)

	if block.NumberU64() > 0 {
		transactions, err = c.PackTransaction(ctx, block, finalBlockHash)
		if err != nil {
			fmt.Println("PackBlockData error when packing Transaction")
			return nil, err
		}
	}

	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: block.Number().Int64(),
			Hash:  finalBlockHash.String(),
		},
		ParentBlockIdentifier: parent,
		Timestamp:             new(big.Int).Mul(block.Time(), big.NewInt(1000)).Int64(),
		Transactions:          transactions,
	}, nil
}

func (c *TomoChainRpcClient) PackTransaction(ctx context.Context, block *tomochaintypes.Block, finalBlockHash tomochaincommon.Hash) ([]*types.Transaction, error) {
	result := []*types.Transaction{}
	blockNumber := block.Number()
	previousBlockNumber := new(big.Int).Sub(blockNumber, tomochaincommon.Big1)
	transactions := block.Transactions()
	sealer := block.Header().Coinbase

	balances := map[tomochaincommon.Address]*big.Int{}

	// balance of sealer of this block
	sealerBalance, err := c.ethClient.BalanceAt(ctx, sealer, previousBlockNumber)
	if err != nil {
		return []*types.Transaction{}, err
	}
	balances[sealer] = sealerBalance

	for _, tx := range transactions {
		var (
			fromBalance, toBalance *big.Int
			ok                     bool
			err                    error
		)
		from := *tx.From()
		to := *tx.To()
		if fromBalance, ok = balances[from]; !ok {
			fromBalance, err = c.ethClient.BalanceAt(ctx, from, previousBlockNumber)
			if err != nil {
				return []*types.Transaction{}, err
			}
			balances[from] = fromBalance
		}
		if toBalance, ok = balances[to]; !ok {
			toBalance, err = c.ethClient.BalanceAt(ctx, to, previousBlockNumber)
			if err != nil {
				return []*types.Transaction{}, err
			}
			balances[to] = toBalance
		}

		// get transaction receipt
		status := ""
		transactionReceipt, err := c.ethClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil || transactionReceipt == nil || transactionReceipt.Status == uint(0) {
			status = common.FAIL
		} else {
			status = common.SUCCESS
		}
		gasUsed := tomochaincommon.Big0
		if transactionReceipt != nil {
			gasUsed = new(big.Int).SetUint64(transactionReceipt.GasUsed)
		}
		fee := new(big.Int).Mul(gasUsed, tx.GasPrice())

		// update new balance
		totalValue := new(big.Int).Add(fee, tx.Value())
		balances[from] = new(big.Int).Sub(balances[from], totalValue)
		balances[to] = new(big.Int).Add(balances[to], tx.Value())
		balances[sealer] = new(big.Int).Add(balances[sealer], fee)

		result = append(result, &types.Transaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: tx.Hash().String(),
			},
			Operations: []*types.Operation{
				// sender
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: 0,
					},
					RelatedOperations: nil,
					Type:              common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_NATIVE_TRANSFER)],
					Status:            status,
					Account: &types.AccountIdentifier{
						Address: (*tx.From()).String(),
					},
					Amount: &types.Amount{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Value:    new(big.Int).Sub(new(big.Int).SetUint64(0), totalValue).String(), // balance change of sender should be negative
						Currency: common.TomoNativeCoin,
					},
					Metadata: map[string]interface{}{
						common.METADATA_NEW_BALANCE: balances[from].String(),
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
					Type:   common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_NATIVE_TRANSFER)],
					Status: status,
					Account: &types.AccountIdentifier{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Address: (*(tx.To())).String(),
					},
					Amount: &types.Amount{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Value:    tx.Value().String(),
						Currency: common.TomoNativeCoin,
					},
					Metadata: map[string]interface{}{
						common.METADATA_NEW_BALANCE: balances[to].String(),
					},
				},

				// fee: send to sealer
				{
					OperationIdentifier: &types.OperationIdentifier{
						Index: 2,
					},
					RelatedOperations: []*types.OperationIdentifier{
						{
							Index: 0,
						},
					},
					Type:   common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_GAS_FEE)],
					Status: status,
					Account: &types.AccountIdentifier{
						Address: sealer.String(),
					},
					Amount: &types.Amount{
						//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
						Value:    fee.String(),
						Currency: common.TomoNativeCoin,
					},
					Metadata: map[string]interface{}{
						common.METADATA_NEW_BALANCE: balances[sealer].String(),
					},
				},
			},
		})
	}

	// at check point, distribute rewards
	if block.NumberU64()%common.Epoch == 0 {
		rewards, err := c.GetBlockReward(ctx, finalBlockHash)
		if err != nil {
			return []*types.Transaction{}, err
		}
		if rewards != nil {
			rewardOperations := []*types.Operation{}
			for _, signer := range rewards {
				for holder, amount := range signer {
					holderAddress := tomochaincommon.HexToAddress(holder)
					if holderBalance, ok := balances[holderAddress]; !ok {
						holderBalance, err = c.ethClient.BalanceAt(ctx, holderAddress, previousBlockNumber)
						if err != nil {
							return []*types.Transaction{}, err
						}
						balances[holderAddress] = holderBalance
					}
					balances[holderAddress] = new(big.Int).Add(balances[holderAddress], amount)

					rewardOperations = append(rewardOperations, &types.Operation{
						OperationIdentifier: &types.OperationIdentifier{
							Index: int64(len(rewardOperations)),
						},
						RelatedOperations: nil,
						Type:              common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_CLAIM_FROM_REWARDING_FUND)],
						Status:            common.SUCCESS,
						Account: &types.AccountIdentifier{
							Address: holder,
						},
						Amount: &types.Amount{
							Value:    amount.String(),
							Currency: common.TomoNativeCoin,
						},
						Metadata: map[string]interface{}{
							common.METADATA_NEW_BALANCE: balances[holderAddress].String(),
						},
					})
				}
			}
			result = append(result, &types.Transaction{
				TransactionIdentifier: &types.TransactionIdentifier{
					Hash: finalBlockHash.String(),
				},
				Operations: rewardOperations,
			})
		}
	}
	return result, nil
}

// GetMempool returns all transactions in mempool
func (c *TomoChainRpcClient) GetMempool(ctx context.Context) ([]tomochaincommon.Hash, error) {
	rpcClient, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	pendingTxs := []*common.RPCTransaction{}
	err = rpcClient.CallContext(ctx, &pendingTxs, common.RPC_METHOD_GET_PENDING_TRANSACTIONS)
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
	err = rpcClient.CallContext(ctx, &pendingTxs, common.RPC_METHOD_GET_PENDING_TRANSACTIONS)
	if err != nil {
		return nil, err
	}

	for _, tx := range pendingTxs {
		if tx.Hash.String() == hash.String() {
			fromBalance, err := c.ethClient.BalanceAt(ctx, tx.From, nil)
			if err != nil {
				return nil, err
			}
			toBalance, err := c.ethClient.BalanceAt(ctx, *(tx.To), nil)
			if err != nil {
				return nil, err
			}
			//TODO: not support internal transaction yet
			return &types.Transaction{
				TransactionIdentifier: &types.TransactionIdentifier{
					Hash: tx.Hash.String(),
				},
				Operations: []*types.Operation{
					// sender
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 0,
						},
						RelatedOperations: nil,
						Type:              common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_NATIVE_TRANSFER)],
						Status:            common.PENDING,
						Account: &types.AccountIdentifier{
							Address: tx.From.String(),
						},
						Amount: &types.Amount{
							//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
							Value:    new(big.Int).Sub(new(big.Int).SetUint64(0), tx.Value.ToInt()).String(),
							Currency: common.TomoNativeCoin,
						},
						Metadata: map[string]interface{}{
							common.METADATA_NEW_BALANCE: new(big.Int).Sub(fromBalance, tx.Value.ToInt()).String(),
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
						Type:   common.TRANSACTION_TYPE_NAME[int32(common.TRANSACTION_TYPE_NATIVE_TRANSFER)],
						Status: common.PENDING,
						Account: &types.AccountIdentifier{
							//TODO: support native transfer only, not support internal transaction (transfer from contract) yet
							Address: (*(tx.To)).String(),
						},
						Amount: &types.Amount{
							//TODO: right for native transfer only, wrong for internal transaction with other tokens or contract transfer
							Value:    tx.Value.ToInt().String(),
							Currency: common.TomoNativeCoin,
						},
						Metadata: map[string]interface{}{
							common.METADATA_NEW_BALANCE: new(big.Int).Add(toBalance, tx.Value.ToInt()).String(),
						},
					},
				},
			}, nil
		}
	}
	return nil, nil
}

// GetBlockReward returns rewards of checkpoint block
func (c *TomoChainRpcClient) GetBlockReward(ctx context.Context, hash tomochaincommon.Hash) (map[string]map[string]*big.Int, error) {
	rpcClient, err := c.ConnectRpc()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	// result structure
	//{
	//	"rewards":{
	//	"signer":{
	//		"holder":amount,
	//		"holder":amount
	//	},
	//	"signer":{
	//		"holder":amount,
	//		"holder":amount
	//	},
	//	"signer":{
	//		"holder":amount,
	//		"holder":amount
	//	},
	//},
	//}
	rewards := map[string]map[string]map[string]*big.Int{}
	if err = rpcClient.CallContext(ctx, &rewards, common.RPC_METHOD_GET_REWARD_BY_HASH, hash); err != nil {
		return nil, err
	}
	if rewards["rewards"] != nil {
		return rewards["rewards"], nil
	}
	return nil, nil
}

func PubToAddress(pubkey []byte) tomochaincommon.Address {
	var address tomochaincommon.Address
	copy(address[:], crypto.Keccak256(pubkey[1:])[12:])
	return address
}

func GetCoinbaseFromHeader(header *tomochaintypes.Header) (tomochaincommon.Address, error) {
	if len(header.Extra) < common.ExtraSeal {
		return tomochaincommon.Address{}, fmt.Errorf("extra-data %d byte suffix signature missing", common.ExtraSeal)
	}
	signature := header.Extra[len(header.Extra)-common.ExtraSeal:]
	pubkey, err := crypto.Ecrecover(posv.SigHash(header).Bytes(), signature)
	if err != nil {
		return tomochaincommon.Address{}, err
	}
	return PubToAddress(pubkey), nil
}
