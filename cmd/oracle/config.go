package main

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/client/common/types"
	"github.com/Carina-labs/HAL9000/utils"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

func SetBechPrefix() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.Seal()
}

func SetInitialDir(krDir string, logDir string) (string, string) {
	ckrDir, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)
	krDir = path.Join(ckrDir, krDir)
	err = os.Mkdir(krDir, 0740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	logDir = path.Join(ckrDir, logDir)
	err = os.Mkdir(logDir, 0740)
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
