package core

import (
	stafiHubXLedgerTypes "github.com/stafiprotocol/stafihub/x/ledger/types"
)

type Message struct {
	Source      RSymbol
	Destination RSymbol
	Reason      Reason
	Content     interface{}
}

type Reason string

const (

	//need send from other chain
	ReasonLiquidityBondResult = Reason("LiquidityBondResult")
	ReasonBondReport          = Reason("BondReport")
	ReasonActiveReport        = Reason("ActiveReport")
	ReasonWithdrawReport      = Reason("WithdrawReport")
	ReasonTransferReport      = Reason("TransferReport")
	ReasonSubmitSignature     = Reason("SubmitSignature")

	ReasonCurrentChainEra  = Reason("CurrentChainEra")
	ReasonNewEra           = Reason("NewEra")
	ReasonBondedPools      = Reason("BondedPools")
	ReasonNewMultisig      = Reason("AsMulti")
	ReasonMultisigExecuted = Reason("MultisigExecuted")
	ReasonGetEraNominated  = Reason("GetEraNominated")

	//need send when got event from stafi chain
	ReasonLiquidityBondEvent = Reason("LiquidityBondEvent")

	//need send when got event from stafi/stafihub chain
	ReasonEraPoolUpdatedEvent    = Reason("EraPoolUpdatedEvent")
	ReasonBondReportedEvent      = Reason("BondReportedEvent")
	ReasonActiveReportedEvent    = Reason("ActiveReportedEvent")
	ReasonWithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	ReasonTransferReportedEvent  = Reason("TransferReportedEvent")
	ReasonNominationUpdatedEvent = Reason("NominationUpdatedEvent")
	ReasonSignatureEnoughEvent   = Reason("SignatureEnoughed")

	ReasonValidatorUpdatedEvent = Reason("ValidatorUpdatedEvent")
)

// msg data used in cosmos
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
