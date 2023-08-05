package main

import (
	"fmt"
	"strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/yaml.v2"

	"yfb-matchup-evaluator/util"

    "github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yahoo"
)

const CONFIG_FILE_PATH = "./config.yaml"

type Config struct {
	YahooClientID string `yaml:"yahoo_client_id"`
	YahooClientSecret string `yaml:"yahoo_client_secret"`
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

// Define a struct to represent the data
type GameData struct {
	Data []struct {
		Date  string   `json:"date"`
		Teams []string `json:"teams"`
	} `json:"data"`
}

// TODO: refactor with mongoUtil
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

func main() {
	r := pat.New()

    config, configErr := readConfig()
    if configErr != nil {
        return
    }

    // requires https://github.com/esplo/docker-local-ssl-termination-proxy/tree/master
	goth.UseProviders(
		yahoo.New(config.YahooClientID, config.YahooClientSecret, "https://localhost"),
	)

	r.Get("/mongoData", func(res http.ResponseWriter, req *http.Request) {
        accessToken, err := util.LoginHandler()
        if err != nil {
            fmt.Println("Error getting access token:", err)
            return
        }

        players, mongoFindErr := util.GetPlayers(accessToken)
        if mongoFindErr != nil {
            fmt.Println("Error:", mongoFindErr)
        } else {
            fmt.Println("Data retrieved successfully!")
        }

        var playersData map[string]interface{}

        playerJsonErr := json.Unmarshal([]byte(players), &playersData)
        if playerJsonErr != nil {
            http.Error(res, "Error parsing players data", http.StatusInternalServerError)
            return
        }

        // Access the "documents" array from playersData
        documents, ok := playersData["documents"].([]interface{})
        if !ok {
            http.Error(res, "Error parsing players data", http.StatusInternalServerError)
            return
        }

        responseJSON, responseErr := json.Marshal(documents)
        if responseErr != nil {
            http.Error(res, "Error encoding response data", http.StatusInternalServerError)
            return
        }

        res.Header().Set("Content-Type", "application/json")
        res.Write(responseJSON)
	})

	r.Get("/schedule", func(res http.ResponseWriter, req *http.Request) {
        gameData := `
        {
            "data": [
                {
                    "date": "10/17/2023",
                    "teams": ["LAL", "DEN", "BOS", "MIA"]
                },
                {
                    "date": "10/18/2023",
                    "teams": ["HOU", "NOP", "SAS", "TOR", "PHX", "PHI", "UTA", "ORL", "ATL", "WAS", "NYK", "CHI", "OKL", "BKN", "MEM", "MIN", "DET", "DAL", "SAC", "POR", "CHA", "CLE", "GSW", "IND"]
                },
                {
                    "date": "10/19/2023",
                    "teams": ["LAL", "MIA", "LAC", "IND"]
                },
                {
                    "date": "10/20/2023",
                    "teams": ["ORL", "PHI", "TOR", "NOP", "WAS", "GSW", "NYK", "DEN", "MEM", "SAS", "HOU", "BOS", "UTA", "CHA", "DET", "ATL", "POR", "PHX", "IND", "BKN", "MIN", "CHI"]
                },
                {
                    "date": "10/21/2023",
                    "teams": ["ORL", "CLE", "CHI", "MEM", "OKL", "IND", "MIL", "DET", "SAC", "DAL", "TOR", "HOU", "DEN", "BOS", "SAS", "LAC", "PHI", "MIA"]
                },
                {
                    "date": "10/22/2023",
                    "teams": ["CHA", "ATL", "WAS", "SAC", "PHX", "POR", "OKL", "UTA", "MIN", "LAC", "LAL", "GSW", "NOP", "CLE"]
                }
            ]
        }`

        res.Header().Set("Content-Type", "application/json")
        res.Write([]byte(gameData))
	})

    r.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
        user, err := gothic.CompleteUserAuth(res, req)
        if err != nil && !strings.Contains(err.Error(), "trying to fetch user information") {
            fmt.Fprintln(res, err)
            return
        }

        // Call the function in yahooUtil.go to fetch the API data.
        apiBody, apiErr := util.GetAPIData(user.AccessToken)
        if apiErr != nil {
            fmt.Println(apiErr)
            return
        }

        // Parse the API data and get the desired JSON
        jsonBytes, jsonErr := util.ParseData(apiBody)
        if jsonErr != nil {
            fmt.Println(jsonErr)
            return
        }

        mongoDeleteErr := util.DeleteDocuments("Cluster0", "yahoo", "rosters")
        if mongoDeleteErr != nil {
            fmt.Println("Error:", mongoDeleteErr)
        } else {
            fmt.Println("Data deleted successfully!")
        }

        mongoInsertErr := util.InsertOneDocument("Cluster0", "yahoo", "rosters", string(jsonBytes))
        if mongoInsertErr != nil {
            fmt.Println("Error:", mongoInsertErr)
        } else {
            fmt.Println("Data inserted successfully!")
        }

        http.Redirect(res, req, "http://localhost:3000/?loggedIn=true", http.StatusFound)
    })

	r.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
	})

	r.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			fmt.Println("/auth/yahoo: res = ", res)
			fmt.Println("/auth/yahoo: gothUser = ", gothUser)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

    // Add a route specifically for handling the redirect from `https://localhost/?code=abc`
    r.Get("/", func(res http.ResponseWriter, req *http.Request) {
        code := req.URL.Query().Get("code")
        state := req.URL.Query().Get("state")

        if code != "" {
            redirectURL := fmt.Sprintf("http://localhost:3000/auth/yahoo/callback/?code=%s&state=%s", code, state)
            http.Redirect(res, req, redirectURL, http.StatusFound)
            return
        }

        fs := http.FileServer(http.Dir("./build"))
        fs.ServeHTTP(res, req)
    })

	port := ":3000"
	fmt.Println("Starting backend server on port", port)
	http.Handle("/", r)
	http.ListenAndServe(port, nil)
}
