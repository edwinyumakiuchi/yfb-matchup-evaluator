import './App.css';
// import React, { useEffect, useState } from 'react';

// npm run build && npm run backend
function App() {
  /* const [roster, setRoster] = useState('');

  useEffect(() => {
    fetch('http://localhost:3000/auth/yahoo/callback')
      .then((response) => response.json())
      .then((rosterData) => {
        setRoster(rosterData);
      })
      .catch((error) => {
        console.error('Error fetching user data:', error);
      });
  }, []);

  if (loggedIn) {
    return (<LoggedIn />);
  } else {
    return (<Home />);
  } */
  return (<Home />);
}

function Home() {
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

/* function LoggedIn() {
  const [roster, setRoster] = useState('');

  useEffect(() => {
    fetch('http://localhost:3000/team')
      .then((response) => response.json())
      .then((rosterData) => {
        setRoster(rosterData);
      })
      .catch((error) => {
        console.error('Error fetching user data:', error);
      });
  }, []);

  return ( */
    {/* <div className="container">
      <div className="col-lg-12">
        <br /> */}
        {/* <span className="pull-right"><a onClick={this.logout}>Log out</a></span> */}
        {/* <p>Matchup evalualtor</p>
        {roster && (
          <div> */}
            {/* Display the rosterData here as a table */}
            {/* <table>
              <thead>
                <tr>
                  <th>Player Name</th>
                </tr>
              </thead>
              <tbody>
                {roster.map((playerName, index) => (
                  <tr key={index}>
                    <td>{playerName}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div> */}
  /* );
} */

export default App;
