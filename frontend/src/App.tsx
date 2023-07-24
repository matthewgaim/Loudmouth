import {useState} from 'react';
import logo from './logo.svg';
import uuid from 'react-uuid';
import './App.css';

function App() {
  const [newComment, setComment] = useState("");
  const [comments, addComments] = useState([]);

  const addComment = async () => {
    const videoid = await getVideoID();
    const timeVal = await getTimeFunc();
    fetch('http://localhost:8080/addcomment', {
      method: 'post',
      mode: 'no-cors',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({
        "videoid": videoid,
        "uuid": uuid(),
        "time": timeVal,
        "comment": newComment
      })
    });
  };
  
  const getVideoID = async () => {
    let videoid = null;
    let [tab] = await chrome.tabs.query({active: true, lastFocusedWindow: true});
    let currURL = tab?.url;
    if (currURL?.startsWith("https://www.netflix.com/watch/")){
      let splitPreId = currURL.split("https://www.netflix.com/watch/");
      let id = splitPreId[1].split("?")[0];
      if(id !== '') videoid = id;
    }
    console.log(`getVideoID: ${videoid}`);
    return videoid
  }
  
  // Used to read properties of Netflix video playing
  const executeScript = (tabId: number, func: any) => new Promise(resolve => {
    chrome.scripting.executeScript({ target: { tabId }, func }, resolve);
  });

  const getTimeFunc = async () => {
      let [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      let vidProps: any;
      vidProps = await executeScript(tab.id || 0,
        () => {
          const videos = document.getElementsByTagName("video");
          return videos[0].currentTime;
        }
      );
      const timeVal = Math.floor(vidProps[0].result);
      return timeVal;
  };

  const getComments = async () => {
    const time = await getTimeFunc();
    const videoid = await getVideoID();
    console.log(time, videoid);
    if(!videoid){
      alert('Can\'t tell what video is playing');
    } else {
      fetch('http://localhost:8080/getcomments', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          time: time,
          videoid: videoid
        })
      })
      .then(response => response.json())
      .then(data => addComments(data));
    }
  }
  
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        {(comments && comments.length) ? comments.map((post: any) => {
          return <div>{post.time} - {post.comment}</div>
        }) : <>
          <h1>No comments available</h1>
        </>}
        <textarea placeholder="Create a comment" value={newComment} onChange={(e)=>setComment(e.target.value)}/>
        <button onClick={addComment}>Add Comment</button>
        <button onClick={getComments}>Get from DB</button>
      </header>
    </div>
  );
}

export default App;
