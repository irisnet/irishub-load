module github.com/irisnet/irishub-load

go 1.12

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/cosmos-sdk v0.38.2
	github.com/grd/statistics v0.0.0-20130405091615-5af75da930c9
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tyler-smith/go-bip39 v1.0.2
	gitlab.bianjie.ai/csrb/csrb v1.0.0-rc0.0.20200610091436-3cb8edf6a576
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/bianjieai/cosmos-sdk v0.39.0-irita-200529
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.33.4-irita-200529
)
