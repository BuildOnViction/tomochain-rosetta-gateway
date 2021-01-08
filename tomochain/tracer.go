// Copyright (c) 2020 TomoChain

package tomochain

import (
	"fmt"
	"github.com/tomochain/tomochain/eth"
	"io/ioutil"
)

// convert raw eth data from client to rosetta

const (
	tracerPath = "tomochain/call_tracer.js"
)

var (
	tracerTimeout = "120s"
)

func loadTraceConfig() (*eth.TraceConfig, error) {
	loadedFile, err := ioutil.ReadFile(tracerPath)
	if err != nil {
		return nil, fmt.Errorf("%w: could not load tracer file", err)
	}

	loadedTracer := string(loadedFile)
	return &eth.TraceConfig{
		Timeout: &tracerTimeout,
		Tracer:  &loadedTracer,
	}, nil
}
