package util

import (
    "fmt"
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