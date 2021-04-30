const app = Vue.createApp({
    name: 'App',
    data: function(){
        return {
            car: "",
        }
    },
    store,
    computed: {
        fuel: function () {
            let v = this.$store.getters.getCarState(this.car, "car-fuel", 0.00)
            return Number((v)).toFixed(1)
        },
        lastLaps: function(){
            return this.$store.getters.getCarState(this.car, "car-last-laps", { "laps": [] }).laps
              .map(x =>  {
                  let v = moment.duration(x["lap-time"] / 1000 / 1000).asSeconds()
                  x["lap-time"] = Number((v)).toFixed(3)
                  return x
              })
        },
        position: function(){
            return this.$store.getters.getCarState(this.car, "car-leaderboard-position", "N/A")
        },
        damage: function(){
            let v = this.$store.getters.getCarState(this.car, "damage", 0)
            return Number((v)).toFixed(1)
        },
        tireWear: function(){
            let v = this.$store.getters.getCarState(this.car, "tire-wear", 0)
            return Number((v)).toFixed(1)
        },
        limbMode: function(){
            return this.$store.getters.getCarState(this.car, "limb-mode", false)
        },
        inPit: function(){
            return this.$store.getters.getCarState(this.car, "car-in-pit", false)
        },
        pitStop: function(){
            return this.$store.getters.getCarState(this.car, "pit-rule-pit-stop-state", "stopped")
        },
        onTrack: function(){
            return this.$store.getters.getCarState(this.car, "car-ontrack", false)
        },
        connectionState: function(){
           return this.$store.getters.connection()
        },
        carCount: function(){
            return this.$store.getters.getCarCount()
        }
    },

    watch: {
        car: function(val){
            this.websocket = websocketConnection({ car: val })
        }
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