const app = Vue.createApp({
    name: 'App',
    data: function() {
        return {

        }
    },
    store,

    mounted: function (){
        this.connection = websocketConnection({ car: 1})
    },

    methods: {
        start: function() {
            this.connection.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "start"
                }
            }))
        },
        stop: function(){
            this.connection.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "stop"
                }
            }))
        },
        pause: function(){
            this.connection.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "pause"
                }
            }))
        }
    }
})

app.use(store)
const vm = app.mount('#app')