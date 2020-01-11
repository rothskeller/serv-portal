<!--
EventsList displays the list of events.
-->

<template lang="pug">
#events-list
  #events-list-title
    | Events in
    |
    select#events-list-year(v-model="year")
      option(v-for="y in years" :key="y" :value="y" v-text="y")
  #events-list-spinner(v-if="loading")
    b-spinner(small)
  #events-list-table(v-else)
    tr
      th Date
      th Event
      th Invited
    tr(v-for="e in events" :key="e.id")
      td(v-text="e.date")
      td: router-link(:to="`/events/${e.id}`" v-text="e.name")
      td
        div(v-for="role in e.roles" v-text="role")
  #events-list-buttons(v-if="canAdd && !loading")
    b-btn(:to="`/events/NEW?year=${year}`") Add Event
</template>

<script>
import range from 'lodash/range'

export default {
  data: () => ({
    year: null,
    years: range(2019, new Date().getFullYear() + 2),
    events: null,
    canAdd: false,
    loading: true,
  }),
  created() {
    this.year = this.$route.params.year || this.$store.state.eventsYear || (new Date().getFullYear())
  },
  watch: {
    year() {
      this.$store.commit('eventsYear', this.year)
      this.load()
    },
  },
  methods: {
    async load() {
      this.loading = true
      const data = (await this.$axios.get(`/api/events?year=${this.year}`)).data
      this.canAdd = data.canAdd
      this.events = data.events
      this.loading = false
    },
  },
}
</script>

<style lang="stylus">
#events-list
  padding 1.5rem 0.75rem
#events-list-title
  display flex
  align-items center
  font-size 1.5rem
#events-list-year
  margin-left 1rem
  font-size 1rem
#events-list-spinner
  margin-top 1.5rem
#events-list-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
#events-list-buttons
  margin-top 1.5rem
</style>
