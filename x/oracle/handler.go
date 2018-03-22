package oracle

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/stake"
)

func NewHandler(keeper Keeper, sk stake.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case OracleMsg:
			return handleOracleMsg(ctx, keeper, sk, msg)
		default:
			errMsg := "Unrecognized oracle Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func getValidator(ctx sdk.Context, sk stake.Keeper, addr sdk.Address) (stake.Validator, bool) {
	valset := sk.GetValidators(ctx)
	for _, val := range valset {
		if val.Address == addr {
			return val, true
		}
	}
	return stake.Validator{}, false
}

func handleOracleMsg(ctx sdk.Context, keeper Keeper, sk stake.Keeper, msg OracleMsg) sdk.Result {
	// Check the signer is a validater
	val, ok := getValidator(ctx, sk, msg.Signer)
	if !ok {
		return ErrNotValidator(msg.Signer)
	}

	oracle := msg.Oracle
	info := keeper.OracleInfo(ctx, oracle)

	// Check the oracle is already processed
	if info.Processed {
		return ErrAlreadyProcessed()
	}

	// TODO: implement valset change later
	// Check if the valset hash is remaining same

	// Add the signer to signer queue
	for _, s := range info.Signers {
		if s == msg.Signer {
			return ErrAlreadySigned()
		}
	}
	info.Signers = append(info.Signers, msg.Signer)
	info.TotalPower += val.VotingPower

	if info.TotalPower >= sm.TotalPower(ctx)*2/3 { // TODO: make "2/3" modifiable
		keeper.Dispatcher.Dispatch(oracle.Type())(ctx, oracle)
	}

}
