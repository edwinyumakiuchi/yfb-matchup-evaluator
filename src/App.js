import './App.css';
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

  useEffect(() => {
    // Make the HTTP request to the backend here
    fetch('/mongoData') // Replace '/api/data' with the appropriate backend API endpoint
      .then((response) => response.json())
      .then((data) => {
        setData(data); // Access the 'documents' array from the data object
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  return (
    <div style={{ marginLeft: '20px' }}>
      {data ? (
        <div>
          <p></p>
          {/* Render the data received from the backend */}
          {data.map((player, index) => (
            <div key={index}>
              <p>Team: {player['Fantasy Team']}</p>
              <table className="bordered-table">
                <thead className="header-row">
                  <tr>
                    <th className="bold centered">PLAYER</th>
                    <th className="bold centered">TEAM</th>
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
                  {player.Roster.map((playerData, index) => (
                    <tr key={index}>
                      <td>{playerData.Player}</td>
                      <td className="bold centered">{playerData.Team}</td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                      <td className="bold centered"></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ))}
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
