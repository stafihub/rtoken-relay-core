// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigPath   = "./config.json"
	defaultKeystorePath = "./keys"
	defaultLogFilePath  = "./log_file"
)

var (
	ChainTypeStafiHub  = "stafiHub"
	ChainTypeSubstrate = "substrate"
	ChainTypeCosmosHub = "cosmosHub"
	ChainTypeSolana    = "solana"
	ChainTypeEthereum  = "ethereum"
	ChainTypeBinance   = "binance"

	ChainTypeSupport = map[string]bool{
		ChainTypeSubstrate: true,
		ChainTypeCosmosHub: true,
		ChainTypeSolana:    true,
		ChainTypeEthereum:  true,
		ChainTypeBinance:   true,
	}
)

type Config struct {
	BlockstorePath string         `json:"blockstorePath"`
	LogFilePath    string         `json:"logFilePath"`
	NativeChain    RawChainConfig `json:"nativeChain"`
	ExternalChain  RawChainConfig `json:"externalChain"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
	Name         string      `json:"name"`
	Rsymbol      string      `json:"rsymbol"`
	EndpointList []string    `json:"endpointList"` // url for rpc endpoint
	KeystorePath string      `json:"keystorePath"`
	Opts         interface{} `json:"opts"`
}

func GetConfig(filePath string) (*Config, error) {
	var cfg Config
	path := defaultConfigPath
	if filePath != "" {
		path = filePath
	}
	err := loadConfig(path, &cfg)
	if err != nil {
		fmt.Println("err loading json file", "err", err.Error())
		return nil, err
	}
	if len(cfg.LogFilePath) == 0 {
		cfg.LogFilePath = defaultLogFilePath
	}
	fmt.Println("Loaded config", "path", path)
	return &cfg, nil
}

func loadConfig(file string, config *Config) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if ext != ".json" {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return json.NewDecoder(f).Decode(&config)
}
