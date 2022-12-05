package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
)

var Sviper *viper.Viper

const (
	ScrtFileName      = ".secret"
	ChainInfoFileName = ".chaininfo"
)

func setDefaultCfgPath(v ...*viper.Viper) {
	home := os.Getenv("HOME")
	pl := [3]string{filepath.Join(home, "config"), "/workspace/config", "."}
	for _, p := range pl {
		if v == nil {
			viper.AddConfigPath(p)
		} else {
			v[0].AddConfigPath(p)
		}
	}
}

func SetScrt(isTest bool) {
	if isTest {
		return
	}
	Sviper = viper.New()
	Sviper.SetConfigName(ScrtFileName)
	Sviper.SetConfigType("yaml")
	setDefaultCfgPath(Sviper)
	err := Sviper.ReadInConfig()
	utils.CheckErr(err, fmt.Sprintf("Can't read %s.yaml", ScrtFileName), types.EXIT)
}

func LoadChainInfo(isTest bool) {
	setDefaultCfgPath()
	viper.SetConfigType("yaml")
	if isTest {
		return
	}
	viper.SetConfigName(ChainInfoFileName)
	err := viper.ReadInConfig()
	utils.CheckErr(err, fmt.Sprintf("Can't read %s.yaml", ChainInfoFileName), types.EXIT)
}
