package util

import (
    "fmt"
	"encoding/json"
	"encoding/xml"

	"yfb-matchup-evaluator/config"
)

type FantasyMatchupContent struct {
    XMLName xml.Name `xml:"fantasy_content"`
    League  League     `xml:"league"`
}

type League struct {
	Scoreboard  Scoreboard  `xml:"scoreboard"`
}

type Scoreboard struct {
	Matchups  []Matchup  `xml:"matchups>matchup"`
}

type Matchup struct {
	MatchupTeams  []MatchupTeam  `xml:"teams>team"`
}

type MatchupTeam struct {
	Name       string     `xml:"name"`
	TeamStats  TeamStats  `xml:"team_stats"`
}

type TeamStats struct {
	Stats  []Stat  `xml:"stats>stat"`
}

type Stat struct {
	StatID     string  `xml:"stat_id"`
	StatValue  string  `xml:"value"`
}

func RetrieveYahooMatchup(accessToken string) ([]byte, error) {
    config, configErr := config.ReadConfig(CONFIG_FILE_PATH)
    if configErr != nil {
        return nil, configErr
    }

    yahooAPIURL := config.YahooLeagueURL + "/" + config.YahooYearID + ".l." + config.YahooLeagueID + "/scoreboard"
    yahooAPIBody, yahooAPIErr := GetAPI(yahooAPIURL, accessToken)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error requesting GET API: %v", yahooAPIErr)
    }

    var fc FantasyMatchupContent
    xmlErr := xml.Unmarshal(yahooAPIBody, &fc)
    if xmlErr != nil {
        return nil, fmt.Errorf("Error while parsing XML: %v", xmlErr)
    }

    resultArray := make([]map[string]interface{}, 0)
    for _, matchup := range fc.League.Scoreboard.Matchups {
        for _, team := range matchup.MatchupTeams {
            statsList := make([]map[string]interface{}, 0)
            for _, stat := range team.TeamStats.Stats {
                statsMap := make(map[string]interface{})
                statsMap["StatID"] = stat.StatID
                statsMap["StatValue"] = stat.StatValue
                statsList = append(statsList, statsMap)
            }
            teamMap := make(map[string]interface{})
            teamMap["MatchupTeam"] = team.Name
            teamMap["Stats"] = statsList
            resultArray = append(resultArray, teamMap)
        }
    }

    // Create the desired JSON structure
    resultJSON := map[string]interface{}{
        "Matchup":       resultArray,
    }

    // Convert the JSON to a formatted string
    jsonBytes, jsonErr := json.MarshalIndent(resultJSON, "", "  ")
    if jsonErr != nil {
        return nil, fmt.Errorf("Error while converting to JSON: %v", jsonErr)
    }

    return jsonBytes, nil
}





