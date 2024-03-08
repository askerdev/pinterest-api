package application

import (
	"github.com/spf13/viper"
)

func MustSetupConfig() {
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
