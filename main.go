package main

import (
	"fmt"
	"sort"
	"strings"
	"strconv"
    "net/http"
    "encoding/json"

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

    hbProjections, err := util.RetrieveMongoData("sample-nba", "projections")
    if err != nil {
        fmt.Println("Error retrieving hashtagbasketball projections:", err)
        return
    }
    fmt.Println("Hashtagbasketball projections retrieved successfully!")

    yahooRosters, err := util.RetrieveMongoData("yahoo", "rosters")
    if err != nil {
        fmt.Println("Error retrieving yahoo rosters:", err)
        return
    }
    fmt.Println("Yahoo rosters retrieved successfully!")
    // fmt.Println("yahooRosters:", string(yahooRosters))

    r.Get("/seasonOutlook", func(res http.ResponseWriter, req *http.Request) {
        var rosters []map[string]interface{}
        if err := json.Unmarshal(yahooRosters, &rosters); err != nil {
            fmt.Println("Error parsing yahooRosters:", err)
            return
        }

        teamStats := make(map[string]map[string]float64)
        for _, roster := range rosters {
            teamName := roster["Fantasy Team"].(string)
            // Initialize teamStats with the desired keys for statistics
            if teamStats[teamName] == nil {
                teamStats[teamName] = make(map[string]float64)

                teamStats[teamName]["fieldGoalMade"] = 0
                teamStats[teamName]["fieldGoalAttempt"] = 0
                teamStats[teamName]["fieldGoal"] = 0
                teamStats[teamName]["freeThrowMade"] = 0
                teamStats[teamName]["freeThrowAttempt"] = 0
                teamStats[teamName]["freeThrow"] = 0
                teamStats[teamName]["threePointMade"] = 0
                teamStats[teamName]["points"] = 0
                teamStats[teamName]["totalRebounds"] = 0
                teamStats[teamName]["assists"] = 0
                teamStats[teamName]["steals"] = 0
                teamStats[teamName]["blocks"] = 0
                teamStats[teamName]["turnovers"] = 0
            }
        }

        // Parse the JSON data from hbProjections
        var projections []map[string]interface{}
        if err := json.Unmarshal(hbProjections, &projections); err != nil {
            fmt.Println("Error parsing hbProjections:", err)
            return
        }

        // Iterate through teamStats and match players with projections
        for _, roster := range rosters {
            teamName := roster["Fantasy Team"].(string)

            for _, playerData := range roster["Roster"].([]interface{}) {
                player := playerData.(map[string]interface{})
                playerName := player["Player"].(string)

                for _, projection := range projections {
                    if projection["name"].(string) == playerName {

                        fieldGoalMade, _ := strconv.ParseFloat(projection["fieldGoalMadeCalculated"].(string), 64)
                        fieldGoalAttempt, _ := strconv.ParseFloat(projection["fieldGoalAttempt"].(string), 64)
                        freeThrowMade, _ := strconv.ParseFloat(projection["freeThrowMadeCalculated"].(string), 64)
                        freeThrowAttempt, _ := strconv.ParseFloat(projection["freeThrowAttempt"].(string), 64)
                        threePointMade, _ := strconv.ParseFloat(projection["threePointMade"].(string), 64)
                        points, _ := strconv.ParseFloat(projection["points"].(string), 64)
                        totalRebounds, _ := strconv.ParseFloat(projection["totalRebounds"].(string), 64)
                        assists, _ := strconv.ParseFloat(projection["assists"].(string), 64)
                        steals, _ := strconv.ParseFloat(projection["steals"].(string), 64)
                        blocks, _ := strconv.ParseFloat(projection["blocks"].(string), 64)
                        turnovers, _ := strconv.ParseFloat(projection["turnovers"].(string), 64)

                        teamStats[teamName]["fieldGoalMade"] += fieldGoalMade
                        teamStats[teamName]["fieldGoalAttempt"] += fieldGoalAttempt
                        teamStats[teamName]["freeThrowMade"] += freeThrowMade
                        teamStats[teamName]["freeThrowAttempt"] += freeThrowAttempt
                        teamStats[teamName]["threePointMade"] += threePointMade
                        teamStats[teamName]["points"] += points
                        teamStats[teamName]["totalRebounds"] += totalRebounds
                        teamStats[teamName]["assists"] += assists
                        teamStats[teamName]["steals"] += steals
                        teamStats[teamName]["blocks"] += blocks
                        teamStats[teamName]["turnovers"] += turnovers
                    }
                }
            }
        }

        // Calculate team projections and convert to JSON
        var teamProjections []map[string]interface{}
        for team, stats := range teamStats {
            teamProjection := make(map[string]interface{})
            teamProjection["Fantasy Team"] = team

            totalFieldGoalMades := stats["fieldGoalMade"]
            totalFieldGoalAttempts := stats["fieldGoalAttempt"]
            totalFreeThrowMades := stats["freeThrowMade"]
            totalFreeThrowAttempts := stats["freeThrowAttempt"]
            totalThreePointMade := stats["threePointMade"]
            totalPoints := stats["points"]
            totalRebounds := stats["totalRebounds"]
            totalAssists := stats["assists"]
            totalSteals := stats["steals"]
            totalBlocks := stats["blocks"]
            totalTurnovers := stats["turnovers"]

            fieldGoalData := make(map[string]interface{})
            fieldGoalData["value"] = formatFloat(totalFieldGoalMades / totalFieldGoalAttempts, 4)
            fieldGoalData["rank"] = 4

            freeThrowData := make(map[string]interface{})
            freeThrowData["value"] = formatFloat(totalFreeThrowMades / totalFreeThrowAttempts, 4)
            freeThrowData["rank"] = 4

            threePointMadeData := make(map[string]interface{})
            threePointMadeData["value"] = formatFloat(totalThreePointMade, 1)
            threePointMadeData["rank"] = 4

            pointsData := make(map[string]interface{})
            pointsData["value"] = formatFloat(totalPoints, 1)
            pointsData["rank"] = 4

            totalReboundsData := make(map[string]interface{})
            totalReboundsData["value"] = formatFloat(totalRebounds, 1)
            totalReboundsData["rank"] = 4

            assistsData := make(map[string]interface{})
            assistsData["value"] = formatFloat(totalAssists, 1)
            assistsData["rank"] = 4

            stealsData := make(map[string]interface{})
            stealsData["value"] = formatFloat(totalSteals, 1)
            stealsData["rank"] = 4

            blocksData := make(map[string]interface{})
            blocksData["value"] = formatFloat(totalBlocks, 1)
            blocksData["rank"] = 4

            turnoversData := make(map[string]interface{})
            turnoversData["value"] = formatFloat(totalTurnovers, 1)
            turnoversData["rank"] = 4

            teamProjection["fieldGoalMade"] = formatFloat(totalFieldGoalMades, 3)
            teamProjection["fieldGoalAttempt"] = formatFloat(totalFieldGoalAttempts, 1)
            teamProjection["fieldGoal"] = fieldGoalData
            teamProjection["freeThrowMade"] = formatFloat(totalFreeThrowMades, 3)
            teamProjection["freeThrowAttempt"] = formatFloat(totalFreeThrowAttempts, 1)
            teamProjection["freeThrow"] = freeThrowData
            teamProjection["threePointMade"] = threePointMadeData
            teamProjection["points"] = pointsData
            teamProjection["totalRebounds"] = totalReboundsData
            teamProjection["assists"] = assistsData
            teamProjection["steals"] = stealsData
            teamProjection["blocks"] = blocksData
            teamProjection["turnovers"] = turnoversData

            teamProjections = append(teamProjections, teamProjection)
        }

        // Define the categories to process
        categories := []string{"fieldGoal", "freeThrow", "threePointMade", "totalRebounds", "assists", "steals", "blocks", "turnovers", "points"}

        for _, category := range categories {
            // Custom sorting function for the current category
            sort.Slice(teamProjections, func(i, j int) bool {
                categoryI := teamProjections[i][category].(map[string]interface{})
                categoryJ := teamProjections[j][category].(map[string]interface{})
                valueI := categoryI["value"].(float64)
                valueJ := categoryJ["value"].(float64)

                if category == "turnovers" {
                    return valueJ > valueI
                }
                return valueI > valueJ
            })

            // Assign ranks based on the sorted order for the current category
            for i := range teamProjections {
                teamProjections[i][category].(map[string]interface{})["rank"] = i + 1
            }
        }

        // Convert to JSON and print or use as needed
        teamProjectionJSON, _ := json.Marshal(teamProjections)
        // fmt.Println("Team Projections:", string(teamProjectionJSON))

        res.Header().Set("Content-Type", "application/json")
        res.Write(teamProjectionJSON)
    })

	r.Get("/yahooRosters", func(res http.ResponseWriter, req *http.Request) {
        res.Header().Set("Content-Type", "application/json")
        res.Write(yahooRosters)
	})

	r.Get("/hbProjections", func(res http.ResponseWriter, req *http.Request) {
        res.Header().Set("Content-Type", "application/json")
        res.Write(hbProjections)
	})

    r.Get("/schedule", func(res http.ResponseWriter, req *http.Request) {
	    hbSchedule, err := util.RetrieveMongoData("sample-nba", "advanced-nba-schedule-grid")
        if err != nil {
            fmt.Println("Error retrieving hashtagbasketball advanced-nba-schedule-grid:", err)
            return
        }
        fmt.Println("Hashtagbasketball advanced-nba-schedule-grid retrieved successfully!")

        res.Header().Set("Content-Type", "application/json")
        res.Write(hbSchedule)
    })

	r.Get("/yahooMatchup", func(res http.ResponseWriter, req *http.Request) {
	    yahooRosters, err := util.RetrieveMongoData("yahoo", "matchup")
        if err != nil {
            fmt.Println("Error retrieving yahoo matchup:", err)
            return
        }

        fmt.Println("Yahoo matchup retrieved successfully!")

        res.Header().Set("Content-Type", "application/json")
        res.Write(yahooRosters)
	})

    r.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
        user, err := gothic.CompleteUserAuth(res, req)
        if err != nil && !strings.Contains(err.Error(), "trying to fetch user information") {
            fmt.Fprintln(res, err)
            return
        }

        util.YahooToMongo("yahoo", "rosters", user.AccessToken)
        util.YahooToMongo("yahoo", "matchup", user.AccessToken)

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

func formatFloat(value float64, decimals int) float64 {
    format := fmt.Sprintf("%%.%df", decimals)
    formattedValue, _ := strconv.ParseFloat(fmt.Sprintf(format, value), 64)
    return formattedValue
}