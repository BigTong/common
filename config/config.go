package config

import (
	"github.com/BigTong/common/log"
	"github.com/spf13/viper"
)

func InitConfigFromFile(fileName string, fileType string) {
	viper.WatchConfig()
	viper.SetConfigType(fileType)
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.SFatal("init config error, %s", err.Error())
	}
}
