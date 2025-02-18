import React from 'react';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import General from './page/General';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<General />} />
      </Routes>
    </Router>
  );
}

export default App;
