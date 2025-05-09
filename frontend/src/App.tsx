// src/App.tsx
import React from 'react';
import './App.css';
import { SmartphoneList } from './components/SmartphoneList';

function App() {
  return (
    <div className="app">
      <header className="app-header">
        <h1>Smartbuy</h1>
      </header>
      <main>
        <SmartphoneList />
      
      </main>
    </div>
  );
}

export default App;