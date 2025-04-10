package keeper_test

import (
	"testing"

	testkeeper "dymechain/testutil/keeper"
	"dymechain/x/dymegovernance/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.DymegovernanceKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
