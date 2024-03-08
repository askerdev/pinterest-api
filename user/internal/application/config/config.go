package config

import (
	"github.com/spf13/viper"
)

func MustSetup() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if viper.GetString("app_mode") != "PRODUCTION" {
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
}
