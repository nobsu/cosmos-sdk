package deriv

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type DerivValidator struct {
	sdk.Validator
	Address sdk.Address
}

func (val DerivValidator) GetOwner() sdk.Address {
	return val.Address
}

type DerivValidatorSet struct {
	key       sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType

	valset sdk.ValidatorSet
}

func NewDerivValidatorSet(cdc *wire.Codec, key sdk.StoreKey, valset sdk.ValidatorSet) DerivValidatorSet {
	return DerivValidatorSet{
		key: key,
		cdc: cdc,

		valset: valset,
	}
}

// IterateValidators implements sdk.ValidatorSet
func (valset DerivValidatorSet) IterateValidators(ctx sdk.Context, fn func(int64, sdk.Validator) bool) {
	valset.valset.IterateValidators(ctx, func(index int64, val sdk.Validator) (stop bool) {
		return fn(index, valset.GetDerivValidator(ctx, val))
	})
}

func (valset DerivValidatorSet) IterateValidatorsBonded(ctx sdk.Context, fn func(int64, sdk.Validator) bool) {
	valset.valset.IterateValidatorsBonded(ctx, func(index int64, val sdk.Validator) (stop bool) {
		return fn(index, valset.GetDerivValidator(ctx, val))
	})
}

func (valset DerivValidatorSet) Validator(ctx sdk.Context, addr sdk.Address) sdk.Validator {
	base := valset.GetBaseValidator(ctx, addr)
	if base == nil {
		return nil
	}

	return DerivValidator{
		Validator: base,
		Address:   addr,
	}
}

func (valset DerivValidatorSet) TotalPower(ctx sdk.Context) sdk.Rat {
	return valset.valset.TotalPower(ctx)
}

// GetDerivBaseKey :: sdk.Address -> sdk.Address
func GetBaseKey(addr sdk.Address) []byte {
	return append([]byte{0x00}, addr...)
}

func (valset DerivValidatorSet) GetBaseValidator(ctx sdk.Context, addr sdk.Address) sdk.Validator {
	store := ctx.KVStore(valset.key)
	base := store.Get(GetBaseKey(addr))
	return valset.valset.Validator(ctx, base)
}

// GetDerivDerivKey :: sdk.Address -> sdk.Address
func GetDerivKey(addr sdk.Address) []byte {
	return append([]byte{0x01}, addr...)
}

func (valset DerivValidatorSet) GetDerivValidator(ctx sdk.Context, val sdk.Validator) sdk.Validator {
	store := ctx.KVStore(valset.key)
	deriv := store.Get(GetDerivKey(val.GetOwner()))
	if deriv == nil {
		return val
	}

	return DerivValidator{
		Validator: val,
		Address:   deriv,
	}
}

func (valset DerivValidatorSet) declareDeriv(ctx sdk.Context, base sdk.Address, deriv sdk.Address) {
	store := ctx.KVStore(valset.key)
	store.Set(GetBaseKey(deriv), base)
	store.Set(GetDerivKey(base), deriv)
}
