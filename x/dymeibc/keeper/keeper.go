package keeper

import (
	"dymechain/x/dymeibc/types"
)

var _ types.QueryServer = Keeper{}
