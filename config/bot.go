package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"

	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/utils"
)

const (
	ActOracle   = "oracle"
	ActStake    = "stake"
	ActRestake  = "restake"
	ActWithdraw = "withdraw"
)

const (
	ControlChain = "nova"
	Gas          = "auto"
	NovaGasPrice = "0unova"
	GasWeight    = 1.1
)

const (
	StdLogFile            = "ctxlog.txt"
	LocalErrlogFile       = "nova_err.txt"
	ExtRedirectErrlogFile = "other_err.txt"
)

func SetInitialDir(krDir string, logDir string) (string, string) {
	cwd, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)

	krDir = path.Join(cwd, "/keyring", krDir)
	err = os.MkdirAll(krDir, 0740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	logDir = path.Join(cwd, logDir)
	err = os.MkdirAll(logDir, 0740)
	if os.IsExist(err) {
		log.Println("** log directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	return krDir, logDir
}

func GetPassphrase(vp *viper.Viper) string {
	pw := vp.GetString("pw")
	pp := fmt.Sprintf("%s\n%s\n", pw, pw)
	return pp
}

func SetAllLogger(logDir, stdLogName, errLogName, errRedirectLogName string, isDisp bool) (*os.File, *os.File, *os.File) {
	var fdLog, fdErr, fdErrExt *os.File
	var err error
	if !isDisp {
		fdLog, err = os.OpenFile(path.Join(logDir, stdLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open logfp", 0)

		// 반환되서 처리할 수 있는 에러 핸들링
		fdErr, err = os.OpenFile(path.Join(logDir, errLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open novaerr", 0)

		// 외부 라이브러리에서 fmt.Fprintf(os.stderr)로 처리하는 애들 핸들링
		fdErrExt, err = os.OpenFile(path.Join(logDir, errRedirectLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open otherErr", 0)

		os.Stderr = fdErrExt
	} else {
		fdLog = os.Stdout
		fdErr = os.Stderr
		fdErrExt = os.Stderr
	}

	return fdLog, fdErr, fdErrExt
}

func InputMnemonic() (mnemonic string) {
	fmt.Println(">>>>>>>>>>>>>> Enter mnemonic (24 words) <<<<<<<<<<<<<<")
	s := bufio.NewScanner(os.Stdin)
	ok := s.Scan()
	if !ok {
		log.Fatalln(" * Unexpected error while setup key")
	}
	mnemonic = s.Text()
	return
}

func CheckBotType(botType string) string {
	switch botType {
	case ActOracle, ActStake, ActRestake, ActWithdraw:
		return botType
	default:
		fmt.Printf(" 🤮 This bot type is not supported. \n\n")
		fmt.Println("Command:")
		fmt.Printf("  hal [action] [flags]\n\n")
		fmt.Println(" [action] : oracle / stake / restake / withdraw")
		fmt.Println(" Use (-h|--help) if you want to see flag usage after set action")
		os.Exit(1)
	}
	return ""
}

func SetupBotBase(f BotCommon, krDir string, ctxOut io.Writer) (ctx client.Context, botInfo keyring.Info, txf tx.Factory) {
	flags := f.GetBase()
	base.SetBechPrefix()
	LoadChainInfo(flags.IsTest)
	NovaInfo := NewChainNetInfo(ControlChain)
	BotScrt := NewBotScrt(NovaInfo.ChainID, "bot_addr", flags.Kn)

	if flags.New {
		SetupBotKey(flags.Kn, krDir, NovaInfo, BotScrt)
		log.Println("🎉 Your keyring has been successfully set.")
		os.Exit(0)
	}

	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)
	os.Stdin = rpipe
	_, err = wpipe.Write([]byte(BotScrt.Passphrase()))
	utils.CheckErr(err, "", 0)

	ctx = base.MakeContext(
		novaapp.ModuleBasics,
		BotScrt.Address(),
		NovaInfo.TmRPC.String(),
		NovaInfo.ChainID,
		krDir,
		keyring.BackendFile,
		rpipe,
		ctxOut,
		false,
	)

	botInfo = base.LoadClientPubInfo(ctx, flags.Kn)
	ctx = base.AddMoreFromInfo(ctx)
	txf = base.MakeTxFactory(ctx, Gas, NovaGasPrice, "", GasWeight)
	return
}

func SetupBotKey(keyname, keyloc string, info *ChainNetInfo, bot BotScrt) {
	ctx := base.MakeContext(
		novaapp.ModuleBasics,
		bot.Address(),
		info.TmRPC.String(),
		info.ChainID,
		keyloc,
		keyring.BackendFile,
		os.Stdin,
		os.Stdout,
		false,
	)

	_ = base.MakeClientWithNewAcc(
		ctx,
		keyname,
		InputMnemonic(),
		sdktypes.FullFundraiserPath,
		hd.Secp256k1,
	)
}
