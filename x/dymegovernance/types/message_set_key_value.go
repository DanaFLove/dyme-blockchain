package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetKeyValue = "set_key_value"

var _ sdk.Msg = &MsgSetKeyValue{}

func NewMsgSetKeyValue(creator string, key string, value string) *MsgSetKeyValue {
	return &MsgSetKeyValue{
		Creator: creator,
		Key:     key,
		Value:   value,
	}
}

func (msg *MsgSetKeyValue) Route() string {
	return RouterKey
}

func (msg *MsgSetKeyValue) Type() string {
	return TypeMsgSetKeyValue
}

func (msg *MsgSetKeyValue) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetKeyValue) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetKeyValue) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
