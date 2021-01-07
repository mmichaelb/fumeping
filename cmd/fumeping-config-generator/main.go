package main

import (
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
	"time"
)

const configPath = "./configs/config.toml"

func main() {
	defaultConfig := &config.Config{
		InfluxDb: config.InfluxDb{
			Organization: "fumeping",
			BucketSyntax: "stats-%s",
			ServerUrl:    "http://localhost:8086/",
			Username:     "admin",
			Password:     "mycrazypassword",
		},
		Destinations: []config.DestinationHost{
			{
				Host:           "mycustomhost",
				Interval:       time.Minute,
				PacketInterval: time.Millisecond * 100,
				Count:          10,
				Timeout:        time.Second * 10,
			},
			{
				Host:           "yetanothercustomhost",
				Interval:       time.Minute,
				PacketInterval: time.Millisecond * 100,
				Count:          10,
				Timeout:        time.Second * 10,
			},
		},
	}
	if err := config.WriteConfig(configPath, defaultConfig); err != nil {
		panic(err)
	}
}
