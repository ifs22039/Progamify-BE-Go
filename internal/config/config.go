package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost            string `mapstructure:"DBHost"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	HuggingfaceApiKey string `mapstructure:"HUGGINGFACE_API_KEY"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	return config
}
