// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stafihub/rtoken-relay-core/common/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Core struct {
	Registry []Chain
	route    *Router
	log      log.Logger
	sysErr   <-chan error
}

func NewCore(logger log.Logger, sysErr <-chan error) *Core {
	return &Core{
		Registry: make([]Chain, 0),
		route:    NewRouter(logger),
		log:      logger,
		sysErr:   sysErr,
	}
}

// AddChain registers the chain in the Registry and calls Chain.SetRouter()
func (c *Core) AddChain(chain Chain) {
	c.Registry = append(c.Registry, chain)
	chain.SetRouter(c.route)
}

// Start will call all registered chains' Start methods and block forever (or until signal is received)
func (c *Core) Start() {
	for _, chain := range c.Registry {
		err := chain.Start()
		if err != nil {
			c.log.Error("failed to start chain", "chain", chain.RSymbol(), "err", err)
			return
		}
		c.log.Info(fmt.Sprintf("Started %s chain", chain.Name()))
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-c.sysErr:
		c.log.Error("FATAL ERROR. Shutting down.", "err", err)
	case sig := <-sigc:
		c.log.Warn("Interrupt received, shutting down now. signal:", sig.String())
	}

	c.route.StopMsgHandler()
	// Signal chains to shutdown
	for _, chain := range c.Registry {
		chain.Stop()
	}
}

func (c *Core) Errors() <-chan error {
	return c.sysErr
}

var sdkContextMutex sync.Mutex

// UseSDKContext uses a custom Bech32 account prefix and returns a restore func
// CONTRACT: When using this function, caller must ensure that lock contention
// doesn't cause program to hang. This function is only for use in codec calls
func UseSdkConfigContext(accountPrefix string) func() {
	// Ensure we're the only one using the global context,
	// lock context to begin function
	sdkContextMutex.Lock()

	// Mutate the sdkConf

	if accountPrefix == "iaa" {
		setIrisPrefix()
	} else {
		setCommonPrefixes(accountPrefix)
	}

	// Return the unlock function, caller must lock and ensure that lock is released
	// before any other function needs to use c.UseSDKContext
	return sdkContextMutex.Unlock
}

func setCommonPrefixes(accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	// Set config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

func setIrisPrefix() {

	// Bech32ChainPrefix defines the prefix of this chain
	Bech32ChainPrefix := "i"

	// PrefixAcc is the prefix for account
	PrefixAcc := "a"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator := "v"

	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus := "c"

	// PrefixPublic is the prefix for public
	PrefixPublic := "p"

	// PrefixAddress is the prefix for address
	PrefixAddress := "a"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr := Bech32ChainPrefix + PrefixAcc + PrefixAddress
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub := Bech32ChainPrefix + PrefixAcc + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr := Bech32ChainPrefix + PrefixValidator + PrefixAddress
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub := Bech32ChainPrefix + PrefixValidator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr := Bech32ChainPrefix + PrefixConsensus + PrefixAddress
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub := Bech32ChainPrefix + PrefixConsensus + PrefixPublic

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}
