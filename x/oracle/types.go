package oracle

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleMsg struct {
	Payload
	Signer sdk.Address
}

func (msg OracleMsg) Get(key interface{}) interface{} {
	return nil
}

func (msg OracleMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg OracleMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Signer}
}

func (msg OracleMsg) GetSigner() sdk.Address {
	return msg.Signer
}

type Payload interface {
	Type() string
	ValidateBasic() sdk.Error
}
