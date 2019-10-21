package cmd

import (
	"fmt"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/irisnet/irishub-load/util/helper/account"
	"github.com/irisnet/irishub-load/util/helper/tx"
	"github.com/spf13/cobra"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/irisnet/irishub-load/types"
	"time"
)

func AirDrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "airDrop",
		Example:                "irishub-load airDrop --config-dir=$HOME",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err         error
				xlsx        *excelize.File
				//faucet_name string
				//faucet_add  string
				airdrop_list []types.AirDropInfo
				faucet_info types.AccountInfoRes
			)
			fmt.Println("Init AirDrop !!!!")
			helper.ReadConfigFile(FlagConfDir)

			if airdrop_list, xlsx, err = helper.ReadAddressList(FlagConfDir); err != nil {
				fmt.Println("ReadAddressList error !!!!")
				return err
			}

			// AirDropAddr 直接卸载conf里面了. 助记词解析出地址，不要了。后续有时间改下。
			/*
				faucet_name = "faucet_" + helper.RandomId()
				fmt.Printf("No.1  Create new faucet account : %s \n", faucet_name)
				if faucet_add, err = account.CreateAccount(faucet_name, constants.KeyPassword, conf.AirDropSeed); err != nil {
					fmt.Println("CreateAccount error !!!!")
					return err
				}
			*/


			/*  先看这里：
				打开lcd ：irislcd-mainnet start --node=http://seed-1.mainnet.irisnet.org:26657 --laddr=tcp://0.0.0.0:1317 --chain-id=irisnet --home=$HOME/.irislcd/ --trust-node
				把data.xlsx和config.json copy到 home目录
				这里把create 注释掉了， 地址有钱， 助记词对就行。 后面可以直接签名，广播。
				AirDropAddr 直接卸载conf里面了，不从助记词读出。
			 */

			// 主网地址发币地址 ： conf.AirDropAddr
			if faucet_info, err = account.GetAccountInfo(conf.AirDropAddr); err != nil {
				fmt.Println("GetAccountInfo error !!!!")
				return err
			}
			sequence, _ := helper.StrToInt(faucet_info.Value.Sequence)
			fmt.Printf("No.2  Get faucet sequence : %d \n", sequence)

			//构造转账交易
			req := types.TransferTxReq{
				Amount:           "", //conf.AirDropAmount,
				ChainID:          conf.AirDropChainId,
				Sequence:         sequence,
				SenderAddr:       conf.AirDropAddr,
				SenderSeed:       conf.AirDropSeed,
				Mode:              "",
			}

			//判断余额
			faucetBalance, _ := account.ParseCoins(helper.IrisattoToIris(faucet_info.Value.Coins))
			minBalance, _ := account.ParseCoins(conf.AirDropAmount)
			minBalance = minBalance+10
			if faucetBalance < minBalance {
				return fmt.Errorf("Faucet balance = %f iris,  not enough balance for stress test! \n", faucetBalance)
			}
			fmt.Printf("No.3  faucetBalance : %f \n", faucetBalance)

			//fmt.Printf("No.4  Transfer balance : %s to  accounts \n",conf.AirDropAmount)
			for i, _ := range airdrop_list{
				if !helper.IsCellEmpty(xlsx, airdrop_list[i]){
					index := helper.IntToStr(airdrop_list[i].Pos)
					fmt.Println("G"+index+" is not empty!")
					continue
				}

				req.RecipientAddr = airdrop_list[i].Address
				req.Amount = airdrop_list[i].Amount

				//转账
				if txRes, err := tx.SendTx(req); err != nil {
					fmt.Println(err.Error())

					airdrop_list[i].Status = "Error"
					airdrop_list[i].Hash = err.Error()
					airdrop_list[i].TransactionTime = ""
					airdrop_list[i].Amount = ""
					helper.WriteAddressList(xlsx, airdrop_list[i])

					break
				} else {
					fmt.Printf("(%d) Send %s to %s ok! \n", i+1,req.Amount, airdrop_list[i].Address)

					req.Sequence++
					airdrop_list[i].Status = ""
					airdrop_list[i].Hash = txRes.Hash
					airdrop_list[i].TransactionTime = time.Now().UTC().Format(time.UnixDate)
					airdrop_list[i].Amount = req.Amount

					helper.WriteAddressList(xlsx, airdrop_list[i])
					time.Sleep(time.Duration(100)*time.Millisecond)
				}
			}

			fmt.Println("-------------\nSend TX ok, wait 15 seconds for Check TX !\n-------------")
			time.Sleep(time.Duration(15)*time.Second)

			fmt.Println("No.5  Check if TX succeeed \n")
			for _, sub := range airdrop_list{
				if sub.Status != "" || sub.Hash == ""{
					continue
				}

				if  err := tx.CheckTx(sub.Hash); err != nil {
					fmt.Println("Check TX Error : "+sub.Address+" "+ sub.Hash)
					sub.Status = "CheckTXError"
				}else{
					fmt.Println(sub.Address+" "+ sub.Hash + " Check TX OK !!")
					sub.Status = "Succeed"
				}

				helper.WriteAddressList(xlsx, sub)
			}

			//if err = helper.SaveAddressList(xlsx, conf.AirDropXlsx); err != nil {
			//	fmt.Println("SaveAddressList error !!!! Save result to tempfile!!")
			//	helper.SaveAddressList(xlsx, conf.AirDropXlsxTemp)
			//	return err
			//}

			fmt.Println("AirDrop end !!!!")
			return nil
		},
	}

	cmd.Flags().AddFlagSet(faucetInitSet)
	cmd.MarkFlagRequired(FlagConfDir)

	return cmd
}


