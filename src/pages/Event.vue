<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
Page(:title="title" menuItem="events" noPadding)
  div.mt-3(v-if="loading")
    b-spinner(small)
  b-card#event-card(v-else-if="canEdit && canAttendance" no-body)
    b-tabs(card)
      b-tab.event-tab-pane(title="Details" no-body)
        EventEdit(:event="event" :roles="roles")
      b-tab.event-tab-pane(title="Attendance" no-body)
        EventAttendance(:event="event" :people="people")
  b-card#event-card(v-else-if="!canEdit && canAttendance" no-body)
    b-tabs(card)
      b-tab.event-tab-pane(title="Details" no-body)
        EventView(:event="event" :roles="roles")
      b-tab.event-tab-pane(title="Attendance" no-body)
        EventAttendance(:event="event" :people="people")
  EventEdit(v-else-if="canEdit" :event="event" :roles="roles")
  EventView(v-else :event="event" :roles="roles")
</template>

<script>
export default {
  data: () => ({
    loading: false,
    title: 'Event',
    canEdit: false,
    canAttendance: false,
    event: null,
    roles: null,
    people: null,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/events/${this.$route.params.id}`)).data
    this.canEdit = data.canEdit
    this.canAttendance = data.canAttendance
    this.event = data.event
    this.title = data.event.id ? `${data.event.date} ${data.event.name}` : 'New Event'
    this.roles = data.roles
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
