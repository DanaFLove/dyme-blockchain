package keeper

import (
	"context"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Stakedyme(goCtx context.Context, msg *types.MsgStakedyme) (*types.MsgStakedymeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("StakedRecord"))

	// check if already in store
	alreadyStaked := store.Get([]byte(msg.Creator))
	if string(alreadyStaked) != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrConflict, "1 DYME Already staked by this wallet")
	}
	store.Set([]byte(msg.Creator), []byte("1"))

	return &types.MsgStakedymeResponse{}, nil
}
