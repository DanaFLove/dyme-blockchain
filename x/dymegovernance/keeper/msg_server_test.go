package keeper_test

import (
	"context"
	"testing"

	keepertest "dymechain/testutil/keeper"
	"dymechain/x/dymegovernance/keeper"
	"dymechain/x/dymegovernance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.DymegovernanceKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
