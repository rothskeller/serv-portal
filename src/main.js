import '@babel/polyfill'
import camelCase from 'lodash/camelCase'
import upperFirst from 'lodash/upperFirst'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/bootstrap-vue'
import './plugins/axios'
import store from './store'
import router from './router'

Vue.config.productionTip = false

// Register every Vue component in the src tree so that they can be used without
// explicit declaration.
const toRequire = require.context('.', true, /.*\.vue$/)
toRequire.keys().forEach(fileName => {
  const config = toRequire(fileName)
  // Get PascalCase name of component
  const name = fileName.split('/').pop().replace(/\.vue$/, '')
  Vue.component(name, config.default || config)
})

// Create the Vue instance.
const vue = new Vue({
  store,
  router,
  render: h => h('Main')
})

// Find out if we're already logged in (via remember-me cookie), and if so, get
// store the session data.
vue.$axios.get('/api/login').then(resp => { store.commit('login', resp.data) }).catch(() => { store.commit('login', null) })

// Mount and start the app.
vue.$mount('#app')
