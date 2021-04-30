const app = Vue.createApp({
    name: 'App',
    data: function() {
        return {

        }
    },
    computed: {
        leaderboard: function(){
            return this.$store.getters.getRaceState("race-leaderboard", { entries: [] }).entries
              .map(x =>  {
                  let v = moment.duration(x["lap"]["lap-time"] / 1000 / 1000).asSeconds()
                  x["lap"]["lap-time"] = Number((v)).toFixed(3)
                  return x
              })
        },
    },
    store,
    mounted: function (){
        this.websocket = websocketConnection({})
    },

    methods: {
        start: function() {
            this.websocket.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "start"
                }
            }))
        },
        stop: function(){
            this.websocket.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "stop"
                }
            }))
        },
        pause: function(){
            this.websocket.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "pause"
                }
            }))
        },
        trackCall: function(){
            this.websocket.send(JSON.stringify({
                race: {
                    name: "race-status",
                    value: "track-call"
                }
            }))
        }
    }
})

app.use(store)
const vm = app.mount('#app')
