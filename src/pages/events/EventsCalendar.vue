<!--
EventsCalendar displays the events in a calendar form.
-->

<template lang="pug">
#events-calendar
  #events-calendar-grid
    .events-calendar-arrow(@click="onYearBackward")
      svg.events-calendar-year-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512")
        path(fill="currentColor" d="M34.5 239L228.9 44.7c9.4-9.4 24.6-9.4 33.9 0l22.7 22.7c9.4 9.4 9.4 24.5 0 33.9L131.5 256l154 154.7c9.3 9.4 9.3 24.5 0 33.9l-22.7 22.7c-9.4 9.4-24.6 9.4-33.9 0L34.5 273c-9.3-9.4-9.3-24.6 0-34zm192 34l194.3 194.3c9.4 9.4 24.6 9.4 33.9 0l22.7-22.7c9.4-9.4 9.4-24.5 0-33.9L323.5 256l154-154.7c9.3-9.4 9.3-24.5 0-33.9l-22.7-22.7c-9.4-9.4-24.6-9.4-33.9 0L226.5 239c-9.3 9.4-9.3 24.6 0 34z")
    .events-calendar-arrow(@click="onMonthBackward")
      svg.events-calendar-month-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512")
        path(fill="currentColor" d="M34.52 239.03L228.87 44.69c9.37-9.37 24.57-9.37 33.94 0l22.67 22.67c9.36 9.36 9.37 24.52.04 33.9L131.49 256l154.02 154.75c9.34 9.38 9.32 24.54-.04 33.9l-22.67 22.67c-9.37 9.37-24.57 9.37-33.94 0L34.52 272.97c-9.37-9.37-9.37-24.57 0-33.94z")
    #events-calendar-month(v-text="month.format('MMMM YYYY')")
    .events-calendar-arrow(@click="onMonthForward")
      svg.events-calendar-month-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512")
        path(fill="currentColor" d="M285.476 272.971L91.132 467.314c-9.373 9.373-24.569 9.373-33.941 0l-22.667-22.667c-9.357-9.357-9.375-24.522-.04-33.901L188.505 256 34.484 101.255c-9.335-9.379-9.317-24.544.04-33.901l22.667-22.667c9.373-9.373 24.569-9.373 33.941 0L285.475 239.03c9.373 9.372 9.373 24.568.001 33.941z")
    #events-calendar-arrow-last.events-calendar-arrow(@click="onYearForward")
      svg.events-calendar-year-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512")
        path(fill="currentColor" d="M477.5 273L283.1 467.3c-9.4 9.4-24.6 9.4-33.9 0l-22.7-22.7c-9.4-9.4-9.4-24.5 0-33.9l154-154.7-154-154.7c-9.3-9.4-9.3-24.5 0-33.9l22.7-22.7c9.4-9.4 24.6-9.4 33.9 0L477.5 239c9.3 9.4 9.3 24.6 0 34zm-192-34L91.1 44.7c-9.4-9.4-24.6-9.4-33.9 0L34.5 67.4c-9.4 9.4-9.4 24.5 0 33.9l154 154.7-154 154.7c-9.3 9.4-9.3 24.5 0 33.9l22.7 22.7c9.4 9.4 24.6 9.4 33.9 0L285.5 273c9.3-9.4 9.3-24.6 0-34z")
    .events-calendar-weekday(v-for="w in ['S','M','T','W','T','F','S']" v-text="w")
    .events-calendar-day(v-for="date in dates" :class="date ? null : 'empty'" @mouseover="onHoverDate(date)" @mouseout="onNoHoverDate" @click="onClickDate(date)")
      div(v-text="date ? date.date() : null")
      .events-calendar-day-dots
        EventTypeDots(:types="typesOn(date)")
      .events-calendar-day-event(v-for="event in eventsOn(date)" :key="event.id")
        EventTypeDots.mr-1(:types="event.types")
        span(v-if="$store.state.touch" v-text="event.name")
        b-link(v-else :to="`/events/${event.id}`" :title="event.name" v-text="event.name")
  #events-calendar-footer(v-if="date")
    #events-calendar-date(v-text="date.format('dddd, MMMM D, YYYY')")
    .events-calendar-event(v-for="e in eventsOn(date)" :key="e.id")
      EventTypeDots.mr-1(:types="e.types")
      b-link(:to="`/events/${e.id}`" v-text="e.name")
    .events-calendar-event(v-if="!eventsOn(date).length") No events scheduled.
