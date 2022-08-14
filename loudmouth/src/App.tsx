import React from 'react';
import {useState} from 'react';
import logo from './logo.svg';
import './App.css';
import { connect, sendMsg } from "./api";

function App() {
  connect();
  const [message, setMessage] = useState('');

  const send = () => {
    console.log(message);
    sendMsg(JSON.stringify({author:"Joe Schmoe", comment:message, title:"Monkey Man", episode:"2"}));
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <form>
          <label>Enter Comment:
            <textarea value={message} onChange={e => setMessage(e.target.value)}/>
          </label>
          <button onClick={send}>Send!</button>
        </form>
      </header>
    </div>
  );
}

export default App;
