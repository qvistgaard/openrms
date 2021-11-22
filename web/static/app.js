// Create a new store instance.

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
      window.setTimeout(function () {
        console.log("Successfully connected to openrms...", this.websocket.OPEN, this.websocket.readyState)
      }, 1000)


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
    console.log(event)
    onMessage(JSON.parse(event.data))
    // store.commit('updateStateFromWebsocket', JSON.parse(event.data))
  }

  this.websocket.onopen = function(event) {
    console.log("Successfully connected to openrms...", this.websocket)
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


