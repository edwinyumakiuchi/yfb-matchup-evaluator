import './App.css';
import React, { useEffect, useState } from 'react';

function App() {
  const [roster, setRoster] = useState('');

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

  const handleYahooLogin = () => {
    // Redirect to the Yahoo authentication endpoint provided by main.go
    window.location.href = 'http://localhost:3000/auth/yahoo';
  };

  return (
    <div className="App">
      <header className="App-header">
        <button onClick={handleYahooLogin}>Login with Yahoo</button>
        {roster && (
          <div>
            {/* Display the rosterData here */}
            <p>Roster Data: {JSON.stringify(roster)}</p>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
