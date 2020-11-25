<!--
EventsCalendar displays the events in a calendar form.
-->

<template lang="pug">
#events-calendar
  #events-calendar-grid
    #events-calendar-heading
      MonthSelect(v-model='monthfmt')
    .events-calendar-weekday(v-for='w in ["S", "M", "T", "W", "T", "F", "S"]', v-text='w')
    .events-calendar-day(
      v-for='date in dates',
      :class='date ? null : "empty"',
      @mouseover='onHoverDate(date)',
      @mouseout='onNoHoverDate',
      @click='onClickDate(date)'
    )
      div(v-text='date ? date.date() : null')
      .events-calendar-day-dots
        EventOrgDot(v-for='e in eventsOn(date)', :key='e.id', :org='e.org')
      .events-calendar-day-event(v-for='event in eventsOn(date)', :key='event.id')
        EventOrgDot.events-calendar-dot-space(:org='event.org')
        span(v-if='touch', v-text='event.name')
        router-link(v-else, :to='`/events/${event.id}`', :title='event.name', v-text='event.name')
  #events-calendar-footer(v-if='showDate')
    #events-calendar-date(v-text='showDate.format("dddd, MMMM D, YYYY")')
    .events-calendar-event(v-for='e in eventsOn(showDate)', :key='e.id')
      EventOrgDot.mr-1(:org='e.org')
      router-link(:to='`/events/${e.id}`', v-text='e.name')
    .events-calendar-event(v-if='!eventsOn(showDate).length') No events scheduled.
</template>

<script lang="ts">
import { defineComponent, inject, Ref, ref, onMounted, watch } from 'vue'
import Cookies from 'js-cookie'
import moment, { Moment } from 'moment-mini'
import axios from '../../plugins/axios'
import { EventOrgDot, MonthSelect } from '../../base'
import setPage from '../../plugins/page'

export type ListEventsVenue = {
  id: number
  name: string
  url: string
}
export type ListEventsEvent = {
  id: number
  name: string
  date: string
  start: string
  venue: null | ListEventsVenue
  org: string
  type: string
}
export type ListEvents = {
  canAdd: boolean
  events: Array<ListEventsEvent>
}

export default defineComponent({
  components: { EventOrgDot, MonthSelect },
  setup(props, { emit }) {
    setPage({ title: 'Events' })

    const touch = inject<Ref<boolean>>('touch')!

    // Choose the month being viewed.
    const month = ref(moment())
    const monthfmt = ref(month.value.format('YYYY-MM'))
    watch(monthfmt, () => {
      console.log('monthfmt', monthfmt.value)
      month.value = moment(monthfmt.value, 'YYYY-MM')
      newMonth()
    })

    // Should we show the events on a particular date under the table?
    const showDate = ref(null as null | Moment)
    const clicked = ref(false)
    function onClickDate(date: Moment) {
      clicked.value = true
      showDate.value = date
    }
    function onHoverDate(date: Moment) {
      if (!clicked.value) showDate.value = date
    }
    function onNoHoverDate() {
      if (!clicked.value) showDate.value = null
    }

    // Load data when a new month is selected.
    let yearLoaded = null as null | number
    let yearEvents = {} as Record<string, Array<ListEventsEvent>>
    const dates = ref([] as Array<null | Moment>)
    async function newMonth() {
      Cookies.set('serv-events-month', month.value.format('YYYY-MM'))
      if (!yearLoaded || yearLoaded != month.value.year()) {
        const data: ListEvents = (await axios.get(`/api/events?year=${month.value.year()}`)).data
        yearEvents = {}
        data.events.forEach((e) => {
          if (!yearEvents[e.date]) yearEvents[e.date] = []
          yearEvents[e.date].push(e)
        })
        yearLoaded = month.value.year()
      }
      dates.value = []
      const start = month.value.clone().startOf('month')
      start.subtract(start.day(), 'days')
      const end = month.value.clone().endOf('month')
      end.add(6 - end.day(), 'days')
      for (let date = start; !date.isAfter(end, 'day'); date = date.clone().add(1, 'day')) {
        dates.value.push(date.isSame(month.value, 'month') ? date : null)
      }
      showDate.value = null
      clicked.value = false
    }
    function eventsOn(date: Moment) {
      return date ? yearEvents[date.format('YYYY-MM-DD')] || [] : []
    }

    // On mount, switch to the last month viewed if any, and remember the fact
    // that this user prefers the calendar view.
    onMounted(() => {
      const monthCookie = Cookies.get('serv-events-month')
      if (monthCookie) month.value = moment(monthCookie, 'YYYY-MM')
      newMonth()
      Cookies.set('serv-events-page', 'calendar', { expires: 3650 })
    })

    return {
      month,
      monthfmt,
      dates,
      onHoverDate,
      onNoHoverDate,
      onClickDate,
      eventsOn,
      touch,
      showDate,
    }
  },
})
</script>

<style lang="postcss">
#events-calendar {
  --dotsMaxWidth: 360px; /* maximum width of the calendar when in dots mode */
  --minTouch: 40px; /* minimum touchable size */
  --arrowSize: var(--minTouch); /* width and height of arrow boxes */
  --printMargin: 1rem; /* margin around calendar when printing */
  display: flex;
  flex-direction: column;
  align-items: center;
  @media print {
    padding: var(--printMargin);
  }
}
#events-calendar-grid {
  display: grid;
  justify-content: center;
  margin-top: 0.5rem;
  max-width: var(--dotsMaxWidth);
  width: 100%;
  grid: auto / repeat(7, 14.2857%);
  @media (min-width: 576px) {
    max-width: none;
    border-left: 1px solid #eee;
  }
  @media print {
    border-top: 1px solid #eee;
  }
}
#events-calendar-heading {
  display: flex;
  grid-area: 1 / 1 / 2 / 8;
}
.events-calendar-weekday {
  margin-top: 0.5rem;
  padding: 0 0 1rem;
  color: #888;
  text-align: center;
  line-height: 1;
  @media (min-width: 576px) {
    border-top: 1px solid #eee;
    border-right: 1px solid #eee;
    border-bottom: 1px solid #eee;
  }
}
.events-calendar-day {
  padding: 0 0 1rem;
  min-height: 3rem;
  text-align: center;
  line-height: 1;
  &:hover {
    background-color: #efefef;
    &.empty {
      background-color: white;
    }
  }
  @media (min-width: 576px) {
    min-height: calc(1rem + 3 * 0.875rem + 1rem);
    border-right: 1px solid #eee;
    border-bottom: 1px solid #eee;
    color: #888;
  }
}
.events-calendar-day-dots {
  @media (min-width: 576px) {
    display: none;
  }
}
.events-calendar-dot-space {
  margin-right: 0.25rem;
}
.events-calendar-day-event {
  display: none;
  overflow: hidden;
  margin: 0 0.25rem;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.875rem;
  line-height: 1.2;
  @media (min-width: 576px) {
    display: block;
    color: #000;
  }
  @media print {
    a {
      color: #000;
      text-decoration: none;
    }
  }
}
#events-calendar-footer {
  padding: 0 0.75rem;
  max-width: var(--dotsMaxWidth);
  width: 100%;
  @media (min-width: 576px) {
    padding: 0.75rem;
    .mouse & {
      display: none;
    }
  }
}
#events-calendar-date {
  font-weight: bold;
}
.events-calendar-event {
  .touch & {
    margin: calc((var(--minTouch) - 1.5rem) / 2) 0;
  }
}
</style>
