package config

import (
	"log"

	"github.com/spf13/viper"
)

type AirbyteConfig struct {
	URL       string `mapstructure:"ABCTLX_AB_URL"`
	ClientId  string `mapstructure:"ABCTLX_AB_CID"`
	ClientKey string `mapstructure:"ABCTLX_AB_CK"`
	Port      int    `mapstructure:"ABCTLX_AB_PORT"`
}

var Data AirbyteConfig

func LoadEnv() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	err := viper.Unmarshal(&Data)
	if err != nil {
		panic("Unable to decode into struct, check your .env types: " + err.Error())
	}
}
