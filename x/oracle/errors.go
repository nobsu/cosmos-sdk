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

func ErrNotValidator(codespace sdk.CodespaceType, address sdk.Address) sdk.Error {
	return newError(codespace, CodeNotValidator, address.String())
}

func ErrAlreadyProcessed(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeAlreadyProcessed, "")
}

func ErrAlreadySigned(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeAlreadySigned, "")
}

// -------------------------
// Helpers

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}
