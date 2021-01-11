package main

import (
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
)

const configPath = "./configs/config.toml"

func main() {
	defaultConfig := &config.Config{
		InfluxDb: config.InfluxDb{
			DatabaseName:   "fumeping",
			ServerUrl:      "http://localhost:8086/",
			AuthEnabled:    true,
			Username:       "admin",
			Password:       "mycrazypassword",
			GatherInterval: 30,
		},
		PingMonitor: config.PingMonitor{
			Timeout:      10,
			PingInterval: 1,
			PayloadSize:  56,
			Destinations: map[string]config.Destination{
				"First": {
					Host: "mycustomhost",
				},
				"Second": {
					Host:    "mycustomipv6host",
					Network: "ipv6",
				},
			},
		},
	}
	if err := config.WriteConfig(configPath, defaultConfig); err != nil {
		panic(err)
	}
}
