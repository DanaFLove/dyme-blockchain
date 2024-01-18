package keeper

import (
	"context"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAdvisors(goCtx context.Context, req *types.QueryGetAdvisorsRequest) (*types.QueryGetAdvisorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var response []byte
	var chosenStore string
	if req.Advisorparam2 == "COUNCIL" {
		chosenStore = types.AdvisoryCouncilStore
	} else {
		chosenStore = types.ElectedAdvisorsStore
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(chosenStore))
	response = store.Get([]byte(req.Advisorparam1))
	return &types.QueryGetAdvisorsResponse{Advisordata: string(response)}, nil

}
