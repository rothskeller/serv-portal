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
        EventTypeDots.mr-1(:types="e.types")
        router-link(:to="`/events/${e.id}`" v-text="e.name")
      .events-list-location
        a(v-if="e.venue && e.venue.url" target="_blank" :href="e.venue.url" v-text="e.venue.name")
        span(v-else-if="e.venue" v-text="e.venue.name")
        span(v-else) TBD
</template>

<script>
import range from 'lodash/range'
import EventTypeDots from '@/base/EventTypeDots'

export default {
  components: { EventTypeDots },
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
  display flex
  flex-wrap wrap
  margin-top 1.5rem
.events-list-heading
  display none
  @media (min-width: 576px)
    display block
    font-weight bold
.events-list-date
  flex none
  margin-top 0.25rem
  width 3.5rem
  font-variant tabular-nums
  @media (min-width: 576px)
    width 10rem
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
  width calc(100vw - 5.5rem)
  text-overflow ellipsis
  white-space nowrap
  @media (min-width: 576px)
    width calc(100vw - 18.5rem)
    // 7rem sidebar
    // 1.5rem tab margins
    // 10rem date/time column
  @media (min-width: 800px)
    width calc(50vw - 9.25rem)
    // above divided by 2
  @media (min-width: 1280px)
    width 30.75rem
    // (80rem-18.5rem)/2
.events-list-location
  display none
  @media (min-width: 800px)
    display block
    flex none
    overflow hidden
    margin-top 0.25rem
    padding-left 0.25rem
    width calc(50vw - 9.25rem)
    // 7rem sidebar
    // 1.5rem tab margins
    // 10rem date/time column
    text-overflow ellipsis
    white-space nowrap
  @media (min-width: 1280px)
    width calc(100vw - 49.75rem)
</style>
