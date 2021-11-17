const store = Vuex.createStore({
    state: {
        leaderboard: [],
        raceStatus: 0,
        raceTimer: 0,
        connection: "disconnected"
    },
    getters: {
        getLeaderboard: state => () => {
            return state.leaderboard
        },
        getRaceTimer: state => () => {
            return state.raceTimer
        },
        getRaceStatus: state => () => {
            return state.raceStatus
        },
        connection: state => () => {
            return state.connection
        }
    },
    mutations: {
        updateLeaderBoard(state, v){
            console.log(v)
            state.leaderboard = v.content.Leaderboard
            state.raceStatus = v.content.RaceStatus
            state.raceTimer = v.content.RaceTimer
            console.log(state.leaderboard)
        },
        connectionState(state, v){
            state.connection = v
        }
    }
})

const app = Vue.createApp({
    name: 'App',
    mixins: [openrms],

    store,

    computed: {
        racetimer: function () {
            return moment.duration(this.$store.getters.getRaceTimer() / 1000 / 1000).asSeconds()
        },
        racestate: function () {
            let s = this.$store.getters.getRaceStatus()
            if (s === 0)  {
                return "Stopped"
            }
            if (s === 1)  {
                return "Paused"
            }
            if (s === 2)  {
                return "Running"
            }
            return "Unknown"
        },
        leaderboard: function(){
            return this.$store.getters.getLeaderboard()
              .map(x =>  {
                  let v = moment.duration(x["last"] / 1000 / 1000).asSeconds()
                  x["last"] = Number((v)).toFixed(3)

                  let b = moment.duration(x["best"] / 1000 / 1000).asSeconds()
                  x["best"] = Number((b)).toFixed(3)

                  let d = moment.duration(x["delta"] / 1000 / 1000).asSeconds()

                  x["isSlow"] = false
                  x["isFast"] = false

                  if (d > 0) {
                      x["isSlow"] = true
                      prefix = "+"
                  } else if (d < 0){
                      x["isFast"] = true
                      prefix = ""
                  } else {
                      prefix = ""
                  }
                  x["delta"] = prefix+Number((d)).toFixed(3)
                  x["fuel"] = Number((x["fuel"])).toFixed(1)

                  return x
              })
        },
    },

    methods: {
        start: function () {
            fetch("http://"+location.host+"/v1/race/start", { method: 'POST' })
        },
        stop: function () {
            fetch("http://"+location.host+"/v1/race/stop", { method: 'POST' })
        },
        pause: function () {
            fetch("http://"+location.host+"/v1/race/pause", { method: 'POST' })
        }
    },

    mounted: function (){
        this.connect(function(event){
            store.commit('updateLeaderBoard', event)
        }, {})
    },
})

app.use(store)
const vm = app.mount('#app')
