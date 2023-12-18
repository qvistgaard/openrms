const app = Vue.createApp({
    name: 'App',
    mixins: [openrms],

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
            return this.$store.getters.getCarState(this.car, "car-web-position", "N/A")
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
        isReady: function(){
            return this.$store.getters.getCarState(this.car, "car-ready", false)
        },
        onTrack: function(){
            return this.$store.getters.getCarState(this.car, "car-ontrack", false)
        },
        connectionState: function(){
           return this.$store.getters.connection()
        },
        raceConfirmed: function(){
            return this.$store.getters.getRaceState("race-confirmation")
        },
        raceState: function(){
            return this.$store.getters.getRaceState("race-state")
        },
        numCars: function(){
            return this.$store.getters.getRaceState("race-cars", 0)
        },
        numCarsReady: function(){
            return this.$store.getters.getRaceState("race-cars-ready", 0)
        },
        countdown: function(){
            return this.$store.getters.getRaceState("race-countdown", 0)
        }
    },

    watch: {
        car: function(val){
            this.websocket = websocketConnection({ car: val })
        }
    },

    methods: {
        ready: function() {
            console.log(this.car)
            this.websocket.sendCarCommand(this.car, "car-ready", true)
        },
    }

})

app.use(store)
const vm = app.mount('#app')