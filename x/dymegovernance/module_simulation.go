package dymegovernance

import (
	"math/rand"

	"dymechain/testutil/sample"
	dymegovernancesimulation "dymechain/x/dymegovernance/simulation"
	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = dymegovernancesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgSetKeyValue = "op_weight_msg_set_key_value"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetKeyValue int = 100

	opWeightMsgStakedyme = "op_weight_msg_stakedyme"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStakedyme int = 100

	opWeightMsgElectAdvisor = "op_weight_msg_elect_advisor"
	// TODO: Determine the simulation weight value
	defaultWeightMsgElectAdvisor int = 100

	opWeightMsgAdviceOnProposal = "op_weight_msg_advice_on_proposal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAdviceOnProposal int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	dymegovernanceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&dymegovernanceGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgSetKeyValue int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetKeyValue, &weightMsgSetKeyValue, nil,
		func(_ *rand.Rand) {
			weightMsgSetKeyValue = defaultWeightMsgSetKeyValue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetKeyValue,
		dymegovernancesimulation.SimulateMsgSetKeyValue(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgStakedyme int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStakedyme, &weightMsgStakedyme, nil,
		func(_ *rand.Rand) {
			weightMsgStakedyme = defaultWeightMsgStakedyme
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStakedyme,
		dymegovernancesimulation.SimulateMsgStakedyme(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgElectAdvisor int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgElectAdvisor, &weightMsgElectAdvisor, nil,
		func(_ *rand.Rand) {
			weightMsgElectAdvisor = defaultWeightMsgElectAdvisor
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgElectAdvisor,
		dymegovernancesimulation.SimulateMsgElectAdvisor(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAdviceOnProposal int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAdviceOnProposal, &weightMsgAdviceOnProposal, nil,
		func(_ *rand.Rand) {
			weightMsgAdviceOnProposal = defaultWeightMsgAdviceOnProposal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAdviceOnProposal,
		dymegovernancesimulation.SimulateMsgAdviceOnProposal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
