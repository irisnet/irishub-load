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
		cmd.FaucetInit(),
		cmd.SignTx(),
		cmd.BroadcastTx(),
		cmd.AirDrop(),
		cmd.SeedTest(),
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



