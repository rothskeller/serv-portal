<!--
EventsList displays the list of events.
-->

<template lang="pug">
#events-list
  #events-list-title
    | Events in
    |
    select#events-list-year(v-model='year')
      option(v-for='y in years', :key='y', :value='y', v-text='y')
  #events-list-spinner(v-if='loading')
    SSpinner
  #events-list-table(v-else)
    .events-list-date.events-list-heading Date
    .events-list-event.events-list-heading Event
    .events-list-location.events-list-heading Location
    template(v-for='e in events')
      .events-list-date
        span.events-list-year(v-text='e.date.substr(0, 5)')
        span(v-text='e.date.substr(5)')
        span.events-list-start(v-text='e.start')
      .events-list-event
        EventOrgDot(:organization='e.organization')
        router-link(:to='`/events/${e.id}`', v-text='e.name')
      .events-list-location
        a(
          v-if='e.venue && e.venue.url',
          target='_blank',
          :href='e.venue.url',
          v-text='e.venue.name'
        )
        span(v-else-if='e.venue', v-text='e.venue.name')
        span(v-else) TBD
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue'
import Cookies from 'js-cookie'
import moment from 'moment-mini'
import axios from '../../plugins/axios'
import { EventOrgDot, SSpinner } from '../../base'
import type { ListEvents, ListEventsEvent } from './EventsCalendar.vue'

export default defineComponent({
  components: { EventOrgDot, SSpinner },
  setup() {
    // The set of years that the user can choose to view.
    const years = [] as Array<number>
    for (let y = 2013; y <= moment().year() + 1; y++) years.push(y)

    // Figure out which year to show.
    let month = Cookies.get('serv-events-month')
    if (!month) {
      month = moment().format('YYYY-MM')
      Cookies.set('serv-events-month', month)
    }
    const year = ref(month.substr(0, 4))

    // Load the year data.
    const events = ref([] as Array<ListEventsEvent>)
    const loading = ref(true)
    async function load() {
      loading.value = true
      const data: ListEvents = (await axios.get(`/api/events?year=${year.value}`)).data
      events.value = data.events
      loading.value = false
      // Store the desired month in case they switch to calendar view.  Of
      // this also remembers the year, in case they return to this page.
      if (parseInt(year.value) == moment().year())
        Cookies.set('serv-events-month', moment().format('YYYY-MM'))
      else if (parseInt(year.value) < moment().year())
        Cookies.set('serv-events-month', `${year.value}-12`)
      else Cookies.set('serv-events-month', `${year.value}-01`)
    }
    load()
    watch(year, load)

    // Remember the fact that this user prefers the list view.
    Cookies.set('serv-events-page', 'list', { expires: 3650 })

    return { year, years, events, loading }
  },
})
</script>

<style lang="postcss">
#events-list {
  padding: 1.5rem 0.75rem;
}
#events-list-title {
  display: flex;
  align-items: center;
  font-size: 1.5rem;
}
#events-list-year {
  margin-left: 1rem;
  font-size: 1rem;
}
#events-list-spinner {
  margin-top: 1.5rem;
}
#events-list-table {
  display: grid;
  margin-top: 1.5rem;
  grid: auto / min-content 1fr;
  @media (min-width: 800px) {
    grid: auto / min-content 2fr 1fr;
  }
  @media (min-width: 1012px) {
    grid: auto / min-content fit-content(30rem) 1fr; /* empirical */
  }
}
.events-list-heading {
  display: none;
  @media (min-width: 576px) {
    display: block;
    font-weight: bold;
  }
}
.events-list-date {
  flex: none;
  margin: 0.25rem 0.75rem 0 0;
  white-space: nowrap;
  font-variant: tabular-nums;
}
.events-list-year {
  display: none;
  @media (min-width: 576px) {
    display: inline;
  }
}
.events-list-start {
  display: none;
  @media (min-width: 576px) {
    display: inline;
    padding-left: 0.25rem;
  }
}
.events-list-event {
  flex: none;
  overflow: hidden;
  margin-top: 0.25rem;
  text-overflow: ellipsis;
  white-space: nowrap;
  & .dot {
    margin-right: 0.25rem;
  }
}
.events-list-location {
  display: none;
  @media (min-width: 800px) {
    display: block;
    flex: none;
    overflow: hidden;
    margin-top: 0.25rem;
    padding-left: 0.25rem;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
