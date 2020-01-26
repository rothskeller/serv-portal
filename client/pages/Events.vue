<!--
Events displays the list of events.
-->

<template lang="pug">
Page(title="Events" menuItem="events" noPadding)
  b-card#events-card(no-body)
    b-tabs(card v-model="tab")
      b-tab.events-tab-pane(title="Calendar" no-body)
        EventsCalendar(@canAdd="canAdd = true")
      b-tab.events-tab-pane(title="List" no-body)
        EventsList(@canAdd="canAdd = true")
      template(v-if="canAdd" v-slot:tabs-end)
        b-nav-item(to="/events/NEW") Add Event
</template>

<script>
import Cookies from 'js-cookie'

export default {
  data: () => ({ tab: 0, canAdd: false }),
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
