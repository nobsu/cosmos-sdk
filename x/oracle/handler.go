package oracle

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
)

func NewHandler(keeper Keeper, sm staking.StakingMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case OracleMsg:
			return handleOracleMsg(ctx, keeper, sm, msg)
		default:
			errMsg := "Unrecognized oracle Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleOracleMsg(ctx sdk.Context, keeper Keeper, sm staking.StakingMapper, msg OracleMsg) sdk.Result {

}
