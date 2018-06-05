package deriv

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

const (
	DefaultCodespace sdk.CodespaceType = 11

	CodeInvalidValidator sdk.CodeType = stake.CodeInvalidValidator
)

func ErrNotValidator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "Not a validator")
}

func ErrEmptyValidator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "Validator address is empty")
}
