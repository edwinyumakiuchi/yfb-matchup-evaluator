package util

import (
    "fmt"
)

func YahooToMongo(database string, collection string, accessToken string) () {
    // Call the function in yahooUtil.go to fetch the API data from Yahoo
    yahooRoster, yahooErr := RetrieveYahooRoster(accessToken)
    if yahooErr != nil {
        fmt.Println(yahooErr)
        return
    }

    mongoDeleteErr := DeleteDocuments("Cluster0", database, collection)
    if mongoDeleteErr != nil {
        fmt.Println("Error:", mongoDeleteErr)
    }

    mongoInsertErr := InsertOneDocument("Cluster0", database, collection, string(yahooRoster))
    if mongoInsertErr != nil {
        fmt.Println("Error:", mongoInsertErr)
    }
}