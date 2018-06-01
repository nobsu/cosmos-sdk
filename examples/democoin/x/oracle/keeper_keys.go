package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// TODO: keys could conflict with each other
// since the user modules provide different codecs
// so disamb bytes are not prefixed
// prefix payload type before its marshalled bytes to fix it

// GetInfoKey returns the key for OracleInfo
func GetInfoKey(p Payload, cdc *wire.Codec) []byte {
	bz := cdc.MustMarshalBinary(p)
	return append([]byte{0x00}, bz...)
}

// GetSignPrefix returns the prefix for signs
func GetSignPrefix(p Payload, cdc *wire.Codec) []byte {
	bz := cdc.MustMarshalBinary(p)
	return append([]byte{0x01}, bz...)
}

// GetSignKey returns the key for sign
func GetSignKey(p Payload, signer sdk.Address, cdc *wire.Codec) []byte {
	return append(GetSignPrefix(p, cdc), signer...)
}
