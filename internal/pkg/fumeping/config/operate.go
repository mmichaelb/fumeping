package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func ReadConfig(filename string) (*Config, error) {
	config := LoadDefaultConfig
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.NewDecoder(file).Decode(&config)
	return &config, err
}

func WriteConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	configFileHeader := fmt.Sprintf("# FumePing Config written on %s\n", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	if _, err = file.Write([]byte(configFileHeader)); err != nil {
		return err
	}
	defer file.Close()
	err = yaml.NewEncoder(file).Encode(config)
	return err
}
