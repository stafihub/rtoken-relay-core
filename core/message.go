package core

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	stafiHubXLedgerTypes "github.com/stafiprotocol/stafihub/x/ledger/types"
)

type OriginalTx string

const (
	OriginalTxDefault       = OriginalTx("default")
	OriginalBond            = OriginalTx("Bond") //bond or unbond
	OriginalUnbond          = OriginalTx("Unbond")
	OriginalWithdrawUnbond  = OriginalTx("WithdrawUnbond")  // used in substrate, because cosmos will auto return unbond token to staker
	OriginalClaimRewards    = OriginalTx("ClaimRewards")    // claim and delegate reward token
	OriginalTransfer        = OriginalTx("Transfer")        //transfer
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

type SubmitSignatures struct {
	Denom      string
	Era        uint32
	Pool       string
	TxType     OriginalTx
	ProposalId []byte
	Signature  [][]byte
	Threshold  uint32
}

func (ss *SubmitSignatures) EncodeToHash() (common.Hash, error) {

	eraBts := make([]byte, 4)
	binary.BigEndian.PutUint32(eraBts, ss.Era)

	packed := make([]byte, 0)
	packed = append(packed, []byte(ss.Denom)...)
	packed = append(packed, eraBts...)
	packed = append(packed, ss.Pool...)
	packed = append(packed, []byte(ss.TxType)...)

	return crypto.Keccak256Hash(packed), nil
}

type EvtSignatureEnough struct {
	Denom      string
	Era        uint32
	Pool       string
	TxType     OriginalTx
	ProposalId []byte
	Signatures [][]byte
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
