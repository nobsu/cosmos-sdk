package oracle

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

func GetInfoKey(p Payload, cdc *wire.Codec) []byte {
	bz := cdc.MustMarshalBinary(p)
	return append([]byte{0x00}, bz...)
}

func GetValidatorsHashKey() []byte {
	return []byte{0x01}
}
