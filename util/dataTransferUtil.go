package util

import (
    "fmt"
    "encoding/json"
    "time"
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
                break
                // return
            }

            mongoInsertErr := InsertOneDocument("Cluster0", database, collection, yahooRoster)
            if mongoInsertErr != nil {
                fmt.Println("Error:", mongoInsertErr)
            }
        }

        // Draft Bid Effect
        for {
            var draftBytes []byte
            var yahooDraftResults string

            draftBytes, yahooErr = RetrieveYahooDraftResults(accessToken)
            if yahooErr != nil {
                fmt.Println(yahooErr)
                break
            } else {
                yahooDraftResults = string(draftBytes)

                team1 := "454.l.47273.t.1"
                team2 := "454.l.47273.t.2"
                teams := []string{team1, team2}

                var results struct {
                    Budgets []struct {
                        AvgBudget string `json:"Avg-Budget"`
                        Budget    string `json:"Budget"`
                        TeamKey   string `json:"TeamKey"`
                        SelfCost  float64 `json:"SelfCost"`
                        AvgCost   float64 `json:"AvgCost"`
                    } `json:"budgets"`
                }

                fmt.Println("yahooDraftResults:", yahooDraftResults)

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
                            myAvgCost := budget.SelfCost
                            leagueAvgCost := budget.AvgCost
                            yahooTeamDraftResults = fmt.Sprintf(`{"priceAdjustment":"%.2f", "avgSelfCost": "%s"}`, float64(myAvgCost - leagueAvgCost), budget.SelfCost)
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