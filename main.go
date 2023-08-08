package main

import (
	"fmt"
	"strings"
    "net/http"

	"yfb-matchup-evaluator/util"
	"yfb-matchup-evaluator/config"

    "github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yahoo"
)

const CONFIG_FILE_PATH = "./config/config.yaml"
const SECRET_CONFIG_FILE_PATH = "./config/secretConfig.yaml"

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

func main() {
	r := pat.New()

    secretConfig, secretConfigErr := config.ReadConfig(SECRET_CONFIG_FILE_PATH)
    if secretConfigErr != nil {
        return
    }

    config, configErr := config.ReadConfig(CONFIG_FILE_PATH)
    if configErr != nil {
        return
    }

    // requires https://github.com/esplo/docker-local-ssl-termination-proxy/tree/master
	goth.UseProviders(
		yahoo.New(secretConfig.YahooClientID, secretConfig.YahooClientSecret, config.YahooRedirectURI),
	)

	r.Get("/yahooRosters", func(res http.ResponseWriter, req *http.Request) {
	    yahooRosters, err := util.RetrieveMongoData("yahoo", "rosters")
        if err != nil {
            fmt.Println("Error retrieving yahoo rosters:", err)
            return
        }

        fmt.Println("Yahoo rosters retrieved successfully!")

        res.Header().Set("Content-Type", "application/json")
        res.Write(yahooRosters)
	})

	r.Get("/hbProjections", func(res http.ResponseWriter, req *http.Request) {
	    hbProjections, err := util.RetrieveMongoData("sample-nba", "projections")
        if err != nil {
            fmt.Println("Error retrieving hashtagbasketball projections:", err)
            return
        }

        fmt.Println("hashtagbasketball projections retrieved successfully!")

        res.Header().Set("Content-Type", "application/json")
        res.Write(hbProjections)
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
                    "teams": ["HOU", "NOP", "SAS", "TOR", "PHX", "PHI", "UTA", "ORL", "ATL", "WAS", "NYK", "CHI", "OKL", "BKN", "MEM", "MIN", "DET", "DAL", "SAC", "POR", "CHA", "CLE", "GSW", "MIL"]
                },
                {
                    "date": "10/19/2023",
                    "teams": ["LAL", "MIA", "LAC", "IND"]
                },
                {
                    "date": "10/20/2023",
                    "teams": ["ORL", "PHI", "TOR", "NOP", "WAS", "GSW", "NYK", "DEN", "MEM", "SAS", "HOU", "BOS", "UTA", "CHA", "DET", "ATL", "POR", "BKN", "MIN", "CHI"]
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
        }

        mongoInsertErr := util.InsertOneDocument("Cluster0", "yahoo", "rosters", string(jsonBytes))
        if mongoInsertErr != nil {
            fmt.Println("Error:", mongoInsertErr)
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
