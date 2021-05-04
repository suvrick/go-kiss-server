
// inject script to game....
(function () {
    var script = document.createElement('script');
    script.src = chrome.runtime.getURL('data/in.js');
    script.onload = function() {
        this.remove();
    };

    (document.head || document.documentElement).appendChild(script);

})()



