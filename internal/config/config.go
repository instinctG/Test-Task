package config

import "github.com/spf13/viper"

type Config struct {
	Port      string `mapstructure:"port"`
	SecretKey string `mapstructure:"secretKey"`
	MongoURL  string `mapstructure:"mongoURL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
