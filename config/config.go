package config

import (
	"github.com/spf13/viper"
)

func InitConfigFromFile(fileName string, fileType string) {
	viper.WatchConfig()
	viper.SetConfigType(fileType)
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		panic("init config error, " + err.Error())
	}
}
