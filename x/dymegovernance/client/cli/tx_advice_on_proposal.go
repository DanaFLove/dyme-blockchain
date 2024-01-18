package cli

import (
	"strconv"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAdviceOnProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advice-on-proposal [proposal-id] [advisory-outcome]",
		Short: "Broadcast message advice-on-proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argProposalId := args[0]
			argAdvisoryOutcome := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAdviceOnProposal(
				clientCtx.GetFromAddress().String(),
				argProposalId,
				argAdvisoryOutcome,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
