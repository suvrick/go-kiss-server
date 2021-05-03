var app = new Vue({
    el: "#app",
    data: {
        self: null,

        email: "",
        password: "",
        msgAuthError: "",

        emailReg: "",
        passwordReg: "",
        passwordReg2: "",
        msgRegError: "",

        bots: [],
        botsContainer: [],
        botsContainerStep: 0,
        
        proxies: [],
        proxiesContainer: [],
        proxiesContainerStep: 0,

        users: [],
        usersContainer: [],
        usersContainerStep: 0,

        toastMsg: ""
    },
    methods: {
        navTab: function(id){
            var triggerEl = document.querySelector('#' + id)
            var tab = new window.bootstrap.Tab(triggerEl)
            tab.show()
        },

        getAllUsers: async function(){
            var result = await this.getFetchData("admin/botovod/all", "GET")
            if (result.code === 200) {
                this.users = result.data.users
                this.updateUsersContainer();
            }
        },

        getAllProxies: async function(){
            var result = await this.getFetchData("proxy/all", "GET")
            if (result.code === 200) {
                this.proxies = result.data.proxies
                this.updateProxiesContainer()
            }
        },

        getAllBots: async function () {
            var result = await this.getFetchData("bots/all", "GET")
            if (result.code === 200) {
                this.bots = result.data.bots
                this.updateBotsContainer()
            }
        },

        updateBotsContainer: function(){
            var start = (this.botsContainerStep * 5)
            var end = (this.botsContainerStep * 5) + 5
            this.botsContainer = this.bots.slice(start, end)
        },
        updateBotsContainerStep: function(step){
            if(this.botsContainerStep === 0 && step === -1){
                return;
            }

            if(this.botsContainerStep * 5 - 5 > this.bots.length && step === 1){
                return;
            }


            this.botsContainerStep += step;
            this.updateBotsContainer();
        },

        updateProxiesContainer: function(){
            var start = (this.proxiesContainerStep * 5)
            var end = (this.proxiesContainerStep * 5) + 5
            this.proxiesContainer = this.proxies.slice(start, end)
        },
        updateProxiesContainerStep: function(step){
            if(this.proxiesContainerStep === 0 && step === -1){
                return;
            }

            if(this.proxiesContainerStep * 5 - 5 > this.proxies.length && step === 1){
                return;
            }


            this.proxiesContainerStep += step;
            this.updateProxiesContainer();
        },

        updateUsersContainer: function(){
            var start = (this.usersContainerStep * 5)
            var end = (this.usersContainerStep * 5) + 5
            this.usersContainer = this.users.slice(start, end)
        },
        updateUsersContainerStep: function(step){
            if(this.usersContainerStep === 0 && step === -1){
                return;
            }

            if(this.usersContainerStep * 5 - 5 > this.users.length && step === 1){
                return;
            }


            this.usersContainerStep += step;
            this.updateUsersContainer();
        },

        addBot: async function(url) {
            var result = await this.getFetchData("/bots/add", "POST", { url: url} )
            if (result.code === 200 ) {

                if (result.error){
                    this.showAlert(result.error)
                    return
                }

                await this.getAllBots();
            }
        },

        removeBot: async function(botID){
            var result = await this.getFetchData("/bots/remove/"+ botID, "GET")

            if (result.code === 200) {

                if (result.error){
                    this.showAlert(result.error)
                    return
                }

                await this.getAllBots();
            }
        },

        getSelf: async function () {
            var result = await this.getFetchData("/user/get", "GET")
            if (result.code === 200) {
                console.log(result)

                if (result.error){
                    this.showAlert(result.error)
                    return
                }

                this.self = result.data.user

                if(this.self.role === 'admin'){
                    //this.getAllProxies();
                    this.getAllUsers();
                }

                this.getAllBots()
            }
        },
        onLogin: async function(){
            
            if(this.email.length === 0 || this.password.length === 0){
                this.msgAuthError = "Попытка отправить невалидные данные"
                return;
            }

            var data = { 
                email: this.email,
                password: this.password
             }

             var response = await fetch('/user/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            
            var result = await response.json();
            this.msgAuthError = result?.error ?? '';
            
            switch(result?.code){
                case 200:
                    window.location.href = '/';
                    return;
                default:
                    return;
            }
        },
        onRegister: async function(){
            
            if(this.emailReg.length === 0 || this.passwordReg.length === 0 || this.passwordReg2.length === 0){
                this.msgRegError = "Попытка отправить невалидные данные"
                return;
            }

            var data = { 
                email: this.emailReg,
                password: this.passwordReg
             }

             var response = await fetch('/user/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            
            var result = await response.json();
            this.msgRegError = result?.error ?? ''

            switch(result?.code){
                case 200:
                    this.navTab('tabLoginBtn')
                    return;
                default:
                    return;
            }
        },
        logout: async function(){
            await this.getFetchData("/user/logout", "GET")
            window.location.href = '/'
        },

        getFetchData: async function (url, method, data) {
            
            var response = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            
            var result = await response.json();
            
            switch(result?.code){
                case 401 :
                    this.navTab('tabLoginBtn')
                    return result;
                case 403 :
                    //No Forbidden
                    this.navTab('tabLoginBtn')
                    return result;
                default:
                    return result
            }
        },

        showAlert: function(msg){

            this.toastMsg = msg;

            var el = document.querySelector('#toast')
            var toast =  new window.bootstrap.Toast(el, { delay: 5000, autohide: true })
            toast.show();
        },

        loadFromFile: async function() {
            var input = document.createElement('input');
            input.type="file";
            input.onchange = ev => {
                const file = ev.target.files[0];
                const reader = new FileReader();
          
                reader.onload = e => {
                    var lines = e.target.result.split("\n")
                    lines.forEach( (url) =>{
                        this.addBot(url)
                    }) 
                }
                reader.readAsText(file);
            }
            input.click()
        }
    },

    created: function () {    
        
        this.getSelf()
    }
})