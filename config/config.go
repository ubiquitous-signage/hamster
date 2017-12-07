package config

import (
	"github.com/spf13/viper"
	"log"
)

func Run() {
	readConfig("wordCloudConfig")
}

func readConfig(fileName string) {
	viper.SetConfigType("toml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}