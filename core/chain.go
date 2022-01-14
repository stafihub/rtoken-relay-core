// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import "github.com/ChainSafe/log15"

type Chain interface {
	SetParams(cfg *ChainConfig, logger log15.Logger, sysErr chan<- error) error
	Start() error // Start chain
	SetRouter(*Router)
	Rsymbol() RSymbol
	Name() string
	Stop()
}

type ChainConfig struct {
	Name            string                 // Human-readable chain name
	Symbol          RSymbol                // symbol
	Endpoint        string                 // url for rpc endpoint
	KeystorePath    string                 // Location of key files
	Insecure        bool                   // Indicated whether the test keyring should be used
	LatestBlockFlag bool                   // If true, overrides blockstore or latest block in config and starts from current block
	Opts            map[string]interface{} // Per chain options
}
