package keeper

import (
	"context"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AdviceOnProposal(goCtx context.Context, msg *types.MsgAdviceOnProposal) (*types.MsgAdviceOnProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.AdvisoryOutcome != types.AdviceStateNone && msg.AdvisoryOutcome != types.AdviceStatePassed && msg.AdvisoryOutcome != types.AdviceStateReturned {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid value for Advisory Outcome")
	}

	// TODO: check if proposal exists and is in a valid state

	// check that only a member of the advisory board can do this transaction
	store1 := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdvisoryCouncilStore))
	checkAdvisor := store1.Get([]byte(msg.Creator))

	if string(checkAdvisor) == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "Only advisors are allowed")
	}
	// Question: does an advice done by only a single member or all members of the advisory council?
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ElectedAdvisorsStore))
	// Proposal advisory state store: [Proposal ID] -> [Proposal advisory state]
	// Advisory states store [advisor wallet] -> ['Advisory Vote']
	store.Set([]byte(msg.ProposalId), []byte(msg.AdvisoryOutcome))
	return &types.MsgAdviceOnProposalResponse{}, nil
}
