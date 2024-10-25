package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAdviceOnProposal = "advice_on_proposal"

var _ sdk.Msg = &MsgAdviceOnProposal{}

func NewMsgAdviceOnProposal(creator string, proposalId string, advisoryOutcome string) *MsgAdviceOnProposal {
	return &MsgAdviceOnProposal{
		Creator:         creator,
		ProposalId:      proposalId,
		AdvisoryOutcome: advisoryOutcome,
	}
}

func (msg *MsgAdviceOnProposal) Route() string {
	return RouterKey
}

func (msg *MsgAdviceOnProposal) Type() string {
	return TypeMsgAdviceOnProposal
}

func (msg *MsgAdviceOnProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAdviceOnProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAdviceOnProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
