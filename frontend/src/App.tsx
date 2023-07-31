import {useState, useEffect} from 'react';
import { getTimeFunc, getVideoID, formatTime } from './helperFuncs';
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

  const getComments = async () => {
    const time = await getTimeFunc();
    const videoid = await getVideoID();
    console.log(time, videoid);
    if(!videoid){
      console.log('Can\'t tell what video is playing');
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

  useEffect(()=>{
    const interval = setInterval(() => {
      getComments();
    }, 5000);
    return () => clearInterval(interval);
  });
  
  return (
    <div className="App">
      <header className="App-header">
        <div className="comments">
          <ul className="divide-y divide-gray-100">
            {comments && comments.length ?
              comments.map((post: any) => {
              return (
                <li id={post.uuid} className="flex justify-between gap-x-6 py-5">
                  <div className="flex gap-x-4">
                    <img className="h-12 w-12 flex-none rounded-full bg-gray-50" src="https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png" alt="" />
                    <div className="min-w-0 flex-auto">
                      <p className="text-sm font-semibold leading-6 comment-text">{post.comment}</p>
                      <p className="mt-1 text-xs leading-5 comment-text">{formatTime(post.time)}</p>
                    </div>
                  </div>
                </li>
            )}) : <>
              <p className='text-3xl font-bold pb-2 text-black'>No comments available</p>
            </>}
          </ul>
        </div>
        <div className="add-comment">
          <textarea className='block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500' placeholder="Create a comment" value={newComment} onChange={(e)=>setComment(e.target.value)}/>
          <button className='rounded-full bg-cyan-700 p-3 comment-text' onClick={addComment}>Add Comment</button>
        </div>
      </header>
    </div>
  );
}

export default App;
