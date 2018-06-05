package deriv

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/wire"
)

var msgCdc = wire.NewCodec()

type MsgDeclareDeriv struct {
	Base  sdk.Address
	Deriv sdk.Address
}

func (msg MsgDeclareDeriv) GetSignBytes() []byte {
	return msgCdc.MustMarshalBinary(msg)
}

func (msg MsgDeclareDeriv) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Base}
}

func (msg MsgDeclareDeriv) Type() string {
	return "deriv"
}

func (msg MsgDeclareDeriv) ValidateBasic() sdk.Error {
	if msg.Base == nil || msg.Deriv == nil {
		return ErrEmptyValidator(DefaultCodespace)
	}
	return nil
}
