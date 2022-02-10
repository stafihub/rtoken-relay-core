package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stafiHubXLedgerTypes "github.com/stafihub/stafihub/x/ledger/types"
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

	//get reason
	ReasonGetPools = Reason("GetPools")
)

// === stafihub -> other chain msg data used in cosmos
type EventEraPoolUpdated struct {
	Denom       string
	LastEra     uint32
	CurrentEra  uint32
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
	PoolUnbond  stafiHubXLedgerTypes.PoolUnbond
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

type EventSignatureEnough struct {
	Denom      string
	Era        uint32
	Pool       string
	TxType     stafiHubXLedgerTypes.OriginalTxType
	ProposalId string
	Signatures [][]byte
	Threshold  uint32
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
	ShotId string
	Action stafiHubXLedgerTypes.BondAction
}

type ProposalActiveReport struct {
	Denom    string
	ShotId   string
	Staked   sdk.Int
	Unstaked sdk.Int
}

type ProposalWithdrawReport struct {
	Denom  string
	ShotId string
}

type ProposalTransferReport struct {
	Denom  string
	ShotId string
}

type ParamSubmitSignature struct {
	Denom     string
	Era       uint32
	Pool      string
	TxType    stafiHubXLedgerTypes.OriginalTxType
	PropId    string
	Signature string
}

// get msg
type ParamGetPools struct {
	Denom string
	Pools chan []string
}
