package sign

import (
	"bytes"
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub-load/util/constants"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/irisnet/irishub-load/types"

	"log"
	"strings"
	"strconv"
	. "github.com/irisnet/irishub-load/conf"


	"github.com/tyler-smith/go-bip39"
	"github.com/irisnet/irishub/crypto/keys/hd"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto"
	"github.com/irisnet/irishub-load/util/helper/account"
	"encoding/json"
	"errors"
)

var (
	Cdc *codec.Codec
)

const (
	amtV    = "1000000000000000"
	feeAmtV = "5000000000000000000"
	denom   = "iris-atto"
	gas     = uint64(20000)
	memo    = ""
)

type StdSignMsg struct {
	ChainID       string      `json:"chain_id"`
	AccountNumber uint64      `json:"account_number"`
	Sequence      uint64      `json:"sequence"`
	Fee           auth.StdFee `json:"fee"`
	Msgs          []sdk.Msg   `json:"msgs"`
	Memo          string      `json:"memo"`
}

func (msg StdSignMsg) Bytes() []byte {
	return auth.StdSignBytes(msg.ChainID, msg.AccountNumber, msg.Sequence, msg.Fee, msg.Msgs, msg.Memo)
}

//custom tx codec
func init() {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	Cdc = cdc
}

func InitAccountSignProcess(fromAddr string, mnemonic string) (types.AccountTestPrivateInfo, error) {
	var Account types.AccountTestPrivateInfo

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return Account, err
	}
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)
	if err != nil {
		return Account, err
	}
	//获取公钥
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()

	//获取sequence id，accountNumber 写入返回的数据结构
	acc, err := account.GetAccountInfo(fromAddr)
	sequence, err := strconv.Atoi(acc.Value.Sequence)
	if err != nil {
		return Account, err
	}
	accountNumber, err := strconv.Atoi(acc.Value.AccountNumber)
	if err != nil {
		return Account, err
	}

	Account.PrivateKey = derivedPriv
	Account.PubKey = pubk.Bytes()
	Account.Addr = fromAddr
	Account.AccountNumber = uint64(accountNumber)
	Account.Sequence = uint64(sequence)
	return Account, err
}

// tx example:
/*{
  "tx": {
    "msg": [
      {
        "type": "irishub/bank/Send",
        "value": {
          "inputs": [
            {
              "address": "faa1lcuw6ewd2gfxap37sejewmta205sgssmv5fnju",
              "coins": [
                {
                  "denom": "iris-atto",
                  "amount": "1000000000000000"
                }
              ]
            }
          ],
          "outputs": [
            {
              "address": "faa1lcuw6ewd2gfxap37sejewmta205sgssmv5fnju",
              "coins": [
                {
                  "denom": "iris-atto",
                  "amount": "1000000000000000"
                }
              ]
            }
          ]
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "iris-atto",
          "amount": "5000000000000000000"
        }
      ],
      "gas": "20000"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AkQeeR40fJhMFkBmh8e/+jcOGYvMOO50YE4trqzaSJ5v"
        },
        "signature": "lCgkVUMEWRWzg36i8TDYqR2yJTyE/VV1CJlet//LEeweF2WcqkN9oEXxcPMonCSSRPWu30+dey8a07VIEmRppA==",
        "account_number": "3",
        "sequence": "3"
      }
    ],
    "memo": ""
  }
}

*/

