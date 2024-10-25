package dymegovernance_test

import (
	"testing"

	keepertest "dymechain/testutil/keeper"
	"dymechain/testutil/nullify"
	"dymechain/x/dymegovernance"
	"dymechain/x/dymegovernance/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DymegovernanceKeeper(t)
	dymegovernance.InitGenesis(ctx, *k, genesisState)
	got := dymegovernance.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}