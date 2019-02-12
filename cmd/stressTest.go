package cmd

import (
	"fmt"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/sign"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/irisnet/irishub-load/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
	"os"
	"bufio"
	"time"
)

func SignTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signtx",
		Example: "irishub-load signtx --config-dir=$HOME/local --tps=1 --duration=1 --account=wenxi",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err                   error
				subFaucetIndex        int
				tps                   int
				duration              int

				PrivateInfo			  types.AccountTestPrivateInfo
				SignedDataString      []string
			)

			helper.ReadConfigFile(FlagConfDir)

			if tps = viper.GetInt(FlagTps);tps <= 0 {
				return fmt.Errorf("tps should > 0 ")
			}
			if duration = viper.GetInt(FlagDuration);duration <= 0 {
				return fmt.Errorf("duration should > 0 ")
			}
			if subFaucetIndex = helper.PraseUser(viper.GetString(FlagAccount));subFaucetIndex == -1 {
				return fmt.Errorf("account %s not exist", viper.GetString(FlagAccount))
			}

			totalTxNum := tps * duration * 60

			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Start signing %d TXs, chain id: %s, Node : %s, sub_faucet : %s \n", 	totalTxNum, conf.ChainId, conf.NodeUrl, conf.SubFaucets[subFaucetIndex].FaucetAddr)

			//读取私钥信息
			if PrivateInfo, err = sign.InitAccountSignProcess(conf.SubFaucets[subFaucetIndex].FaucetAddr, conf.SubFaucets[subFaucetIndex].Seed); err!=nil {
				return fmt.Errorf("Get private info error : %s", err.Error())
			}
			//生成签名数据
			if SignedDataString, err = sign.GenSignTxByTend(totalTxNum, subFaucetIndex, conf.ChainId, conf.SubFaucets, PrivateInfo); err!=nil {
				return fmt.Errorf("GenSignTx error : %s", err.Error())
			}

			log.Printf("Finish signing TXs \n")

			//把签名数据写入文件
			conf.Output = helper.GetPath(conf.Output)
			filePath := fmt.Sprintf("%v/SignedTXs", conf.Output)
			if err = helper.CreateFolder(conf.Output) ; err != nil {
				return err
			}
			if err = helper.WriteFile(filePath, []byte(strings.Join(SignedDataString, "\n"))); err != nil {
				return err
			}

			log.Printf("Finish WriteFile into %s \n",filePath )

			return nil
		},
	}

	cmd.Flags().AddFlagSet(signTXFlagSet)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTps)
	cmd.MarkFlagRequired(FlagDuration)
	cmd.MarkFlagRequired(FlagAccount)

	return cmd
}

func BroadcastTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast",
		Example: "irishub-load broadcast --config-dir=$HOME/local --tps=1 --interval=5",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err                   error
				tps                   int
				interval              int
			)
			helper.ReadConfigFile(FlagConfDir)

			if tps = viper.GetInt(FlagTps);tps <= 0 {
				return fmt.Errorf("tps should > 0 ")
			}
			if interval = viper.GetInt(FlagInterval);interval < 5 {
				return fmt.Errorf("duration should >= 5 ")
			}

			TXsForOneBlock := tps * interval  //发完后还没有到达 commitInterval， 则等待

			conf.Output = helper.GetPath(conf.Output+"/SignedTXs")
			file, err := os.OpenFile(conf.Output, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("can't find directory in %v\n", conf.Output)
			}
			defer file.Close()
			sc := bufio.NewScanner(file)

			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Start broadcasting TXs!!")

			count := 0
			timeTemp := time.Now()
			for sc.Scan() { //sc.Scan()默认以 \n 分隔
				count++
				_, err = sign.BroadcastTx(sc.Text())
				if err != nil  {
					//返回500
					for {
						time.Sleep(time.Millisecond * 500)
						_, err = sign.BroadcastTx(sc.Text())
						if err == nil {
							break
						}
					}
				}

				if count % TXsForOneBlock == 0 {
					if timeTemp.Add(time.Second * time.Duration(interval)).After(time.Now()) {
						time.Sleep(timeTemp.Add(time.Second * time.Duration(interval)).Sub(time.Now()))
					}
					log.Printf("Broadcast %d TXs\n", count)
					timeTemp = time.Now()
				}
			}

			if err := sc.Err(); err != nil{
				fmt.Println("An error has happened, when we run buf scanner")
				return err
			}


			log.Printf("End broadcasting TXs\n", )
			return nil
		},
	}

	cmd.Flags().AddFlagSet(broadcastTXFlagSet)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTps)
	cmd.MarkFlagRequired(FlagInterval)

	return cmd
}

