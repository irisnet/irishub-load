package tx

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-load/sign"
	"github.com/irisnet/irishub-load/types"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/irisnet/irishub-load/util/helper"
	"strings"
)

/////////////////////////////////////////

func SendTx(req types.TransferTxReq) (types.TransferTxRes, error) {
	var (
		err              error
		transferTxInfo   types.TransferTxRes
		PrivateInfo      types.AccountTestPrivateInfo
		SignedDataString string

		response []byte
	)

	if PrivateInfo, err = sign.InitAccountSignProcess(req.SenderAddr, req.SenderSeed); err != nil {
		return transferTxInfo, fmt.Errorf("Get private info error : %s", err.Error())
	}

	if SignedDataString, err = sign.GenSingleSignTxByTend(req, PrivateInfo); err != nil {
		return transferTxInfo, fmt.Errorf("GenSignTx error : %s", err.Error())
	}

	//如果签名不通过 可以这里打印出来看下有没有问题 是不是chainid错了
	//fmt.Println(SignedDataString)

	if response, err = sign.BroadcastTx(SignedDataString, req.Mode); err != nil {
		fmt.Println(string(response))
		return transferTxInfo, fmt.Errorf("BroadcastTx error : %s", err.Error())
	}

	if err := json.Unmarshal(response, &transferTxInfo); err != nil {
		return transferTxInfo, err
	}

	return transferTxInfo, nil
}

func CheckTx(tx string) error {
	uri := fmt.Sprintf(constants.UriTxs, tx)
	statusCode, resByte, err := helper.HttpClientGetData(uri)

	if err != nil {
		return err
	}

	if statusCode == constants.StatusCodeOk {
		if strings.Contains(string(resByte), "hash") && !strings.Contains(string(resByte), "failed") {
			return nil
		}

		return fmt.Errorf(string(resByte))
	} else {
		return fmt.Errorf("status code is not ok, code: %v", statusCode)
	}

}
