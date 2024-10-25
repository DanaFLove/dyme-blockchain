package types

const (
	// ModuleName defines the module name
	ModuleName = "dymegovernance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dymegovernance"

	// Initial Stakes Amount
	InitialStake = "1000000udyme"

	TimeFormat        = "2006-01-02 15:04:05"
	ProposalChannelId = "1078657716953821285"

	KeyRingBackend = "test"

	// Data store names
	ElectedAdvisorsStore = "ElectedAdvisors"
	AdvisoryCouncilStore = "AdvisoryCouncil"
	ProposalAdviceStore  = "ProposalAdvice"

	AdviceStateNone     = "NO_ADVICE"
	AdviceStateReturned = "RETURNED"
	AdviceStatePassed   = "ADVICE"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
