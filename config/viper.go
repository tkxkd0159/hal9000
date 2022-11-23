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
	ScrtFileName          = ".secret"
	ChainInfoFileName     = ".chaininfo"
	ChainInfoTestFileName = ChainInfoFileName + ".test"
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

func init() {
	Sviper = setScrt()
}

func setScrt() *viper.Viper {
	sViper := viper.New()
	sViper.SetConfigName(ScrtFileName)
	sViper.SetConfigType("yaml")
	setDefaultCfgPath(sViper)
	err := sViper.ReadInConfig()
	utils.CheckErr(err, fmt.Sprintf("Can't read %s.yaml", ScrtFileName), types.EXIT)

	return sViper
}

func LoadChainInfo(isTest bool) {
	setDefaultCfgPath()
	viper.SetConfigType("yaml")
	if isTest {
		viper.SetConfigName(ChainInfoTestFileName)
		err := viper.ReadInConfig()
		utils.CheckErr(err, fmt.Sprintf("Can't read %s.yaml", ChainInfoTestFileName), types.EXIT)
	} else {
		viper.SetConfigName(ChainInfoFileName)
		err := viper.ReadInConfig()
		utils.CheckErr(err, fmt.Sprintf("Can't read %s.yaml", ChainInfoFileName), types.EXIT)
	}
}
