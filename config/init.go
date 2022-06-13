package config

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
)

var Sviper *viper.Viper

func init() {
	setChainInfo()
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

func setChainInfo() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(".chaininfo.yml")
	err := viper.ReadInConfig()
	utils.CheckErr(err, "Can't read .chaininfo.yml", types.EXIT)
}
