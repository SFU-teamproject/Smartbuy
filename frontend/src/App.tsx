// src/App.tsx
import React from 'react';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { SmartphoneList } from './components/SmartphoneList';
import { SmartphoneDetail } from './components/SmartphoneDetail';

function App() {
  return (
    <Router>
    <div className="app">
      <header className="app-header">
        <h1>Smartbuy</h1>
      </header>
      <main>
        <Routes>
          <Route path="/" element={<SmartphoneList />} />
          <Route path="/smartphones/:id" element={<SmartphoneDetail />} />
        </Routes>
      </main>
    </div>
  </Router>
  );
}

export default App;