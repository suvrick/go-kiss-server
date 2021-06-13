
var app = new Vue({

    el: "#app",
    data: {
        self: null,

        frameUrl: "",
        host: "ws://localhost:8080/ws",
        //host: "wss://suvricksoft.ru/ws",

        bots: [],
        botsContainer: [],
        botsContainerStep: 0,

        selectedBots: [],
        selectAllFlag: false,


        toastMsg: "",
        progress: true,

        task_id: "1",
        target_id: "",
        count: 0,

        hideTask: false,

        currentBot: null,

        client: null,

        tasks: 0,
    },
    watch: {
        task_id: function (n, o) {
            this.hideTask = !["1", "6", "7"].includes(this.task_id)
        }
    },
    computed: {
    },
    methods: {
        showDetailBot(bot) {
            this.currentBot = bot;

            var myModal = new bootstrap.Modal(document.getElementById('detailDialog'), {
                keyboard: false
            })
            myModal.show()
        },
        toggleSelected() {
            this.selectAllFlag = !this.selectAllFlag
            this.selectedBots = []
            if (this.selectAllFlag) {
                this.bots.forEach(b => {
                    this.selectedBots.push(b)
                })
            }
        },
        updateSelectedBots: function () {
            console.log(this.selectedBots)
            for (let i = 0; i < this.selectedBots.length; i++) {
                const b = this.selectedBots[i];
                this.updateBot(b.UID)
            }
        },
        deleteBotsHandle() {
            for (let i = 0; i < this.selectedBots.length; i++) {
                const b = this.selectedBots[i];
                this.removeBot(b.UID)
            }
        },
        addBotClick: function () {

            var m = document.getElementById("botAddDialogClose")
            m.click();

            console.log(m)

            if (this.frameUrl === '') {
                return
            }

            this.addBot(this.frameUrl)
            this.frameUrl = '';
            this.showAlert("adding new bot")
        },
        addTaskClick: function () {

            document.getElementById("taskAddDialogClose").click();
            console.log(this.task_id)
            switch (this.task_id) {
                case "1": {
                    this.loadFromFile()
                    return
                }
                case "2": {
                    //ROSE
                    // BUY(6); good_id:I, cost:I, target_id:I, data:I, price_type:B, count:I, hash:S, params: S
                    // [2, 1, 48232366, 9845, 0, 3, "5a2410809e4a9e24ad7ce07f89dd2a18", ""]
                    let prize = {
                        good_id: 2,
                        cost: 1,
                        target_id: parseInt(this.target_id, 10) ?? 0,
                        data: 9845,
                        price_type: 0,
                        count: parseInt(this.count, 10) ?? 1,
                        hash: "5a2410809e4a9e24ad7ce07f89dd2a18",
                        params: "{\"category\": 70, \"screen\": 4}"
                    }
                    this.prizeBotSend(prize)
                    return
                }
                case "3": {
                    //Зажигалка
                    //[2, 1, 42083206, 10120, 0, 1, "5a3f246110761d3e9d5d344af9b5aaa0", "{"category": 70, "screen": 4}"]
                    let prize = {
                        good_id: 2,
                        cost: 1,
                        target_id: parseInt(this.target_id, 10) ?? 0,
                        data: 10120,
                        price_type: 0,
                        count: parseInt(this.count, 10) ?? 1,
                        hash: "5a3f246110761d3e9d5d344af9b5aaa0",
                        params: "{\"category\": 70, \"screen\": 4}"
                    }
                    this.prizeBotSend(prize)
                    return
                }
                case "4": {
                    //HEART
                    //1, 3, 49009358, 0, 0, 1
                    let prize = {
                        good_id: 1,
                        cost: 3,
                        target_id: parseInt(this.target_id, 10) ?? 0,
                        data: 0,
                        price_type: 0,
                        count: parseInt(this.count, 10) ?? 1,
                        hash: "",
                        params: ""
                    }
                    this.prizeBotSend(prize)
                    return
                }
                case "5": {
                    //Фляжка
                    //[2, 1, 48273340, 10125, 0, 4, "dcddfb78c71cb85d2e7cd978d46f2ee5", "{"category": 70, "screen": 4}"]
                    let prize = {
                        good_id: 2,
                        cost: 1,
                        target_id: parseInt(this.target_id, 10) ?? 0,
                        data: 10125,
                        price_type: 0,
                        count: parseInt(this.count, 10) ?? 1,
                        hash: "dcddfb78c71cb85d2e7cd978d46f2ee5",
                        params: "{\"category\": 70, \"screen\": 4}"
                    }
                    this.prizeBotSend(prize)
                    return
                }
                case "6": {
                    this.removeBotSend()
                    return
                }
                case "7": {
                    this.updateBotSend()
                    return
                }
                case "8": {
                    this.viewBotSend()
                    return
                }
            }


        },
        toggleRow: function (bot) {
            let r = this.selectedBots.find(b => b.UID === bot.UID)
            if (r) {
                this.selectedBots = this.selectedBots.filter(b => b.UID != r.UID)
            } else {
                this.selectedBots.push(bot)
            }
        },
        updateBotsContainer: function () {
            var start = (this.botsContainerStep * 5)
            var end = (this.botsContainerStep * 5) + 5
            this.botsContainer = this.bots.slice(start, end)
        },
        updateBotsContainerStep: function (step) {
            if (this.botsContainerStep === 0 && step === -1) {
                return;
            }

            if (this.botsContainerStep * 5 - 5 > this.bots.length && step === 1) {
                return;
            }


            this.botsContainerStep += step;
            this.updateBotsContainer();
        },
        showAlert: function (msg) {
            console.log(msg)
            // this.toastMsg = msg;
            // var el = document.querySelector('#toast')
            // var toast = new window.bootstrap.Toast(el, { delay: 5000, autohide: true })
            // toast.show();
        },


        addBotSend(url) {
            var packet = {
                type: ClientPacketType.ADD_BOT_SEND,
                data: {
                    url: url
                }
            }

            var cmd = JSON.stringify(packet)
            this.client.send(cmd)
        },
        updateBotSend() {
            this.selectedBots.forEach(b => {
                var packet = {
                    type: ClientPacketType.UPDATE_BOT_SEND,
                    data: {
                        uid: b.UID
                    }
                }

                var cmd = JSON.stringify(packet)
                this.client.send(cmd)
            })
        },
        removeBotSend() {

            this.selectedBots.forEach(b => {
                var packet = {
                    type: ClientPacketType.REMOVE_BOT_SEND,
                    data: {
                        uid: b.UID
                    }
                }

                var cmd = JSON.stringify(packet)
                this.client.send(cmd)
            })
        },
        allBotSend() {
            var packet = {
                type: ClientPacketType.ALL_BOT_SEND,
                data: {}
            }

            var cmd = JSON.stringify(packet)
            this.client.send(cmd)
        },
        prizeBotSend(prize) {

                var packet = {
                    type: ClientPacketType.PRIZE_BOT_SEND,
                    data: {
                        uids: this.selectedBots.map( o => { return o.UID }),
                        ...prize
                    }
                }
                var cmd = JSON.stringify(packet)
                this.client.send(cmd)
        },
        viewBotSend(){
            var packet = {
                type: ClientPacketType.VIEW_BOT_SEND,
                data: {
                    uids: this.selectedBots.map( o => { return o.UID }),
                    target_id: parseInt(this.target_id, 10)
                }
            }
            var cmd = JSON.stringify(packet)
            this.client.send(cmd)
        },

        addBotRecv(data) {

            if (data.bot === null) {
                return
            }

            this.bots.push(data.bot)
            this.bots = this.bots.sort(function (a, b) {
                if (a.UID > b.UID) {
                    return 1;
                }
                if (a.UID < b.UID) {
                    return -1;
                }
                return 0;
            });

            this.updateBotsContainer()
        },
        updateBotRecv(data) {
            this.bots.forEach((item, i)=>{
                if(item.UID === data.bot.UID) {
                    this.bots[i] = data.bot
                }
            })

            this.selectedBots = this.selectedBots.filter(b => b.UID != data.bot.UID)
            this.updateBotsContainer()
        },
        removeBotRecv(data) {
            this.bots = this.bots.filter(b => b.UID != data.uid)
            this.selectedBots = this.selectedBots.filter(b => b.UID != data.uid)
            this.updateBotsContainer()
        },
        allBotRecv(data) {
            this.bots = data.bots.sort(function (a, b) {
                if (a.UID > b.UID) {
                    return 1;
                }
                if (a.UID < b.UID) {
                    return -1;
                }
                return 0;
            });

            this.updateBotsContainer()
        },

        selfRecv(data) {
            this.self = data.user
            this.allBotSend()
        },

        addTaskRecv() {
            this.tasks++
        },
        removeTaskRecv() {
            this.tasks--
        },

        loadFromFile: async function () {
            var input = document.createElement('input');
            input.type = "file";
            input.onchange = ev => {
                const file = ev.target.files[0];
                const reader = new FileReader();

                reader.onload = e => {
                    var lines = e.target.result.split("\n")
                    lines.forEach((url) => {
                        this.addBotSend(url)
                    })
                }
                reader.readAsText(file);
            }
            input.click()
        },
        parsePacket(pack) {
            let p = JSON.parse(pack)
            console.log(p)
            switch (p.type) {

                case ServerPacketType.SELF_RECV:
                    this.selfRecv(p.data)
                    break;

                case ServerPacketType.ADD_BOT_RECV:
                    this.addBotRecv(p.data)
                    break;
                case ServerPacketType.ALL_BOT_RECV:
                    this.allBotRecv(p.data)
                    break;
                case ServerPacketType.REMOVE_BOT_RECV:
                    this.removeBotRecv(p.data)
                    break;
                case ServerPacketType.UPDATE_BOT_RECV:
                    this.updateBotRecv(p.data)
                    break;

                case ServerPacketType.ADD_TASK_RECV:
                    this.addTaskRecv()
                    break;
                case ServerPacketType.REMOVE_TASK_RECV:
                    this.removeTaskRecv()
                    break;
                case ServerPacketType.ERROR_RECV:
                    this.showAlert(p.data.error)
                    break;
                default:
                    break;
            }
        },
        initSocket() {

            this.client = new WebSocket(this.host);

            this.client.onopen = (e) => {
                console.log("socket open")
            };

            this.client.onmessage = (event) => {
                this.parsePacket(event.data)
            };

            this.client.onclose = function (event) {
                console.log("socket close")
            };

            this.client.onerror = function (error) {
                console.log("socket error:", error)
                location.href = '/login'
            };
        }
    },

    created: function () {
        this.initSocket()
    }
})

const PrizeType = Object.freeze({
    ROSE: "ROSE",
    ROSE_FREE: "ROSE_FREE",
    HEART: "HEART",
    HEART_FREE: "HEART_FREE"
});

const ClientPacketType = Object.freeze({
    ADD_BOT_SEND: "ADD_BOT_SEND",
    ALL_BOT_SEND: "ALL_BOT_SEND",
    REMOVE_BOT_SEND: "REMOVE_BOT_SEND",
    UPDATE_BOT_SEND: "UPDATE_BOT_SEND",

    PRIZE_BOT_SEND: "PRIZE_BOT_SEND",
    VIEW_BOT_SEND: "VIEW_BOT_SEND"
});

const ServerPacketType = Object.freeze({
    SELF_RECV: "SELF_RECV",

    ADD_BOT_RECV: "ADD_BOT_RECV",
    ALL_BOT_RECV: "ALL_BOT_RECV",
    REMOVE_BOT_RECV: "REMOVE_BOT_RECV",
    UPDATE_BOT_RECV: "UPDATE_BOT_RECV",

    ADD_TASK_RECV: "ADD_TASK_RECV",
    REMOVE_TASK_RECV: "REMOVE_TASK_RECV",

    ERROR_RECV: "ERROR_RECV"
});