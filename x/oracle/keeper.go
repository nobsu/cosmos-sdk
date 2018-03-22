package oracle

import (
	"github.com/cosmos/cosmos-sdk/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	key        sdk.StoreKey
	cdc        *wire.Codec
	dispatcher Dispatcher
}

func (keeper Keeper) Dispatcher() Dispatcher {
	return keeper.dispatcher
}
