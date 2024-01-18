package simulation

import (
	"math/rand"

	"dymechain/x/dymegovernance/keeper"
	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgStakedyme(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgStakedyme{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the Stakedyme simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Stakedyme simulation not implemented"), nil, nil
	}
}
