package config

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
)

func SetEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	utils.HandleErr(err, "Can't read .env", types.EXIT)
}

func SetNetInfo() *viper.Viper {
	netViper := viper.New()
	netViper.SetConfigType("yaml")
	netViper.SetConfigFile(".netinfo.yml")
	err := netViper.ReadInConfig()
	utils.HandleErr(err, "Can't read .netinfo.yml", types.EXIT)

	return netViper
}
