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

// 创建4个用户，并向他们分别转账测试所需的最低iris数量
func FaucetInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "init",
		Example:                "irishub-load init --config-dir=$HOME/config.json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err         error
				faucet_name string
				faucet_add  string
				faucet_info types.AccountInfoRes
			)
			fmt.Println("Init start !!!!")

			// 从config.json中读取参数
			helper.ReadConfigFile(FlagConfDir)
			faucet_name = "faucet_" + helper.RandomId()
			fmt.Printf("No.1  Create new faucet account : %s \n", faucet_name)

			// recover 水龙头用户， 这个账户（faa1lcuw6ewd2gfxap37sejewmta205sgssmv5fnju）需要在程序运行前手动转入一定数量的测试代币，
			// 然后再从这个账户把iris分别转到其他4个账户
			if faucet_add, err = account.CreateAccount(faucet_name, constants.KeyPassword, conf.FaucetSeed); err != nil {
				fmt.Println("CreateAccount error !!!!")
				return err
			}

			// 读取水龙头账户信息
			if faucet_info, err = account.GetAccountInfo(faucet_add); err != nil {
				fmt.Println("GetAccountInfo error !!!!"+err.Error())
				return err
			}
			sequence, _ := helper.StrToInt(faucet_info.Value.Sequence)
			fmt.Printf("No.2  Get faucet sequence : %d \n", sequence)

			//判断faucet的余额是否大于4倍min-balance（测试所需单个账户的最小余额）
			faucetBalance, _ := account.ParseCoins(helper.IrisattoToIris(faucet_info.Value.Coins[0].Amount))
			minBalance, _ := account.ParseCoins(conf.MinBalance)
			if faucetBalance < minBalance*4 {
				return fmt.Errorf("Faucet balance = %f iris,  not enough balance for stress test! \n", faucetBalance)
			}

			//构造转账交易
			req := types.TransferTxReq{
				Amount: conf.MinBalance,
				Recipient: "",
				BaseTx: types.BaseTx{
					LocalAccountName: faucet_name,
					Password:         constants.KeyPassword,
					ChainID:          conf.ChainId,
					AccountNumber:    faucet_info.Value.AccountNumber,
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

					//判断余额是否充足，充足则不转账，continue
					if balance, _ := account.ParseCoins(helper.IrisattoToIris(accInfo.Value.Coins[0].Amount)); balance >= minBalance {
						fmt.Printf("Enough balance for %s, balance %f iris >= minBalance : %f iris \n", subFaucet.FaucetAddr,balance,minBalance)
						continue
					}
				}

				req.Recipient = subFaucet.FaucetAddr

				//给指定地址转账minBalance ,500000iris
				if msg, err := tx.SendTx(req,  faucet_info.Value.Address, true); err != nil {
					fmt.Println(msg)
					return err
				} else {
					fmt.Printf("Send %s to %s succeed! \n", conf.MinBalance, subFaucet.FaucetAddr)
					//如果转账成功， sequence+1
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



