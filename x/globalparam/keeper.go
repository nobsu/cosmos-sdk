package globalparam

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type keeper struct {
	cdc *wire.Codec
	key sdk.StoreKey
}

func NewKeeper(key sdk.StoreKey, cdc *wire.Codec) keeper {
	return keeper{
		cdc: cdc,
		key: key,
	}
}

func (k keeper) get(ctx sdk.Context, key string, ptr interface{}) error {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(key))
	return k.cdc.UnmarshalBinary(bz, ptr)
}

func (k keeper) getBytes(ctx sdk.Context, key string) []byte {
	store := ctx.KVStore(k.key)
	return store.Get([]byte(key))
}

func (k keeper) set(ctx sdk.Context, key string, param interface{}) error {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(key))
	if bz != nil {
		ptrty := reflect.PtrTo(reflect.TypeOf(param))
		ptr := reflect.New(ptrty).Interface()

		if k.cdc.UnmarshalBinary(bz, ptr) != nil {
			return fmt.Errorf("Type mismatch with stored param and provided param")
		}
	}

	bz, err := k.cdc.MarshalBinary(param)
	if err != nil {
		return err
	}
	store.Set([]byte(key), bz)

	return nil
}

func (k keeper) Getter() Getter {
	return Getter{k}
}

func (k keeper) Setter() Setter {
	return Setter{k}
}

// Getter exposes methods related with only getting params
type Getter struct {
	k keeper
}

func (k Getter) Get(ctx sdk.Context, key string, ptr interface{}) error {
	return k.k.get(ctx, key, ptr)
}

func (k Getter) GetBytes(ctx sdk.Context, key string) []byte {
	return k.k.getBytes(ctx, key)
}

// Setter exposes all methods including Set
type Setter struct {
	k keeper
}

func (k Setter) Get(ctx sdk.Context, key string, ptr interface{}) error {
	return k.k.get(ctx, key, ptr)
}

func (k Setter) GetBytes(ctx sdk.Context, key string) []byte {
	return k.k.getBytes(ctx, key)
}

func (k Setter) Set(ctx sdk.Context, key string, param interface{}) error {
	return k.k.set(ctx, key, param)
}
