package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
	YahooClientID     string `yaml:"yahoo_client_id"`
	YahooClientSecret string `yaml:"yahoo_client_secret"`

    YahooRedirectURI  string `yaml:"yahoo_redirect_uri"`
    YahooTeamURL      string `yaml:"yahoo_team_url"`
    YahooLeagueURL    string `yaml:"yahoo_league_url"`
    YahooYearID       string `yaml:"yahoo_year_id"`
    YahooLeagueID     string `yaml:"yahoo_league_id"`
    YahooTeamID       string `yaml:"yahoo_team_id"`

	MongoKey          string `yaml:"mongo_key"`
	MongoUsername     string `yaml:"mongo_username"`
	MongoPassword     string `yaml:"mongo_password"`
}

func ReadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
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