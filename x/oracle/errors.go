package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Oracle errors reserve 1101-1199
	CodeNotValidator     sdk.CodeType = 1101
	CodeAlreadyProcessed sdk.CodeType = 1102
	CodeAlreadySigned    sdk.CodeType = 1103
	CodeUnknownRequest   sdk.CodeType = sdk.CodeUnknownRequest
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeNotValidator:
		return "Oracle is not signed by a validator"
	case CodeAlreadyProcessed:
		return "Oracle is already processed"
	case CodeAlreadySigned:
		return "Oracle is already signed by this signer"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func ErrNotValidator(address sdk.Address) sdk.Error {
	return newError(CodeNotValidator, address.String())
}

func ErrAlreadyProcessed() sdk.Error {
	return newError(CodeAlreadyProcessed, "")
}

func ErrAlreadySigned() sdk.Error {
	return newError(CodeAlreadySigned, "")
}

// -------------------------
// Helpers

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}
