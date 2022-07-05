package config

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
)

var Sviper *viper.Viper

func init() {
	Sviper = setEnv()
}

func setEnv() *viper.Viper {
	//viper.SetConfigFile(".env")
	//err := viper.ReadInConfig()
	//utils.CheckErr(err, "Can't read .env", types.EXIT)

	sViper := viper.New()
	sViper.SetConfigType("yaml")
	sViper.SetConfigFile(".secret.yml")
	err := sViper.ReadInConfig()
	utils.CheckErr(err, "Can't read .secret.yml", types.EXIT)

	return sViper
}

func SetChainInfo(isTest bool) {
	viper.SetConfigType("yaml")
	if isTest {
		viper.SetConfigFile(".chaininfo.test.yml")
		err := viper.ReadInConfig()
		utils.CheckErr(err, "Can't read .chaininfo.test.yml", 0)
	} else {
		viper.SetConfigFile(".chaininfo.yml")
		err := viper.ReadInConfig()
		utils.CheckErr(err, "Can't read .chaininfo.yml", types.EXIT)
	}
}
