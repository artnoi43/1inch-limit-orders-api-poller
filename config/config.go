package config

import (
	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/get"
	"github.com/spf13/viper"
)

type Config struct {
	Chain    enums.Chain `mapstructure:"chain"`
	Get      get.Config  `mapstructure:"get"`
	MaxGuard int         `mapstructure:"max_guard"`
}

func Load() (*Config, error) {
	viper.SetDefault("chain", enums.ChainEthereum)
	viper.SetDefault("max_guard", 10)
	viper.SetDefault("get.limit", 500)
	viper.SetDefault("get.interval", 2)
	viper.SetDefault("get.limit_order_url", "https://limit-orders.1inch.io/v2.0/1/limit-order/all")

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		// Config file not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// WriteConfig() just won't create new file if doesn't exist
			viper.SafeWriteConfig()
		} else {
			return nil, err
		}
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}
