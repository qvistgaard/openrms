const app = Vue.createApp({
    name: 'App',
    data: function() {
        return {

        }
    },
    computed: {
        leaderboard: function(){
            /*
            return this.$store.getters.getRaceState("race-leaderboard", { "laps": [] }).laps
              .map(x =>  {
                  let v = moment.duration(x["lap-time"] / 1000 / 1000).asSeconds()
                  x["lap-time"] = Number((v)).toFixed(3)
                  return x
              })

             */
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
        this.connection = websocketConnection({})
    },
})

app.use(store)
const vm = app.mount('#app')


app.config.globalProperties.$filters = {
    lapTime(value) {
        return moment.duration(value / 1000 / 1000).asSeconds()
    }
}
