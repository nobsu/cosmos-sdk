package deriv

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(valset DerivValidatorSet) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeclareDeriv:
			return handleMsgDeclareDeriv(ctx, valset, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in deriv module").Result()
		}
	}
}

func handleMsgDeclareDeriv(ctx sdk.Context, valset DerivValidatorSet, msg MsgDeclareDeriv) sdk.Result {
	val := valset.valset.Validator(ctx, msg.Base)
	if val == nil {
		return ErrNotValidator(valset.codespace).Result()
	}

	valset.declareDeriv(ctx, msg.Base, msg.Deriv)
	return sdk.Result{}
}
