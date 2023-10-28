package util

import (
    "fmt"
	"sort"
    "strconv"
	// "reflect"
    "encoding/json"
)

type TeamScore struct {
    Wins   float64
    Losses float64
    Ties   float64
}

func GetTeamProjection(yahooRosters []byte, hbProjections []byte) []byte {
    var rosters []map[string]interface{}
    if err := json.Unmarshal(yahooRosters, &rosters); err != nil {
        fmt.Println("Error parsing yahooRosters:", err)
        return nil
    }

    teamStats := make(map[string]map[string]float64)
    for _, roster := range rosters {
        teamName := roster["Fantasy Team"].(string)
        // Initialize teamStats with the desired keys for statistics
        if teamStats[teamName] == nil {
            teamStats[teamName] = make(map[string]float64)

            teamStats[teamName]["fieldGoalMade"] = 0
            teamStats[teamName]["fieldGoalAttempt"] = 0
            teamStats[teamName]["fieldGoal"] = 0
            teamStats[teamName]["freeThrowMade"] = 0
            teamStats[teamName]["freeThrowAttempt"] = 0
            teamStats[teamName]["freeThrow"] = 0
            teamStats[teamName]["threePointMade"] = 0
            teamStats[teamName]["points"] = 0
            teamStats[teamName]["totalRebounds"] = 0
            teamStats[teamName]["assists"] = 0
            teamStats[teamName]["steals"] = 0
            teamStats[teamName]["blocks"] = 0
            teamStats[teamName]["turnovers"] = 0
        }
    }

    // Parse the JSON data from hbProjections
    var projections []map[string]interface{}
    if err := json.Unmarshal(hbProjections, &projections); err != nil {
        fmt.Println("Error parsing hbProjections:", err)
        return nil
    }

    // Iterate through teamStats and match players with projections
    for _, roster := range rosters {
        teamName := roster["Fantasy Team"].(string)

        for _, playerData := range roster["Roster"].([]interface{}) {
            player := playerData.(map[string]interface{})
            playerName := player["Player"].(string)

            for _, projection := range projections {
                if projection["name"].(string) == playerName {

                    fieldGoalMade, _ := strconv.ParseFloat(projection["fieldGoalMadeCalculated"].(string), 64)
                    fieldGoalAttempt, _ := strconv.ParseFloat(projection["fieldGoalAttempt"].(string), 64)
                    freeThrowMade, _ := strconv.ParseFloat(projection["freeThrowMadeCalculated"].(string), 64)
                    freeThrowAttempt, _ := strconv.ParseFloat(projection["freeThrowAttempt"].(string), 64)
                    threePointMade, _ := strconv.ParseFloat(projection["threePointMade"].(string), 64)
                    points, _ := strconv.ParseFloat(projection["points"].(string), 64)
                    totalRebounds, _ := strconv.ParseFloat(projection["totalRebounds"].(string), 64)
                    assists, _ := strconv.ParseFloat(projection["assists"].(string), 64)
                    steals, _ := strconv.ParseFloat(projection["steals"].(string), 64)
                    blocks, _ := strconv.ParseFloat(projection["blocks"].(string), 64)
                    turnovers, _ := strconv.ParseFloat(projection["turnovers"].(string), 64)

                    teamStats[teamName]["fieldGoalMade"] += fieldGoalMade
                    teamStats[teamName]["fieldGoalAttempt"] += fieldGoalAttempt
                    teamStats[teamName]["freeThrowMade"] += freeThrowMade
                    teamStats[teamName]["freeThrowAttempt"] += freeThrowAttempt
                    teamStats[teamName]["threePointMade"] += threePointMade
                    teamStats[teamName]["points"] += points
                    teamStats[teamName]["totalRebounds"] += totalRebounds
                    teamStats[teamName]["assists"] += assists
                    teamStats[teamName]["steals"] += steals
                    teamStats[teamName]["blocks"] += blocks
                    teamStats[teamName]["turnovers"] += turnovers
                }
            }
        }
    }

    // Calculate team projections and convert to JSON
    var teamProjections []map[string]interface{}
    for team, stats := range teamStats {
        teamProjection := make(map[string]interface{})
        teamProjection["Fantasy Team"] = team

        totalFieldGoalMades := stats["fieldGoalMade"]
        totalFieldGoalAttempts := stats["fieldGoalAttempt"]
        totalFreeThrowMades := stats["freeThrowMade"]
        totalFreeThrowAttempts := stats["freeThrowAttempt"]
        totalThreePointMade := stats["threePointMade"]
        totalPoints := stats["points"]
        totalRebounds := stats["totalRebounds"]
        totalAssists := stats["assists"]
        totalSteals := stats["steals"]
        totalBlocks := stats["blocks"]
        totalTurnovers := stats["turnovers"]

        fieldGoalData := make(map[string]interface{})
        fieldGoalData["value"] = formatFloat(totalFieldGoalMades / totalFieldGoalAttempts, 4)
        fieldGoalData["rank"] = 4

        freeThrowData := make(map[string]interface{})
        freeThrowData["value"] = formatFloat(totalFreeThrowMades / totalFreeThrowAttempts, 4)
        freeThrowData["rank"] = 4

        threePointMadeData := make(map[string]interface{})
        threePointMadeData["value"] = formatFloat(totalThreePointMade, 1)
        threePointMadeData["rank"] = 4

        pointsData := make(map[string]interface{})
        pointsData["value"] = formatFloat(totalPoints, 1)
        pointsData["rank"] = 4

        totalReboundsData := make(map[string]interface{})
        totalReboundsData["value"] = formatFloat(totalRebounds, 1)
        totalReboundsData["rank"] = 4

        assistsData := make(map[string]interface{})
        assistsData["value"] = formatFloat(totalAssists, 1)
        assistsData["rank"] = 4

        stealsData := make(map[string]interface{})
        stealsData["value"] = formatFloat(totalSteals, 1)
        stealsData["rank"] = 4

        blocksData := make(map[string]interface{})
        blocksData["value"] = formatFloat(totalBlocks, 1)
        blocksData["rank"] = 4

        turnoversData := make(map[string]interface{})
        turnoversData["value"] = formatFloat(totalTurnovers, 1)
        turnoversData["rank"] = 4

        teamProjection["fieldGoalMade"] = formatFloat(totalFieldGoalMades, 3)
        teamProjection["fieldGoalAttempt"] = formatFloat(totalFieldGoalAttempts, 1)
        teamProjection["fieldGoal"] = fieldGoalData
        teamProjection["freeThrowMade"] = formatFloat(totalFreeThrowMades, 3)
        teamProjection["freeThrowAttempt"] = formatFloat(totalFreeThrowAttempts, 1)
        teamProjection["freeThrow"] = freeThrowData
        teamProjection["threePointMade"] = threePointMadeData
        teamProjection["points"] = pointsData
        teamProjection["totalRebounds"] = totalReboundsData
        teamProjection["assists"] = assistsData
        teamProjection["steals"] = stealsData
        teamProjection["blocks"] = blocksData
        teamProjection["turnovers"] = turnoversData

        teamProjections = append(teamProjections, teamProjection)
    }

    // Define the categories to process
    categories := []string{"fieldGoal", "freeThrow", "threePointMade", "totalRebounds", "assists", "steals", "blocks", "turnovers", "points"}

    for _, category := range categories {
        // Custom sorting function for the current category
        sort.Slice(teamProjections, func(i, j int) bool {
            categoryI := teamProjections[i][category].(map[string]interface{})
            categoryJ := teamProjections[j][category].(map[string]interface{})
            valueI := categoryI["value"].(float64)
            valueJ := categoryJ["value"].(float64)

            if category == "turnovers" {
                return valueJ > valueI
            }
            return valueI > valueJ
        })

        // Assign ranks based on the sorted order for the current category
        for i := range teamProjections {
            teamProjections[i][category].(map[string]interface{})["rank"] = i + 1
        }
    }

    // Convert to JSON and print or use as needed
    teamProjectionJSON, _ := json.Marshal(teamProjections)
    // fmt.Println("Team Projections:", string(teamProjectionJSON))

    return teamProjectionJSON
}

