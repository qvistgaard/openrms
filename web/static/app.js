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
      for(const item of v.cars) {
        const id = item.id;
        if (typeof state.cars[id] === 'undefined') {
          state.cars[id] = {}
        }
        for (const change of item.changes) {
          state.cars[id][change.name] = { value: change.value }
        }
      }
      for(const item of v.race){
        for (const change of item.changes) {
          state.race[change.name] = change.value
        }
      }
    },
    connectionState(state, v){
      state.connection = v
    }
  }
})

function websocketConnection(params) {
  var query = Object.keys(params)
      .map(k => encodeURIComponent(k) + '=' + encodeURIComponent(params[k]))
      .join('&');

  console.log("Starting connection to WebSocket Server")

  this.connection = new WebSocket("ws://localhost:8080/ws?"+query)
  this.connection.onmessage = function(event) {
    store.commit('updateStateFromWebsocket', JSON.parse(event.data))

  }

  this.connection.onopen = function(event) {
    console.log("Successfully connected to the echo websocket server...")
    store.commit('connectionState', "connected")
  }
  this.connection.onerror = function(event) {
    console.log("Error in connection")
    store.commit('connectionState', "error")
  }
  this.connection.onclose = function(event) {
    console.log("Closed Connection")
    store.commit('connectionState', "closed")
  }

  return this.connection
}


