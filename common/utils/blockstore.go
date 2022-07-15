package utils

import (
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"strconv"
)

func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}

func parseStartBlock(startBlock string) uint64 {
	//if blk, ok := cfg.Opts[""]; ok {
	blk, err := strconv.ParseUint(startBlock, 10, 32)
	if err != nil {
		panic(err)
	}
	return blk
}