</template>

<script>
import moment from 'moment'

export default {
  data: () => ({
    month: moment(),
    dates: [],
    year: null,
    events: null,
    date: null,
    clicked: false,
  }),
  mounted() {
    this.newMonth()
  },
  methods: {
    eventsOn(date) {
      return date ? this.events[date.format('YYYY-MM-DD')] || [] : []
    },
    typesOn(date) {
      const types = {}
      this.eventsOn(date).forEach(e => { e.types.forEach(t => { types[t] = true }) })
      return Object.keys(types)
    },
    async newMonth() {
      if (!this.year || this.year != this.month.year()) {
        const data = (await this.$axios.get(`/api/events?year=${this.month.year()}`)).data
        if (data.canAdd) this.$emit('canAdd')
        const events = {}
        data.events.forEach(e => {
          if (!events[e.date]) events[e.date] = []
          events[e.date].push(e)
        })
        this.events = events
      }
      const dates = []
      const start = this.month.clone().startOf('month')
      start.subtract(start.day(), 'days')
      const end = this.month.clone().endOf('month')
      end.add(6 - end.day(), 'days')
      for (let date = start; !date.isAfter(end, 'day'); date = date.clone().add(1, 'day')) {
        dates.push(date.isSame(this.month, 'month') ? date : null)
      }
      this.dates = dates
      this.date = null
      this.clicked = false
    },
    groupToClass(group) { return group.toLowerCase().replace(' ', '-') },
    onClickDate(date) { this.clicked = true, this.date = date },
    onHoverDate(date) { if (!this.clicked) this.date = date },
    onMonthBackward() {
      this.month.subtract(1, 'month')
      this.newMonth()
    },
    onMonthForward() {
      this.month.add(1, 'month')
      this.newMonth()
    },
    onNoHoverDate() { if (!this.clicked) this.date = null },
    onYearBackward() {
      this.month.subtract(1, 'year')
      this.newMonth()
    },
    onYearForward() {
      this.month.add(1, 'year')
      this.newMonth()
    },
  },
}
</script>

<style lang="stylus">
sidebarWidth = 7rem // width of the sidebar that appears on screens wider than labelThreshold
smallSideMargin = 0.5rem // side margin of the calendar on screens narrower than dotsMaxWidth
dotsMaxWidth = 360px // maximum width of the calendar when in dots mode
labelThreshold = 576px // threshold where we switch from dots mode to label mode
minTouch = 40px // minimum touchable size
arrowSize = minTouch // width and height of arrow boxes
printMargin = 1rem // margin around calendar when printing
#events-calendar
  margin 0 auto
  max-width dotsMaxWidth
  @media (min-width: labelThreshold)
    max-width none
#events-calendar-grid
  display flex
  flex-wrap wrap
  justify-content center
  margin-top 0.5rem
  @media (min-width: labelThreshold)
    border-left 1px solid #eee
  @media print
    margin printMargin
    border-top 1px solid #eee
.events-calendar-arrow
  display flex
  justify-content center
  align-items center
  width arrowSize
  height arrowSize
  cursor pointer
  user-select none
  &:hover
    background-color #efefef
#events-calendar-arrow-last
  @media (min-width: labelThreshold)
    border-right 1px solid #eee
.events-calendar-year-arrow
  width 1rem
  @media print
    display none
.events-calendar-month-arrow
  width 0.625rem
  @media print
    display none
