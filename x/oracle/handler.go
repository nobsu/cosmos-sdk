package oracle

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/stake"
)

type Handler func(ctx sdk.Context, p Payload) sdk.Error

func (keeper Keeper) Handle(h Handler, ctx sdk.Context, o Oracle) sdk.Result {
	sk := keeper.sk

	signer := o.GetSigner()

	// Check the signer is a validater
	val, ok := getValidator(ctx, keeper.sk, signer)
	if !ok {
		return ErrNotValidator(signer).Result()
	}

	info := keeper.OracleInfo(ctx, o)

	// Check the oracle is already processed
	if info.Processed {
		return ErrAlreadyProcessed().Result()
	}

	// TODO: implement valset change later
	// Check if the valset hash is remaining same

	// Add the signer to signer queue
	for _, s := range info.Signers {
		if bytes.Equal(s, signer) {
			return ErrAlreadySigned().Result()
		}
	}
	info.Signers = append(info.Signers, signer)
	info.Power = info.Power.Add(val.Power)

	supermaj := sdk.NewRat(2, 3)
	totalPower := sk.GetPool(ctx).BondedShares
	if info.Power.GT(totalPower.Mul(supermaj)) { // TODO: make "2/3" modifiable
		cctx, write := ctx.CacheContext()
		err := h(cctx, o)
		info.Processed = true
		keeper.setInfo(ctx, o, info)
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeOK,
				Log:  err.ABCILog(),
			}
		}
		write()
		return sdk.Result{}
	}
	keeper.setInfo(ctx, o, info)
	return sdk.Result{}
}

func getValidator(ctx sdk.Context, sk stake.Keeper, addr sdk.Address) (stake.Validator, bool) {
	valset := sk.GetValidators(ctx)
	for _, val := range valset {
		if bytes.Equal(val.Address, addr) {
			return val, true
		}
	}
	return stake.Validator{}, false
}
