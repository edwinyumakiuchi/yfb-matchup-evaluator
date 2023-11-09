import './App.css';
import React from 'react';
import { useState, useEffect } from 'react';

// npm run build && npm run backend
function App() {
  const urlParams = new URLSearchParams(window.location.search);
  const loggedIn = urlParams.get('loggedIn') === 'true';

  if (loggedIn === true) {
    return (<Home />);
  } else {
    return (<LogIn />);
  }
}

function LogIn() {
  const handleYahooLogin = () => {
    // Redirect to the Yahoo authentication endpoint provided by main.go
    window.location.href = 'http://localhost:3000/auth/yahoo';
  };

  return (
    <div className="App">
      <header className="App-header">
        <p>Yahoo Fantasy Basketball: matchup evaluator</p>
        <button onClick={handleYahooLogin}>Login with Yahoo</button>
      </header>
    </div>
  );
}

function Home() {
  const [data, setData] = useState(null);
  const [teams, setTeams] = useState(null);
  const [selfTeamName, setSelfTeamName] = useState(null);
  const [projectionData, setProjectionData] = useState(null);
  const [gameData, setGameData] = useState(null);
  const [matchupData, setMatchupData] = useState(null);
  const [seasonOutlookData, setSeasonOutlookData] = useState(null);
  const [teamRankingData, setTeamRankingData] = useState(null);

  useEffect(() => {
    fetch('/yahooRosters')
      .then((response) => response.json())
      .then((data) => {
        setTeams(data);
        for (let i = 0; i < data.length; i++) {
            const team = data[i];
            if (team.hasOwnProperty("isSelfTeam") && team["isSelfTeam"] === true) {
                setData([team]);
                setSelfTeamName(team["Fantasy Team"]);
                break; // Found the team, exit the loop
            }
        }
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/hbProjections')
      .then((response) => response.json())
      .then((data) => {
        setProjectionData(data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/schedule')
      .then((response) => response.json())
      .then((data) => {
        setGameData(data[0]);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/yahooMatchup')
      .then((response) => response.json())
      .then((data) => {
        setMatchupData(data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/seasonOutlook')
      .then((response) => response.json())
      .then((data) => {
        setSeasonOutlookData(data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/scoreProjection')
      .then((response) => response.json())
      .then((data) => {
        setTeamRankingData(data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  const statIdToLabel = {
    "9004003": "FGM/FGA",
    "5": "FG",
    "9007006": "FTM/FTA",
    "8": "FT",
    "10": "3PM",
    "12": "PTS",
    "15": "REB",
    "16": "AST",
    "17": "ST",
    "18": "BLK",
    "19": "TO"
  };

  return (
    <>
    <div style={{ marginLeft: '20px' }}>
      {data && projectionData && gameData && matchupData && selfTeamName ? (
        <div>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">DATE</th>
                <th className="bold centered">PLAYER</th>
                <th className="bold centered">TEAM</th>
                <th className="bold centered">POS</th>
                <th className="bold centered">FGM</th>
                <th className="bold centered">FGA</th>
                <th className="bold centered">FG%</th>
                <th className="bold centered">FTM</th>
                <th className="bold centered">FTA</th>
                <th className="bold centered">FT%</th>
                <th className="bold centered">3PTM</th>
                <th className="bold centered">PTS</th>
                <th className="bold centered">REB</th>
                <th className="bold centered">AST</th>
                <th className="bold centered">ST</th>
                <th className="bold centered">BLK</th>
                <th className="bold centered">TO</th>
              </tr>
            </thead>
            <tbody>
              {gameData.data.map((game) => {
                const playersForDate = data
                  .flatMap((player) => player.Roster)
                  .filter((playerData) => game.teams.includes(playerData.Team));

                return (
                  <React.Fragment key={game.date}>
                    <tr>
                      <td className="bold centered" rowSpan={playersForDate.length + 1}>
                        {game.date}
                      </td>
                    </tr>

                    {playersForDate.map((playerData, index) => {
                      const matchedPlayer = projectionData.find((player) => player.name === playerData.Player);

                      return (
                        <tr key={index}>
                          <td>{playerData.Player}</td>
                          <td className="bold centered">{playerData.Team}</td>
                          <td className="bold centered">{playerData.Position}</td>
                          {matchedPlayer ? (
                            <>
                              <td className="bold centered">{matchedPlayer.fieldGoalMadeCalculated}</td>
                              <td className="bold centered">{matchedPlayer.fieldGoalAttempt}</td>
                              <td className="bold centered">{matchedPlayer.fieldGoal}</td>
                              <td className="bold centered">{matchedPlayer.freeThrowMadeCalculated}</td>
                              <td className="bold centered">{matchedPlayer.freeThrowAttempt}</td>
                              <td className="bold centered">{matchedPlayer.freeThrow}</td>
                              <td className="bold centered">{matchedPlayer.threePointMade}</td>
                              <td className="bold centered">{matchedPlayer.points}</td>
                              <td className="bold centered">{matchedPlayer.totalRebounds}</td>
                              <td className="bold centered">{matchedPlayer.assists}</td>
                              <td className="bold centered">{matchedPlayer.steals}</td>
                              <td className="bold centered">{matchedPlayer.blocks}</td>
                              <td className="bold centered">{matchedPlayer.turnovers}</td>
                            </>
                          ) : (
                            // Render cells with 0 when matchedPlayer is not found
                            <>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                              <td className="bold centered">0</td>
                            </>
                          )}
                        </tr>
                      );
                    })}
                  </React.Fragment>
                );
              })}

              {/* TOTAL row */}
              {/* <tr>
                <td className="bold centered">TOTAL</td>
                <td className="bold centered"></td>
                <td className="bold centered"></td>
                <td className="bold centered"></td>
                <td className="bold centered">{Number(calculateField('fieldGoalMadeCalculated')).toFixed(3)}</td>
                <td className="bold centered">{Number(calculateField('fieldGoalAttempt')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateFieldGoalPercentage()/100).toFixed(5)}</td>
                <td className="bold centered">{Number(calculateField('freeThrowMadeCalculated')).toFixed(3)}</td>
                <td className="bold centered">{Number(calculateField('freeThrowAttempt')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateFreeThrowPercentage()/100).toFixed(5)}</td>
                <td className="bold centered">{Number(calculateField('threePointMade')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('points')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('totalRebounds')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('assists')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('steals')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('blocks')).toFixed(1)}</td>
                <td className="bold centered">{Number(calculateField('turnovers')).toFixed(1)}</td>
              </tr> */}

              <tr>
                <td>Current matchup</td>
                <td></td>
                <td></td>
                <td></td>
                {matchupData[0]?.Matchup.map((matchup, matchupIndex) => {
                  if (matchup.MatchupTeam === selfTeamName) {
                    const statCells = matchup.Stats.map((stat, index) => {
                      // Check if StatValue is in the format "x/y"
                      const statParts = stat.StatValue.split('/');

                      if (statParts.length === 2) {
                        // If it's in the "x/y" format, create separate <td> for "x" and "y"
                        return (
                          <React.Fragment key={index}>
                            <td className="bold centered">{statParts[0]}</td>
                            <td className="bold centered">{statParts[1]}</td>
                          </React.Fragment>
                        );
                      } else if (!stat.StatValue.includes('.')) {
                        // If not in "x/y" format, create a single <td> with the original value
                        return <td className="bold centered" key={index}>{stat.StatValue}</td>;
                      } else {
                        return <td></td>;
                      }
                    });

                    return (
                      <React.Fragment key={matchupIndex}>
                        {statCells}
                      </React.Fragment>
                    );
                  }
                })}
              </tr>

            </tbody>
          </table>
        </div>
      ) : (
        <>
          <p>Failed to fetch data from the backend.</p>
        </>
      )}
    </div>

    <div style={{ marginLeft: '20px' }}>
      {matchupData && selfTeamName ? (
        <div>
          <h2>Matchup Data</h2>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">Team</th>
                {matchupData[0]?.Matchup[0]?.Stats.map((stat, index) => (
                  <th className="bold centered" key={index}>{statIdToLabel[stat.StatID]}</th>
                ))}
                <th className="bold centered">Score</th>
              </tr>
            </thead>
            <tbody>
              {(() => {
                const teamStats = {}; // Object to store win and loss counts for each team

                // Function to calculate ratio and update win-loss records
                const calculateRatioAndUpdateRecords = (team, selfStat, winLossFieldIndex) => {
                  const teamStat = team.Stats[winLossFieldIndex]?.StatValue;
                  const [teamValue, teamTotal] = teamStat.split('/').map(Number);
                  const [selfValue, selfTotal] = selfStat.split('/').map(Number);
                  const teamRatio = teamValue / teamTotal;
                  const selfRatio = selfValue / selfTotal;

                  if (!(team.MatchupTeam in teamStats)) {
                    teamStats[team.MatchupTeam] = {
                      winNumber: 0,
                      lossNumber: 0,
                      tieNumber: 0,
                    };
                  }

                  if (selfRatio > teamRatio) {
                    teamStats[team.MatchupTeam].winNumber++;
                  } else if (selfRatio < teamRatio) {
                    teamStats[team.MatchupTeam].lossNumber++;
                  } else {
                    teamStats[team.MatchupTeam].tieNumber++;
                  }
                };

                // Function to calculate non-ratio category and update win-loss records
                const calculateAndUpdateRecords = (team, selfStat, winLossFieldIndex) => {
                  const teamStat = team.Stats[winLossFieldIndex]?.StatValue;

                  if (!(team.MatchupTeam in teamStats)) {
                    teamStats[team.MatchupTeam] = {
                      winNumber: 0,
                      lossNumber: 0,
                      tieNumber: 0,
                    };
                  }

                  if (parseInt(selfStat) > parseInt(teamStat)) {
                    teamStats[team.MatchupTeam].winNumber++;
                  } else if (parseInt(selfStat) < parseInt(teamStat)) {
                    teamStats[team.MatchupTeam].lossNumber++;
                  } else {
                    teamStats[team.MatchupTeam].tieNumber++;
                  }
                };

                // Function to calculate turnover and update win-loss records
                const calculateToAndUpdateRecords = (team, selfStat, winLossFieldIndex) => {
                  const teamStat = team.Stats[winLossFieldIndex]?.StatValue;

                  if (!(team.MatchupTeam in teamStats)) {
                    teamStats[team.MatchupTeam] = {
                      winNumber: 0,
                      lossNumber: 0,
                      tieNumber: 0,
                    };
                  }

                  if (parseInt(selfStat) < parseInt(teamStat)) {
                    teamStats[team.MatchupTeam].winNumber++;
                  } else if (parseInt(selfStat) > parseInt(teamStat)) {
                    teamStats[team.MatchupTeam].lossNumber++;
                  } else {
                    teamStats[team.MatchupTeam].tieNumber++;
                  }
                };

                matchupData[0]?.Matchup.forEach((team) => {
                  if (team.MatchupTeam !== selfTeamName) {
                    const selfFg = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[0]?.StatValue;
                    const selfFt = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[2]?.StatValue;
                    const selfTpm = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[4]?.StatValue;
                    const selfPts = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[5]?.StatValue;
                    const selfReb = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[6]?.StatValue;
                    const selfAst = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[7]?.StatValue;
                    const selfStl = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[8]?.StatValue;
                    const selfBlk = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[9]?.StatValue;
                    const selfTo = matchupData[0]?.Matchup.find(t => t.MatchupTeam === selfTeamName)?.Stats[10]?.StatValue;

                    calculateRatioAndUpdateRecords(team, selfFg, 0); // FGM/FGA
                    calculateRatioAndUpdateRecords(team, selfFt, 2); // FTM/FTA

                    calculateAndUpdateRecords(team, selfTpm, 4); // 3PM
                    calculateAndUpdateRecords(team, selfPts, 5); // PTS
                    calculateAndUpdateRecords(team, selfReb, 6); // REB
                    calculateAndUpdateRecords(team, selfAst, 7); // AST
                    calculateAndUpdateRecords(team, selfStl, 8); // STL
                    calculateAndUpdateRecords(team, selfBlk, 9); // BLK

                    calculateToAndUpdateRecords(team, selfTo, 10); // TO
                  }
                });

                return matchupData[0]?.Matchup.map((matchup, index) => (
                  <tr key={index}>
                    <td className="bold centered">{matchup.MatchupTeam}</td>
                    {matchup.Stats.map((stat, statIndex) => (
                      <td className="bold centered" key={statIndex}>{stat.StatValue}</td>
                    ))}
                    <td className="bold centered">
                      {matchup.MatchupTeam === selfTeamName ? (
                        ''
                      ) : (
                        `${teamStats[matchup.MatchupTeam].winNumber}-${teamStats[matchup.MatchupTeam].lossNumber}-${teamStats[matchup.MatchupTeam].tieNumber}`
                      )}
                    </td>
                  </tr>
                ));
              })()}
            </tbody>
          </table>
        </div>
      ) : null}
    </div>

    <div style={{ marginLeft: '20px' }}>
      {data && projectionData ? (
        <div>
          <h2>Season Outlook</h2>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">FANTASY TEAM</th>
                <th className="bold centered">FGM</th>
                <th className="bold centered">FGA</th>
                <th className="bold centered">FG%</th>
                <th className="bold centered">FTM</th>
                <th className="bold centered">FTA</th>
                <th className="bold centered">FT%</th>
                <th className="bold centered">3PTM</th>
                <th className="bold centered">PTS</th>
                <th className="bold centered">REB</th>
                <th className="bold centered">AST</th>
                <th className="bold centered">ST</th>
                <th className="bold centered">BLK</th>
                <th className="bold centered">TO</th>
                <th className="bold centered">SCORE</th>
              </tr>
            </thead>
            <tbody>
              {teams.map((team, teamIndex) => {
                const teamPlayers = team.Roster;
                const teamName = team['Fantasy Team'];

                const fieldGoalMadeSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.fieldGoalMadeCalculated) : sum;
                }, 0);
                // const fieldGoalMadeAverage = fieldGoalMadeSum / teamPlayers.length;

                const fieldGoalAttemptSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.fieldGoalAttempt) : sum;
                }, 0);
                // const fieldGoalAttemptAverage = fieldGoalAttemptSum / teamPlayers.length;

                const fieldGoal = fieldGoalMadeSum / fieldGoalAttemptSum

                const freeThrowMadeSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.freeThrowMadeCalculated) : sum;
                }, 0);
                // const freeThrowMadeAverage = freeThrowMadeSum / teamPlayers.length;

                const freeThrowAttemptSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.freeThrowAttempt) : sum;
                }, 0);
                // const freeThrowAttemptAverage = freeThrowAttemptSum / teamPlayers.length;

                const freeThrow = freeThrowMadeSum / freeThrowAttemptSum

                const threePointMadeSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.threePointMade) : sum;
                }, 0);
                // const threePointMadeAverage = threePointMadeSum / teamPlayers.length;

                const pointsSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.points) : sum;
                }, 0);
                // const pointsAverage = pointsSum / teamPlayers.length;

                const totalReboundsSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.totalRebounds) : sum;
                }, 0);
                // const totalReboundsAverage = totalReboundsSum / teamPlayers.length;

                const assistsSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.assists) : sum;
                }, 0);
                // const assistsAverage = assistsSum / teamPlayers.length;

                const stealsSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.steals) : sum;
                }, 0);
                // const stealsAverage = stealsSum / teamPlayers.length;

                const blocksSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.blocks) : sum;
                }, 0);
                // const blocksAverage = blocksSum / teamPlayers.length;

                const turnoversSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.turnovers) : sum;
                }, 0);
                // const turnoversAverage = turnoversSum / teamPlayers.length;

                let wins = 0;
                let losses = 0;
                let ties = 0;
                let score;

                if (teamName !== selfTeamName) {
                  const selfTeam = teams.find(team => team['Fantasy Team'] === selfTeamName);

                  if (selfTeam) {
                    // Calculate FG% for the self team
                    const selfFieldGoalMadeSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.fieldGoalMadeCalculated) : sum;
                    }, 0);
                    // const selfFieldGoalMadeAverage = selfFieldGoalMadeSum / selfTeam.Roster.length;

                    const selfFieldGoalAttemptSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.fieldGoalAttempt) : sum;
                    }, 0);
                    // const selfFieldGoalAttemptAverage = selfFieldGoalAttemptSum / selfTeam.Roster.length;

                    const selfFieldGoal = selfFieldGoalMadeSum / selfFieldGoalAttemptSum;

                    if (selfFieldGoal > fieldGoal) {
                      wins++;
                    } else if (selfFieldGoal < fieldGoal) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfFreeThrowMadeSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.freeThrowMadeCalculated) : sum;
                    }, 0);
                    // const selfFreeThrowMadeAverage = selfFreeThrowMadeSum / selfTeam.Roster.length;

                    const selfFreeThrowAttemptSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.freeThrowAttempt) : sum;
                    }, 0);
                    // const selfFreeThrowAttemptAverage = selfFreeThrowAttemptSum / selfTeam.Roster.length;

                    const selfFreeThrow = selfFreeThrowMadeSum / selfFreeThrowAttemptSum;

                    if (selfFreeThrow > freeThrow) {
                      wins++;
                    } else if (selfFreeThrow < freeThrow) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfThreePointMadeSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.threePointMade) : sum;
                    }, 0);
                    // const selfThreePointMadeAverage = selfThreePointMadeSum / selfTeam.Roster.length;

                    if (selfThreePointMadeSum > threePointMadeSum) {
                      wins++;
                    } else if (selfThreePointMadeSum < threePointMadeSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfPointsSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.points) : sum;
                    }, 0);
                    // const selfPointsAverage = selfPointsSum / selfTeam.Roster.length;

                    if (selfPointsSum > pointsSum) {
                      wins++;
                    } else if (selfPointsSum < pointsSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfReboundsSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.totalRebounds) : sum;
                    }, 0);
                    // const selfReboundsAverage = selfReboundsSum / selfTeam.Roster.length;

                    if (selfReboundsSum > totalReboundsSum) {
                      wins++;
                    } else if (selfReboundsSum < totalReboundsSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfAssistsSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.assists) : sum;
                    }, 0);
                    // const selfAssistsAverage = selfAssistsSum / selfTeam.Roster.length;

                    if (selfAssistsSum > assistsSum) {
                      wins++;
                    } else if (selfAssistsSum < assistsSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfStealsSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.steals) : sum;
                    }, 0);
                    // const selfStealsAverage = selfStealsSum / selfTeam.Roster.length;

                    if (selfStealsSum > stealsSum) {
                      wins++;
                    } else if (selfStealsSum < stealsSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfBlocksSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.blocks) : sum;
                    }, 0);
                    // const selfBlocksAverage = selfBlocksSum / selfTeam.Roster.length;

                    if (selfBlocksSum > blocksSum) {
                      wins++;
                    } else if (selfBlocksSum < blocksSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    const selfTurnoversSum = selfTeam.Roster.reduce((sum, playerData) => {
                      const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                      return matchingProjection ? sum + parseFloat(matchingProjection.turnovers) : sum;
                    }, 0);
                    // const selfTurnoversAverage = selfTurnoversSum / selfTeam.Roster.length;

                    if (selfTurnoversSum < turnoversSum) {
                      wins++;
                    } else if (selfTurnoversSum > turnoversSum) {
                      losses++;
                    } else {
                      ties++;
                    }

                    score = `${wins}-${losses}-${ties}`;
                  }
                }

                return (
                  <React.Fragment key={teamIndex}>
                    <tr>
                      <td className="bold centered">{teamName}</td>
                      <td className="bold centered">{fieldGoalMadeSum.toFixed(3)}</td>
                      <td className="bold centered">{fieldGoalAttemptSum.toFixed(1)}</td>
                      <td className="bold centered">{fieldGoal.toFixed(4)}</td>
                      <td className="bold centered">{freeThrowMadeSum.toFixed(3)}</td>
                      <td className="bold centered">{freeThrowAttemptSum.toFixed(1)}</td>
                      <td className="bold centered">{freeThrow.toFixed(4)}</td>
                      <td className="bold centered">{threePointMadeSum.toFixed(1)}</td>
                      <td className="bold centered">{pointsSum.toFixed(1)}</td>
                      <td className="bold centered">{totalReboundsSum.toFixed(1)}</td>
                      <td className="bold centered">{assistsSum.toFixed(1)}</td>
                      <td className="bold centered">{stealsSum.toFixed(1)}</td>
                      <td className="bold centered">{blocksSum.toFixed(1)}</td>
                      <td className="bold centered">{turnoversSum.toFixed(1)}</td>
                      <td className="bold centered">{score}</td>
                    </tr>
                  </React.Fragment>
                );
              })}

            </tbody>
          </table>
        </div>
      ) : (
        <>
          <p>Failed to fetch data from the backend.</p>
        </>
      )}
    </div>

    <div style={{ marginLeft: '20px' }}>
      {seasonOutlookData ? (
        <div>
          <h2>Ranking per Category</h2>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">FANTASY TEAM</th>
                <th className="bold centered">FG%</th>
                <th className="bold centered">FT%</th>
                <th className="bold centered">3PTM</th>
                <th className="bold centered">PTS</th>
                <th className="bold centered">REB</th>
                <th className="bold centered">AST</th>
                <th className="bold centered">ST</th>
                <th className="bold centered">BLK</th>
                <th className="bold centered">TO</th>
              </tr>
            </thead>
            <tbody>
              {teams.map((team, teamIndex) => {
                const teamName = team["Fantasy Team"];

                // Find the corresponding element in seasonOutlookData
                const matchingTeamData = seasonOutlookData.find(
                  (data) => data["Fantasy Team"] === teamName
                );

                const fieldGoalRank = matchingTeamData ? matchingTeamData["fieldGoal"]["rank"] : '';
                const freeThrowRank = matchingTeamData ? matchingTeamData["freeThrow"]["rank"] : '';
                const threePointMadeRank = matchingTeamData ? matchingTeamData["threePointMade"]["rank"] : '';
                const pointsRank = matchingTeamData ? matchingTeamData["points"]["rank"] : '';
                const totalReboundsRank = matchingTeamData ? matchingTeamData["totalRebounds"]["rank"] : '';
                const assistsRank = matchingTeamData ? matchingTeamData["assists"]["rank"] : '';
                const stealsRank = matchingTeamData ? matchingTeamData["steals"]["rank"] : '';
                const blocksRank = matchingTeamData ? matchingTeamData["blocks"]["rank"] : '';
                const turnoversRank = matchingTeamData ? matchingTeamData["turnovers"]["rank"] : '';

                return (
                  <React.Fragment key={teamIndex}>
                    <tr>
                      <td className="bold centered">{teamName}</td>
                      <td className="bold centered">{fieldGoalRank}</td>
                      <td className="bold centered">{freeThrowRank}</td>
                      <td className="bold centered">{threePointMadeRank}</td>
                      <td className="bold centered">{pointsRank}</td>
                      <td className="bold centered">{totalReboundsRank}</td>
                      <td className="bold centered">{assistsRank}</td>
                      <td className="bold centered">{stealsRank}</td>
                      <td className="bold centered">{blocksRank}</td>
                      <td className="bold centered">{turnoversRank}</td>
                    </tr>
                  </React.Fragment>
                );
              })}

            </tbody>
          </table>
        </div>
      ) : (
        <>
          <p>Failed to fetch data from the backend.</p>
        </>
      )}
    </div>

    <div style={{ marginLeft: '20px' }}>
      {teamRankingData ? (
        <div>
          <h2>Team Ranking</h2>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">FANTASY TEAM</th>
                {teamRankingData.map((teamData, index) => {
                  const teamName = teamData['Fantasy Team'];
                  return (
                    <th key={index} className="bold centered">
                      {teamName}
                    </th>
                  );
                })}
              </tr>
            </thead>
            <tbody>
              {teamRankingData.map((teamData, teamIndex) => {
                const teamName = teamData['Fantasy Team'];
                return (
                  <tr key={teamIndex}>
                    <td className="bold centered">{teamName}</td>
                    {teamRankingData.map((opponentData, opponentIndex) => {
                      const opponentName = opponentData['Fantasy Team'];
                      const wins = teamData[opponentName] ? teamData[opponentName].Wins : '-';
                      const losses = teamData[opponentName] ? teamData[opponentName].Losses : '-';
                      const ties = teamData[opponentName] ? teamData[opponentName].Ties : '-';
                      return (
                        <td key={opponentIndex} className="bold centered">
                          {teamName === opponentName ? '-' : `${losses}-${wins}-${ties}`}
                        </td>
                      );
                    })}
                  </tr>
                );
              })}
              <tr>
                <td className="bold centered">Total Score</td>
                {teamRankingData.map((teamData, index) => {
                  const averageTotalScore = teamData['Average'] ? teamData['Average'].TotalScore : '-';
                  return (
                    <td key={index} className="bold centered">
                      {averageTotalScore.toFixed(3)}
                    </td>
                  );
                })}
              </tr>
            </tbody>
          </table>
        </div>
      ) : (
        <>
          <p>Failed to fetch data from the backend.</p>
        </>
      )}
    </div>

    </>
  );
}

export default App;
