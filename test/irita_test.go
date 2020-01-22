package test

import (
	"encoding/json"
	"fmt"
	"github.com/bartekn/go-bip39"
	"github.com/bianjieai/irita/app"
	"github.com/bianjieai/irita/cli_test"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub-load/util/helper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

var (
	startCoins = sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(10000)), sdk.NewCoin("stake", sdk.NewInt(10000000000)))

	BIP44Prefix        = "44'/118'/"
	FullFundraiserPath = BIP44Prefix + "0'/0/0"

	number =300
	cdc = app.MakeCodec()
)

func Test_tx(t *testing.T) {
	f := NewFixtures(t)
	to := f.KeyAddress(fmt.Sprintf("key_%d", 0))
	var signTxs string
	var privateList [][32]byte
	privPath := filepath.Join(f.IritaCLIHome, "priv.keys")
	data, err := ioutil.ReadFile(privPath)
	json.Unmarshal(data,&privateList)


	for i := 0; i < len(privateList); i++ {
		priv := secp256k1.PrivKeySecp256k1(privateList[i])
		address := sdk.AccAddress(priv.PubKey().Address())
		from := address.String()
		_, tx, _ := f.TxSend(from, to, sdk.NewCoin("iris", sdk.NewInt(1)), "--generate-only")
		var stdTx auth.StdTx
		cdc.MustUnmarshalJSON([]byte(tx),&stdTx)

		account := f.QueryAccount(address)

		stdSignMsg := auth.StdSignMsg{
			AccountNumber:account.AccountNumber,
			Sequence:account.Sequence,
			ChainID:f.ChainID,
			Fee:stdTx.Fee,
			Msgs:stdTx.GetMsgs(),
			Memo:stdTx.GetMemo(),
		}

		sigBytes,err := priv.Sign(stdSignMsg.Bytes())
		require.NoError(t,err)
		stdTx.Signatures = append(stdTx.Signatures, auth.StdSignature{
			PubKey:priv.PubKey(),
			Signature:sigBytes,
		})

		signedTx,_ := cdc.MarshalJSON(stdTx)
		signedStr := string(signedTx)
		signedStr = strings.Replace(signedStr, `"type":"cosmos-sdk/StdTx","value"`, `"tx"`, -1)
		signedStr = signedStr[:len(signedStr)-1] + `,"mode":"async"}` //立即返回，block / sync / async
		signTxs = signTxs + signedStr + "\n"
		require.NoError(t, err)
	}
	signPath := filepath.Join(f.IritaCLIHome, "signTxs.json")
	err = helper.WriteFile(signPath, []byte(signTxs))
	require.NoError(t, err)

}

func Test_init(t *testing.T) {
	t.Parallel()
	_ = InitFixtures(t,number)
}

func InitFixtures(t *testing.T,n int) (f *clitest.Fixtures) {
	f = NewFixtures(t)

	// reset test state
	f.UnsafeResetAll()

	f.KeysAdd(fmt.Sprintf("key_%d", 0))

	// ensure that CLI output is in JSON format
	f.CLIConfig("output", "json")

	// NOTE: GDInit sets the ChainID
	f.GDInit(fmt.Sprintf("key_%d", 0),"--chain-id=test")
	f.AddGenesisAccount(f.KeyAddress(fmt.Sprintf("key_%d", 0)), startCoins)


	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("broadcast-mode", "block")
	f.CLIConfig("trust-node", "true")

	var privateList [][32]byte-

	// start an account with tokens
	for i := 0; i < n; i++ {
		entropySeed, err := bip39.NewEntropy(256)
		require.NoError(t,err)
		mnemonic, err := bip39.NewMnemonic(entropySeed)
		require.NoError(t,err)
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
		require.NoError(t,err)

		masterPriv, ch := hd.ComputeMastersFromSeed(seed)
		derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, FullFundraiserPath)
		require.NoError(t,err)

		//获取公钥
		pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()
		f.AddGenesisAccount(sdk.AccAddress(pubk.Address()), startCoins)
		privateList =append(privateList,derivedPriv)

		fmt.Printf("generated number : %d\n",i+1)
	}

	privsStr,err := json.Marshal(privateList)
	require.NoError(t,err)
	privPath := filepath.Join(f.IritaCLIHome, "priv.keys")
	err = helper.WriteFile(privPath, []byte(privsStr))
	require.NoError(t, err)

	f.GenTx(fmt.Sprintf("key_%d", 0))
	f.CollectGenTxs()
	fmt.Print("generated success")
	return
}

// NewFixtures creates a new instance of Fixtures with many vars set
func NewFixtures(t *testing.T) *clitest.Fixtures {
	tmpDir := "irita_test"

	return &clitest.Fixtures{
		T:              t,
		ChainID:        "test",
		RootDir:        tmpDir,
		IritadBinary:   "irita",
		IritaCLIBinary: "iritacli",
		IritadHome:     filepath.Join(tmpDir, ".irita"),
		IritaCLIHome:   filepath.Join(tmpDir, ".iritacli"),
		RPCAddr:        "localhost:26657",
		P2PAddr:        "localhost:26656",
		Port:           "26657",
	}
}
