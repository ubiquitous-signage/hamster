package config

import (
	"github.com/spf13/viper"
	"log"
)

func Run(env string) {
	if env == "prod" {
		readConfig("config.prod")
		log.Println("starting with prodcution settings...")
	} else {
		readConfig("config.dev")
		log.Println("starting with development settings...")
	}
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