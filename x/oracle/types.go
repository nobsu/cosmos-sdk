package oracle

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

// implements sdk.Msg
type OracleMsg struct {
	Oracle types.Oracle
	Signer sdk.Address
}

func (msg OracleMsg) Type() string {
	return "oracle"
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

func (msg OracleMsg) ValidateBasic() sdk.Error {
	return msg.Oracle.ValidateBasic()
}

func (msg OracleMsg) GetSigners() []sdk.Address {
	return []sdk.Address{Oracle.Signer}
}
