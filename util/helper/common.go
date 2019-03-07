package helper

import (
	"bufio"
	"bytes"
	"github.com/irisnet/irishub-load/conf"
	"github.com/irisnet/irishub-load/util/constants"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"runtime"
	"strings"
	"math/rand"
	"time"
	"github.com/spf13/viper"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
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

func RandomId() string{
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return strconv.Itoa(r.Intn(89999)+10000)
}

func PraseUser(name string) int {
	switch name {
	case "wenxi":
		return 0
	case "silei":
		return 1
	case "haoyang":
		return 2
	case "jiacheng":
		return 3
	default:
		return -1
	}
}

func ReadConfigFile(dir string) error{
	confDir := viper.GetString(dir)
	viper.SetConfigName("config")
	viper.AddConfigPath(confDir)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.UnmarshalKey("Node", &conf.NodeUrl)
	viper.UnmarshalKey("Output", &conf.Output)
	viper.UnmarshalKey("ChainId", &conf.ChainId)
	viper.UnmarshalKey("MinBalance", &conf.MinBalance)
	viper.UnmarshalKey("FaucetSeed", &conf.FaucetSeed)
	viper.UnmarshalKey("SubFaucets", &conf.SubFaucets)

	viper.UnmarshalKey("AirDropSeed", &conf.AirDropSeed)
	viper.UnmarshalKey("AirDropGas", &conf.AirDropGas)
	viper.UnmarshalKey("AirDropFee", &conf.AirDropFee)
	viper.UnmarshalKey("AirDropAmount", &conf.AirDropAmount)
	viper.UnmarshalKey("AirDropXlsx", &conf.AirDropXlsx)

	return nil
}

/////////////////////////////////

func ReadAddressList(dir string) ([]types.AirDropInfo, *excelize.File, error){
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
			if j == 1 && i>=1 {
				airdrop_info.Address = colCell
				airdrop_info.Pos     = i+1
				airdrop_list = append(airdrop_list, airdrop_info)
			}
		}
	}

	if err = SaveAddressList(xlsx); err != nil {
		fmt.Println("SaveAddressList error !!!!")
		return airdrop_list, xlsx, err
	}

	return airdrop_list,xlsx, nil
}

func WriteAddressList(xlsx *excelize.File, airDropinfo types.AirDropInfo) {
	index := IntToStr(airDropinfo.Pos)
	xlsx.SetCellValue("Sheet1", "G"+index, airDropinfo.Status)
	xlsx.SetCellValue("Sheet1", "H"+index, airDropinfo.Hash)
	xlsx.SetCellValue("Sheet1", "I"+index, airDropinfo.TransactionTime)
	xlsx.SetCellValue("Sheet1", "J"+index, airDropinfo.Amount)
}

func IsCellEmpty(xlsx *excelize.File, airDropinfo types.AirDropInfo) bool {
	index := IntToStr(airDropinfo.Pos)
	return xlsx.GetCellValue("Sheet1", "G"+index) == ""
}

func SaveAddressList(xlsx *excelize.File) error{
	err := xlsx.SaveAs(conf.AirDropXlsx)
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

func GetPath(in string) string{
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