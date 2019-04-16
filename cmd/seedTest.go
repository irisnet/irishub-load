package cmd

import (
//	"github.com/irisnet/irishub-load/conf"
//	"github.com/irisnet/irishub-load/sign"
//	"github.com/irisnet/irishub-load/util/helper"
//	"github.com/irisnet/irishub-load/types"
	"github.com/spf13/cobra"
//	"github.com/spf13/viper"
//	"log"
	//"strings"

	"github.com/irisnet/irishub-load/util/helper"
	"fmt"

	"github.com/tyler-smith/go-bip39"
	"github.com/irisnet/irishub/crypto/keys/hd"
	//"github.com/tendermint/tendermint/crypto/secp256k1"

	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
	"github.com/irisnet/irishub-load/conf"
	"os"
	"bufio"
	"encoding/json"
	"errors"
)

type SeedAccountInfo struct {
	PrivateKey    string
	PubKey        string
	Addr          string
}

type InputAccountInfo struct {
	PriveteKeyLen  int         `json:"priveteKeyLen"`
	Secret         string      `json:"secret"`
	Address        string      `json:"address"`
	PrivateKey     string      `json:"privateKey"`
	PublicKey      string      `json:"publicKey"`
}

func SeedTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seedtest",
		Example: "irishub-load seedtest --config-dir=$HOME/local",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("start \n")

			var (
				err                   error
				inputInfo             InputAccountInfo
				seedInfo			  SeedAccountInfo
			)

			helper.ReadConfigFile(FlagConfDir)


			file, err := os.OpenFile("D:/seedtest.txt", os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("can't find directory in %v\n", conf.Output)
			}
			defer file.Close()
			sc := bufio.NewScanner(file)


			fmt.Println("Start reading seed ....")

			count := 0
			error_count := 0
			for sc.Scan() {
				js := sc.Text()
				err = json.Unmarshal([]byte(js), &inputInfo)
				if err != nil {
					return fmt.Errorf("can't prase input json\n %s",err.Error())
				}

				//读取私钥信息
				if seedInfo, err = GetAccountInfoFromSeed(inputInfo.Secret); err!=nil {
					return fmt.Errorf("Get private info error : %s", err.Error())
				}

				if  err = Comparedata(inputInfo, seedInfo); err!=nil {
					//return fmt.Errorf("Result not equal : %s", err.Error())
					error_count++
				}

				count++
				fmt.Println("Compare " , count , " seeds ok!")
			}


			fmt.Println("total errors " , error_count)
			return nil
		},
	}

	cmd.Flags().AddFlagSet(seedtestFlagSet)
	cmd.MarkFlagRequired(FlagConfDir)

	return cmd
}

func Comparedata(inputInfo InputAccountInfo, seedInfo SeedAccountInfo) error {
	if inputInfo.PriveteKeyLen == 31 {
		inputInfo.PrivateKey = "00"+inputInfo.PrivateKey
	}

	if inputInfo.PrivateKey != seedInfo.PrivateKey {
		fmt.Println(inputInfo.PrivateKey , " != " , seedInfo.PrivateKey)
		return errors.New("PrivateKey not equal")
	}

	if inputInfo.Address != seedInfo.Addr {
		return errors.New("Address not equal")
	}

	if inputInfo.PublicKey != seedInfo.PubKey {
		return errors.New("PrivateKey not equal")
	}

	return nil
}


func GetAccountInfoFromSeed(mnemonic string) (SeedAccountInfo, error) {
	var account SeedAccountInfo

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return account, err
	}
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)
	if err != nil {
		return account, err
	}
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()

	account.PrivateKey = strings.ToUpper(hex.EncodeToString(derivedPriv[:]))
	account.PubKey = helper.ConvertFromHex("fap", hex.EncodeToString(pubk.Bytes()))
	account.Addr = helper.ConvertFromHex("faa", pubk.Address().String()) //pubk.Address().String()

	return account, err
}

