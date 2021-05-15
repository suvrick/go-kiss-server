

/*************************************************************************************************


        По всем техническим вопросам пишите в телеграм @suvrick


************************************************************************************************/

var isAutoKiss = false;
var isHidePopup = false;
var isAutoSaveKick = false;
var isHideToggleMenuPopupBtn = false;

var menuAppElement = {}
var screenGameElement = {}

var timerHidePopup = 0;
var timerMenuToggle = 0;
var timerConnectToGame = 90;      //счетчик попыток иницилизации

var selfID = 0;
var catchPacketID = [28, 29, 308];  // пакеты 

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

    if (!isAutoKiss) {
        return;
    }

    //console.log(">>>>>>>>", buffer)

    switch (buffer.type) {
        /* BOTTLE_LEADER */
        case 28:

            var leader = buffer[0]

            if (leader === selfID) {
                setInterval(() => {
                    Main.connection.sendData(28, 0);
                }, 5000);
            }

            break;
        /* BOTTLE_ROLL */
        case 29:

            var leader = buffer[0]
            var rolled = buffer[1]

            if (leader === selfID || rolled === selfID) {

                setInterval(() => {
                    Main.connection.sendData(29, 1);
                }, 8000);
            }

            break;
        /* KICKS */
        case 308:

            var kickID = buffer[0][0][0]
            var kickID2 = buffer[0][0][1]

            //console.log(">>>>>>>>>>>>", buffer, kickID, kickID2)

            if (kickID != selfID) {
                return;
            }

            //Если выключены автоспасения, выходим
            if (!isAutoSaveKick) {
                return;
            }

            setInterval(() => {
                Main.connection.sendData(30, selfID);
            }, 7000);
            break;
    }
}



// Иницилизация меню menuAppElement
function createPopupMenu() {

    menuAppElement = document.createElement("div")
    menuAppElement.classList.add("autokiss_container")
    menuAppElement.innerHTML = `    
        <h3>Helper KissMe FREE</h3>
        <span onclick="menuAppElement.style.display='none'" ></span>

        <ul class="autokiss_menu" >    
            <li>
                <span class="on"></span>
                <label id="autoKissBtn" >Автопоцелуи</label>
            </li>
            <li>
                <span class="on"></span>
                <label id="autoSaveBtn" >Автоспасения</label>
            </li>
            <li>
                <span class="on"></span>
                <label id="hidePopupBtn" >Скрыть всплыв.окна</label>
            </li>       
            <li>
                <span class="on"></span>
                <label id="hideToggleMenuPopupBtn" >Скрыть кнопку меню</label>
            </li>
            <li>
                <span class="on"></span>
                <label id="unlockGuestBtn" >Разблокировать гостей </label>
            </li>
        </ul>

        <div class="footer" >
            <p>
                По всем техническим вопросам пишем мне в телеграм  <a target="_blank" href="https://t.me/suvrick">@suvrick</a>
            </p>
        </div>
    `

    //push buttons container
    screenGameElement.appendChild(menuAppElement);

    // set event to menu item
    for (let b of document.querySelectorAll(".autokiss_menu > li")) {
        b.addEventListener('click', function () {

            this.children[0].classList.toggle("on")
            this.children[0].classList.toggle("off")

            var itemID = this.children[1].id;

            switch (itemID) {
                case "autoKissBtn":
                    isAutoKiss = !isAutoKiss;
                    isAutoKiss ?  setTopLine() : delTopMark();  
                    break;
                case "autoSaveBtn":
                    isAutoSaveKick = !isAutoSaveKick;
                    break;
                case "hidePopupBtn":
                    isHidePopup = !isHidePopup;
                    isHidePopup ? hidePopup() : clearInterval(timerPopup)
                    break;
                case "hideToggleMenuPopupBtn":
                    isHideToggleMenuPopupBtn = !isHideToggleMenuPopupBtn;
                    isHideToggleMenuPopupBtn ? hideMenuToggle(this.children[1]) : clearInterval(timerMenuToggle)
                    break;
                case "unlockGuestBtn":
                    unlockGuest()
                    break;
                default:
                    break;
            }
        })
    }
}

// Иницилизация переключателя hideToggleMenuPopupBtn для menuAppElement 
// Переключатель hideToggleMenuPopupBtn отвечает за отображения меню приложения
// Переключатель помещается в контайнер с классом `header-buttons`
function createToggleBtnMenu() {

    var headerButtonsElement = document.getElementsByClassName("header-buttons")[0]

    var hideToggleMenuPopupBtn = document.createElement("span");
    hideToggleMenuPopupBtn.classList.add("btn_toggle_menu");

    hideToggleMenuPopupBtn.addEventListener("click", function () {
        menuAppElement.style.display != 'block' ? menuAppElement.style.display = 'block' : menuAppElement.style.display = 'none';
    });

    hideToggleMenuPopupBtn.addEventListener("mouseover", function () {
        hideToggleMenuPopupBtn.style.opacity = 1;
    });

    headerButtonsElement.appendChild(hideToggleMenuPopupBtn);
}

//Разблокировка гостей (снятия маски рамытия)
function unlockGuest() {
    var items = this.document.getElementsByClassName("guest")
    for (let i = 0; i < items.length; i++) {
        items[i].setAttribute("is-unlocked", true)
    }
}

//Скрываем всплывающие окна в игре
function hidePopup() {

    timerPopup = setInterval( () => {
        var items = this.document.getElementsByClassName("popup")

        if (items === undefined || items === null)
            return;

        for (let i = 0; i < items.length; i++) {

            if (items[i].parentElement.className.includes("dialog") || items[i].parentElement.id.includes("dialog")) {
                items[i].parentElement.remove()
            }

            if (items[i].parentElement.parentElement.className.includes("dialog") || items[i].parentElement.parentElement.id.includes("dialog")) {
                items[i].parentElement.parentElement.remove()
            }

        }
    }, 500)
}

// Скрывает главную кнопку
function hideMenuToggle(el) {
    timerMenuToggle = setInterval(() => {
        if (isHideToggleMenuPopupBtn)
            el.style.opacity = 0.1;
    }, 3000)
}

//Метка над контайнером игры
function setTopLine() {
    document.getElementsByTagName("body")[0].style.borderTop = "3px solid yellow";
}

//Снятия метки (линии) над контайнером игры
function delTopMark() {
    document.getElementsByTagName("body")[0].style.borderTop = "0px solid yellow";
}

//Иницилизация 
function initialize() {

    timerConnectToGame--;

    if (!this.hasOwnProperty("Main")) {
        console.log("Не смог найти обЪект Main")
        if (timerConnectToGame === 0) {
            alert(msgError)
            clearInterval(timerInit)
        }
        return;
    }

    //Пытаемся найти главный контейнер с игрой
    screenGameElement = document.getElementById("screen_game");
    if (screenGameElement === undefined || screenGameElement === null) {
        console.log("Не смог найти главный контайнер #screen_game")
        if (timerConnectToGame === 0) {
            alert(msgError)
            clearInterval(timerInit)
        }
        return;
    }

    clearInterval(timerInit)

    // get self id 
    selfID = Main.self.id;

    // catch packet by server id
    Main.connection.listen(receiveDataMain, catchPacketID);


    createToggleBtnMenu();
    createPopupMenu();
}

var timerInit = setInterval(initialize, 1000);