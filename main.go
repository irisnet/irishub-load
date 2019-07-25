package main

import (
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub-load/cmd"
	"github.com/spf13/viper"
)

func main() {
	//创建一个cmd 对象
	rootCmd := NewRootCmd()

	//添加回调函数
	rootCmd.AddCommand(

		//压力测试
		// 1）创建账户   irishub-load init --config-dir=$HOME
		//
		// 2）签名交易   irishub-load signtx --config-dir=$HOME --tps=100 --duration=1 --account=user0
		//              其中：总交易量=tps*duration*60
		//
		// 3）广播交易   irishub-load broadcast --config-dir=$HOME/local --tps=50
		//              其中： 每秒广播tps个交易， 如果在1s时间段完成， 则sleep直到1s时间结束，
		//                    再开始新一轮广播， 循环直至所有的交易广播完毕。

		cmd.FaucetInit(),
		cmd.SignTx(),
		cmd.BroadcastTx(),

		//其他feature，和压测无关
		cmd.AirDrop(),
		cmd.SeedTest(),
		cmd.RandTest(),
	)

	//PreRun bind viper
	executor := prepareMainCmd(rootCmd)

	//rootCmd.Execute 会先执行prerun
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "irishub-load",
		Short: "",
	}

	return rootCmd
}

func prepareMainCmd(cmd *cobra.Command) *cobra.Command {
	//绑定preRun
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return bindFlagsLoadViper(cmd)
	}

	return cmd
}

func bindFlagsLoadViper(rootCmd *cobra.Command) error {
	// viper bind flag
	// 后续viper.GetString(FlagConfDir)等处要绑定flag后才可以使用
	viper.BindPFlags(rootCmd.Flags())

	//bind各命令的flag , rootCmd.Commands() = 空？
	//for _, c := range rootCmd.Commands() {
	//	viper.BindPFlags(c.Flags())
	//}

	return nil
}



