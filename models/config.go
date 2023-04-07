package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EnvConfig struct {
	DBHost                     string `mapstructure:"POSTGRES_HOST"`
	DBUserName                 string `mapstructure:"POSTGRES_USER"`
	DBUserPassword             string `mapstructure:"POSTGRES_PASSWORD"`
	DBName                     string `mapstructure:"POSTGRES_DB"`
	DBPort                     string `mapstructure:"POSTGRES_PORT"`
	ServerHost                 string `mapstructure:"SERVER_HOST"`
	ServerPort                 string `mapstructure:"SERVER_PORT"`
	UseRedis                   bool   `mapstructure:"USE_REDIS"`
	RedisDefaultAddr           string `mapstructure:"REDIS_DEFAULT_ADDR"`
	RedisPassword              string `mapstructure:"REDIS_PASSWORD"`
	JWTSecretKey               string `mapstructure:"JWT_SECRET"`
	JWTAccessExpirationMinutes int    `mapstructure:"JWT_ACCESS_EXPIRATION_MINUTES"`
	JWTRefreshExpirationDays   int    `mapstructure:"JWT_REFRESH_EXPIRATION_DAYS"`
	Mode                       string `mapstructure:"MODE"`
}

func (config *EnvConfig) Validate() error {
	return validation.ValidateStruct(config,
		validation.Field(&config.DBPort, is.Port),
		validation.Field(&config.DBHost, validation.Required),
		validation.Field(&config.DBUserPassword, validation.Required),
		validation.Field(&config.DBName, validation.Required),
		validation.Field(&config.UseRedis, validation.In(true, false)),
		validation.Field(&config.RedisDefaultAddr),

		validation.Field(&config.JWTSecretKey, validation.Required),
		validation.Field(&config.JWTAccessExpirationMinutes, validation.Required),
		validation.Field(&config.JWTRefreshExpirationDays, validation.Required),

		validation.Field(&config.Mode, validation.In("debug", "release")),
	)
}
