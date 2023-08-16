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
  const [projectionData, setProjectionData] = useState(null);
  const [gameData, setGameData] = useState(null);
  const [matchupData, setMatchupData] = useState(null);

  useEffect(() => {
    fetch('/yahooRosters')
      .then((response) => response.json())
      .then((data) => {
        setData(data);
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
        setGameData(data);
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
      {matchupData ? (
        <div>
          <h2>Matchup Data</h2>
          <table className="bordered-table">
            <thead className="header-row">
              <tr>
                <th className="bold centered">Team</th>
                {matchupData[0]?.Matchup[0]?.Stats.map((stat, index) => (
                  <th className="bold centered" key={index}>{statIdToLabel[stat.StatID]}</th>
                ))}
              </tr>
            </thead>
            <tbody>
                {matchupData[0]?.Matchup.map((matchup, index) => (
                  <tr key={index}>
                    <td className="bold centered">{matchup.MatchupTeam}</td>
                    {matchup.Stats.map((stat, statIndex) => {
                      const statValue = stat.StatValue;
                      return (
                        <td className="bold centered" key={statIndex}>{statValue}</td>
                      );
                    })}
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      ) : null}
    </div>
    </>
  );
}

export default App;
