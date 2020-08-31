// Copyright (c) 2020 TomoChain

package main

import (
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	"github.com/tomochain/tomochain-rosetta-gateway/config"
	"github.com/tomochain/tomochain-rosetta-gateway/services"
	tc "github.com/tomochain/tomochain-rosetta-gateway/tomochain-client"
	"log"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	ConfigPath = "ConfigPath"
)

// NewBlockchainRouter returns a Mux http.Handler from a collection of
// Rosetta service controllers.
func NewBlockchainRouter(client tc.TomoChainClient) (http.Handler, error) {
	asserter, err := asserter.NewServer(common.SupportedOperationTypes(),
		false,
		[]*types.NetworkIdentifier{
			{
				Blockchain: client.GetConfig().NetworkIdentifier.Blockchain,
				Network:    client.GetConfig().NetworkIdentifier.Network,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), asserter)
	accountAPIController := server.NewAccountAPIController(services.NewAccountAPIService(client), asserter)
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client), asserter)
	mempoolAPIController := server.NewMempoolAPIController(services.NewMempoolAPIService(client), asserter)
	constructionAPIController := server.NewConstructionAPIController(services.NewConstructionAPIService(client), asserter)
	r := server.NewRouter(networkAPIController, accountAPIController, blockAPIController, mempoolAPIController, constructionAPIController)
	return server.CorsMiddleware(server.LoggerMiddleware(r)), nil
}

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("ERROR: Failed to parse config: %v\n", err)
	}
	client, err := tc.NewTomoChainClient(cfg)
	if err != nil {
		log.Fatalf("ERROR: Failed to prepare TomoChain RPC client: %v\n", err)
	}

	// Start the server.
	router, err := NewBlockchainRouter(client)
	if err != nil {
		log.Fatalf("ERROR: Failed to init router: %v\n", err)
	}
	log.Println("listen", "0.0.0.0:"+cfg.Server.Port)
	if err := http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router); err != nil {
		log.Fatalf("TomoChain Rosetta Gateway server exited with error: %v\n", err)
	}
}
