// Copyright (c) 2020 TomoChain

package main

import (
	"github.com/fatih/color"
	"github.com/tomochain/tomochain-rosetta-gateway/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
