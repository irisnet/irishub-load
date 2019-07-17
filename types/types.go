package types

import (
	"github.com/irisnet/irishub/types"
)

////////////////////////////////////////////////////
type PublicKey struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type Coin struct {
	Denom   string `json:"denom"`
	Amount  string `json:"amount"`
}

////////////////////////////////////////////////////

type AccountInfo struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	AccountNumber    string `json:"account_number"`
	Sequence         string `json:"sequence"`
	Address          string `json:"address"`
	AccountName      string `json:"account_name"`
	Seed             string `json:"seed"`
}

type AccountInfoRes struct {
	Type  string           `json:"type"`
	Value AccountInfoValue `json:"value"`
}

type AccountInfoValue struct {
	Address       string   `json:"address"`
	Coins         []Coin    `json:"coins"`
	PublicKey     PublicKey   `json:"public_key"`
	AccountNumber string   `json:"account_number"`
	Sequence      string   `json:"sequence"`
}

////////////////////////////////////////////////////

type AirDropInfo struct {
	Pos                   int
	Address               string
	Status                string
	Hash                  string
	TransactionTime       string
	Amount                string
}

////////////////////////////////////////////////////

type TransferTxRes struct {
	Check_tx   CheckTx `json:"check_tx"`
	Deliver_tx DeliverTx `json:"deliver_tx"`
	Hash       string `json:"hash"`
	Height     string `json:"height"`
}

type CheckTx struct {
	GasWanted string `json:"gasWanted"`
	GasUsed   string `json:"gasUsed"`
}

type DeliverTx struct {
	Log         string `json:"log"`
	GasWanted   string `json:"gasWanted"`
	GasUsed    string `json:"gasUsed"`
	Tags       []Tag `json:"tags"`
}

type Tag struct {
	Key         string `json:"key"`
	Value      string `json:"value"`
}


////////////////////////////////////////////////////

type TransferTxReq struct {
	Amount    string
	ChainID   string
	Sequence  int
	SenderAddr  string
	SenderSeed  string
	RecipientAddr  string
	Mode         string
}



type BaseTx struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	AccountNumber    string `json:"account_number"`
	Sequence         string `json:"sequence"`
	Gas              string `json:"gas"`
	Fees             string `json:"fee"`
	Memo             string `json:"memo"`
}

////////////////////////////////////////////////////

type ErrorRes struct {
	RestAPI      string `json:"rest api"`
	Code         int    `json:"code"`
	ErrorMessage string `json:"err message"`
}

type KeyCreateReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

type KeyCreateRes struct {
	Address string `json:"address"`
	Seed    string `json:"seed"`
}

////////////////////////////////////////////////////

type TxDataRes struct {
	Type  string `json:"type"`
	Value PostTx `json:"value"`
}

type PostTx struct {
	Msgs       []TxDataInfo   `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

type StdFee struct {
	Amount types.Coins `json:"amount"`
	Gas    string      `json:"gas"`
}

type StdSignature struct {
	PubKey        PubKey `json:"pub_key"` // optional
	Signature     string `json:"signature"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

type PubKey struct {
	Type  string `json:"type"` // optional
	Value string `json:"value"`
}

type TxDataInfo struct {
	Type  string      `json:"type"`
	Value TxDataValue `json:"value"`
}

type TxDataValue struct {
	Input  []InOutPutData `json:"inputs"`
	Output []InOutPutData `json:"outputs"`
}

type InOutPutData struct {
	Address string      `json:"address"`
	Amount  types.Coins `json:"coins"`
}

type TxBroadcast struct {
	Tx PostTx `json:"tx"`
}

////////////////////////////////////////////////////

type AccountTestPrivateInfo struct {
	PrivateKey    [32]byte
	PubKey        []byte
	Addr          string
	AccountNumber uint64
	Sequence      uint64
}


