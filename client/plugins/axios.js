import Vue from 'vue'
import axios from 'axios'
import router from '../router'
import store from '../store'

const config = {
  xsrfCookieName: 'auth',
}

const _axios = axios.create(config);

_axios.interceptors.response.use(
  response => response,
  error => {
    if (
      (error.response && error.response.status === 401) &&
      (router && router.currentRoute && !router.currentRoute.matched.some(record => record.meta.allow401)) &&
      (store.state.started)) {
      store.commit('login', null)
      console.log(error.response && error.response.status, router && router.currentRoute, error.request)
      router.replace({ path: '/login', query: { redirect: router.currentRoute.fullPath } })
    }
    return Promise.reject(error);
  }
)

Plugin.install = function (Vue, options) {
  Vue.axios = _axios
  window.axios = _axios
  Object.defineProperties(Vue.prototype, {
    axios: { get() { return _axios } },
    $axios: { get() { return _axios } },
  })
}

Vue.use(Plugin)

export default Plugin