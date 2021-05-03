
var _this = {};

var isDebug = true;
var urlData = "";
var urlInit = "";

var selfID = 0;
var roomID = 0;

var isGame = false;
var isMenuShow = false;
var autoSaveKick = false;
var autoMoveToRoom = true;
var isShowAlert = false;
var isHidePopup = false;
var isHideGeneralBtn = false;

var menu = {}
var screenGame = {}
var headerButtons = {}

var hideGeneralBtn = {}
var btnClose = {}
var autoKissBtn = {}
var autoSaveBtn = {}
var hidePopupBtn = {}
var unlockGuestBtn = {}

var timerPopup = 0;





function getData(data) {

    fetch(urlData + "/" + selfID, {
        method: "POST",
        body: data,
        mode: 'no-cors',
        headers: {
            'Accept': 'application/json, text/plain, */*',
            'Content-Type': 'application/json'
          },
    }).then(response => response.json())
    .then(xhr => {
        if (xhr.status === 200) {
            var result = JSON.parse(xhr.responseText);

            if (result.status === 403) {
                showAlert();
                return;
            }

            if (result.status === 200) {
                setTimeout(callHandler.bind(_this, result), result.delay);
                return;
            }

            console.log("unknow status:", result)
        }
    })


    // var xhr = new XMLHttpRequest();
    // xhr.open("POST", urlData + "/" + selfID, true);
    // xhr.setRequestHeader('Content-Type', 'application/octet-stream');
    // xhr.setRequestHeader('Content-Type', 'application/json');
    // xhr.setRequestHeader('Content-Type', '*');


    // xhr.onload = function () {

       
    // };
    // xhr.send(data);
}

function callHandler(result) {

    if(result.code === 28 || result.code === 30) {
        _this.Main.connection.sendData(result.code, result.data[0]);
        return;
    }

    _this.Main.connection.sendData(result.code, result.data);


}

function showAlert() {

    if (isShowAlert)
        return;

    var div = document.createElement("div")
    div.style.position = "absolute";
    div.style.bottom = "0"
    div.style.width = "300px"
    div.style.padding = "10px";
    div.style.background = "white";

    var p = document.createElement("span");
    p.style.padding = " 0 10px 10px 10px";
    p.style.display = "block";
    p.innerText = "Программа не зарегистрирована.Тестовый пириод закончился.\nНапишите мне в телеграмм @help_auto_kiss для приобретения программы";

    var head = document.createElement("span");
    head.innerText = "Bottle Auto Kiss Helper";
    head.style.display = "block";
    head.style.fontWeight = "bold";
    head.style.paddingLeft = "10px"

    var btnClose = document.createElement("span");
    btnClose.style.display = "block";
    btnClose.style.position = "absolute";
    btnClose.style.right = "10px";
    btnClose.style.top = "10px";
    btnClose.style.cursor = "pointer";
    btnClose.innerText = "x";
    btnClose.addEventListener("click", function () {
        screenGame.removeChild(div)
    })

    div.appendChild(head);
    div.appendChild(p);
    div.appendChild(btnClose);

    screenGame.appendChild(div)

    isShowAlert = true;
}

function unlock(){
    var s = document.getElementsByClassName("recaptcha-checkbox goog-inline-block recaptcha-checkbox-unchecked rc-anchor-checkbox")[0];
    if(s === undefined || s === null) {
        console.log("nullable object capcha")
        return;
    }

    s.click();
}

function receiveDataMain(buffer) {
    //console.log(buffer)
    if(buffer.type === 25){
        roomID = buffer[0]
        console.log("RoomID:", roomID)
        return;
    }


    if(!isGame)
    return; 

    var arr = new ArrayBuffer(buffer.bytesLength + 6 + 10);
    var data = new DataView(arr, 0, buffer.bytesLength + 6 + 10);

    data.setInt32(0, buffer.id, true);
    data.setInt16(4, buffer.type, true);

    if (buffer.type === 29) {
        data.setInt32(6, buffer[0], true);
        data.setInt32(10, buffer[1], true);
       // data.setInt32(14, buffer[2], true);
       // data.setInt32(18, buffer[3], true);
    }


    if(buffer.type === 27 ){
        console.log(buffer[0],selfID , autoMoveToRoom)
        if(buffer[0] === selfID && autoMoveToRoom) {
            console.log("tru move to roomID:", roomID)

            if(roomID === 0) {
                return;
            }
            
            Main.connection.sendData(202, roomID)
            return;
        }
    }

    if (buffer.type === 28) {
        data.setInt32(6, buffer[0], true);
    }

    if (buffer.type === 308) {

        if(buffer.bytesLength < 3) {
            console.log("308 >>>>>>>> ", buffer)
            return;
        }

        var kickID = buffer[0][0][0]
        var kickID2 = buffer[0][0][1]

        if (kickID != selfID ){
            return;
        }

        if(!autoSaveKick) {
            return;
        }

        data.setInt32(6, kickID, true);
        data.setInt32(10, kickID2, true);
        //console.log("kickIDS:", kickID, kickID2)
        //console.log("autosavekick send");
    }

    getData(data.buffer);

}

function unlockGuest(){
    var items = _this.document.getElementsByClassName("guest")
    for (let i = 0; i < items.length; i++) {
        items[i].setAttribute("is-unlocked", true)
    }
}

