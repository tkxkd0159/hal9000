package cmd

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

func SetInitialDir(krDir string, logDir string) (string, string) {
	ckrDir, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)

	krDir = path.Join(ckrDir, krDir)
	err = os.MkdirAll(krDir, 0740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	logDir = path.Join(ckrDir, logDir)
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
