package util

import (
    "fmt"
    "strconv"
	"encoding/json"
	"encoding/xml"

	"yfb-matchup-evaluator/config"
)

type FantasyContent struct {
    XMLName xml.Name `xml:"fantasy_content"`
    Team    Team     `xml:"team"`
}

type FantasyLeagueContent struct {
    XMLName       xml.Name      `xml:"fantasy_content"`
    DraftLeague   DraftLeague   `xml:"league"`
}

type DraftLeague struct {
    XMLName        xml.Name        `xml:"league"`
    DraftResults   []DraftResult   `xml:"draft_results>draft_result"`
}

type DraftResult struct {
    XMLName       xml.Name      `xml:"draft_result"`
    Cost          int           `xml:"cost"`
    TeamKey       string        `xml:"team_key"`
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

func RetrieveYahooRoster(accessToken string, teamID int) ([]byte, error) {
    config, configErr := config.ReadConfig(CONFIG_FILE_PATH)
    if configErr != nil {
        return nil, configErr
    }

    yahooAPIURL := config.YahooTeamURL + "/" + config.YahooYearID + ".l." + config.YahooLeagueID + ".t." + strconv.Itoa(teamID) + "/roster"
    yahooAPIBody, yahooAPIErr := GetAPI(yahooAPIURL, accessToken)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error requesting GET API: %v", yahooAPIErr)
    }

    var fc FantasyContent
    xmlErr := xml.Unmarshal(yahooAPIBody, &fc)
    if xmlErr != nil {
        return nil, fmt.Errorf("Error while parsing XML: %v", xmlErr)
    }

    // Extract the player names and team abbreviation from the roster
    playersWithTeam := make([]map[string]string, len(fc.Team.Roster.Players))
    for i, player := range fc.Team.Roster.Players {
        playerData := map[string]string{
            "Player":     player.Name,
            "Team":       player.TeamAbbr,
            "Position":   player.Position,
        }
        playersWithTeam[i] = playerData
    }

    isSelfTeam := false
    if strconv.Itoa(teamID) == config.YahooTeamID {
        isSelfTeam = true
    }

    // Create the desired JSON structure
    resultJSON := map[string]interface{}{
        "Roster":       playersWithTeam,
        "Fantasy Team": fc.Team.Name,
        "isSelfTeam": isSelfTeam,
    }

    // Convert the JSON to a formatted string
    jsonBytes, jsonErr := json.MarshalIndent(resultJSON, "", "  ")
    if jsonErr != nil {
        return nil, fmt.Errorf("Error while converting to JSON: %v", jsonErr)
    }

    return jsonBytes, nil
}

func RetrieveYahooDraftResults(accessToken string) ([]byte, error) {
    config, configErr := config.ReadConfig(CONFIG_FILE_PATH)
    if configErr != nil {
        return nil, configErr
    }

    yahooAPIURL := config.YahooLeagueURL + config.YahooYearID + ".l." + config.YahooLeagueID + "/draftresults"
    yahooAPIBody, yahooAPIErr := GetAPI(yahooAPIURL, accessToken)
    if yahooAPIErr != nil {
        return nil, fmt.Errorf("Error requesting GET API: %v", yahooAPIErr)
    }

    var fc FantasyLeagueContent
    xmlErr := xml.Unmarshal(yahooAPIBody, &fc)
    if xmlErr != nil {
        return nil, fmt.Errorf("Error while parsing XML: %v", xmlErr)
    }

    draftData := make(map[int]map[string]string)
    for i, draftResult := range fc.DraftLeague.DraftResults {
        // Create a nested map for Cost and TeamKey
        data := map[string]string{
            "Cost":    strconv.Itoa(draftResult.Cost),
            "TeamKey": draftResult.TeamKey,
        }
        draftData[i] = data // Assign the nested map to the outer map using the index as key
    }

	// Resulting map to store aggregated data
	aggregatedData := make(map[string]map[string]string)
	totalCost := 0

	for _, data := range draftData {
		costStr := data["Cost"]
		teamKey := data["TeamKey"]

		// Convert cost to integer
		cost, err := strconv.Atoi(costStr)
		if err != nil {
			fmt.Println("Error converting cost:", err)
			continue
		}
		totalCost += cost

		// Check if the TeamKey already exists in the aggregatedData
		if _, exists := aggregatedData[teamKey]; exists {
			// If exists, add to the existing cost
			existingCost, _ := strconv.Atoi(aggregatedData[teamKey]["Budget"])
			aggregatedData[teamKey]["Budget"] = strconv.Itoa(existingCost + cost)
		} else {
			// If not exists, create a new entry
			aggregatedData[teamKey] = map[string]string{
				"Budget":    costStr,
				"TeamKey": teamKey,
			}
		}
	}

	for teamKey, budget := range aggregatedData {
		ownCost, _ := strconv.Atoi(budget["Budget"])

		// Calculate total cost excluding the current team
		otherTotal := totalCost - ownCost
		otherCount := len(aggregatedData) - 1 // Subtracting one if this team is included

		var avgCost float64
		if otherCount > 0 {
			avgCost = float64(otherTotal) / float64(otherCount)
		}

		// Add Avg-Cost to the current budget entry
		budget["Avg-Budget"] = fmt.Sprintf("%.2f", avgCost)
		aggregatedData[teamKey] = budget // Update the aggregated data
	}

	// Prepare the final output in the desired format
	teamBudgets := make([]map[string]string, 0)
	for _, teamData := range aggregatedData {
		teamBudgets = append(teamBudgets, teamData)
	}

	for _, budget := range teamBudgets {
		// Update Avg-Cost
		avgCost, _ := strconv.ParseFloat(budget["Avg-Budget"], 64)
		budget["Avg-Budget"] = fmt.Sprintf("%.2f", 200.0-avgCost)

		// Update Cost
		cost, _ := strconv.Atoi(budget["Budget"])
		budget["Budget"] = strconv.Itoa(200 - cost)
	}

	finalResult := map[string]interface{}{
		"budgets": teamBudgets, // Key name can be changed as needed
	}

	budgetJsonData, err := json.Marshal(finalResult)
	if err != nil {
		return nil, fmt.Errorf("Error converting to JSON:", err)
	}

    return budgetJsonData, nil
}