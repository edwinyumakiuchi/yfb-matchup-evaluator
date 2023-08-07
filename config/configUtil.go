package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

const CONFIG_FILE_PATH = "./config/config.yaml"

type Config struct {
	YahooClientID string `yaml:"yahoo_client_id"`
	YahooClientSecret string `yaml:"yahoo_client_secret"`

	MongoKey string `yaml:"mongo_key"`
	MongoUsername string `yaml:"mongo_username"`
	MongoPassword string `yaml:"mongo_password"`
}

func ReadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}