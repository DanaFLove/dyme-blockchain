package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Elect the msg sender as an advisor
func (k msgServer) ElectAdvisor(goCtx context.Context, msg *types.MsgElectAdvisor) (*types.MsgElectAdvisorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ElectedAdvisorsStore))
	store.Set([]byte(msg.Creator), []byte(fmt.Sprintf("{%d}---{%d}", 0, 0)))

	// add advisor to Advisory council directly (TEST ONLY, this should happen after sufficient number of YES votes on elected advisor)

	store2 := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdvisoryCouncilStore))
	store2.Set([]byte(msg.Creator), []byte(strconv.Itoa(int(time.Now().Unix()))))

	return &types.MsgElectAdvisorResponse{}, nil
}
