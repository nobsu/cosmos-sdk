package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func GetInfoKey(p Payload, cdc *wire.Codec) []byte {
	bz := cdc.MustMarshalBinary(p)
	return append([]byte{0x00}, bz...)
}

func GetIsSignedPrefix(p Payload, cdc *wire.Codec) []byte {
	bz := cdc.MustMarshalBinary(p)
	return append([]byte{0x01}, bz...)
}

func GetIsSignedKey(p Payload, signer sdk.Address, cdc *wire.Codec) []byte {
	return append(GetIsSignedPrefix(p, cdc), signer...)
}
