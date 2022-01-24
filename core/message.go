package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stafiHubXLedgerTypes "github.com/stafiprotocol/stafihub/x/ledger/types"
)

type OriginalTx string

const (
	OriginalTxDefault       = OriginalTx("default")
	OriginalBond            = OriginalTx("Bond") //bond or unbond
	OriginalUnbond          = OriginalTx("Unbond")
	OriginalWithdrawUnbond  = OriginalTx("WithdrawUnbond")
	OriginalWithdrawReward  = OriginalTx("WithdrawReward")
	OriginalTransfer        = OriginalTx("Transfer")        //transfer
	OriginalClaimRewards    = OriginalTx("ClaimRewards")    // claim
	OriginalUpdateValidator = OriginalTx("UpdateValidator") // redelegate
)

type Message struct {
	Source      RSymbol
	Destination RSymbol
	Reason      Reason
	Content     interface{}
}

type Reason string

const (
	//send from other chain
	ReasonNewEra           = Reason("NewEra")
	ReasonExeLiquidityBond = Reason("ExeLiquidityBond")
	ReasonBondReport       = Reason("BondReport")
	ReasonActiveReport     = Reason("ActiveReport")
	ReasonWithdrawReport   = Reason("WithdrawReport")
	ReasonTransferReport   = Reason("TransferReport")
	ReasonSubmitSignature  = Reason("SubmitSignature")

	ReasonCurrentChainEra  = Reason("CurrentChainEra")
	ReasonBondedPools      = Reason("BondedPools")
	ReasonNewMultisig      = Reason("AsMulti")
	ReasonMultisigExecuted = Reason("MultisigExecuted")
	ReasonGetEraNominated  = Reason("GetEraNominated")

	//send when got event from stafi chain
	ReasonLiquidityBondEvent = Reason("LiquidityBondEvent")

	//send when got event from stafi/stafihub chain
	ReasonEraPoolUpdatedEvent    = Reason("EraPoolUpdatedEvent")
	ReasonBondReportedEvent      = Reason("BondReportedEvent")
	ReasonActiveReportedEvent    = Reason("ActiveReportedEvent")
	ReasonWithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	ReasonTransferReportedEvent  = Reason("TransferReportedEvent")
	ReasonNominationUpdatedEvent = Reason("NominationUpdatedEvent")
	ReasonSignatureEnoughEvent   = Reason("SignatureEnoughed")

	ReasonValidatorUpdatedEvent = Reason("ValidatorUpdatedEvent")
)

// === stafihub -> other chain msg data used in cosmos
type EventEraPoolUpdated struct {
	Denom       string
	LastEra     string
	CurrentEra  string
	ShotId      string
	LasterVoter string
	Snapshot    stafiHubXLedgerTypes.BondSnapshot
}

type EventBondReported struct {
	Denom       string
	ShotId      string
	LasterVoter string
	Snapshot    stafiHubXLedgerTypes.BondSnapshot
}

type EventActiveReported struct {
	Denom       string
	ShotId      string
	LasterVoter string
	Snapshot    stafiHubXLedgerTypes.BondSnapshot
}

type EventWithdrawReported struct {
	Denom       string
	ShotId      string
	LasterVoter string
	Snapshot    stafiHubXLedgerTypes.BondSnapshot
	PoolUnbond  stafiHubXLedgerTypes.PoolUnbond
}

type EventTransferReported struct {
	Denom       string
	ShotId      string
	LasterVoter string
}

// === other chain -> stafihub msg data used in cosmos
type ProposalExeLiquidityBond struct {
	Denom     string
	Bonder    string
	Pool      string
	Blockhash string
	Txhash    string
	Amount    sdk.Int
}

type ProposalSetChainEra struct {
	Denom string
	Era   uint32
}

type ProposalBondReport struct {
	Denom  string
	ShotId []byte
	Action stafiHubXLedgerTypes.BondAction
}

type ProposalActiveReport struct {
	Denom    string
	ShotId   []byte
	Staked   sdk.Int
	Unstaked sdk.Int
}

type ProposalWithdrawReport struct {
	Denom  string
	ShotId []byte
}

type ProposalTransferReport struct {
	Denom  string
	ShotId []byte
}
