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
    .events-list-date.events-list-heading Date
    .events-list-event.events-list-heading Event
    .events-list-location.events-list-heading Location
    template(v-for="e in events")
      .events-list-date
        span.events-list-year(v-text="e.date.substr(0, 5)")
        span(v-text="e.date.substr(5)")
        span.events-list-start(v-text="e.start")
      .events-list-event
        EventOrgDot.mr-1(:organization="e.organization")
        router-link(:to="`/events/${e.id}`" v-text="e.name")
      .events-list-location
        a(v-if="e.venue && e.venue.url" target="_blank" :href="e.venue.url" v-text="e.venue.name")
        span(v-else-if="e.venue" v-text="e.venue.name")
        span(v-else) TBD
</template>

<script>
import range from 'lodash/range'
import EventOrgDot from '@/base/EventOrgDot'

export default {
  components: { EventOrgDot },
  data: () => ({
    year: null,
    years: range(2013, new Date().getFullYear() + 2),
    events: null,
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
      if (data.canAdd) this.$emit('canAdd')
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
  display grid
  margin-top 1.5rem
  grid auto / min-content 1fr
  @media (min-width: 800px)
    grid auto / min-content 2fr 1fr
  @media (min-width: 1012px)
    grid unquote('auto / min-content fit-content(30rem) 1fr') // empirical
.events-list-heading
  display none
  @media (min-width: 576px)
    display block
    font-weight bold
.events-list-date
  flex none
  margin 0.25rem 0.75rem 0 0
  white-space nowrap
  font-variant tabular-nums
.events-list-year
  display none
  @media (min-width: 576px)
    display inline
.events-list-start
  display none
  @media (min-width: 576px)
    display inline
    padding-left 0.25rem
.events-list-event
  flex none
  overflow hidden
  margin-top 0.25rem
  text-overflow ellipsis
  white-space nowrap
.events-list-location
  display none
  @media (min-width: 800px)
    display block
    flex none
    overflow hidden
    margin-top 0.25rem
    padding-left 0.25rem
    text-overflow ellipsis
    white-space nowrap
</style>
