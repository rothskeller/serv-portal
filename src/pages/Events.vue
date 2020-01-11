<!--
Events displays the list of events.
-->

<template lang="pug">
Page(title="Events" menuItem="events" noPadding)
  b-card#events-card(no-body)
    b-tabs(card v-model="tab")
      b-tab.events-tab-pane(title="Calendar" no-body)
        EventsCalendar
      b-tab.events-tab-pane(title="List" no-body)
        EventsList
</template>

<script>
import Cookies from 'js-cookie'

export default {
  data: () => ({ tab: 0 }),
  mounted() {
    this.tab = Cookies.get('serv-events-list') === '1' ? 1 : 0
  },
  watch: {
    tab() {
      Cookies.set('serv-events-list', this.tab, { expires: 3650 })
    },
  },
}
</script>

<style lang="stylus">
#events-card
  height calc(100vh - 40px)
  border none
  .card-header
    @media print
      display none
.events-tab-pane
  overflow-y auto
  height calc(100vh - 3.25rem - 42px)
</style>
