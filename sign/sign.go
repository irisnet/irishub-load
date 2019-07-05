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
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()

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
	amount, ok := sdk.NewIntFromString(amtV)
	coins := sdk.Coins{{Denom: denom, Amount: amount}}

	input := bank.Input{Address: from, Coins: coins}

	feea, ok := sdk.NewIntFromString(feeAmtV)
	feeAmt := sdk.Coins{{Denom: denom, Amount: feea}}
	if !ok {
		return nil, errors.New("err in String to int")
	}
	var msgs []sdk.Msg

	fee := auth.StdFee{Amount: feeAmt, Gas: gas}

	var signedData []string
	priv := secp256k1.PrivKeySecp256k1(accountPrivate.PrivateKey)
	sigChan := make(chan auth.StdTx, 30)
	counter := 0
	sequence := accountPrivate.Sequence
	for i := 0; i < testNum; i++ {
		to, err := sdk.AccAddressFromBech32(subFaucets[counter].FaucetAddr)
		if err != nil {
			return nil, errors.New("err 2 in address to String")
		}
		output := bank.Output{Address: to, Coins: coins}
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
		go genSignedDataByTend(priv, sigMsg, accountPrivate, sigChan, msgs, fee)
		tx := <-sigChan
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
		signedData = append(signedData, string(postTxBytes))
		sequence = sequence + 1
		counter = counter + 1
		if counter == len(subFaucets) {
			counter = 0
		}
	}
	return signedData, nil
}

func genSignedDataByTend(priv secp256k1.PrivKeySecp256k1, sigMsg StdSignMsg, accountPrivate types.AccountTestPrivateInfo, sigChan chan auth.StdTx, msgs []sdk.Msg, fee auth.StdFee) {
	sigBz := sigMsg.Bytes()
	sigByte, err := priv.Sign(sigBz)
	if err != nil {
		return
	}
	sig := auth.StdSignature{
		PubKey:        priv.PubKey(),
		Signature:     sigByte,
		AccountNumber: accountPrivate.AccountNumber,
		Sequence:      sigMsg.Sequence,
	}

	sigs := []auth.StdSignature{sig}
	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	sigChan <- tx
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

