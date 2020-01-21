package cmd

import (
	"bufio"
	"fmt"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/sign"
	"github.com/irisnet/irishub-load/types"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

func SignTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "signtx",
		Example: "irishub-load signtx --config-dir=$HOME/local --tps=1 --duration=1 --account=user0",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err            error
				subFaucetIndex int
				tps            int
				duration       int

				PrivateInfo      types.AccountTestPrivateInfo
				SignedDataString []string
			)

			// 从config.json中读取参数
			helper.ReadConfigFile(FlagConfDir)

			//先读指定账户私钥， 再生成签名数据(交易总数=tps*duration*60)
			if tps = viper.GetInt(FlagTps); tps <= 0 {
				return fmt.Errorf("tps should > 0 ")
			}
			if duration = viper.GetInt(FlagDuration); duration <= 0 {
				return fmt.Errorf("duration should > 0 ")
			}
			//读取发送交易子账户的序号
			if subFaucetIndex = helper.PraseUser(viper.GetString(FlagAccount)); subFaucetIndex == -1 {
				return fmt.Errorf("account %s not exist", viper.GetString(FlagAccount))
			}

			//总交易数
			totalTxNum := tps * duration * 60

			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Start signing %d TXs, chain id: %s, Node : %s, sub_faucet : %s \n", totalTxNum, conf.ChainId, conf.NodeUrl, conf.SubFaucets[subFaucetIndex].FaucetAddr)

			//用助记词恢复账户，并读取私钥，sequence id，account number
			if PrivateInfo, err = sign.InitAccountSignProcess(conf.SubFaucets[subFaucetIndex].FaucetAddr, conf.SubFaucets[subFaucetIndex].Seed); err != nil {
				return fmt.Errorf("Get private info error : %s", err.Error())
			}

			//生成所有签名后的交易
			if SignedDataString, err = sign.GenSignTxByTend(totalTxNum, subFaucetIndex, conf.ChainId, conf.SubFaucets, PrivateInfo); err != nil {
				return fmt.Errorf("GenSignTx error : %s", err.Error())
			}

			log.Printf("Finish signing TXs \n")

			//把签名后数据写入文件，以便后续用broadcast调用
			conf.Output = helper.GetPath(conf.Output)
			filePath := fmt.Sprintf("%v/SignedTXs", conf.Output)
			if err = helper.CreateFolder(conf.Output); err != nil {
				return err
			}
			if err = helper.WriteFile(filePath, []byte(strings.Join(SignedDataString, "\n"))); err != nil {
				return err
			}

			log.Printf("Finish WriteFile into %s \n", filePath)

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
		Use:     "broadcast",
		Example: "irishub-load broadcast --config-dir=$HOME/local --tps=1",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				err      error
				tps      int
				interval int
			)
			// 从config.json中读取参数
			helper.ReadConfigFile(FlagConfDir)

			//每轮广播tps*interval个交易， 如果在interval时间段完成， 则sleep直到interval时间结束，
			//再开始新一轮广播， 循环直至所有的交易广播完毕。
			if tps = viper.GetInt(FlagTps); tps <= 0 {
				return fmt.Errorf("tps should > 0 ")
			}

			interval = 1

			//一个块所发送的交易总数。发完后，如果时间还没有到，则等待。
			TXsForOneBlock := tps * interval

			//读取签名后的交易文件
			conf.Output = helper.GetPath(conf.Output + "/SignedTXs")
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

			//逐条广播交易，sc.Scan()默认以 \n 分隔
			for sc.Scan() {
				count++
				_, err = sign.BroadcastTx(sc.Text(), "async=true")

				//如果遇到网络拥堵（lcd返回500）
				//则每隔半秒检查一次，直至网络恢复
				if err != nil {
					for {
						time.Sleep(time.Millisecond * 500)
						_, err = sign.BroadcastTx(sc.Text(), "async=true")
						if err == nil {
							break
						}
					}
				}

				//每隔TXsForOneBlock条检查一下，如果时间未用完则等待，并且打印当前广播进度。
				if count%TXsForOneBlock == 0 {
					if timeTemp.Add(time.Millisecond * time.Duration(interval*1000)).After(time.Now()) {
						time.Sleep(timeTemp.Add(time.Millisecond * time.Duration(interval*1000)).Sub(time.Now()))
					}
					log.Printf("Broadcast %d TXs\n", count)
					timeTemp = time.Now()
				}
			}

			if err := sc.Err(); err != nil {
				fmt.Println("An error has happened, when we run buf scanner")
				return err
			}

			log.Printf("End broadcasting TXs\n")
			return nil
		},
	}

	cmd.Flags().AddFlagSet(broadcastTXFlagSet)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTps)

	return cmd
}
