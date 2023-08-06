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

  useEffect(() => {
    fetch('/mongoData')
      .then((response) => response.json())
      .then((data) => {
        setData(data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    fetch('/mongoProjectionData')
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

  return (
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
                      // Update the player names in projectionData before searching for the matched player
                      const updatedProjectionData = projectionData.map((player) => {
                        if (player.name === "PJ Washington") {
                          return { ...player, name: "P.J. Washington" };
                        } else {
                          return player;
                        }
                      });

                      // Find the player from updatedProjectionData that matches the current "Player" in playerData
                      const matchedPlayer = updatedProjectionData.find((player) => player.name === playerData.Player);

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
                            // Render empty cells when matchedPlayer is not found
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
            </tbody>
          </table>
        </div>
      ) : (
        <>
          <p>Failed to fetch data from the backend.</p>
        </>
      )}
    </div>
  );
}

export default App;
