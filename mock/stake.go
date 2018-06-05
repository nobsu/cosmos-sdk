package mock

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
)

type Validator struct {
	Address sdk.Address
	Power   sdk.Rat
}

func (v Validator) GetStatus() sdk.BondStatus {
	return sdk.Bonded
}

func (v Validator) GetOwner() sdk.Address {
	return v.Address
}

func (v Validator) GetPubKey() crypto.PubKey {
	return nil
}

func (v Validator) GetPower() sdk.Rat {
	return v.Power
}

func (v Validator) GetBondHeight() int64 {
	return 0
}

type ValidatorSet struct {
	Validators []Validator
}

// IterateValidators implements sdk.ValidatorSet
func (vs *ValidatorSet) IterateValidators(ctx sdk.Context, fn func(index int64, Validator sdk.Validator) bool) {
	for i, val := range vs.Validators {
		if fn(int64(i), val) {
			break
		}
	}
}

// IterateValidatorsBonded implements sdk.ValidatorSet
func (vs *ValidatorSet) IterateValidatorsBonded(ctx sdk.Context, fn func(index int64, Validator sdk.Validator) bool) {
	vs.IterateValidators(ctx, fn)
}

// Validator implements sdk.ValidatorSet
func (vs *ValidatorSet) Validator(ctx sdk.Context, addr sdk.Address) sdk.Validator {
	for _, val := range vs.Validators {
		if bytes.Equal(val.Address, addr) {
			return val
		}
	}
	return nil
}

// TotalPower implements sdk.ValidatorSet
func (vs *ValidatorSet) TotalPower(ctx sdk.Context) sdk.Rat {
	res := sdk.ZeroRat()
	for _, val := range vs.Validators {
		res = res.Add(val.Power)
	}
	return res
}

func (vs *ValidatorSet) AddValidator(val Validator) {
	vs.Validators = append(vs.Validators, val)
}

func (vs *ValidatorSet) RemoveValidator(addr sdk.Address) {
	pos := -1
	for i, val := range vs.Validators {
		if bytes.Equal(val.Address, addr) {
			pos = i
			break
		}
	}
	if pos == -1 {
		return
	}
	vs.Validators = append(vs.Validators[:pos], vs.Validators[pos+1:]...)
}
