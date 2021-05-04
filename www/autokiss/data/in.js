

/*


        По всем техническим вопросам пишите в телеграмм @suvrick



*/


/************************************ UI ***************************/

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

/************************************END UI****************************/

var _this = {};
var selfID = 0;
var tryTimeConnectToGame = 90; //счетчик попыток иницилизации
var catchPacketID = [28, 29, 308]; // пакеты 
var msgError = 'Уп-с. Критическая ошибка. Сообщите об этом мне в телеграм @suvrick';

/*

    Ядро
    Тут и происходит все волшебство программы 
    Перехват пакетов
	

    Client 
    BOTTLE_LEAVE(27);	
    BOTTLE_ROLL(28); speed:I
    BOTTLE_KISS(29); answer:B
    BOTTLE_SAVE(30); target_id:I
    BOTTLE_KICK(31); player_id:I

    Server 
    BOTTLE_LEADER(28); leader_id:I
    BOTTLE_ROLL(29); leader:I, rolled_id:I, speed:I, time:I,
    BOTTLE_KISS(30); player:I, answer:B
    KICKS(308); [target_id:I, actor_id:I, kick_left:I]
    KICK_SAVE(309); target_id:I, actor_id:I

*/

function receiveDataMain(buffer) {

    if(!isGame){
        return; 
    }

    //console.log(">>>>>>>>", buffer)

    switch(buffer.type) {
        /* BOTTLE_LEADER */
        case 28: 

            var leader = buffer[0]

            if (leader === selfID) {
                setInterval(()=>{
                    Main.connection.sendData(28, 0);
                }, 5000);
            }

        break;
        /* BOTTLE_ROLL */
        case 29: 
            
            var leader = buffer[0]
            var rolled = buffer[1]

            if (leader === selfID || rolled === selfID) {
              
                setInterval(()=>{
                    Main.connection.sendData(29, 1);
                }, 8000);
            }

        break;
        /* KICKS */
        case 308 :
           
            var kickID = buffer[0][0][0]
            var kickID2 = buffer[0][0][1]

            //console.log(">>>>>>>>>>>>", buffer, kickID, kickID2)
    
            if (kickID != selfID ){
                return;
            }
    
            //Если выключены автоспасения, выходим
            if(!autoSaveKick) {
                return;
            }

            setInterval(()=>{
                Main.connection.sendData(30, selfID);
            }, 7000);
        break;
    }
}

//Разблокировка гостей (снятия маски рамытия)
function unlockGuest(){
    var items = _this.document.getElementsByClassName("guest")
    for (let i = 0; i < items.length; i++) {
        items[i].setAttribute("is-unlocked", true)
    }
}

//Скрываем всплывающие окна в игре
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

//Метка над контайнером игры
function setTopLine() {
    document.getElementsByTagName("body")[0].style.borderTop = "3px solid yellow";
}

//Снятия метки (линии) над контайнером игры
function delTopMark() {
    document.getElementsByTagName("body")[0].style.borderTop = "0px solid yellow";
}

// UI
function createPopupMenu(){
    menu = document.createElement("div")
    menu.classList.add("menu")

    menu.innerHTML = `
    
    <ul class="kissme_help" >
        <li>
            <h3>Helper KissMe FREE</h3>
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

// Обработка нажания кнопок
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

 //Иницилизация 
function initialize() {

    tryTimeConnectToGame--;

    if (!this.hasOwnProperty("Main")) {
        console.log("Не смог найти обЪект Main")
        if (tryTimeConnectToGame === 0) {
            alert(msgError)
            clearInterval(timerInit)
        }
        return;
    }

    //Пыьаемся найти главный контейнер с игрой
    screenGame = document.getElementById("screen_game");
    if(screenGame === undefined || screenGame === null) {
        console.log("Не смог найти главный контайнер #screen_game")
        if (tryTimeConnectToGame === 0) {
            alert(msgError)
            clearInterval(timerInit)
        }
        return;
    }


    //Ищим блок с кнопками...    
    headerButtons = document.getElementsByClassName("header-buttons")[0]
    if(headerButtons === undefined || headerButtons === null) {
        console.log("Не смог найти класс .header-buttons")
        if (tryTimeConnectToGame === 0) {
            alert(msgError)
            clearInterval(timerInit)
        }
        return;
    }

    clearInterval(timerInit)

    _this = this;
    selfID = Main.self.id;

    // catch packet by server id
    Main.connection.listen(receiveDataMain, catchPacketID);


    createPopupMenu();
    addBtn();
}

var timerInit = setInterval(initialize, 1000);