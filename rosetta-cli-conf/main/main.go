package main

import (
	"encoding/json"
	"fmt"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cast"
	"github.com/tomochain/tomochain-rosetta-gateway/common"
	tomochaincommon "github.com/tomochain/tomochain/common"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultInputFile = "genesis.json"
	DefaultOutputFile = "bootstrap_balances.json"
)
type BootstrapBalanceItem struct {
	AccountIdentifier *RosettaTypes.AccountIdentifier `json:"account_identifier"`
	Currency *RosettaTypes.Currency                   `json:"currency"`
	Value string                                      `json:"value"`
}

// input: genesis file
// output: bootstrap balance
func main() {

	path, _ := os.Getwd()

	genesisFile := DefaultInputFile
	if len(os.Args) > 1 {
		genesisFile = os.Args[1]
	}
	data, err := ioutil.ReadFile(filepath.Join(path, genesisFile))

	genesis := make(map[string]interface{})
	if err == nil {
		err = json.Unmarshal(data, &genesis)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	bootstrapBalances := []*BootstrapBalanceItem{}

	if allocBalances, ok := genesis["alloc"]; ok {
		wallets := allocBalances.(map[string]interface{})
		for addr, data := range wallets {
			walletData := data.(map[string]interface{})
			if hexBalance, ok := walletData["balance"] ; ok {
				balance, good := new(big.Int).SetString(strings.TrimPrefix(cast.ToString(hexBalance), "0x"), 16)
				if !good {
					fmt.Println("Cannot parse balance of address ", addr, err)
					return
				}
				bootstrapBalances = append(bootstrapBalances, &BootstrapBalanceItem{
					AccountIdentifier: &RosettaTypes.AccountIdentifier{
						Address:    tomochaincommon.HexToAddress(addr).Hex(),
					},
					Currency:         common.TomoNativeCoin ,
					Value:             balance.String(),
				})
			}
		}
	}

	output, err := json.Marshal(bootstrapBalances)
	if err != nil {
		fmt.Println("Unable to marshal bootstrapBalances", err)
		return
	}
	outputFile := DefaultOutputFile
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	err = ioutil.WriteFile(filepath.Join(path, outputFile), output, 0644)

	if err != nil {
		fmt.Println("Unable to write output file", outputFile, err)
	}
}