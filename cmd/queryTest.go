package cmd

import (
	"fmt"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/irisnet/irishub-load/util/helper/account"
	"github.com/spf13/cobra"
)

func QueryTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "querytest",
		Example: "irishub-load querytest --config-dir=$HOME",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
			//tx_info types.AccountInfoRes
			//err         error
			)

			fmt.Println("Start queryTest test!")
			helper.ReadConfigFile(FlagConfDir)

			//tx_info, err = account.GetTxInfo()

			account.GetTxInfo()
			return nil
		},
	}

	cmd.Flags().AddFlagSet(faucetInitSet)
	cmd.MarkFlagRequired(FlagConfDir)

	return cmd
}