#events-calendar-month
  display block
  align-self center
  width 'calc(100vw - 4 * %s - 2 * %s)' % (arrowSize smallSideMargin)
  text-align center
  font-size 1.2rem // lets longest month and year fit at 320px width
  @media (min-width: dotsMaxWidth)
    width 'calc(%s - 4 * %s)' % (dotsMaxWidth arrowSize) // 4*2.5rem arrows
  @media (min-width: labelThreshold)
    width 'calc(100vw - %s - 1px - 4 * %s)' % (sidebarWidth arrowSize)
  @media print
    width 'calc(100vw - 2 * %s - 1px - 4 * %s)' % (printMargin arrowSize)
.events-calendar-weekday
  margin-top 0.5rem
  padding 0 0 1rem
  width 'calc((100vw - (2 * %s)) / 7 - 0.02px)' % smallSideMargin
  // Different browsers round calculations different ways.  Firefox rounds to
  // 1/60th of a pixel.  So 4/7, which should be 0.428, gets rounded up to
  // 0.433.  That could cause Sunday to wrap to the next line, which we need to
  // avoid.  By subtracting 0.02px, the right margin could be off by up to
  // 0.14px, but at least it won't wrap.  Chrome uses 1/64th of a pixel, same
  // story.
  color #888
  text-align center
  line-height 1
  @media (min-width: dotsMaxWidth)
    width 'calc(%s / 7 - 0.02px)' % dotsMaxWidth
  @media (min-width: labelThreshold)
    width 'calc((100vw - %s - 1px) / 7 - 0.02px)' % sidebarWidth
    border-top 1px solid #eee
    border-right 1px solid #eee
    border-bottom 1px solid #eee
  @media print
    width 'calc((100vw - 2 * %s - 1px) / 7 - 0.02px)' % printMargin
.events-calendar-day
  padding 0 0 1rem
  min-height 3rem
  width 'calc((100vw - (2 * %s)) / 7 - 0.02px)' % smallSideMargin
  text-align center
  line-height 1
  &:hover
    background-color #efefef
    &.empty
      background-color white
  @media (min-width: dotsMaxWidth)
    width 'calc(%s / 7 - 0.02px)' % dotsMaxWidth
  @media (min-width: labelThreshold)
    min-height calc(1rem + 3 * 0.875rem + 1rem)
    width 'calc((100vw - %s - 1px) / 7 - 0.02px)' % sidebarWidth
    border-right 1px solid #eee
    border-bottom 1px solid #eee
    color #888
  @media print
    width 'calc((100vw - 2 * %s - 1px) / 7 - 0.02px)' % printMargin
.events-calendar-day-dots
  @media (min-width: labelThreshold)
    display none
.events-calendar-dot
  display inline-block
  margin 0 0 0 2px
  width 10px
  height 10px
  border-radius 50%
  &:first-child
    margin 0
.events-calendar-serv
  background-color #a9a9a9
.events-calendar-cert-d
  background-color #469990
.events-calendar-cert-t
  border 2px solid #469990
  background-color white
.events-calendar-listos
  background-color #f032e6
.events-calendar-outreach
  background-color #800000
.events-calendar-pep
  border 2px solid #f032e6
  background-color white
.events-calendar-sares
  background-color #f58231
.events-calendar-scc-ares
  border 2px solid #f58231
  background-color white
.events-calendar-snap
  background-color #4363d8
  // colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/
.events-calendar-day-event
  display none
  overflow hidden
  margin 0 0.25rem
  text-align left
  text-overflow ellipsis
  white-space nowrap
  font-size 0.875rem
  line-height 1.2
  @media (min-width: labelThreshold)
    display block
    color #000
  @media print
    a
      color #000
      text-decoration none
#events-calendar-footer
  margin 0 0.75rem
  @media (min-width: labelThreshold)
    margin 0.75rem
    .mouse &
      display none
#events-calendar-date
  font-weight bold
.events-calendar-event
  .touch &
    margin 'calc((%s - 1.5rem) / 2) 0' % minTouch
</style>
