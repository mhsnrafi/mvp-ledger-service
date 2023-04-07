package services

import (
	"github.com/spf13/viper"
	"ledger-service/models"
)

var Config *models.EnvConfig

func LoadConfig() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("SERVER_PORT", "8000")
	v.SetDefault("MODE", "debug")
	v.SetConfigType("dotenv")
	v.SetConfigName(".env.local")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}

	if err := Config.Validate(); err != nil {
		panic(err)
	}
}
