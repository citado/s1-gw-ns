package config

import (
	"time"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/lora/api"
)

// nolint: gomnd
func Default() Config {
	return Config{
		Tries: 10,
		LoRaServer: api.Config{
			URL:      "http://127.0.0.1:8080",
			Username: "admin",
			Password: "admin",
		},
		App: app.Config{
			Addr:  "127.0.0.1",
			Port:  1883,
			Total: 10,
			Delay: 1 * time.Second,
		},
	}
}
