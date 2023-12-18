const app = Vue.createApp({
    mixins: [openrms],

    name: 'App',
    computed: {
        cars: function(){
            return Object.keys(this.$store.state.cars)
        },
        countdown: function(){
            return this.$store.getters.getRaceState("race-countdown", 0)
        },
        raceConfirmed: function(){
            return this.$store.getters.getRaceState("race-confirmation")
        },
        raceState: function(){
            return this.$store.getters.getRaceState("race-state")
        },
    },
    store,

    mounted: function (){
        this.connect()
    },

    methods: {
        isReady: function (car){
            return this.carState(car, "car-ready", false)
        },
        position: function (car){
            return this.carState(car, "car-web-position", "N/A")
        },
        inPit: function (car){
            return this.carState(car, "car-in-pit", false)
        },
        fuel: function (car) {
            let v = this.carState(car, "car-fuel", 0.00)
            return Number((v))
        },
        lap: function (car) {
            return this.carState(car, "car-lap", {})["lap-number"]
        },
        time: function (car) {
            let lt = this.carState(car, "car-lap", {})["lap-time"]
            return moment.duration(lt / 1000 / 1000).asSeconds().toFixed(3)
        },

        isStarted: function () {
            return this.$store.getters.getRaceState("race-state", "none") === "started"
        },
        isRunning: function () {
            return this.$store.getters.getRaceState("race-state", "none") === "running"
        },
        isConfirmed: function () {
            return this.$store.getters.getRaceState("race-confirmation", "unconfirmed") === "confirmed"
        }
    }
})

app.use(store)
const vm = app.mount('#app')
