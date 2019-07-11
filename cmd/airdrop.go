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
	"github.com/360EntSecGroup-Skylar/excelize"
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
				faucet_name string
				faucet_add  string
				airdrop_list []types.AirDropInfo
				faucet_info types.AccountInfoRes
				//record_list  map[string]string
			)
			fmt.Println("Init AirDrop !!!!")
			helper.ReadConfigFile(FlagConfDir)

			//fmt.Println("Read transfer Record !!!!")
			//if record_list, err = helper.ReadRecord(); err != nil {
			//	return err
			//}

			if airdrop_list, xlsx, err = helper.ReadAddressList(FlagConfDir); err != nil {
				fmt.Println("ReadAddressList error !!!!")
				return err
			}

			faucet_name = "faucet_" + helper.RandomId()
			fmt.Printf("No.1  Create new faucet account : %s \n", faucet_name)

			if faucet_add, err = account.CreateAccount(faucet_name, constants.KeyPassword, conf.AirDropSeed); err != nil {
				fmt.Println("CreateAccount error !!!!")
				return err
			}

			if faucet_info, err = account.GetAccountInfo(faucet_add); err != nil {
				fmt.Println("GetAccountInfo error !!!!")
				return err
			}
			sequence, _ := helper.StrToInt(faucet_info.Value.Sequence)
			fmt.Printf("No.2  Get faucet sequence : %d \n", sequence)

			req := types.TransferTxReq{
				Amount: conf.AirDropAmount,
				Recipient: "",//faucet_info.Value.Address,
				BaseTx: types.BaseTx{
					LocalAccountName: faucet_name,
					Password:         constants.KeyPassword,
					ChainID:          conf.ChainId,
					AccountNumber:    faucet_info.Value.AccountNumber,
					Sequence:         helper.IntToStr(sequence),
					Gas:              conf.AirDropGas,
					Fees:             conf.AirDropFee,
					Memo:            "airdrop token",
				},
			}

			//判断余额
			faucetBalance, _ := account.ParseCoins(helper.IrisattoToIris(faucet_info.Value.Coins))
			minBalance, _ := account.ParseCoins(conf.AirDropAmount)
			minBalance = minBalance+10
			if faucetBalance < minBalance {
				return fmt.Errorf("Faucet balance = %f iris,  not enough balance for stress test! \n", faucetBalance)
			}
			fmt.Printf("No.3  faucetBalance : %f \n", faucetBalance)

			fmt.Printf("No.4  Transfer balance : %s to  accounts \n",conf.AirDropAmount)
			for i, _ := range airdrop_list{
				if !helper.IsCellEmpty(xlsx, airdrop_list[i]){
					index := helper.IntToStr(airdrop_list[i].Pos)
					fmt.Println("G"+index+" is not empty!")
					continue
				}

				////查重
				//if (record_list[airdrop_list[i].Address] != "") {
				//	fmt.Println("Duplicated transfer : "+req.Recipient+" to "+airdrop_list[i].Address)
				//	airdrop_list[i].Status = "Duplicated"
				//	airdrop_list[i].Hash = ""
				//	airdrop_list[i].TransactionTime = ""
				//	airdrop_list[i].Amount = ""
				//	helper.WriteAddressList(xlsx, airdrop_list[i])
				//	continue
				//}

				//随机金额
				if conf.AirDropRandom {
					if req.Amount, err = account.RandomCoin(conf.AirDropAmount); err != nil {
						return err
					}
				}

				req.Recipient = airdrop_list[i].Address

				//转账
				if txRes, err := tx.SendTx(req, faucet_info.Value.Address, false); err != nil {
					fmt.Println(err.Error())

					airdrop_list[i].Status = "Error"
					airdrop_list[i].Hash = err.Error()
					airdrop_list[i].TransactionTime = ""
					airdrop_list[i].Amount = ""
					helper.WriteAddressList(xlsx, airdrop_list[i])

					break
				} else {
					fmt.Printf("(%d) Send %s to %s ok! \n", i+1,req.Amount, airdrop_list[i].Address)

					sequence++
					req.BaseTx.Sequence = helper.IntToStr(sequence)
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

			if err = helper.SaveAddressList(xlsx, conf.AirDropXlsx); err != nil {
				fmt.Println("SaveAddressList error !!!! Save result to tempfile!!")
				helper.SaveAddressList(xlsx, conf.AirDropXlsxTemp)
				return err
			}

			fmt.Println("AirDrop end !!!!")
			return nil
		},
	}

	cmd.Flags().AddFlagSet(faucetInitSet)
	cmd.MarkFlagRequired(FlagConfDir)

	return cmd
}


