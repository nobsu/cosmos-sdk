package oracle

import (
	"github.com/cosmos/cosmos-sdk/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

type Keeper struct {
	key sdk.StoreKey
	cdc *wire.Codec

	sk stake.Keeper
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, sk stake.Keeper) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		sk:  sk,
	}
}

type OracleInfo struct {
	// ValidatorsHash []byte
	Signers   []sdk.Address
	Power     sdk.Rat
	Processed bool
}

func (keeper Keeper) OracleInfo(ctx sdk.Context, oracle Oracle) (res OracleInfo) {
	store := ctx.KVStore(keeper.key)

	key, err := keeper.cdc.MarshalBinary(oracle)
	if err != nil {
		panic(err)
	}

	bz := store.Get(key)

	if bz == nil {
		return OracleInfo{
			Signers:   []sdk.Address{},
			Power:     sdk.ZeroRat,
			Processed: false,
		}
	}

	if err = keeper.cdc.UnmarshalBinary(bz, &res); err != nil {
		panic(err)
	}

	return
}

func (keeper Keeper) setInfo(ctx sdk.Context, oracle Oracle, info OracleInfo) {
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
