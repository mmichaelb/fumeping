package main

import (
	"github.com/mmichaelb/fumeping/internal/pkg/fumeping/config"
)

const configPath = "./configs/config.toml"

func main() {
	if err := config.WriteConfig(configPath, &config.DefaultConfig); err != nil {
		panic(err)
	}
}
