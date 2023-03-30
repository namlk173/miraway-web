package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	ServerPort         string `mapstructure:"SERVER_PORT"`
	ContextTimeout     int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBName             string `mapstructure:"DB_NAME"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPass             string `mapstructure:"DB_PASS"`
	AccessTokenExpiry  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiry int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Can't find .env file", err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalln("Can't load data from file", err)
	}

	return &env
}
