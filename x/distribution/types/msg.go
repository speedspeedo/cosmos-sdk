//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to identify transaction types
const MsgType = "distr"

// Verify interface at compile time
var _, _ sdk.Msg = &MsgModifyWithdrawAddress{}, &MsgWithdrawDelegatorRewardsAll{}
var _, _ sdk.Msg = &MsgWithdrawDelegationReward{}, &MsgWithdrawValidatorRewardsAll{}

//______________________________________________________________________

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgModifyWithdrawAddress struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	WithdrawAddr  sdk.AccAddress `json:"delegator_addr"`
}

func NewMsgModifyWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) MsgModifyWithdrawAddress {
	return MsgModifyWithdrawAddress{
		DelegatorAddr: delAddr,
		WithdrawAddr:  withdrawAddr,
	}
}

func (msg MsgModifyWithdrawAddress) Type() string { return MsgType }
func (msg MsgModifyWithdrawAddress) Name() string { return "withdraw_delegation_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgModifyWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgModifyWithdrawAddress) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgModifyWithdrawAddress) ValidateBasic() sdk.Error {
	if msg.DelegatorAddr == nil {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.WithdrawAddr == nil {
		return ErrNilWithdrawAddr(DefaultCodespace)
	}
	return nil
}

//______________________________________________________________________

// msg struct for delegation withdraw for all of the delegator's delegations
type MsgWithdrawDelegatorRewardsAll struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
}

func NewMsgWithdrawDelegationRewardsAll(delAddr sdk.AccAddress) MsgWithdrawDelegatorRewardsAll {
	return MsgWithdrawDelegatorRewardsAll{
		DelegatorAddr: delAddr,
	}
}

func (msg MsgWithdrawDelegatorRewardsAll) Type() string { return MsgType }
func (msg MsgWithdrawDelegatorRewardsAll) Name() string { return "withdraw_delegation_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorRewardsAll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorRewardsAll) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgWithdrawDelegatorRewardsAll) ValidateBasic() sdk.Error {
	if msg.DelegatorAddr == nil {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

//______________________________________________________________________

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegationReward struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
}

func NewMsgWithdrawDelegationReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) MsgWithdrawDelegationReward {
	return MsgWithdrawDelegationReward{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
	}
}

func (msg MsgWithdrawDelegationReward) Type() string { return MsgType }
func (msg MsgWithdrawDelegationReward) Name() string { return "withdraw_delegation_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegationReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegationReward) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgWithdrawDelegationReward) ValidateBasic() sdk.Error {
	if msg.DelegatorAddr == nil {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddr == nil {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

//______________________________________________________________________

// msg struct for validator withdraw
type MsgWithdrawValidatorRewardsAll struct {
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
}

func NewMsgWithdrawValidatorRewardsAll(valAddr sdk.ValAddress) MsgWithdrawValidatorRewardsAll {
	return MsgWithdrawValidatorRewardsAll{
		ValidatorAddr: valAddr,
	}
}

func (msg MsgWithdrawValidatorRewardsAll) Type() string { return MsgType }
func (msg MsgWithdrawValidatorRewardsAll) Name() string { return "withdraw_validator_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorRewardsAll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorRewardsAll) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgWithdrawValidatorRewardsAll) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr == nil {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}
