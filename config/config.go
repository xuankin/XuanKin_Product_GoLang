package config

import "github.com/spf13/viper"

type Config struct {
	DbSource             string `mapstructure:"DB_SOURCE"`
	ServerAddress        string `mapstructure:"SERVER_ADDRESS"`
	RedisAddr            string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	RedisDB              int    `mapstructure:"REDIS_DB"`
	BaseUrl              string `mapstructure:"BASE_URL"`
	ElasticsearchAddress string `mapstructure:"ELASTICSEARCH_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
