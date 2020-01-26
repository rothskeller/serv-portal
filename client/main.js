import '@babel/polyfill'
import camelCase from 'lodash/camelCase'
import upperFirst from 'lodash/upperFirst'
import 'mutationobserver-shim'
import Vue from 'vue'
import * as VueGoogleMaps from 'vue2-google-maps'
import './plugins/bootstrap-vue'
import './plugins/axios'
import store from './store'
import router from './router'
import Main from './Main.vue'
Vue.config.productionTip = false

Vue.use(VueGoogleMaps, { load: { key: 'AIzaSyDYiDjdYhCKZnM4qbK68KZRjKZqJiQ1dZw' } })

// Find out if we're already logged in (via cookie), and if so, get and store
// the session data.
Vue.axios.get('/api/login')
  .then(resp => { store.commit('login', resp.data) })
  .catch(() => { store.commit('login', null) })
  .finally(() => {
    new Vue({
      store,
      router,
      render: h => h(Main)
    }).$mount('#app')
  })
