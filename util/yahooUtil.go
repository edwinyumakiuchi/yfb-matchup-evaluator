package util

import (
    "fmt"
    "io"
    "net/http"
	"encoding/json"
	"encoding/xml"
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

func GetAPIData(accessToken string) ([]byte, error) {
    apiURL := "https://fantasysports.yahooapis.com/fantasy/v2/team/418.l.33024.t.4/roster"
    apiReq, apiErr := http.NewRequest("GET", apiURL, nil)
    if apiErr != nil {
        return nil, fmt.Errorf("Error creating request: %v", apiErr)
    }
    apiReq.Header.Set("Authorization", "Bearer "+accessToken)

    client := &http.Client{}
    apiResp, apiErr := client.Do(apiReq)
    if apiErr != nil {
        return nil, fmt.Errorf("Error making request: %v", apiErr)
    }
    defer apiResp.Body.Close()

    if apiResp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API responded with status code %d", apiResp.StatusCode)
    }

    apiBody, apiErr := io.ReadAll(apiResp.Body)
    if apiErr != nil {
        return nil, fmt.Errorf("Error reading response body: %v", apiErr)
    }

    return apiBody, nil
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