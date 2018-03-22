package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Oracle interface {
	Type() string
	ValidateBasic() sdk.Error
}