func GenSignTxByTend(testNum int, fromIndex int, chainId string, subFaucets []SubFaucet, accountPrivate types.AccountTestPrivateInfo) ([]string, error) {

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)

	from, err := sdk.AccAddressFromBech32(subFaucets[fromIndex].FaucetAddr)

	if err != nil {
		return nil, errors.New("err in address to String")
	}
	//每个交易转账1000000000000000iris-atto = 0.001iris
	amount, ok := sdk.NewIntFromString(amtV)
	coins := sdk.Coins{{Denom: denom, Amount: amount}}

	//构造转账的结构
	input := bank.Input{Address: from, Coins: coins}

	//构造fee
	feea, ok := sdk.NewIntFromString(feeAmtV)
	feeAmt := sdk.Coins{{Denom: denom, Amount: feea}}
	if !ok {
		return nil, errors.New("err in String to int")
	}
	var msgs []sdk.Msg

	fee := auth.StdFee{Amount: feeAmt, Gas: gas}

	//把所有签名后的交易都写入这个数据结构
	var signedData []string
	priv := secp256k1.PrivKeySecp256k1(accountPrivate.PrivateKey)
	counter := 0
	sequence := accountPrivate.Sequence

	//一共构造、签名testNum条交易，每条交易sequence都要+1
	//分别转账到4个不同的账户， 用counter实现4个账户地址的循环
	for i := 0; i < testNum; i++ {
		to, err := sdk.AccAddressFromBech32(subFaucets[counter].FaucetAddr)
		if err != nil {
			return nil, errors.New("err 2 in address to String")
		}
		output := bank.Output{Address: to, Coins: coins}
		//构造"msg"
		msgs = []sdk.Msg{bank.MsgSend{
			Inputs:  []bank.Input{input},
			Outputs: []bank.Output{output},
		}}

		sigMsg := StdSignMsg{
			ChainID:       chainId,
			AccountNumber: accountPrivate.AccountNumber,
			Sequence:      sequence,
			Memo:          "",
			Msgs:          msgs,
			Fee:           fee,
		}

		//签名单条交易
		tx := genSignedDataByTend(priv, sigMsg, accountPrivate, msgs, fee)
		bz, _ := cdc.MarshalJSON(tx)
		var signedTx types.TxDataRes
		err = json.Unmarshal(bz, &signedTx)
		if err != nil {
			log.Printf("%v: sign tx failed: %v\n", "SignByTend", err)
			return nil, err
		}

		postTx := types.TxBroadcast{
			Tx: signedTx.Value,
		}
		postTx.Tx.Msgs[0].Type = "irishub/bank/Send"
		postTxBytes, err := json.Marshal(postTx)
		//log.Printf("%s\n", postTxBytes)
		if err != nil {
			log.Printf("%v: cdc marshal json fail: %v\n", "", err)
			return nil, err
		}

		//把创建的签名后的交易逐条写入
		signedData = append(signedData, string(postTxBytes))
		sequence = sequence + 1
		counter = counter + 1
		if counter == len(subFaucets) {
			counter = 0
		}
	}
	return signedData, nil
}

func genSignedDataByTend(priv secp256k1.PrivKeySecp256k1, sigMsg StdSignMsg, accountPrivate types.AccountTestPrivateInfo, msgs []sdk.Msg, fee auth.StdFee) (auth.StdTx) {
	sigBz := sigMsg.Bytes()
	sigByte, err := priv.Sign(sigBz)
	if err != nil {
		return auth.StdTx{}
	}
	sig := auth.StdSignature{
		PubKey:        priv.PubKey(),
		Signature:     sigByte,
		AccountNumber: accountPrivate.AccountNumber,
		Sequence:      sigMsg.Sequence,
	}

	sigs := []auth.StdSignature{sig}
	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}

func BroadcastTx(txBody string) ([]byte, error) {
	reqBytes := []byte(txBody)

	reqBuffer := bytes.NewBuffer(reqBytes)
	uri := constants.UriTxBroadcast + "?async=true"
	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	// handle response
	if err != nil {
		return nil, err
	}
	if statusCode != constants.StatusCodeOk {
		return nil, fmt.Errorf("unexcepted status code: %v", statusCode)
	}

	if strings.Contains(string(resBytes), "invalid") {
		log.Printf("%s\n", resBytes)
		return nil, fmt.Errorf("unexcepted information: %v", "invalid")
	}

	a := strings.Contains(string(resBytes), "check_tx") && strings.Contains(string(resBytes), "deliver_tx")
	b := strings.Contains(string(resBytes), "hash") && strings.Contains(string(resBytes), "height")
	if a && b {
		return resBytes, nil
	} else {
		log.Printf("check_tx check broadcast information failed\n")
		return nil, fmt.Errorf("check broadcast information failed")
	}
	return resBytes, nil
}

