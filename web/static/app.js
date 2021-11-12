// Create a new store instance.
/*

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
      console.log(v)

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
*/


const openrms = {
/*  data: function(){
    return {
      car: "",
      setStateCarId: "null",
      setStateCarState: "null",
      setStateCarValue: "null",
      setCourseState: "null",
      setCourseValue: "null"
    }
  },*/
  methods: {
    connect: function (onMessage, params = {}){
      this.websocket = websocketConnection(onMessage, params)
    },

    formSetCarState: function(){
      this.carCommand(this.setStateCarId, this.setStateCarState, this.setStateCarValue)
      console.log("test")
      return false
    },
    formSetCourseState: function(){
      this.raceCommand(this.setCourseState, this.setCourseValue)
      console.log("test")
      return false
    },

    carState: function(car, state, d){
      return this.$store.getters.getCarState(car, state, d)
    },
    raceState: function(state, d){
      return this.$store.getters.getRaceState(state, d)
    },

    raceCommand: function (state, value) {
      this.websocket.sendRaceCommand( state, value)
    },
    carCommand: function (car, state, value) {
      this.websocket.sendCarCommand( car, state, value)
    },

    start: function() {
      this.raceCommand( "race-state", "started")
    },
    stop: function(){
      this.websocket.sendRaceCommand("race-state", "stopped")
    },
    pause: function(){
      this.websocket.sendRaceCommand("race-status", "pause")
    },
    trackCall: function(){
      this.websocket.sendRaceCommand("race-status", "track-call")
    }
  }
}



function websocketConnection(onMessage, params) {
  var query = Object.keys(params)
      .map(k => encodeURIComponent(k) + '=' + encodeURIComponent(params[k]))
      .join('&');

  console.log("Starting connection to WebSocket Server", params)

  this.websocket = new WebSocket("ws://"+location.host+"/ws?"+query)
  this.websocket.onmessage = function(event) {
    onMessage(JSON.parse(event.data))
    // store.commit('updateStateFromWebsocket', JSON.parse(event.data))
  }

  this.websocket.onopen = function(event) {
    console.log("Successfully connected to the echo webserver server...")
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
      websocketConnection(onMessage, params);
    }, 1000);
  }
  return this.websocket
}


