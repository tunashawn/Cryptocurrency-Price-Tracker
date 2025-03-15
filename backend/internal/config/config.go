package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

func GetConfig(cfg any) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return fmt.Errorf("read config from env %v", err)
	}

	return nil
}

type PriceTrackerConfig struct {
	BinanceWebSocketURL string `envconfig:"BINANCE_WEBSOCKET_URL" required:"true"`
	WsFetchInterval     int64  `envconfig:"WEBSOCKET_FETCH_INTERVAL" default:"1m"`
}

type SqliteConfig struct {
	SqliteURL string `envconfig:"SQLITE_URL" required:"true"`
}
