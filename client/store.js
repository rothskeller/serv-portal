import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    me: null,
    started: false,
    eventsYear: null,
    touch: false,
    page: {},
  },
  mutations: {
    login(state, data) {
      state.me = data
      state.started = true
    },
    eventsYear(state, year) { state.eventsYear = year },
    setPage(state, page) { state.page = page },
    setTouch(state, touch) { state.touch = touch },
  },
  actions: {
  }
})

try {
  const mql = window.matchMedia(
    '(pointer:coarse),(hover:none),(-moz-touch-enabled:1)'
  )
  store.commit('setTouch', mql.matches)
  mql.addListener(evt => {
    store.commit('setTouch', evt.matches)
  })
} catch (err) {
  console.error('Touch detection:', err)
}

export default store
