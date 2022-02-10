// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"github.com/ChainSafe/log15"
	"github.com/stafihub/rtoken-relay-core/common/config"
)

type Chain interface {
	Initialize(cfg *config.RawChainConfig, logger log15.Logger, sysErr chan<- error) error
	Start() error // Start chain
	SetRouter(*Router)
	RSymbol() RSymbol
	Name() string
	Stop()
}