// Calculate the value for a player compared to another player
func CalculateScore(team, opponentTeam map[string]interface{}) (float64, float64, float64) {
	wins := 0.0
	losses := 0.0
	ties := 0.0

	for category := range team {
		if category == "Fantasy Team" {
			continue
		}

		switch categoryStat := team[category].(type) {
		case map[string]interface{}:
			if categoryValue, ok := categoryStat["value"].(float64); ok {
				opponentCategoryStat := opponentTeam[category].(map[string]interface{})["value"].(float64)
				if categoryValue > opponentCategoryStat {
				    if category != "turnovers" {
				        wins++
				    } else {
				        losses++
				    }
				} else if categoryValue < opponentCategoryStat {
				    if category != "turnovers" {
				        losses++
				    } else {
				        wins++
				    }
				} else {
				    ties++
				}
			}
		case float64:
			// float: non-category field, skip
			continue
		}
	}
	return wins, losses, ties
}

func CalculateAverages(scoreProjectionJSON []map[string]interface{}) []map[string]interface{} {
	for _, entry := range scoreProjectionJSON {
		// Create a map for "Average" data
		averageData := make(map[string]float64)

		// Initialize count for players
		count := 0

		// Iterate through the players and accumulate wins, losses, and ties
		for opponentTeam, score := range entry {
			if opponentTeam != "Fantasy Team" {
				if playerData, isPlayerData := score.(TeamScore); isPlayerData {
					win := playerData.Wins
					loss := playerData.Losses
					tie := playerData.Ties

					averageData["Wins"] += win
					averageData["Losses"] += loss
					averageData["Ties"] += tie

					count++
				}
			}
		}

		// Calculate averages
		averageData["Wins"] /= float64(count)
		averageData["Wins"] = float64(int(averageData["Wins"]*1000)) / 1000
		averageData["Losses"] /= float64(count)
		averageData["Losses"] = float64(int(averageData["Losses"]*1000)) / 1000
		averageData["Ties"] /= float64(count)
		averageData["Ties"] = float64(int(averageData["Ties"]*1000)) / 1000
		averageData["TotalScore"] = averageData["Wins"] + (averageData["Ties"] * 0.5)

		// Add the "Average" data to the entry
		entry["Average"] = averageData
	}

	sort.Slice(scoreProjectionJSON, func(i, j int) bool {
		averageI := scoreProjectionJSON[i]["Average"].(map[string]float64)
		averageJ := scoreProjectionJSON[j]["Average"].(map[string]float64)
		return averageI["TotalScore"] > averageJ["TotalScore"]
	})

	return scoreProjectionJSON
}

func formatFloat(value float64, decimals int) float64 {
    format := fmt.Sprintf("%%.%df", decimals)
    formattedValue, _ := strconv.ParseFloat(fmt.Sprintf(format, value), 64)
    return formattedValue
}