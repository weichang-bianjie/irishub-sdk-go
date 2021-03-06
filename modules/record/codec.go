package record

import (
	"github.com/irisnet/irishub-sdk-go/codec"
	"github.com/irisnet/irishub-sdk-go/codec/types"
	cryptocodec "github.com/irisnet/irishub-sdk-go/crypto/codec"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRecord{},
	)
}
