package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"

	"yfb-matchup-evaluator/util"

    "github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/yahoo"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func main() {
	r := pat.New()

    // prereq: https://github.com/esplo/docker-local-ssl-termination-proxy/blob/master/Dockerfile
	goth.UseProviders(
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "https://localhost"),
	)

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

        mongoErr := util.InsertOneDocument("Cluster0", "yahoo", "rosters", string(jsonBytes))
        if mongoErr != nil {
            fmt.Println("Error:", mongoErr)
        } else {
            fmt.Println("Data inserted successfully!")
        }

        // Send the roster data as a JSON response
        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(http.StatusOK)
        res.Write(jsonBytes)
    })

    r.Get("/home", func(res http.ResponseWriter, req *http.Request) {
        // Send the roster data as a JSON response
        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(http.StatusOK)
        // res.Write(jsonBytes)
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
