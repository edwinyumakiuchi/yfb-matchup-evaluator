package util

import (
    "fmt"
    "encoding/json"
    "time"
    "strconv"
)

const leagueSize = 12

func YahooToMongo(database string, collection string, accessToken string) () {
    var yahooRoster string
    var yahooMatchup string
    var yahooErr error

    mongoDeleteErr := DeleteDocuments("Cluster0", database, collection)
    if mongoDeleteErr != nil {
        fmt.Println("Error:", mongoDeleteErr)
    }

    if collection == "rosters" {
        var rosterBytes []byte
        for i := 1; i <= leagueSize; i++ {
            rosterBytes, yahooErr = RetrieveYahooRoster(accessToken, i)
            yahooRoster = string(rosterBytes)

            if yahooErr != nil {
                fmt.Println(yahooErr)
                return
            }

            mongoInsertErr := InsertOneDocument("Cluster0", database, collection, yahooRoster)
            if mongoInsertErr != nil {
                fmt.Println("Error:", mongoInsertErr)
            }
        }

        for {
            var draftBytes []byte
            var yahooDraftResults string

            draftBytes, yahooErr = RetrieveYahooDraftResults(accessToken)
            yahooDraftResults = string(draftBytes)

            team1 := "428.l.57655.t.4"
            team2 := "428.l.57655.t.8"
            teams := []string{team1, team2}

            var results struct {
                Budgets []struct {
                    AvgBudget string `json:"Avg-Budget"`
                    Budget    string `json:"Budget"`
                    TeamKey   string `json:"TeamKey"`
                } `json:"budgets"`
            }

            if err := json.Unmarshal([]byte(yahooDraftResults), &results); err != nil {
                fmt.Println("Error unmarshaling JSON:", err)
                return
            }

            for _, team := range teams {
                mongoDeleteErr := DeleteDocuments("Cluster0", database, "draft-" + team)
                if mongoDeleteErr != nil {
                    fmt.Println("Error:", mongoDeleteErr)
                }

                var yahooTeamDraftResults string
                for _, budget := range results.Budgets {
                    if budget.TeamKey == team {
                        // Construct the desired output
                        myBudget, _ := strconv.ParseFloat(budget.Budget, 64)
                        avgBudget, _ := strconv.ParseFloat(budget.AvgBudget, 64)
                        yahooTeamDraftResults = fmt.Sprintf(`{"priceAdjustment":"%.2f", "avgBudget": "%s"}`, float64((myBudget - avgBudget) / leagueSize), budget.AvgBudget)
                        break
                    }
                }

                fmt.Println("Inserting: ", yahooTeamDraftResults)
                mongoInsertErr := InsertOneDocument("Cluster0", database, "draft-" + team, yahooTeamDraftResults)
                if mongoInsertErr != nil {
                    fmt.Println("Error:", mongoInsertErr)
                }
            }

            time.Sleep(10 * time.Second)
        }
    } else {
        var matchupBytes []byte
        matchupBytes, yahooErr = RetrieveYahooMatchup(accessToken)
        yahooMatchup = string(matchupBytes)

        if yahooErr != nil {
            fmt.Println(yahooErr)
            return
        }

        mongoInsertErr := InsertOneDocument("Cluster0", database, collection, yahooMatchup)
        if mongoInsertErr != nil {
            fmt.Println("Error:", mongoInsertErr)
        }
    }
}