package oracle

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Handler func(ctx sdk.Context, p Payload) sdk.Error

func (keeper Keeper) update(ctx sdk.Context, val sdk.Validator, valset sdk.ValidatorSet, p Payload, info OracleInfo) OracleInfo {
	info.Power = info.Power.Add(val.GetPower())

	supermaj := sdk.NewRat(2, 3)
	totalPower := valset.TotalPower(ctx)
	if !info.Power.GT(totalPower.Mul(supermaj)) {
		return info
	}

	hash := ctx.BlockHeader().ValidatorsHash
	if !bytes.Equal(hash, info.Hash) {
		newinfo := OracleInfo{
			Power:     sdk.ZeroRat(),
			Hash:      hash,
			Processed: false,
		}
		prefix := GetSignPrefix(p, keeper.cdc)
		store := ctx.KVStore(keeper.key)
		iter := sdk.KVStorePrefixIterator(store, prefix)
		for ; iter.Valid(); iter.Next() {
			if valset.Validator(ctx, iter.Value()) != nil {
				store.Delete(iter.Key())
				continue
			}
			newinfo.Power = newinfo.Power.Add(val.GetPower())
		}
		if newinfo.Power.GT(totalPower.Mul(supermaj)) {
			newinfo.Processed = true
		}
		return newinfo
	}

	info.Processed = true
	return info
}

func (keeper Keeper) Handle(h Handler, ctx sdk.Context, o OracleMsg, codespace sdk.CodespaceType) sdk.Result {
	valset := keeper.valset

	signer := o.Signer
	payload := o.Payload

	// Check the signer is a validater
	val := valset.Validator(ctx, signer)
	if val == nil {
		return ErrNotValidator(codespace, signer).Result()
	}

	info := keeper.OracleInfo(ctx, payload)

	// Check the oracle is already processed
	if info.Processed {
		return ErrAlreadyProcessed(codespace).Result()
	}

	// Check double signing
	if keeper.signed(ctx, payload, signer) {
		return ErrAlreadySigned(codespace).Result()
	}

	keeper.sign(ctx, payload, signer)

	info = keeper.update(ctx, val, valset, payload, info)
	if info.Processed {
		info = OracleInfo{Processed: true}
	}

	keeper.setInfo(ctx, payload, info)

	if info.Processed {
		keeper.clearSigns(ctx, payload)
		cctx, write := ctx.CacheContext()
		err := h(cctx, payload)
		if err != nil {
			return sdk.Result{
				Code: sdk.ABCICodeOK,
				Log:  err.ABCILog(),
			}
		}
		write()

	}

	return sdk.Result{}
}
