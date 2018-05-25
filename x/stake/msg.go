package stake

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	crypto "github.com/tendermint/go-crypto"
)

// name to idetify transaction types
const MsgType = "stake"

// XXX remove: think it makes more sense belonging with the Params so we can
// initialize at genesis - to allow for the same tests we should should make
// the ValidateBasic() function a return from an initializable function
// ValidateBasic(bondDenom string) function
const StakingToken = "steak"

//Verify interface at compile time
var _, _, _, _ sdk.Msg = &MsgDeclareCandidacy{}, &MsgEditCandidacy{}, &MsgDelegate{}, &MsgUnbond{}

var msgCdc = wire.NewCodec()

func init() {
	wire.RegisterCrypto(msgCdc)
}

//______________________________________________________________________

// MsgDeclareCandidacy - struct for unbonding transactions
type MsgDeclareCandidacy struct {
	Description
	ValidatorAddr sdk.Address   `json:"address"`
	PubKey        crypto.PubKey `json:"pubkey"`
	Bond          sdk.Coin      `json:"bond"`
}

func NewMsgDeclareCandidacy(validatorAddr sdk.Address, pubkey crypto.PubKey,
	bond sdk.Coin, description Description) MsgDeclareCandidacy {
	return MsgDeclareCandidacy{
		Description:   description,
		ValidatorAddr: validatorAddr,
		PubKey:        pubkey,
		Bond:          bond,
	}
}

//nolint
func (msg MsgDeclareCandidacy) Type() string              { return MsgType } //TODO update "stake/declarecandidacy"
func (msg MsgDeclareCandidacy) GetSigners() []sdk.Address { return []sdk.Address{msg.ValidatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgDeclareCandidacy) GetSignBytes() []byte {
	return msgCdc.MustMarshalBinary(msg)
}

// quick validity check
func (msg MsgDeclareCandidacy) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr == nil {
		return ErrValidatorEmpty(DefaultCodespace)
	}
	if msg.Bond.Denom != StakingToken {
		return ErrBadBondingDenom(DefaultCodespace)
	}
	if msg.Bond.Amount <= 0 {
		return ErrBadBondingAmount(DefaultCodespace)
	}
	empty := Description{}
	if msg.Description == empty {
		return newError(DefaultCodespace, CodeInvalidInput, "description must be included")
	}
	return nil
}

//______________________________________________________________________

// MsgEditCandidacy - struct for editing a validator
type MsgEditCandidacy struct {
	Description
	ValidatorAddr sdk.Address `json:"address"`
}

func NewMsgEditCandidacy(validatorAddr sdk.Address, description Description) MsgEditCandidacy {
	return MsgEditCandidacy{
		Description:   description,
		ValidatorAddr: validatorAddr,
	}
}

//nolint
func (msg MsgEditCandidacy) Type() string              { return MsgType } //TODO update "stake/msgeditcandidacy"
func (msg MsgEditCandidacy) GetSigners() []sdk.Address { return []sdk.Address{msg.ValidatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgEditCandidacy) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgEditCandidacy) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr == nil {
		return ErrValidatorEmpty(DefaultCodespace)
	}
	empty := Description{}
	if msg.Description == empty {
		return newError(DefaultCodespace, CodeInvalidInput, "Transaction must include some information to modify")
	}
	return nil
}

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddr sdk.Address `json:"address"`
	ValidatorAddr sdk.Address `json:"address"`
	Bond          sdk.Coin    `json:"bond"`
}

func NewMsgDelegate(delegatorAddr, validatorAddr sdk.Address, bond sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
		Bond:          bond,
	}
}

//nolint
func (msg MsgDelegate) Type() string              { return MsgType } //TODO update "stake/msgeditcandidacy"
func (msg MsgDelegate) GetSigners() []sdk.Address { return []sdk.Address{msg.DelegatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgDelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddr == nil {
		return ErrBadDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddr == nil {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.Bond.Denom != StakingToken {
		return ErrBadBondingDenom(DefaultCodespace)
	}
	if msg.Bond.Amount <= 0 {
		return ErrBadBondingAmount(DefaultCodespace)
	}
	return nil
}

//______________________________________________________________________

// MsgUnbond - struct for unbonding transactions
type MsgUnbond struct {
	DelegatorAddr sdk.Address `json:"address"`
	ValidatorAddr sdk.Address `json:"address"`
	Shares        string      `json:"shares"`
}

func NewMsgUnbond(delegatorAddr, validatorAddr sdk.Address, shares string) MsgUnbond {
	return MsgUnbond{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
		Shares:        shares,
	}
}

//nolint
func (msg MsgUnbond) Type() string              { return MsgType } //TODO update "stake/msgeditcandidacy"
func (msg MsgUnbond) GetSigners() []sdk.Address { return []sdk.Address{msg.DelegatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgUnbond) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgUnbond) ValidateBasic() sdk.Error {
	if msg.DelegatorAddr == nil {
		return ErrBadDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddr == nil {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.Shares != "MAX" {
		rat, err := sdk.NewRatFromDecimal(msg.Shares)
		if err != nil {
			return ErrBadShares(DefaultCodespace)
		}
		if rat.IsZero() || rat.LT(sdk.ZeroRat()) {
			return ErrBadShares(DefaultCodespace)
		}
	}
	return nil
}

//______________________________________________________________________

// MsgUnrevoke - struct for unrevoking revoked validator
type MsgUnrevoke struct {
	ValidatorAddr sdk.Address `json:"address"`
}

func NewMsgUnrevoke(validatorAddr sdk.Address) MsgUnrevoke {
	return MsgUnrevoke{
		ValidatorAddr: validatorAddr,
	}
}

func (msg MsgUnrevoke) Type() string              { return MsgType }
func (msg MsgUnrevoke) GetSigners() []sdk.Address { return []sdk.Address{msg.ValidatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgUnrevoke) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgUnrevoke) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr == nil {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	return nil
}
