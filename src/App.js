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
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Make the HTTP request to the backend here
    fetch('/mongoData') // Replace '/api/data' with the appropriate backend API endpoint
      .then((response) => response.json())
      .then((data) => {
        setData(data); // Access the 'documents' array from the data object
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <p>Yahoo Fantasy Basketball: matchup evaluator</p>
        {loading ? (
          <p>Loading...</p>
        ) : (
          data ? (
            <div>
              {/* Render the data received from the backend */}
              {data.map((player) => (
                <div key={player.ID}>
                  <p>Team Name: {player['Team name']}</p>
                  <table>
                    <thead>
                      <tr>
                        <th>Player Names</th>
                      </tr>
                    </thead>
                    <tbody>
                      {/* Render each player's name separately */}
                      {player.Roster.map((playerName, index) => (
                        <tr key={index}>
                          <td>{playerName}</td>
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
              {data && <p>{JSON.stringify(data)}</p>}
            </>
          )
        )}
      </header>
    </div>
  );
}

export default App;
