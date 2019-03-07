package tx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-load/types"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/irisnet/irishub-load/util/helper"
)

/////////////////////////////////////////
func SendTx(req types.TransferTxReq, dstAddress string) (types.TransferTxRes, error) {
	var (
		transferTxInfo types.TransferTxRes
	)

	uri := fmt.Sprintf(constants.UriTransfer, dstAddress)+"?commit=true"

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return transferTxInfo, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)
	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	if err != nil {
		return transferTxInfo, err
	}

	//fmt.Println(string(resBytes))
	if statusCode == constants.StatusCodeOk {
		fmt.Printf("Send %s to %s ok! \n",req.Amount,dstAddress)
		if err := json.Unmarshal(resBytes, &transferTxInfo); err != nil {
			return transferTxInfo, err
		}

		return transferTxInfo, nil
	} else {
		return transferTxInfo, fmt.Errorf(string(resBytes))
	}
}

func ChechTx(sender string, recipient string)(bool,error) {
	uri := fmt.Sprintf(constants.UriTxs, sender, recipient)
	statusCode, resByte, err := helper.HttpClientGetData(uri)

	if err != nil {
		return false, err
	}

	if statusCode == constants.StatusCodeOk {
		return string(resByte) != "[]", nil
	} else {
		return false, fmt.Errorf("status code is not ok, code: %v", statusCode)
	}
}