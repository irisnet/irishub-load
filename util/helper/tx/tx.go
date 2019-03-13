package tx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-load/types"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/irisnet/irishub-load/util/helper"
	"strings"
)

/////////////////////////////////////////

func SendTx(req types.TransferTxReq, dstAddress string, sync bool) (types.TransferTxRes, error) {
	var (
		transferTxInfo types.TransferTxRes
		op string
	)

	if sync {
		op= "?commit=true"
	} else {
		op= ""
	}

	uri := fmt.Sprintf(constants.UriTransfer, dstAddress)+op

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
		if err := json.Unmarshal(resBytes, &transferTxInfo); err != nil {
			return transferTxInfo, err
		}

		return transferTxInfo, nil
	} else {
		return transferTxInfo, fmt.Errorf(string(resBytes))
	}
}

func ChechTx(tx string) error {
	uri := fmt.Sprintf(constants.UriTxs, tx)
	statusCode, resByte, err := helper.HttpClientGetData(uri)

	if err != nil {
		return err
	}

	if statusCode == constants.StatusCodeOk {
		if strings.Contains(string(resByte), "hash"){
			return  nil
		}

		return  fmt.Errorf(string(resByte))
	} else {
		return fmt.Errorf("status code is not ok, code: %v", statusCode)
	}


}

//判断是否已经转账过  上链检查  支持一个全节点配置index_all_tags = true

//if duplicated, err := tx.ChechTx(req.Sender, sub.Address); err == nil {
//if duplicated {
//fmt.Println("Duplicated transfer : "+req.Sender+" to "+sub.Address)
//sub.Status = "Duplicated"
//sub.Hash = ""
//sub.TransactionTime = ""
//sub.Amount = ""
//helper.WriteAddressList(xlsx, sub)
//continue
//}
//}else {
//fmt.Println(err.Error())
//break
//}

//func ChechTx2(sender string, recipient string)(bool,error) {
//	uri := fmt.Sprintf(constants.UriTxs, sender, recipient)
//	statusCode, resByte, err := helper.HttpClientGetData(uri)
//
//	if err != nil {
//		return false, err
//	}
//
//	if statusCode == constants.StatusCodeOk {
//		return string(resByte) != "[]", nil
//	} else {
//		return false, fmt.Errorf("status code is not ok, code: %v", statusCode)
//	}
//}