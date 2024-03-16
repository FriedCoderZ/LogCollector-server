package config

import (
	"log"

	"github.com/spf13/viper"
)

// 读取配置文件config
type Config struct {
	Crypto CryptoConfig
}

var (
	config Config
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	viper.Unmarshal(&config)
}

func GetConfig() Config {
	return config
}
