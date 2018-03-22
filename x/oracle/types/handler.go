package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Handler func(sdk.Context, Oracle) sdk.Error
