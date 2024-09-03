package keeper

import (
	"context"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetKeyValue(goCtx context.Context, req *types.QueryGetKeyValueRequest) (*types.QueryGetKeyValueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("DockerStore"))

	value := store.Get([]byte(req.Key))

	return &types.QueryGetKeyValueResponse{Value: string(value)}, nil
}
