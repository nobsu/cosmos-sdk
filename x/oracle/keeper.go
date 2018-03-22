package oracle

import (
	"github.com/cosmos/cosmos-sdk/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

type Keeper struct {
	key        sdk.StoreKey
	cdc        *wire.Codec
	dispatcher types.Dispatcher
}

func (keeper Keeper) Dispatcher() types.Dispatcher {
	return keeper.dispatcher
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		dispatcher: types.NewDispatcher(),
	}
}

type OracleInfo struct {
	// ValidatorsHash []byte
	Signers    []sdk.Address
	TotalPower uint64
	Processed  bool
}

func (keeper Keeper) OracleInfo(ctx sdk.Context, oracle types.Oracle) (res OracleInfo) {
	store := ctx.KVStore(keeper.key)

	bz, err := keeper.cdc.MarshalBinary(oracle)
	if err != nil {
		panic(err)
	}

	if bz == nil {
		return OracleInfo{
			Signers:    []sdk.Address{},
			TotalPower: 0,
			Processed:  false,
		}
	}

	if err = keeper.cdc.UnmarshalBinary(bz, &res); err != nil {
		panic(err)
	}

	return
}

func (keeper Keeper) setInfo(ctx sdk.Context, oracle types.Oracle, info OracleInfo) {
	store := ctx.KVStore(keeper.key)

	k, err := keeper.cdc.MarshalBinary(oracle)
	if err != nil {
		panic(err)
	}

	v, err := keeper.cdc.MarshalBinary(info)
	if err != nil {
		panic(err)
	}

	store.Set(k, v)
}
