package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stafiHubXLedgerTypes "github.com/stafihub/stafihub/x/ledger/types"
	stafiHubXRValidatorTypes "github.com/stafihub/stafihub/x/rvalidator/types"
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
	ReasonNewEra                       = Reason("NewEra")
	ReasonExeLiquidityBond             = Reason("ExeLiquidityBond")
	ReasonExeNativeAndLsmLiquidityBond = Reason("ExeLiquidityNativeAndLsmBond")
	ReasonBondReport                   = Reason("BondReport")
	ReasonActiveReport                 = Reason("ActiveReport")
	ReasonTransferReport               = Reason("TransferReport")
	ReasonSubmitSignature              = Reason("SubmitSignature")
	ReasonRValidatorUpdateReport       = Reason("RValidatorUpdateReport")
	ReasonInterchainTx                 = Reason("InterchainTx")

	//send when got event from stafihub chain
	ReasonEraPoolUpdatedEvent    = Reason("EraPoolUpdatedEvent")
	ReasonBondReportedEvent      = Reason("BondReportedEvent")
	ReasonActiveReportedEvent    = Reason("ActiveReportedEvent")
	ReasonTransferReportedEvent  = Reason("TransferReportedEvent")
	ReasonSignatureEnoughEvent   = Reason("SignatureEnoughed")
	ReasonRValidatorUpdatedEvent = Reason("RValidatorUpdatedEvent")
	ReasonRValidatorAddedEvent   = Reason("RValidatorAddedEvent")
	ReasonRParamsChangedEvent    = Reason("RParamsChangedEvent")
	ReasonInitPoolEvent          = Reason("InitPoolEvent")
	ReasonRemovePoolEvent        = Reason("RemovePoolEvent")

	//get reason
	ReasonGetPools              = Reason("GetPools")
	ReasonGetSignatures         = Reason("GetSignatures")
	ReasonGetBondRecord         = Reason("GetBondRecord")
	ReasonGetInterchainTxStatus = Reason("GetInterchainTxStatus")
)

// === stafihub -> other chain msg data used in cosmos
type EventEraPoolUpdated struct {
	Denom      string
	LastEra    uint32
	CurrentEra uint32
	ShotId     string
	Snapshot   stafiHubXLedgerTypes.BondSnapshot
}

type EventBondReported struct {
	Denom    string
	ShotId   string
	Snapshot stafiHubXLedgerTypes.BondSnapshot
}

type EventActiveReported struct {
	Denom      string
	ShotId     string
	Snapshot   stafiHubXLedgerTypes.BondSnapshot
	PoolUnbond []*stafiHubXLedgerTypes.Unbonding
}
type EventTransferReported struct {
	Denom  string
	ShotId string
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

type EventRValidatorUpdated struct {
	Denom          string
	PoolAddress    string
	Era            uint32
	OldAddress     string
	NewAddress     string
	CycleVersion   uint64
	CycleNumber    uint64
	CycleSeconds   uint64
	BlockTimestamp int64
}
type EventRValidatorAdded struct {
	Denom        string
	PoolAddress  string
	Era          uint32
	AddedAddress string
}

type EventRParamsChanged struct {
	Denom      string
	GasPrice   string `json:"gasPrice"`
	EraSeconds uint32 `json:"eraSeconds"`
	LeastBond  string `json:"leastBond"`
	Offset     int32  `json:"offset"`
}

type EventInitPool struct {
	Denom             string
	PoolAddress       string
	WithdrawalAddress string
	HostChannelId     string
	Validators        []string
}

type EventRemovePool struct {
	Denom       string
	PoolAddress string
}

// === other chain -> stafihub msg data used in cosmos
type ProposalExeLiquidityBond struct {
	Denom  string
	Bonder string
	Pool   string
	Txhash string
	Amount sdk.Int
	State  stafiHubXLedgerTypes.LiquidityBondState
}

type ProposalExeNativeAndLsmLiquidityBond struct {
	Denom            string
	Bonder           string
	Pool             string
	Txhash           string
	NativeBondAmount sdk.Int
	LsmBondAmount    sdk.Int
	State            stafiHubXLedgerTypes.LiquidityBondState
	Msgs             []sdk.Msg
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

type ProposalRValidatorUpdateReport struct {
	Denom        string
	PoolAddress  string
	CycleVersion uint64
	CycleNumber  uint64
	Status       stafiHubXRValidatorTypes.UpdateRValidatorStatus
}

type ProposalInterchainTx struct {
	Denom  string
	Pool   string
	Era    uint32
	TxType stafiHubXLedgerTypes.OriginalTxType
	Factor uint32
	Msgs   []sdk.Msg
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

type ParamGetSignatures struct {
	Denom  string
	Era    uint32
	Pool   string
	TxType stafiHubXLedgerTypes.OriginalTxType
	PropId string
	Sigs   chan []string
}

type ParamGetBondRecord struct {
	Denom      string
	TxHash     string
	BondRecord chan stafiHubXLedgerTypes.BondRecord
}

type ParamGetInterchainTxStatus struct {
	PropId string
	Status chan stafiHubXLedgerTypes.InterchainTxStatus
}
