// Copyright (c) 2020 TomoChain

package configuration

import (
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/tomochain/tomochain-rosetta-gateway/tomochain"
	"github.com/tomochain/tomochain/params"
	"math/big"
	"os"
	"strconv"
)

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// TomoEnv is an optional environment variable
	// used to connect tomochain-rosetta to an already
	// running tomo node.
	TomoEnv = "TOMO"

	// DefaultTomoURL is the default URL for
	// a running geth node. This is used
	// when GethEnv is not populated.
	DefaultTomoURL = "http://localhost:8545"

	// Mainnet is the TomoChain Mainnet.
	Mainnet string = "MAINNET"

	// Testnet is TomoChain Public Testnet
	Testnet string = "TESTNET"

	// Devnet is TomoChain network for development
	Devnet string = "DEVNET"
)

var (
	// MiddlewareVersion is the version of tomochain-rosetta.
	MiddlewareVersion = "0.0.2"
)

type Mode string

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	TomoURL                string
	RemoteTomo             bool
	Port                   int
	TomoArguments          string

	Params *params.ChainConfig
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*Configuration, error) {
	config := &Configuration{}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: tomochain.Blockchain,
			Network:    tomochain.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = tomochain.MainnetGenesisBlockIdentifier
		config.Params = params.TomoMainnetChainConfig
		config.TomoArguments = tomochain.MainnetTomoArguments
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: tomochain.Blockchain,
			Network:    tomochain.TestnetNetwork,
		}
		config.GenesisBlockIdentifier = &types.BlockIdentifier{
			Hash: "",
			Index: tomochain.GenesisBlockIndex,
		}
		testnetChainConfig := params.TomoMainnetChainConfig
		testnetChainConfig.ChainId = new(big.Int).SetUint64(cast.ToUint64(tomochain.TestnetNetwork))
		config.Params = testnetChainConfig
		config.TomoArguments = tomochain.TestnetTomoArguments
	case Devnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: tomochain.Blockchain,
			Network:    tomochain.DevnetNetwork,
		}
		config.GenesisBlockIdentifier = &types.BlockIdentifier{
			Hash: "0x8853b9408238a5f3d33fce226114b2ac85c487a9e1978979c30b31e2403ed575",
			Index: tomochain.GenesisBlockIndex,
		}
		devnetChainConfig := params.TomoMainnetChainConfig
		devnetChainConfig.ChainId = new(big.Int).SetUint64(cast.ToUint64(tomochain.DevnetNetwork))
		config.Params = devnetChainConfig
		config.TomoArguments = tomochain.DevnetTomoArguments
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.TomoURL = DefaultTomoURL
	envGethURL := os.Getenv(TomoEnv)
	if len(envGethURL) > 0 {
		config.RemoteTomo = true
		config.TomoURL = envGethURL
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}
