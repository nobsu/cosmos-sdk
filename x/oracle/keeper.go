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
	Power     sdk.Rat
	Hash      []byte
	Processed bool
}

func EmptyOracleInfo(ctx sdk.Context) OracleInfo {
	return OracleInfo{
		Power:     sdk.ZeroRat(),
		Hash:      ctx.BlockHeader().ValidatorsHash,
		Processed: false,
	}
}

func (keeper Keeper) OracleInfo(ctx sdk.Context, p Payload) (res OracleInfo) {
	store := ctx.KVStore(keeper.key)

	key := GetInfoKey(p, keeper.cdc)

	bz := store.Get(key)

	if bz == nil {
		return EmptyOracleInfo(ctx)
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

func (keeper Keeper) sign(ctx sdk.Context, p Payload, signer sdk.Address) {
	store := ctx.KVStore(keeper.key)
	key := GetSignKey(p, signer, keeper.cdc)
	store.Set(key, signer)
}

func (keeper Keeper) signed(ctx sdk.Context, p Payload, signer sdk.Address) bool {
	store := ctx.KVStore(keeper.key)
	key := GetSignKey(p, signer, keeper.cdc)
	return store.Has(key)
}

func (keeper Keeper) clearSigns(ctx sdk.Context, p Payload) {
	store := ctx.KVStore(keeper.key)
	prefix := GetSignPrefix(p, keeper.cdc)

	iter := sdk.KVStorePrefixIterator(store, prefix)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
	iter.Close()
}
