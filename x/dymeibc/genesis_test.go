package dymeibc_test

import (
	"testing"

	keepertest "dymechain/testutil/keeper"
	"dymechain/testutil/nullify"
	"dymechain/x/dymeibc"
	"dymechain/x/dymeibc/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DymeibcKeeper(t)
	dymeibc.InitGenesis(ctx, *k, genesisState)
	got := dymeibc.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
