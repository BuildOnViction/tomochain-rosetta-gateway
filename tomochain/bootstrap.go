package tomochain

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
	"strings"
)

type BootstrapBalanceItem struct {
	AccountIdentifier *RosettaTypes.AccountIdentifier `json:"account_identifier"`
	Currency *RosettaTypes.Currency                   `json:"currency"`
	Value string                                      `json:"value"`
}

// inputFile: genesis file
// outputFile: bootstrap balance
func GenerateBootstrapFile(inputFile, outputFile string) error {
	data, err := ioutil.ReadFile(inputFile)

	genesis := make(map[string]interface{})
	if err == nil {
		err = json.Unmarshal(data, &genesis)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
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
					return err
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

	output, err := json.MarshalIndent(bootstrapBalances, "", "	")
	if err != nil {
		fmt.Println("Unable to marshal bootstrapBalances", err)
		return err
	}
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	err = ioutil.WriteFile(outputFile, output, 0644)

	if err != nil {
		fmt.Println("Unable to write output file", outputFile, err)
	}
	return nil
}

