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

	/*  先看这里：
	0. 编译主网iris版本， ide 编辑配置中 airDrop --config-dir=$HOME

	1. 打开lcd（注意要主网版本）
		irislcd start --node=http://seed-1.mainnet.irisnet.org:26657 --laddr=tcp://0.0.0.0:1317 --chain-id=irisnet --home=$HOME/.irislcd/ --trust-node

	明细查账：     https://www.irisplorer.io/#/txs/transfers
	查水龙头余额： https://www.irisplorer.io/#/address/iaa13nzsae74qype65rshc0wyvhk9s0l3uecwf8y93

	2. 把主网版本的data和config在home目录（/Users/sherlock/）下，改下config中空投的助记词 strike-drop：

	3. 确保 iaa13nzsae74qype65rshc0wyvhk9s0l3uecwf8y93 有钱

	4. 把create account删除了（lcd不支持key了），后续没有用了。实在要用可以考虑iks。
	   助记词可以生成签名后的交易。 简单查账的话可以把水龙头的iaa在config里面事先写好。

	5. 完成后会把交易hash和结果 写回data.xlsl

	6. 另外注意下：收款地址不要有 iaa1mvfej6hvkuplcvm9aa2qdeh54npvdnshzcjpat 这个是交易所地址，需要备注memo， 先核查。
	*/

	cmd := &cobra.Command{
		Use:     "airDrop",
		Example: "irishub-load airDrop --config-dir=$HOME",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err          error
				xlsx         *excelize.File
				airdrop_list []types.AirDropInfo
				faucet_info  types.AccountInfoRes
			)
			fmt.Println("No.1  Init AirDrop !!!!")
			helper.ReadConfigFile(FlagConfDir)

			//从excel中读取金额和地址列表
			if airdrop_list, xlsx, err = helper.ReadAddressList(FlagConfDir); err != nil {
				fmt.Println("ReadAddressList error !!!!")
				return err
			}

			// 读出水龙头信息， iaa地址：conf.AirDropAddr
			if faucet_info, err = account.GetAccountInfo(conf.AirDropAddr); err != nil {
				fmt.Println("GetAccountInfo error !!!!")
				return err
			}
			sequence := faucet_info.Value.Sequence
			fmt.Printf("No.2  Get faucet sequence : %d \n", sequence)

			//构造转账交易
			req := types.TransferTxReq{
				Amount:     conf.AirDropAmount, //如果不是固定值，后面会重新赋值
				ChainID:    conf.AirDropChainId,
				Sequence:   sequence,
				SenderAddr: conf.AirDropAddr,
				SenderSeed: conf.AirDropSeed,
				Mode:       "", //调试的时候可以用 "commit=true",
			}

			//判断余额
			faucetBalance, _ := account.ParseCoins(helper.IrisattoToIris(faucet_info.Value.Coins))
			minBalance, _ := account.ParseCoins(conf.AirDropAmount)
			minBalance = minBalance + 10
			if faucetBalance < minBalance {
				return fmt.Errorf("Faucet balance = %f iris,  not enough balance for stress test! \n", faucetBalance)
			}
			fmt.Printf("No.3  faucetBalance : %f \n", faucetBalance)

			for i := range airdrop_list {
				req.Amount = airdrop_list[i].Amount //如果是固定值，这里可以注释掉
				req.RecipientAddr = airdrop_list[i].Address

				//转账
				if txRes, err := tx.SendTx(req); err != nil {
					//转账失败，break
					fmt.Println(err.Error())

					airdrop_list[i].Status = "Error"
					airdrop_list[i].Hash = err.Error()
					helper.WriteAddressList(xlsx, airdrop_list[i])
					break
				} else {
					//转账成功
					fmt.Printf("(%d) Send %s iris to %s ok! \n", i+1, req.Amount, airdrop_list[i].Address)

					req.Sequence++
					airdrop_list[i].Status = ""
					airdrop_list[i].Hash = txRes.Hash

					//把数据都写到内存列表中
					helper.WriteAddressList(xlsx, airdrop_list[i])
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			}

			fmt.Println("-------------\nSend TX ok, wait 15 seconds for Check TX !\n-------------")
			time.Sleep(time.Duration(15) * time.Second)

			fmt.Println("No.4  Check if TX succeeed \n")
			for i, sub := range airdrop_list {
				if sub.Status != "" || sub.Hash == "" {
					continue
				}

				if err := tx.CheckTx(sub.Hash); err != nil {
					fmt.Printf("(%d) Check TX Error : "+sub.Address+" "+sub.Hash+"\n", i+1)
					sub.Status = "CheckTXError"
				} else {
					fmt.Printf("(%d) "+sub.Address+" "+sub.Hash+" Check TX OK !!\n", i+1)
					sub.Status = "Succeed"
				}

				//最终把所有数据都写到内存列表中
				helper.WriteAddressList(xlsx, sub)
			}

			//hash写c列 结果写d列， 注意data.xlsx要关闭
			if err = helper.SaveAddressList(xlsx, conf.AirDropXlsx); err != nil {
				fmt.Println("SaveAddressList error !!!! ")
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
