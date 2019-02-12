package tx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
)

/////////////////////////////////////////
func SendTx(req types.TransferTxReq, dstAddress string) ([]byte, error) {
	uri := fmt.Sprintf(constants.UriTransfer, dstAddress)+"?commit=true"

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)
	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	if err != nil {
		return nil, err
	}

	//fmt.Println(string(resBytes))

	if statusCode == constants.StatusCodeOk {
		fmt.Printf("Send %s to %s ok! \n",req.Amount,dstAddress)
		return resBytes, nil
	} else {
		return nil, fmt.Errorf(string(resBytes))
	}
}