function hidePopup() {
    
    timerPopup = setInterval(function(){
        var items = _this.document.getElementsByClassName("popup")

        if(items === undefined || items === null)
            return;

        for (let i = 0; i < items.length; i++) {

            if(items[i].parentElement.className.includes("dialog") || items[i].parentElement.id.includes("dialog")){
                console.log("delete popup", items[i].parentElement.className)
                items[i].parentElement.remove() 
            }

            if(items[i].parentElement.parentElement.className.includes("dialog") || items[i].parentElement.parentElement.id.includes("dialog")){
                console.log("delete popup", items[i].parentElement.parentElement.className)
                items[i].parentElement.parentElement.remove() 
            }

        }
    }, 1000)
}

function setTopLine() {
    document.getElementsByTagName("body")[0].style.borderTop = "3px solid yellow";
}

function delTopMark() {
    document.getElementsByTagName("body")[0].style.borderTop = "0px solid yellow";
}

function createPopupMenu(){
    menu = document.createElement("div")
    menu.classList.add("menu")

    menu.innerHTML = `
    
    <ul>
        <li>
            <h3>Helper KissMe [v3.3]</h3>
            <span id="btnClose">x</span>
        <li>       
        <li>
            <label id="hideGeneralBtn" >Скрыть кнопку меню</label>
        </li>
        <li>
            <label id="autoKissBtn" >Автопоцелуи (вкл)</label>
        </li>
        <li>
            <label id="autoSaveBtn" >Автоспасения (вкл)</label>
        </li>
        <li>
            <label id="hidePopupBtn" >Скрыть всплыв.окна (вкл)</label>
        </li>       
        <li>
            <label id="unlockGuestBtn" >Разблокировать гостей </label>
        </li>
    </ul>
    `

    screenGame.appendChild(menu);
}

function addBtn() {

    hideGeneralBtn = document.getElementById("hideGeneralBtn")
    btnClose = document.getElementById("btnClose")
    autoKissBtn = document.getElementById("autoKissBtn")
    autoSaveBtn = document.getElementById("autoSaveBtn")
    hidePopupBtn = document.getElementById("hidePopupBtn")
    unlockGuestBtn = document.getElementById("unlockGuestBtn")

     var btn = document.createElement("span")
     btn.innerText = "+";
     btn.classList.add("btn")
     btn.addEventListener("click", function(){
        isMenuShow = !isMenuShow;
        if(isMenuShow){
            menu.style.display = "block";
        } else {
            menu.style.display = "none";
        }

     })
 
     btn.addEventListener("mouseover", function(){
         btn.style.opacity = 1;
     })
 
     headerButtons.appendChild(btn)

     hideGeneralBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";

        isHideGeneralBtn = !isHideGeneralBtn;
        if(isHideGeneralBtn) {
            hideGeneralBtn.innerText = "Показать кнопку меню"
        } else {
            hideGeneralBtn.innerText = "Скрыть кнопку меню"
        }
     })

     btnClose.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
     })

    autoSaveBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
        autoSaveKick = !autoSaveKick;

        if(autoSaveKick) {
            autoSaveBtn.innerText = "Автоспасения (выкл)"
        } else {
            autoSaveBtn.innerText = "Автоспасения (вкл)"
        }
     })

     autoKissBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
        isGame = !isGame;
        if(isGame) {
            autoKissBtn.innerText = "Автопоцелуи (выкл)"
            setTopLine();
        } else {
            autoKissBtn.innerText = "Автопоцелуи (вкл)"
            delTopMark();
        }
     })

     hidePopupBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
        isHidePopup = !isHidePopup;
        if(isHidePopup) {
            console.log("hide popup")
            hidePopupBtn.innerText = "Скрыть всплыв.окна (выкл)";
            hidePopup();
        } else {
            console.log("show popup")
            hidePopupBtn.innerText = "Скрыть всплыв.окна (вкл)";
            clearInterval(timerPopup);
        }
     })

     unlockGuestBtn.addEventListener("click",function(){
        unlockGuest();
     })

 
     setInterval(()=>{
         if(isHideGeneralBtn)
            btn.style.opacity = 0.1;
     }, 3000)
 }

function init() {

    if (!this.hasOwnProperty("Main"))
        return;

    screenGame = document.getElementById("screen_game");
    if(screenGame === undefined || screenGame === null)
        return;

    headerButtons = document.getElementsByClassName("header-buttons")[0]
    if(headerButtons === undefined || headerButtons === null)
        return;

    clearInterval(timerInit)

    if(isDebug) {
        urlData = "http://localhost:8080/autokiss/who";
        urlInit = "http://localhost:8080/autokiss/init";
    } else {
        urlData = "https://suvricksoft.ru/autokiss/who";
        urlInit = "https://suvricksoft.ru/autokiss/init";
    }

    _this = this;
    selfID = Main.self.id;
    Main.connection.listen(receiveDataMain, [25, 28, 29, 308]);

    var xhr = new XMLHttpRequest();
    xhr.open("GET", urlInit + "/" + selfID, true);
    xhr.send();


    createPopupMenu();
    addBtn();
}

var timerInit = setInterval(init, 1000)