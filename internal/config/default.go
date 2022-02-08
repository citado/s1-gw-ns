package config

import (
	"time"

	"github.com/citado/s1-gw-ns/internal/app"
	"github.com/citado/s1-gw-ns/internal/lora"
)

// nolint: gomnd
func Default() Config {
	return Config{
		Tries: 10,
		App: app.Config{
			Addr:  "127.0.0.1",
			Port:  1883,
			Total: 10,
			Delay: 1 * time.Second,
		},
		Gateways: []lora.Config{
			{
				MAC: "b827ebffff70c80a",
				Keys: lora.Keys{
					NetworkSKey:     "DB56B6C3002A4763A79E64573C629D97",
					ApplicationSKey: "94B49CD7BC621BC46571D019640804AA",
				},
				Device: lora.Device{
					Addr: "26011CF6",
				},
			},
		},
	}
}