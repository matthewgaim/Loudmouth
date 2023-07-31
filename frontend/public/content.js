function getTime () {
    return $('video').prop('currentTime');
}
  
chrome.runtime.onMessage.addListener(
    function (request, sender, sendResponse) {
        if (request.action === 'getTime') {
            sendResponse({
                url: window.location.href,
                time: getTime()
            });
        }
    }
);