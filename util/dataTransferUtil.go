package util

import (
    "fmt"
    "encoding/json"
    "time"
    "strconv"

    "yfb-matchup-evaluator/config"
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
            if yahooErr != nil && yahooErr.Error() != "Draft has not started." {
                fmt.Println(yahooErr)
                break
            } else {
                yahooDraftResults = string(draftBytes)

                config, _ := config.ReadConfig(CONFIG_FILE_PATH)
                team1 := config.YahooYearID + ".l." + config.YahooLeagueID + ".t." + config.YahooTeamID
                teams := []string{team1}

                var results struct {
                    Budgets []struct {
                        AvgBudget string `json:"Avg-Budget"`
                        Budget    float64 `json:"Budget"`
                        TeamKey   string `json:"TeamKey"`
                        SelfCost  float64 `json:"SelfCost"`
                        AvgCost   string `json:"AvgCost"`
                    } `json:"budgets"`
                }

                if yahooErr == nil {
                    if err := json.Unmarshal([]byte(yahooDraftResults), &results); err != nil {
                        fmt.Println("Error unmarshaling yahooDraftResults:", err)
                        return
                    }
                }

                for _, team := range teams {
                    mongoDeleteErr := DeleteDocuments("Cluster0", database, "draft-" + team)
                    if mongoDeleteErr != nil {
                        fmt.Println("Error:", mongoDeleteErr)
                    }

                    var yahooTeamDraftResults string
                    if yahooErr != nil && yahooErr.Error() == "Draft has not started." {
                        yahooTeamDraftResults = fmt.Sprintf(`{"priceAdjustment":"%.2f", "avgSelfCost": "%s"}`, float64(0), "15")
                    } else {
                        for _, budget := range results.Budgets {
                            if budget.TeamKey == team {
                                myAvgCost := budget.SelfCost
                                leagueAvgCost, _ := strconv.ParseFloat(budget.AvgCost, 64)
                                yahooTeamDraftResults = fmt.Sprintf(`{"priceAdjustment":"%.2f", "avgSelfCost": "%s"}`, float64(myAvgCost - leagueAvgCost), budget.SelfCost)
                                break
                            }
                        }
                    }

                    type DraftResults struct {
                    	PriceAdjustment string `json:"priceAdjustment"`
                    	AvgSelfCost     string `json:"avgSelfCost"`
                    }

                    var results DraftResults
                    err := json.Unmarshal([]byte(yahooTeamDraftResults), &results)
                    if err != nil {
                        fmt.Println("Error unmarshalling yahooTeamDraftResults:", err)
                        return
                    }
                    fmt.Println("Price Adjustment: ", results.PriceAdjustment)

                    /* mongoInsertErr := InsertOneDocument("Cluster0", database, "draft-" + team, yahooTeamDraftResults)
                    if mongoInsertErr != nil {
                        fmt.Println("Error:", mongoInsertErr)
                    } */
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