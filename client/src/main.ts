import { createApp } from 'vue'
import store from './plugins/store'
import { checkLogin } from './plugins/login'
import router from './router'
import Main from './Main.vue'
import './global.css'

const app = createApp(Main)
app.use(store)

checkLogin().then(() => {
    app.use(router)
    app.mount('#app')
})
