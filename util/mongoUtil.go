package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
)

const API_ENDPOINT = "https://us-west-2.aws.data.mongodb-api.com/app/data-natmv/endpoint/data/v1/action/"
const API_INSERTONE_ENDPOINT = "insertOne"
const CONFIG_FILE_PATH = "./config.yaml"

type Config struct {
	MongoKey string `yaml:"mongo_key"`
}

type InsertRequest struct {
	DataSource string     `json:"dataSource"`
	Database   string     `json:"database"`
	Collection string     `json:"collection"`
	Document   map[string]interface{} `json:"document"`
}

func readConfig() (*Config, error) {
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

func InsertOneDocument(dataSource, database, collection string, documentJSON string) error {
    config, err := readConfig()
	if err != nil {
		return err
	}

	var doc map[string]interface{}
	err = json.Unmarshal([]byte(documentJSON), &doc)
	if err != nil {
		return err
	}

	insertData := InsertRequest{
		DataSource: dataSource,
		Database:   database,
		Collection: collection,
		Document: doc,
	}

	requestBody, err := json.Marshal(insertData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", API_ENDPOINT + API_INSERTONE_ENDPOINT, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.MongoKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Collection %s: Error storing data", collection)
	}

	return nil
}
