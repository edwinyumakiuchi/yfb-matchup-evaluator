package util

import (
    "fmt"
)

func YahooToMongo(database string, collection string, accessToken string) () {
    var yahooRoster string
    var yahooMatchup string
    var yahooErr error

    if collection == "rosters" {
        var rosterBytes []byte
        rosterBytes, yahooErr = RetrieveYahooRoster(accessToken)
        yahooRoster = string(rosterBytes)
    } else {
        var matchupBytes []byte
        matchupBytes, yahooErr = RetrieveYahooMatchup(accessToken)
        yahooMatchup = string(matchupBytes)
    }

    if yahooErr != nil {
        fmt.Println(yahooErr)
        return
    }

    mongoDeleteErr := DeleteDocuments("Cluster0", database, collection)
    if mongoDeleteErr != nil {
        fmt.Println("Error:", mongoDeleteErr)
    }

    if collection == "rosters" {
        mongoInsertErr := InsertOneDocument("Cluster0", database, collection, yahooRoster)
        if mongoInsertErr != nil {
            fmt.Println("Error:", mongoInsertErr)
        }
    } else {
        mongoInsertErr := InsertOneDocument("Cluster0", database, collection, yahooMatchup)
        if mongoInsertErr != nil {
            fmt.Println("Error:", mongoInsertErr)
        }
    }
}