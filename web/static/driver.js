const app = Vue.createApp({
    name: 'App',
    data: function() {
        return {

        }
    },
    store,

    mounted: function (){
        websocketConnection({ car: 1})
    }
})

app.use(store)
const vm = app.mount('#app')