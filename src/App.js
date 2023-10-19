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

  const calculateField = (field) => {
    if (!gameData || !data || !projectionData) return 0;

    return gameData.data.reduce((total, game) => {
      const playersForDate = data
        .flatMap((player) => player.Roster)
        .filter((playerData) => game.teams.includes(playerData.Team));

      const sumForField = playersForDate.reduce((sum, playerData) => {
        const matchedPlayer = projectionData.find((player) => player.name === playerData.Player);
        return matchedPlayer ? sum + +matchedPlayer[field] : sum;
      }, 0);

      return total + sumForField;
    }, 0);
  };

  const calculateFieldGoalPercentage = () => {
    const fieldGoalsMade = calculateField('fieldGoalMadeCalculated');
    const fieldGoalAttempts = calculateField('fieldGoalAttempt');
    const percentage = (fieldGoalsMade / fieldGoalAttempts) * 100;
    return isNaN(percentage) ? '0.00' : percentage.toFixed(3);
  };

  const calculateFreeThrowPercentage = () => {
    const freeThrowsMade = calculateField('freeThrowMadeCalculated');
    const freeThrowAttempts = calculateField('freeThrowAttempt');
    const percentage = (freeThrowsMade / freeThrowAttempts) * 100;
    return isNaN(percentage) ? '0.00' : percentage.toFixed(3);
  };

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
      {data && projectionData ? (
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
              <tr>
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

                  if (selfStat > teamStat) {
                    teamStats[team.MatchupTeam].winNumber++;
                  } else if (selfStat < teamStat) {
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

                  if (selfStat < teamStat) {
                    teamStats[team.MatchupTeam].winNumber++;
                  } else if (selfStat > teamStat) {
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
                const fieldGoalMadeAverage = fieldGoalMadeSum / teamPlayers.length;

                const assistsSum = teamPlayers.reduce((sum, playerData) => {
                  const matchingProjection = projectionData.find(projection => projection.name === playerData.Player);
                  return matchingProjection ? sum + parseFloat(matchingProjection.assists) : sum;
                }, 0);
                const assistsAverage = assistsSum / teamPlayers.length;

                return (
                  <React.Fragment key={teamIndex}>
                    <tr>
                      <td className="bold centered">{teamName}</td>
                      <td className="bold centered">{fieldGoalMadeAverage.toFixed(3)}</td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered">{assistsAverage.toFixed(3)}</td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
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



    </>
  );
}

export default App;
