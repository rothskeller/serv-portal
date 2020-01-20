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

// Find out if we're already logged in (via cookie), and if so, get and store
// the session data.
Vue.axios.get('/api/login')
  .then(resp => { store.commit('login', resp.data) })
  .catch(() => { store.commit('login', null) })
  .finally(() => {
    new Vue({
      store,
      router,
      render: h => h('Main')
    }).$mount('#app')
  })
