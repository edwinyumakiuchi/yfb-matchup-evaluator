package util

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"yfb-matchup-evaluator/config"
)

const API_APP_ENDPOINT = "https://us-west-2.aws.data.mongodb-api.com/app/data-natmv/"
const API_AUTH_ENDPOINT = "https://realm.mongodb.com/api/client/v2.0/app/data-natmv/auth/providers/local-userpass/login"
const API_ACTION_ENDPOINT = "endpoint/data/v1/action/"
const API_DELETEMANY_ENDPOINT = "deleteMany"
const API_INSERTONE_ENDPOINT = "insertOne"
const API_FIND_ENDPOINT = "find"

const CONFIG_FILE_PATH = "./config/config.yaml"
const SECRET_CONFIG_FILE_PATH = "./config/secretConfig.yaml"

// Define a struct to represent the login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	DeviceID     string `json:"device_id"`
}

type DeleteRequest struct {
	DataSource string               `json:"dataSource"`
	Database   string               `json:"database"`
	Collection string               `json:"collection"`
	Filter     map[string]interface{} `json:"filter"`
}

type InsertRequest struct {
	DataSource string                 `json:"dataSource"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Document   map[string]interface{} `json:"document"`
}

type MongoPlayer struct {
	ID        string            `json:"_id"`
	Names     map[string]string `json:"names"`
	TeamName  string            `json:"Team name"`
	Roster    []string          `json:"Roster"`
}

func (p *MongoPlayer) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID    string          `json:"_id"`
		Names map[string][]string `json:"-"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.ID = aux.ID

	// Convert the map of dynamic field names to a map with single player names
	p.Names = make(map[string]string)
	for fieldName, names := range aux.Names {
		if len(names) > 0 {
			p.Names[fieldName] = names[0]
		}
	}

	return nil
}

// Define a handler to handle the login API call
func LoginHandler() (string, error) {
    secretConfig, secretConfigErr := config.ReadConfig(SECRET_CONFIG_FILE_PATH)
	if secretConfigErr != nil {
		return "", secretConfigErr
	}

	credentials := map[string]string{
		"username": secretConfig.MongoUsername,
		"password": secretConfig.MongoPassword,
	}

	// Convert the credentials to JSON
	data, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	// Make the API call to your MongoDB endpoint
	resp, err := http.Post(API_AUTH_ENDPOINT, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

	// Decode the response JSON into the LoginResponse struct
	var loginResp LoginResponse
    err = json.Unmarshal(responseBody, &loginResp)
    if err != nil {
        return "", err
    }

	return loginResp.AccessToken, nil
}

func DeleteDocuments(dataSource, database, collection string) error {
    secretConfig, secretConfigErr := config.ReadConfig(SECRET_CONFIG_FILE_PATH)
	if secretConfigErr != nil {
		return secretConfigErr
	}

	deleteData := DeleteRequest{
		DataSource: dataSource,
		Database:   database,
		Collection: collection,
		Filter:     map[string]interface{}{},
	}

	requestBody, err := json.Marshal(deleteData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", API_APP_ENDPOINT + API_ACTION_ENDPOINT + API_DELETEMANY_ENDPOINT, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", secretConfig.MongoKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Collection %s: Error deleting data", collection)
	}

	return nil
}

func InsertOneDocument(dataSource, database, collection string, documentJSON string) error {
    secretConfig, secretConfigErr := config.ReadConfig(SECRET_CONFIG_FILE_PATH)
	if secretConfigErr != nil {
		return secretConfigErr
	}

	var doc map[string]interface{}
	var err error
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

	req, err := http.NewRequest("POST", API_APP_ENDPOINT + API_ACTION_ENDPOINT + API_INSERTONE_ENDPOINT, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", secretConfig.MongoKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Collection %s: Error storing data", collection)
	}

	return nil
}

// Define a function to fetch players data
func GetPlayers(database string, collection string, accessToken string) (string, error) {
	// Define the request payload
	payload := map[string]interface{}{
		"dataSource": "Cluster0",
		"database":   database,
		"collection": collection,
		"filter":     map[string]interface{}{},
	}

	// Convert the payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make the API call to your MongoDB endpoint
	req, err := http.NewRequest("POST", API_APP_ENDPOINT + API_ACTION_ENDPOINT + API_FIND_ENDPOINT, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + accessToken)

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

    // Read the response body
    responseBody, err := ioutil.ReadAll(resp.Body)

    return string(responseBody), nil
}

func RetrieveMongoData(database string, collection string) ([]byte, error) {
    mongoAccessToken, loginErr := LoginHandler()
    if loginErr != nil {
        fmt.Println("Error retrieving Mongo access token:", loginErr)
        return nil, loginErr
    }

    players, mongoFindErr := GetPlayers(database, collection, mongoAccessToken)
    if mongoFindErr != nil {
        fmt.Println("Error:", mongoFindErr)
        return nil, mongoFindErr
    }

    var data map[string]interface{}

    jsonErr := json.Unmarshal([]byte(players), &data)
    if jsonErr != nil {
        fmt.Println("Error parsing data:", jsonErr)
        return nil, jsonErr
    }

    documents, ok := data["documents"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("Error parsing documents data")
    }

    encodedData, encodingErr := json.Marshal(documents)
    if encodingErr != nil {
        fmt.Println("Error encoding response data:", encodingErr)
        return nil, encodingErr
    }

    return encodedData, nil
}