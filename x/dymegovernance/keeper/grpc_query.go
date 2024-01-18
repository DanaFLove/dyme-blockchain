package keeper

import (
	"dymechain/x/dymegovernance/types"
)

var _ types.QueryServer = Keeper{}
