package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
	// TODO include option to always include prefix bytes.
	cdc.RegisterConcrete(SendMsg{}, "cosmos-sdk/SendMsg", nil)
	cdc.RegisterConcrete(IssueMsg{}, "cosmos-sdk/IssueMsg", nil)

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature] types.
}
