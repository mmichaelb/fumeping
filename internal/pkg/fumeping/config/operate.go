package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

func ReadConfig(filename string) (*Config, error) {
	config := &Config{}
	_, err := toml.DecodeFile(filename, config)
	return config, err
}

func WriteConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	configFileHeader := fmt.Sprintf("# Config written on %s\n", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	if _, err = file.Write([]byte(configFileHeader)); err != nil {
		return err
	}
	defer file.Close()
	err = toml.NewEncoder(file).Encode(config)
	return err
}
