import React, {Component, useState, useEffect} from 'react';
import logo from './logo.svg';
import './App.css';
import { connect, sendMsg } from "./network";
var i = 1;
function App() {
  const [message, setMessage] = useState("");
  const [comments, setComments] = useState("");

  useEffect(() => {
    connect((msg: any) => {
      if (msg) {
        console.log(comments);
        setComments(comments + '\n' + msg.author + ": " + msg.comment);
      }
    });
  }, []);

  const send = () => {
    console.log(message);
    sendMsg(JSON.stringify({author:"Joe Schmoe", comment:message, title:"Monkey Man", episode:"2"}));
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
          <label>Enter Comment:
            <textarea value={message} onChange={e => setMessage(e.target.value)}/>
          </label>
          <button onClick={send}>Send!</button>
        <div>{comments}</div>
      </header>
    </div>
  );
}

export default App;
