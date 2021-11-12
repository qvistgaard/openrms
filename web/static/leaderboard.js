const store = Vuex.createStore({
    state: {
        leaderboard: [],
        connection: "disconnected"
    },
    getters: {
        getLeaderboard: state => () => {
            return state.leaderboard
        },
        connection: state => () => {
            return state.connection
        }
    },
    mutations: {
        updateLeaderBoard(state, v){
            console.log(v)
            state.leaderboard = v.content
            console.log(state.leaderboard)
        },
        connectionState(state, v){
            console.log(v)
            state.connection = v
        }
    }
})

const app = Vue.createApp({
    name: 'App',
    mixins: [openrms],

    store,

    computed: {
        leaderboard: function(){
            console.log("RELOAD")

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

    mounted: function (){
        this.connect(function(event){
            store.commit('updateLeaderBoard', event)
        }, {})
    },
})

app.use(store)
const vm = app.mount('#app')
