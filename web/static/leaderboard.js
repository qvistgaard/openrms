const app = Vue.createApp({
    name: 'App',
    data: function() {
        return {

        }
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
