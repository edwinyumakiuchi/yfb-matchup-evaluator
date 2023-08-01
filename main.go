package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"
	"io"

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

	goth.UseProviders(
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "https://localhost"),
	)

    r.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
        user, err := gothic.CompleteUserAuth(res, req)
        if err != nil && !strings.Contains(err.Error(), "trying to fetch user information") {
            fmt.Fprintln(res, err)
            return
        }

        // Build the URL with the access token and the specific API endpoint.
        apiURL := "https://fantasysports.yahooapis.com/fantasy/v2/team/418.l.33024.t.4/roster"
        apiReq, apiErr := http.NewRequest("GET", apiURL, nil)
        if apiErr != nil {
            fmt.Println("Error creating request:", apiErr)
            return
        }
        apiReq.Header.Set("Authorization", "Bearer "+user.AccessToken)

        client := &http.Client{}
        apiResp, apiErr := client.Do(apiReq)
        if apiErr != nil {
            fmt.Println("Error making request:", apiErr)
            return
        }
        defer apiResp.Body.Close()

        if apiResp.StatusCode != http.StatusOK {
            fmt.Printf("API responded with status code %d\n", apiResp.StatusCode)
            return
        }

        apiBody, apiErr := io.ReadAll(apiResp.Body)
        // fmt.Println("Response Body:", string(apiBody))

        // Send the roster data as a JSON response
        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(http.StatusOK)
        res.Write(apiBody)
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
