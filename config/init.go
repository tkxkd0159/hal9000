package config

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
)

var Nviper *viper.Viper
var Sviper *viper.Viper

func init() {
	Sviper = SetEnv()
	Nviper = SetNetInfo()
}

func SetEnv() *viper.Viper {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	utils.HandleErr(err, "Can't read .env", types.EXIT)

	sViper := viper.New()
	sViper.SetConfigType("yaml")
	sViper.SetConfigFile(".secret.yml")
	err = sViper.ReadInConfig()
	utils.HandleErr(err, "Can't read .secret.yml", types.EXIT)

	return sViper
}

func SetNetInfo() *viper.Viper {
	netViper := viper.New()
	netViper.SetConfigType("yaml")
	netViper.SetConfigFile(".netinfo.yml")
	err := netViper.ReadInConfig()
	utils.HandleErr(err, "Can't read .netinfo.yml", types.EXIT)

	return netViper
}
