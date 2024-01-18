package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakedyme = "stakedyme"

var _ sdk.Msg = &MsgStakedyme{}

func NewMsgStakedyme(creator string) *MsgStakedyme {
	return &MsgStakedyme{
		Creator: creator,
	}
}

func (msg *MsgStakedyme) Route() string {
	return RouterKey
}

func (msg *MsgStakedyme) Type() string {
	return TypeMsgStakedyme
}

func (msg *MsgStakedyme) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStakedyme) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStakedyme) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
