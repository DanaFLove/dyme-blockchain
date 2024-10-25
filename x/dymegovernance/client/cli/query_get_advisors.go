package cli

import (
	"strconv"

	"dymechain/x/dymegovernance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetAdvisors() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-advisors [advisorparam-1] [advisorparam-2]",
		Short: "Query get-advisors",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAdvisorparam1 := args[0]
			reqAdvisorparam2 := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetAdvisorsRequest{

				Advisorparam1: reqAdvisorparam1,
				Advisorparam2: reqAdvisorparam2,
			}

			res, err := queryClient.GetAdvisors(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
