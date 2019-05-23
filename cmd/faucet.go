package cmd

import (
	"fmt"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/irisnet/irishub-load/util/helper/account"
	"github.com/irisnet/irishub-load/util/helper/tx"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub-load/types"
)

func FaucetInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "init",
		Example:                "irishub-load init --config-dir=$HOME/local",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err         error
				faucet_name string
				faucet_add  string
				faucet_info types.AccountInfoRes
			)
			fmt.Println("Init start !!!!")
			helper.ReadConfigFile(FlagConfDir)
			faucet_name = "faucet_" + helper.RandomId()
			fmt.Printf("No.1  Create new faucet account : %s \n", faucet_name)

			if faucet_add, err = account.CreateAccount(faucet_name, constants.KeyPassword, conf.FaucetSeed); err != nil {
				fmt.Println("CreateAccount error !!!!")
				return err
			}

			if faucet_info, err = account.GetAccountInfo(faucet_add); err != nil {
				fmt.Println("GetAccountInfo error !!!!")
				return err
			}
			sequence, _ := helper.StrToInt(faucet_info.Sequence)
			fmt.Printf("No.2  Get faucet sequence : %d \n", sequence)

			//判断余额是否大于4倍min-balance,测试所需的最小余额
			faucetBalance, _ := account.ParseCoins(faucet_info.Coins[0])
			minBalance, _ := account.ParseCoins(conf.MinBalance)
			if faucetBalance < minBalance*4 {
				return fmt.Errorf("Faucet balance = %f iris,  not enough balance for stress test! \n", faucetBalance)
			}

			req := types.TransferTxReq{
				Amount: conf.MinBalance,
				Sender: faucet_info.Address,
				BaseTx: types.BaseTx{
					LocalAccountName: faucet_name,
					Password:         constants.KeyPassword,
					ChainID:          conf.ChainId,
					AccountNumber:    faucet_info.AccountNumber,
					Sequence:         helper.IntToStr(sequence),
					Gas:              constants.MockDefaultGas,
					Fees:             constants.MockDefaultFee,
					Memo:             fmt.Sprintf("transfer token"),
				},
			}

			//分别给4个账户转账
			fmt.Printf("No.3  Transfer balance : %s to 4 accounts \n",conf.MinBalance)
			for _, subFaucet := range conf.SubFaucets{
				if accInfo, err := account.GetAccountInfo(subFaucet.FaucetAddr); err == nil {

					//判断余额是否充足，充足则continue
					if balance, _ := account.ParseCoins(accInfo.Coins[0]); balance >= minBalance {
						fmt.Printf("Enough balance for %s , balance %f iris >= minBalance : %f iris \n", subFaucet.FaucetAddr, balance, minBalance)
						continue
					}
				}

				//转账最小测试金额
				if msg, err := tx.SendTx(req, subFaucet.FaucetAddr, true); err != nil {
					fmt.Println(msg)
					return err
				} else {
					sequence++
					req.BaseTx.Sequence = helper.IntToStr(sequence)
				}
			}

			fmt.Println("Init end !!!!")
			return nil
		},
	}

	cmd.Flags().AddFlagSet(faucetInitSet)
	cmd.MarkFlagRequired(FlagConfDir)

	return cmd
}



