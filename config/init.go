package config

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	utils.HandleErr(err)
}

func SetNetInfo() *viper.Viper {
	netViper := viper.New()
	netViper.SetConfigType("yaml")
	netViper.SetConfigFile(".netinfo.yml")
	err := netViper.ReadInConfig()
	utils.HandleErr(err)

	return netViper
}
