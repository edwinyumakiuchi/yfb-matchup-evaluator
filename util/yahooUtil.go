package util

import (
    "fmt"
    "strconv"
	"encoding/json"
	"encoding/xml"

	"yfb-matchup-evaluator/config"
)

/* type FantasyYearContent struct {
    XMLName xml.Name `xml:"fantasy_content"`
    YearGame    YearGame     `xml:"game"`
}

type YearGame struct {
    XMLName xml.Name `xml:"game"`
    YearGameKey    int     `xml:"game_key"`
} */

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
    Name           string          `xml:"name"`
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

    resultJSON := map[string]interface{}{
        "Roster":       playersWithTeam,
        "Fantasy Team": fc.Team.Name,
        "isSelfTeam": isSelfTeam,
    }

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

    /* yahooLeagueAPIURL := "https://fantasysports.yahooapis.com/fantasy/v2/game/nba"
    fmt.Println("yahooLeagueAPIURL: ", yahooLeagueAPIURL)
    yahooLeagueAPIBody, yahooLeagueAPIErr := GetAPI(yahooLeagueAPIURL, accessToken)
    if yahooLeagueAPIErr != nil {
        return nil, fmt.Errorf("Error requesting GET API: %v", yahooLeagueAPIErr)
    }
    var yfc FantasyYearContent
    xmlYearErr := xml.Unmarshal(yahooLeagueAPIBody, &yfc)
    if xmlYearErr != nil {
        return nil, fmt.Errorf("Error while parsing XML: %v", xmlYearErr)
    }
    fmt.Println("yfc: ", yfc) */

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

    // DELETEME
    // fc.DraftLeague.DraftResults = fc.DraftLeague.DraftResults[:len(fc.DraftLeague.DraftResults)-154]
    // fmt.Println("Modified fc:", fc)

    draftData := make(map[int]map[string]string)
    for i, draftResult := range fc.DraftLeague.DraftResults {
        if draftResult.Cost != 0 {
            data := map[string]string{
                "Cost":    strconv.Itoa(draftResult.Cost),
                "TeamKey": draftResult.TeamKey,
            }
            draftData[i] = data
        }
    }
    if len(draftData) == 0 {
        return nil, fmt.Errorf("Draft has not started.")
    }

	var aggregatedData = make(map[string]map[string]interface{})
	leagueSize, _ := strconv.Atoi(config.LeagueSize)
    numPlayerPerTeam, _ := strconv.Atoi(config.RosterSize)
	totalCost := 0
	// total number of players remaining to be drafted
	totalPlayerSlotRemaining := leagueSize * numPlayerPerTeam

	for _, data := range draftData {
		costStr := data["Cost"]
		teamKey := data["TeamKey"]

		cost, err := strconv.Atoi(costStr)
		if err != nil {
			fmt.Println("Error converting cost:", err)
			continue
		}
		totalCost += cost
		totalPlayerSlotRemaining--

		// Check if the TeamKey already exists
		if _, exists := aggregatedData[teamKey]; exists {
			// If exists, add to the existing cost
			existingCost, _ := aggregatedData[teamKey]["Budget"].(int)
			aggregatedData[teamKey]["Budget"] = existingCost + cost
			aggregatedData[teamKey]["NumPlayerLeft"] = aggregatedData[teamKey]["NumPlayerLeft"].(int) - 1
		} else {
			// If not exists, create a new entry
            aggregatedData[teamKey] = map[string]interface{}{
                "Budget":       cost,  // Budget spent
                "TeamKey":      teamKey,
                "NumPlayerLeft": numPlayerPerTeam - 1,
            }
		}
	}

	for i := 1; i <= leagueSize; i++ {
		teamKey := fmt.Sprintf(config.YahooYearID + ".l." + config.YahooLeagueID + ".t.%d", i)

		// Check if the team key exists
		if _, exists := aggregatedData[teamKey]; !exists {
			// If it does not exist, create a new entry
			aggregatedData[teamKey] = map[string]interface{}{
				"Budget":        0,
				"NumPlayerLeft": numPlayerPerTeam,
				"TeamKey":      teamKey,
			}
		}
	}

	if (totalPlayerSlotRemaining == 0) {
	    return nil, fmt.Errorf("Remaining player slot is 0, draft has ended!")
	}

	for teamKey, budget := range aggregatedData {
	    if (budget["NumPlayerLeft"].(int) == 0) {
	        totalCost -= budget["Budget"].(int)
	        continue
	    }

		ownCost, _ := budget["Budget"].(int)

		// Calculate total cost excluding the current team
		otherTotal := (200 * (leagueSize - 1)) - (totalCost - ownCost)
		// Total number of players remaining to be drafted for the rest of the league
		otherPlayerSlotRemaining := totalPlayerSlotRemaining - budget["NumPlayerLeft"].(int)
        if otherPlayerSlotRemaining == 0 {
            return nil, fmt.Errorf("Only one team remaining to draft.")
        }

		var avgCost float64
		avgCost = float64(otherTotal) / float64(otherPlayerSlotRemaining)

		// Average cost per player for the rest of the league
		budget["AvgCost"] = fmt.Sprintf("%.2f", avgCost)
		aggregatedData[teamKey] = budget
	}

	teamBudgets := make([]map[string]interface{}, 0)
	for _, teamData := range aggregatedData {
	    if (teamData["NumPlayerLeft"].(int) != 0) {
		    teamBudgets = append(teamBudgets, teamData)
		}
	}

	for _, budget := range teamBudgets {
		cost, _ := budget["Budget"]
		ownPlayerLeft, _ := budget["NumPlayerLeft"]

        // Average cost per player for self
		var avgSelfCost float64
		avgSelfCost = (200 - float64(cost.(int))) / float64(ownPlayerLeft.(int))

		budget["SelfCost"] = avgSelfCost
	}

	finalResult := map[string]interface{}{
		"budgets": teamBudgets,
	}

	budgetJsonData, err := json.Marshal(finalResult)
	if err != nil {
		return nil, fmt.Errorf("Error converting to JSON:", err)
	}

    return budgetJsonData, nil
}