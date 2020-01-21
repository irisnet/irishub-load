package cmd

import (
	"github.com/spf13/pflag"
)

const (
	FlagConfDir = "config-dir"

	FlagTps      = "tps"
	FlagDuration = "duration"
	FlagAccount  = "account"
)

var (
	faucetInitSet      = pflag.NewFlagSet("", pflag.ContinueOnError)
	signTXFlagSet      = pflag.NewFlagSet("", pflag.ContinueOnError)
	broadcastTXFlagSet = pflag.NewFlagSet("", pflag.ContinueOnError)
)

func init() {
	faucetInitSet.String(FlagConfDir, "", "directory for save config data")

	signTXFlagSet.String(FlagConfDir, "", "directory of config file")
	signTXFlagSet.Int(FlagTps, 0, "max tps per second")
	signTXFlagSet.Int(FlagDuration, 0, "time duration during the pressure test")
	signTXFlagSet.String(FlagAccount, "wenxi", "test account index in distribute systems")

	broadcastTXFlagSet.String(FlagConfDir, "", "directory of config file")
	broadcastTXFlagSet.Int(FlagTps, 0, "max tps per second")
}
