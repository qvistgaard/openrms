// Create a new store instance.
const store = Vuex.createStore({
  state () {
    return {
      cars: { },
      race: { },
      connection: "disconnected"
    }
  },

  getters: {
    getCarState: (state) => (car, n, d = null) => {
      if (typeof state.cars[car] !== "undefined"){
        if (typeof state.cars[car][n] !== "undefined"){
          if (typeof state.cars[car][n].value !== "undefined"){
            return state.cars[car][n].value
          }
        }
      }
      return d
    },
    getRaceState: (state) => (n, d) => {
      if (typeof state.race[n] !== "undefined"){
        if (typeof state.race[n].value !== "undefined") {
          return state.race[n].value
        }
      }
      return d
    },
    getCarCount: state => () => {
      return Object.keys(state.cars).length
    },

    connection: state => () => {
      return state.connection
    }
  },

  mutations: {
    updateStateFromWebsocket(state, v) {
      let s = state
      for(const item of v.cars) {
        const id = item.id;
        if (typeof s.cars[id] === 'undefined') {
          s.cars[id] = {}
        }
        for (const change of item.changes) {
          s.cars[id][change.name] = { value: change.value }
        }
        state.cars = {
          ...state.cars,
          [id]: { ...s.cars[id] }
        }

      }
      for(const item of v.race){
        for (const change of item.changes) {
          s.race[change.name] = { value: change.value }
        }
      }
      state.race = {
        ...s.race
      }
    },
    connectionState(state, v){
      console.log(v)
      state.connection = v
    }
  }
})

function websocketConnection(params) {
  var query = Object.keys(params)
      .map(k => encodeURIComponent(k) + '=' + encodeURIComponent(params[k]))
      .join('&');

  console.log("Starting connection to WebSocket Server")

  this.websocket = new WebSocket("ws://localhost:8080/ws?"+query)
  this.websocket.onmessage = function(event) {
    store.commit('updateStateFromWebsocket', JSON.parse(event.data))
  }

  this.websocket.onopen = function(event) {
    console.log("Successfully connected to the echo websocket server...")
    store.commit('connectionState', "connected")
  }
  this.websocket.onerror = function(event) {
    console.log("Error in connection")
    store.commit('connectionState', "error")
  }
  this.websocket.onclose = function(event) {
    console.log("Closed Connection")
    store.commit('connectionState', "closed")
    setTimeout(function() {
      websocketConnection(params);
    }, 1000);
  }

  return this.websocket
}


