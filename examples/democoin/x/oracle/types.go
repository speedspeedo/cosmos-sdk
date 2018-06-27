package oracle

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg - struct for voting on payloads
type Msg struct {
	Payload
	Signer sdk.Address
}

// GetSignBytes implements sdk.Msg
func (msg Msg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners implements sdk.Msg
func (msg Msg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Signer}
}

// Payload defines inner data for actual execution
type Payload interface {
	Type() string
	ValidateBasic() sdk.Error
}
