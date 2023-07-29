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

const formatTime = (time: string) => {
    let intTime = Number(time)
    const hour = Math.floor(intTime/3600)
    const minute = Math.floor(intTime/60)
    const second = Math.floor(intTime%60)
    return `${(hour ? hour : '00')}:${(minute ? minute : '00')}:${(second ? second : '00')}`;
}

export {getTimeFunc, getVideoID, formatTime}