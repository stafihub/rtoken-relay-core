package main

import (
	"errors"
	"os"
	"strconv"

	log "github.com/ChainSafe/log15"
	cosmosChain "github.com/stafiprotocol/cosmos-relay-sdk/chain"
	"github.com/stafiprotocol/rtoken-relay-core/config"
	"github.com/stafiprotocol/rtoken-relay-core/core"
	stafiHubChain "github.com/stafiprotocol/stafi-hub-relay-sdk/chain"
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

	for _, chainConfig := range cfg.Chains {
		var newChain core.Chain
		logger := log.Root().New("chain", chainConfig.Name)

		switch chainConfig.Type {
		case config.ChainTypeStafiHub:
			newChain = stafiHubChain.NewChain()
			newChain.Initialize(&chainConfig, logger, sysErr)
		case config.ChainTypeCosmosHub:
			newChain = cosmosChain.NewChain()
			newChain.Initialize(&chainConfig, logger, sysErr)
		default:
			return errors.New("unsupport Chain Type")
		}

		c.AddChain(newChain)
	}

	c.Start()

	return nil
}
