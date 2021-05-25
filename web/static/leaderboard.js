const app = Vue.createApp({
    name: 'App',
    mixins: [openrms],

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
        this.connect()
    },
})

app.use(store)
const vm = app.mount('#app')
