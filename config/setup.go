package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/tkxkd0159/HAL9000/utils"
)

func SetInitialDir(keyname string, logSubdir string) (string, string) {
	cwd, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)

	krDir := path.Join(cwd, "/keyring", keyname)
	err = os.MkdirAll(krDir, 0o740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	logDir := path.Join(cwd, logSubdir)
	err = os.MkdirAll(logDir, 0o740)
	if os.IsExist(err) {
		log.Println("** log directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	return krDir, logDir
}

func SetAllLogger(logDir, stdLogName, errLogName, errRedirectLogName string, isDisp bool) (*os.File, *os.File, *os.File) {
	var fdLog, fdErr, fdErrExt *os.File
	var err error
	if !isDisp {
		fdLog, err = os.OpenFile(path.Join(logDir, stdLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		utils.CheckErr(err, "cannot open logfp", 0)

		// 반환되서 처리할 수 있는 에러 핸들링
		fdErr, err = os.OpenFile(path.Join(logDir, errLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		utils.CheckErr(err, "cannot open novaerr", 0)

		// 외부 라이브러리에서 fmt.Fprintf(os.stderr)로 처리하는 애들 핸들링 (redirect stderr)
		fdErrExt, err = os.OpenFile(path.Join(logDir, errRedirectLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
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
	case ActOracle, ActStake, ActAutoStake, ActWithdraw, ActAutoClaim:
		return botType
	default:
		fmt.Printf(" 🤮 This bot type is not supported. \n\n")
		fmt.Println("Command:")
		fmt.Printf("  hal [action] [flags]\n\n")
		fmt.Println(" [action] : oracle / stake / restake / withdraw / autoclaim")
		fmt.Println(" Use (-h|--help) if you want to see flag usage after set action")
		os.Exit(1)
	}
	return ""
}

func CheckTesterType(botType string) string {
	switch botType {
	case ActOracle, ActStake, ActAutoStake, ActWithdraw:
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
