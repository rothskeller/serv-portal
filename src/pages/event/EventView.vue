<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
#event-view
  #event-view-name(v-text="event.name")
  #event-view-groups(v-if="event.servGroups.length" v-text="event.servGroups.join(', ')")
  #event-view-date-time(v-text="dateTimeFmt")
  #event-view-venue
    #event-view-venue-name(v-text="event.venue ? event.venue.name : 'Location TBD'")
    #event-view-venue-address(v-if="event.venue")
      span(v-text="event.venue.address")
      span(v-if="event.venue.city" v-text="`, ${event.venue.city}`")
      span#event-view-venue-map(v-if="event.venue.url")
        |  (
        a(target="_blank" :href="event.venue.url") map
        | )
  #event-view-details(v-if="event.details" v-html="event.details")
</template>

<script>
import moment from 'moment'

export default {
  props: {
    event: Object,
  },
  computed: {
    dateTimeFmt() {
      const date = moment(this.event.date, 'YYYY-MM-DD')
      const start = moment(this.event.start, 'HH:mm')
      const end = moment(this.event.end, 'HH:mm')
      if (start.format('a') !== end.format('a'))
        return `${date.format('dddd, MMMM D, YYYY')}\n${start.format('h:mma')} to ${end.format('h:mma')}`
      else
        return `${date.format('dddd, MMMM D, YYYY')}\n${start.format('h:mm')} to ${end.format('h:mma')}`
    },
  },
}
</script>

<style lang="stylus">
#event-view
  margin 1.5rem 0.75rem
#event-view-name
  font-weight bold
  font-size 1.25rem
  line-height 1.2
#event-view-groups
  color #888
#event-view-date-time
  margin-top 0.75rem
  white-space pre-line
  line-height 1.2
#event-view-venue
  margin-top 0.75rem
  line-height 1.2
#event-view-venue-address
  font-size 0.875rem
#event-view-details
  margin-top 0.75rem
  max-width 40rem
  white-space pre-line
  line-height 1.2
</style>
