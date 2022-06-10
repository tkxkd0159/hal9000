package client

import (
	"github.com/Carina-labs/HAL9000/config"
	"github.com/spf13/viper"
)

var NV *viper.Viper

func init() {
	NV = config.SetNetInfo()
}
