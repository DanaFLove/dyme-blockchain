package keeper

import (
	"context"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetKeyValue(goCtx context.Context, msg *types.MsgSetKeyValue) (*types.MsgSetKeyValueResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("DockerStore"))

	store.Set([]byte(msg.Key), []byte(msg.Value))

	return &types.MsgSetKeyValueResponse{}, nil
}
