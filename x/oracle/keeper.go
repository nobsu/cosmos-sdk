package oracle

import (
	"github.com/cosmos/cosmos-sdk/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	key sdk.StoreKey
	cdc *wire.Codec

	valset sdk.ValidatorSet
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec, valset sdk.ValidatorSet) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,

		valset: valset,
	}
}

type OracleInfo struct {
	// ValidatorsHash []byte
	Signers   []sdk.Address
	Power     sdk.Rat
	Processed bool
}

func EmptyOracleInfo() OracleInfo {
	return OracleInfo{
		Signers:   []sdk.Address{},
		Power:     sdk.ZeroRat(),
		Processed: false,
	}
}

func (keeper Keeper) OracleInfo(ctx sdk.Context, p Payload) (res OracleInfo) {
	store := ctx.KVStore(keeper.key)

	key := GetInfoKey(p, keeper.cdc)

	bz := store.Get(key)

	if bz == nil {
		return EmptyOracleInfo()
	}

	keeper.cdc.MustUnmarshalBinary(bz, &res)

	return
}

func (keeper Keeper) setInfo(ctx sdk.Context, p Payload, info OracleInfo) {
	store := ctx.KVStore(keeper.key)

	key := GetInfoKey(p, keeper.cdc)

	bz := keeper.cdc.MustMarshalBinary(info)

	store.Set(key, bz)
}
