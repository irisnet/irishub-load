package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/spf13/viper"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-load/types"
)

func CheckFileExist(filePath string) (bool, error) {
	exists := true
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			// unknown err
			return false, err
		}
	}
	return exists, nil
}

func CreateFolder(folderPath string) error {
	folderExist, err := CheckFileExist(folderPath)
	if err != nil {
		return err
	}

	if !folderExist {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteFile(filePath string, content []byte) error {
	os.Remove(filePath)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	fileWrite := bufio.NewWriter(file)
	_, err = fileWrite.Write(content)
	if err != nil {
		return err
	}
	fileWrite.Flush()
	return nil
}

/////////////////////////////////////

func IntToStr(amt int) string {
	return strconv.Itoa(amt)
}

func StrToInt(amt string) (int, error) {
	return strconv.Atoi(amt)
}

func RandomId() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return strconv.Itoa(r.Intn(899999) + 100000)
}

func PraseUser(name string) int {
	switch name {
	case "user0":
		return 0
	case "user1":
		return 1
	case "user2":
		return 2
	case "user3":
		return 3
	case "user4":
		return 4
	default:
		return -1
	}
}

func ReadConfigFile(dir string) error {
	confDir := viper.GetString(dir)
	viper.SetConfigName("config_csrb") // config.json
	viper.AddConfigPath(confDir)  // $HOME
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	viper.UnmarshalKey("Node", &conf.NodeUrl)
	viper.UnmarshalKey("Output", &conf.Output)
	viper.UnmarshalKey("ChainId", &conf.ChainId)
	viper.UnmarshalKey("MinBalance", &conf.MinBalance)
	viper.UnmarshalKey("FaucetAddr", &conf.FaucetAddr)
	viper.UnmarshalKey("FaucetSeed", &conf.FaucetSeed)
	viper.UnmarshalKey("SubFaucets", &conf.SubFaucets)

	viper.UnmarshalKey("AirDropChainId", &conf.AirDropChainId)
	viper.UnmarshalKey("AirDropAddr", &conf.AirDropAddr)
	viper.UnmarshalKey("AirDropSeed", &conf.AirDropSeed)
	viper.UnmarshalKey("AirDropAmount", &conf.AirDropAmount)
	viper.UnmarshalKey("AirDropXlsx", &conf.AirDropXlsx)

	return nil
}

/////////////////////////////////
func ReadAddressList(dir string) ([]types.AirDropInfo, *excelize.File, error) {
	//fmt.Println("ReadAddressList() !!!!")

	var (
		airdrop_list []types.AirDropInfo
		airdrop_info types.AirDropInfo
	)

	xlsx, err := excelize.OpenFile(conf.AirDropXlsx)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	rows := xlsx.GetRows("Sheet1")
	for i, row := range rows {
		for j, colCell := range row {
			//读取金额 ，第1列， 跳过首行
			if j == 0 && i >= 1 {
				airdrop_info.Amount = colCell
			}

			//读取地址 ，第2列， 跳过首行
			if j == 1 && i >= 1 {
				airdrop_info.Address = colCell
				airdrop_info.Pos = i + 1
				airdrop_list = append(airdrop_list, airdrop_info)
			}
		}
	}

	return airdrop_list, xlsx, nil
}

func WriteAddressList(xlsx *excelize.File, airDropinfo types.AirDropInfo) {
	index := IntToStr(airDropinfo.Pos)
	xlsx.SetCellValue("Sheet1", "C"+index, airDropinfo.Hash)
	xlsx.SetCellValue("Sheet1", "D"+index, airDropinfo.Status)
}

//判断单元格是否为空，现在不用了。
func IsCellEmpty(xlsx *excelize.File, airDropinfo types.AirDropInfo) bool {
	index := IntToStr(airDropinfo.Pos)
	result := xlsx.GetCellValue("Sheet1", "G"+index)
	return result == ""
}

func SaveAddressList(xlsx *excelize.File, file string) error {
	err := xlsx.SaveAs(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/////////////////////////////////

func HttpClientPostJsonData(uri string, requestBody *bytes.Buffer) (int, []byte, error) {
	url := conf.NodeUrl + uri
	res, err := http.Post(url, constants.HeaderContentTypeJson, requestBody)

	if err != nil {
		return 0, nil, err
	}

	if res == nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resByte, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, resByte, nil
}

func HttpClientGetData(uri string) (int, []byte, error) {
	res, err := http.Get(conf.NodeUrl + uri)
	defer res.Body.Close()

	if err != nil {
		return 0, nil, err
	}
	if res == nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, resByte, nil
}

func GetPath(in string) string {
	if strings.HasPrefix(in, "$HOME") {
		in = UserHomeDir() + in[5:]
	}

	return in
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

/////////////////////////////////

func IrisattoToIris(coins []types.Coin) string {
	coin := types.Coin{"0", "0"}
	for _, subCoin := range coins {
		if subCoin.Denom == "iris-atto" {
			coin = subCoin
			break
		}
	}

	m := big.NewInt(math.MaxInt64)
	n, _ := new(big.Int).SetString(coin.Amount, 10)
	decimal := big.NewInt(1000000000000000000)
	m.Div(n, decimal)

	return m.String()
}

//支持小数，如1.23iris
func IrisToIrisatto(amount string) sdk.Coins {
	amtStr := strings.Replace(amount, constants.Denom, "", -1)
	dec, _ := sdk.NewDecFromStr(amtStr)
	decimal := sdk.NewDec(1000000000000000000)
	n := dec.Mul(decimal).TruncateInt() //做乘法，然后去掉小数后面取整数 ，fmt.Printf(" : %s \n", n.String())
	return sdk.Coins{{Denom: "iris-atto", Amount: n}}
}
