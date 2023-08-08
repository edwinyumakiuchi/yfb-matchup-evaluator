package util

import (
    "fmt"
    "io"
    "net/http"
	"encoding/json"
	"encoding/xml"

	"yfb-matchup-evaluator/config"
)

type FantasyContent struct {
    XMLName xml.Name `xml:"fantasy_content"`
    Team    Team     `xml:"team"`
}

type Team struct {
    XMLName xml.Name `xml:"team"`
    Name    string   `xml:"name"`
    Roster  Roster   `xml:"roster"`
}

type Roster struct {
    XMLName xml.Name  `xml:"roster"`
    Players []Player  `xml:"players>player"`
}

type Player struct {
    XMLName xml.Name `xml:"player"`
    Name    string   `xml:"name>full"`
    TeamAbbr string   `xml:"editorial_team_abbr"`
    Position string   `xml:"display_position"`
}

func RetrieveYahooRoster(accessToken string) ([]byte, error) {
    config, configErr := config.ReadConfig(CONFIG_FILE_PATH)
    if configErr != nil {
        return nil, configErr
    }

    yahooAPIURL := config.YahooTeamURL + "/" + config.YahooYearID + ".l." + config.YahooLeagueID + ".t." + config.YahooTeamID + "/roster"
    yahooAPIReq, yahooAPIErr := http.NewRequest("GET", yahooAPIURL, nil)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error creating request: %v", yahooAPIErr)
    }
    yahooAPIReq.Header.Set("Authorization", "Bearer " + accessToken)

    client := &http.Client{}
    yahooAPIResp, yahooAPIErr := client.Do(yahooAPIReq)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error making request: %v", yahooAPIErr)
    }
    defer yahooAPIResp.Body.Close()

    if yahooAPIResp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API responded with status code %d", yahooAPIResp.StatusCode)
    }

    yahooAPIBody, yahooAPIErr := io.ReadAll(yahooAPIResp.Body)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error reading response body: %v", yahooAPIErr)
    }

    return yahooAPIBody, nil
}

func ParseData(apiBody []byte) ([]byte, error) {
    var fc FantasyContent
    xmlErr := xml.Unmarshal(apiBody, &fc)
    if xmlErr != nil {
        return nil, fmt.Errorf("Error while parsing XML: %v", xmlErr)
    }

    // Extract the player names and team abbreviation from the roster
    playersWithTeam := make([]map[string]string, len(fc.Team.Roster.Players))
    for i, player := range fc.Team.Roster.Players {
        playerData := map[string]string{
            "Player": player.Name,
            "Team":   player.TeamAbbr,
            "Position":   player.Position,
        }
        playersWithTeam[i] = playerData
    }

    // Create the desired JSON structure
    resultJSON := map[string]interface{}{
        "Roster":       playersWithTeam,
        "Fantasy Team": fc.Team.Name,
    }

    // Convert the JSON to a formatted string
    jsonBytes, jsonErr := json.MarshalIndent(resultJSON, "", "  ")
    if jsonErr != nil {
        return nil, fmt.Errorf("Error while converting to JSON: %v", jsonErr)
    }

    return jsonBytes, nil
}