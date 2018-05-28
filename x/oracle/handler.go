package oracle

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Handler func(ctx sdk.Context, p Payload) sdk.Error

func updateInfo(ctx sdk.Context, val sdk.Validator, valset sdk.ValidatorSet, info OracleInfo) OracleInfo {
	info.Signers = append(info.Signers, val.GetOwner())
	info.Power = info.Power.Add(val.GetPower())

	supermaj := sdk.NewRat(2, 3)
	totalPower := valset.TotalPower(ctx)
	if !info.Power.GT(totalPower.Mul(supermaj)) {
		return info
	}

	hash := ctx.BlockHeader().ValidatorsHash
	if !bytes.Equal(hash, info.Hash) {
		newinfo := OracleInfo{
			Signers:   []sdk.Address{},
			Power:     sdk.ZeroRat(),
			Hash:      hash,
			Processed: false,
		}
		for _, s := range info.Signers {
			val := valset.Validator(ctx, s)
			if val != nil {
				newinfo.Signers = append(newinfo.Signers, val.GetOwner())
				newinfo.Power = newinfo.Power.Add(val.GetPower())
			}
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

	// Add the signer to signer queue
	for _, s := range info.Signers {
		if bytes.Equal(s, signer) {
			return ErrAlreadySigned(codespace).Result()
		}
	}

	info = updateInfo(ctx, val, valset, info)
	if info.Processed {
		info = OracleInfo{Processed: true}
	}

	keeper.setInfo(ctx, payload, info)

	if info.Processed {
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
