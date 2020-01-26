<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
Page(:title="title" menuItem="events" noPadding)
  div.mt-3(v-if="loading")
    b-spinner(small)
  b-card#event-card(v-else-if="tabs" no-body)
    b-tabs(card)
      b-tab.event-tab-pane(v-if="!newe" title="Details" no-body)
        EventView(:event="event")
      b-tab.event-tab-pane(v-if="canEdit" title="Edit" no-body)
        EventEdit(:event="event" :groups="groups" :venues="venues" :types="types")
      b-tab.event-tab-pane(v-if="canAttendance" title="Attendance" no-body)
        EventAttendance(:event="event" :people="people")
  EventEdit(v-else-if="canEdit" :event="event" :groups="groups" :venues="venues" :types="types")
  EventView(v-else :event="event")
</template>

<script>
export default {
  data: () => ({
    loading: false,
    title: 'Event',
    canEdit: false,
    canAttendance: false,
    event: null,
    groups: null,
    types: null,
    venues: null,
    people: null,
  }),
  computed: {
    newe() { return this.$route.params.id === 'NEW' },
    tabs() {
      return (this.newe ? 0 : 1) + (this.canEdit ? 1 : 0) + (this.canAttendance ? 1 : 0) > 1
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/events/${this.$route.params.id}`)).data
    this.canEdit = data.canEdit
    this.canAttendance = data.canAttendance
    this.event = data.event
    this.title = data.event.id ? `${data.event.date} ${data.event.name}` : 'New Event'
    this.groups = data.groups
    this.venues = data.venues
    this.types = data.types
    this.people = data.people
    this.loading = false
  },
}
</script>

<style lang="stylus">
#event-card
  height calc(100vh - 40px)
  border none
.event-tab-pane
  overflow-y auto
  height calc(100vh - 3.25rem - 42px)
</style>
