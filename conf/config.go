package conf

var (
	NodeUrl string
	ChainId string
	Output  string

	FaucetAddr string
	FaucetSeed string
	MinBalance string
	SubFaucets []SubFaucet

	AirDropChainId string
	AirDropAddr    string
	AirDropSeed    string
	AirDropAmount  string
	AirDropXlsx    string
)

type SubFaucet struct {
	FaucetName     string `json:"faucet_name" mapstructure:"faucet_name"`
	FaucetPassword string `json:"faucet_password" mapstructure:"faucet_password"`
	FaucetAddr     string `json:"faucet_addr" mapstructure:"faucet_addr"`
	Seed           string `json:"seed" mapstructure:"seed"`
}
