const app = Vue.createApp({
    name: 'App',
    mixins: [openrms],

    computed: {
        leaderboard: function(){
            console.log("RELOAD")

            return this.$store.getters.getRaceState("race-leaderboard", { entries: [] }).entries
              .map(x =>  {
                  let v = moment.duration(x["lap"]["lap-time"] / 1000 / 1000).asSeconds()
                  x["lap"]["lap-time"] = Number((v)).toFixed(3)

                  let b = moment.duration(x["best"]["lap-time"] / 1000 / 1000).asSeconds()
                  x["best"]["lap-time"] = Number((b)).toFixed(3)

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
