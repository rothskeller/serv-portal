<!--
Events displays the list of events.
-->

<template lang="pug">
Page(title="Events" menuItem="events")
  #events-title
    | Events
    select#events-year(v-model="year")
      option(v-for="y in years" :key="y" :value="y" v-text="y")
  #events-spinner(v-if="loading")
    b-spinner(small)
  #events-table(v-else)
    tr
      th Date
      th Event
      th Invited
      th
    tr(v-for="e in events" :key="e.id")
      td(v-text="e.date")
      td: router-link(:to="`/events/${e.id}`" v-text="e.name")
      td
        div(v-for="role in e.roles" v-text="role")
      td: router-link(v-if="e.canAttendance" :to="`/events/${e.id}/attendance`") Attendance
  #events-buttons(v-if="canAdd && !loading")
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
#events
  padding 1.5rem 0.75rem
#events-title
  display flex
  align-items center
  font-size 1.5rem
#events-year
  margin-left 1rem
  font-size 1rem
#events-spinner
  margin-top 1.5rem
#events-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
#events-buttons
  margin-top 1.5rem
</style>
