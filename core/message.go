package core

type Message struct {
	Source      RSymbol
	Destination RSymbol
	Reason      Reason
	Content     interface{}
}

type Reason string

const (
	ReasonLiquidityBondResult = Reason("LiquidityBondResult")

	ReasonCurrentChainEra  = Reason("CurrentChainEra")
	ReasonNewEra           = Reason("NewEra")
	ReasonBondedPools      = Reason("BondedPools")
	ReasonInformChain      = Reason("InformChain")
	ReasonActiveReport     = Reason("ActiveReport")
	ReasonNewMultisig      = Reason("AsMulti")
	ReasonMultisigExecuted = Reason("MultisigExecuted")
	ReasonGetEraNominated  = Reason("GetEraNominated")
	ReasonSubmitSignature  = Reason("SubmitSignature")

	//need send when got event from stafi/stafihub chain
	ReasonLiquidityBondEvent = Reason("LiquidityBondEvent")

	ReasonEraPoolUpdatedEvent    = Reason("EraPoolUpdatedEvent")
	ReasonBondReportedEvent      = Reason("BondReportedEvent")
	ReasonActiveReportedEvent    = Reason("ActiveReportedEvent")
	ReasonWithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	ReasonTransferReportedEvent  = Reason("TransferReportedEvent")
	ReasonNominationUpdatedEvent = Reason("NominationUpdatedEvent")
	ReasonSignatureEnoughEvent   = Reason("SignatureEnoughed")

	ReasonValidatorUpdatedEvent = Reason("ValidatorUpdatedEvent")
)
