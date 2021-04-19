// Create a new store instance.
const store = Vuex.createStore({
  state () {
    return {
      cars: { },
      race: { }
    }
  },
  mutations: {
    updateStateFromWebsocket(state, v) {
      console.log(v)
      for(const item of v.cars) {
        const id = item.id;
        if (typeof state.cars[id] === 'undefined') {
          state.cars[id] = {}
          console.log(state.cars, id)
        }
        for (const change of item.changes) {
          console.log(change)
          state.cars[id][change.name] = { value: change.value }
        }
        console.log(state)
      }
    }
  }
})


const app = Vue.createApp({
  name: 'App',
  data: function() {
    return {

    }
  },
  store,

  created: function() {



    console.log("Starting connection to WebSocket Server")

    this.connection = new WebSocket("ws://localhost:8080/ws")
    this.connection.onmessage = function(event) {
      store.commit('updateStateFromWebsocket', JSON.parse(event.data))
    }

    this.connection.onopen = function(event) {
      console.log(event)
      console.log("Successfully connected to the echo websocket server...")
    }
  }
})

app.config.globalProperties.$filters = {
  lapTime(value) {
    return value / (1000 * 1000 * 1000)+"s"
  }
}

app.use(store)
const vm = app.mount('#app')