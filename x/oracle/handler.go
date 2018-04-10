package oracle

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/stake"
)

type Handler func(ctx sdk.Context, o Oracle) sdk.Error

func (keeper Keeper) Handle(h Handler, ctx sdk.Context, msg OracleMsg) sdk.Result {
	sk := keeper.sk

	// Check the signer is a validater
	val, ok := getValidator(ctx, keeper.sk, msg.Signer)
	if !ok {
		return ErrNotValidator(msg.Signer).Result()
	}

	oracle := msg.Oracle
	info := keeper.OracleInfo(ctx, oracle)

	// Check the oracle is already processed
	if info.Processed {
		return ErrAlreadyProcessed().Result()
	}

	// TODO: implement valset change later
	// Check if the valset hash is remaining same

	// Add the signer to signer queue
	for _, s := range info.Signers {
		if bytes.Equal(s, msg.Signer) {
			return ErrAlreadySigned().Result()
		}
	}
	info.Signers = append(info.Signers, msg.Signer)
	info.Power = info.Power.Add(val.Power)

	supermaj := sdk.NewRat(2, 3)
	totalPower := sk.GetPool(ctx).BondedShares
	if info.Power.GT(totalPower.Mul(supermaj)) { // TODO: make "2/3" modifiable
		cctx, write := ctx.CacheContext()
		err := h(cctx, oracle)
		info.Processed = true
		keeper.setInfo(ctx, oracle, info)
		if err != nil {
			return sdk.Result{
				Code: sdk.CodeOK,
				Log:  err.ABCILog(),
			}
		}
		write()
		return sdk.Result{}
	}
	keeper.setInfo(ctx, oracle, info)
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
