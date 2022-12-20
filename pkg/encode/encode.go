package encode

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

// Codec creates a new Codec
func Codec() codec.Codec {
	interfaceRegistry := types.NewInterfaceRegistry()
	return codec.NewProtoCodec(interfaceRegistry)
}
