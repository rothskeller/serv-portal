import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    me: null,
    started: false,
    eventsYear: null,
  },
  mutations: {
    login(state, data) {
      state.me = data
      state.started = true
    },
    eventsYear(state, year) { state.eventsYear = year },
  },
  actions: {
  }
})
