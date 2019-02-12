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