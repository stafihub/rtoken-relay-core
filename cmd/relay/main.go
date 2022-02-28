package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	log "github.com/ChainSafe/log15"
	cosmosChain "github.com/stafihub/cosmos-relay-sdk/chain"
	"github.com/stafihub/rtoken-relay-core/common/config"
	"github.com/stafihub/rtoken-relay-core/common/core"
	stafiHubChain "github.com/stafihub/stafi-hub-relay-sdk/chain"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var mainFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
}

var generateFlags = []cli.Flag{
	config.KeystorePathFlag,
	config.NetworkFlag,
}

var bncGenerateFlags = []cli.Flag{
	config.KeystorePathFlag,
	config.BncNetwork,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage reth keystore",
	Description: "The accounts command is used to manage the relay keystore.\n" +
		"\tMake sure the keystore dir is exist before generating\n" +
		"\tTo generate a substrate keystore: relay accounts gensub\n" +
		"\tTo generate a ethereum keystore: relay accounts geneth\n" +
		"\tTo generate a bc chain keystore: relay accounts genbc\n" +
		"\tTo list keys: chainbridge accounts list",
	Subcommands: []*cli.Command{
		{
			Action:      handleGenerateSubCmd,
			Name:        "gensub",
			Usage:       "generate subsrate keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the substrate keystore.",
		},
		{
			Action:      handleGenerateEthCmd,
			Name:        "geneth",
			Usage:       "generate ethereum keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.",
		},
		{
			Action:      handleGenerateBcCmd,
			Name:        "genbc",
			Usage:       "generate binance chain keystore",
			Flags:       bncGenerateFlags,
			Description: "The generate subcommand is used to generate the binance chain keystore.",
		},
	},
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2022 Stafi Protocol Authors"
	app.Name = "reley"
	app.Usage = "relay"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2022"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
	}

	app.Flags = append(app.Flags, mainFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	logger := log.Root()
	var lvl log.Lvl
	if lvlToInt, err := strconv.Atoi(ctx.String(config.VerbosityFlag.Name)); err == nil {
		lvl = log.Lvl(lvlToInt)
	} else if lvl, err = log.LvlFromString(ctx.String(config.VerbosityFlag.Name)); err != nil {
		return err
	}

	logger.SetHandler(log.MultiHandler(
		log.LvlFilterHandler(
			lvl,
			log.StreamHandler(os.Stdout, log.LogfmtFormat())),
		log.Must.FileHandler("relay_log.json", log.JsonFormat()),
		log.LvlFilterHandler(
			log.LvlError,
			log.Must.FileHandler("relay_log_errors.json", log.JsonFormat()))))

	return nil
}

func run(ctx *cli.Context) error {
	err := startLogger(ctx)
	if err != nil {
		return err
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Used to signal core shutdown due to fatal error
	sysErr := make(chan error)
	c := core.NewCore(sysErr)

	// init stafiHub
	stafiHubChainConfig := cfg.NativeChain
	bts, err := json.Marshal(stafiHubChainConfig.Opts)
	if err != nil {
		return err
	}
	option := stafiHubChain.ConfigOption{}
	err = json.Unmarshal(bts, &option)
	if err != nil {
		return err
	}
	option.CaredSymbol = cfg.ExternalChain.Rsymbol
	option.BlockstorePath = cfg.BlockstorePath

	stafiHubChainConfig.Opts = option
	stafiHubChain := stafiHubChain.NewChain()
	logger := log.Root().New("chain", stafiHubChainConfig.Name)
	err = stafiHubChain.Initialize(&stafiHubChainConfig, logger, sysErr)
	if err != nil {
		return err
	}
	c.AddChain(stafiHubChain)

	// init externa chain
	chainConfig := cfg.ExternalChain
	var newChain core.Chain
	// prepare r params
	rParams, err := stafiHubChain.GetRParams(chainConfig.Rsymbol)
	if err != nil {
		return err
	}
	bts, err = json.Marshal(chainConfig.Opts)
	if err != nil {
		return err
	}
	cosmosOption := cosmosChain.ConfigOption{}
	err = json.Unmarshal(bts, &cosmosOption)
	if err != nil {
		return err
	}
	eraSeconds, err := strconv.Atoi(rParams.RParams.EraSeconds)
	if err != nil {
		return err
	}
	if len(cosmosOption.PoolNameSubKey) == 0 {
		return fmt.Errorf("no pool and subkey")
	}

	// prepare pools and threshold
	poolRes, err := stafiHubChain.GetPools(rParams.RParams.Denom)
	if err != nil {
		return err
	}
	cosmosOption.PoolAddressThreshold = make(map[string]uint32)
	for _, poolAddress := range poolRes.GetAddrs() {
		poolDetail, err := stafiHubChain.GetPoolDetail(rParams.RParams.Denom, poolAddress)
		if err != nil {
			return err
		}
		if poolDetail.Detail.Threshold <= 0 {
			return fmt.Errorf("pool threshold is zero in stafihub, pool: %s", poolAddress)
		}
		cosmosOption.PoolAddressThreshold[poolAddress] = poolDetail.Detail.Threshold
	}

	cosmosOption.EraSeconds = eraSeconds
	cosmosOption.GasPrice = rParams.RParams.GasPrice
	cosmosOption.TargetValidators = rParams.RParams.Validators
	cosmosOption.LeastBond = rParams.RParams.LeastBond.BigInt()
	cosmosOption.BlockstorePath = cfg.BlockstorePath

	chainConfig.Opts = cosmosOption
	newChain = cosmosChain.NewChain()
	externalChainLogger := log.Root().New("chain", chainConfig.Name)
	err = newChain.Initialize(&chainConfig, externalChainLogger, sysErr)
	if err != nil {
		return fmt.Errorf("newChain.Initialize failed: %s", err)
	}
	c.AddChain(newChain)
	c.Start()
	return nil
}
