<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
#event-view
  b-form-group(label="Event name" label-for="event-name" label-cols-sm="auto" label-class="event-edit-label")
    b-input#event-name(plaintext :value="event.name")
  b-form-group(label="Event date" label-for="event-date" label-cols-sm="auto" label-class="event-edit-label")
    b-input#event-date(type="date" plaintext :value="event.date")
  b-form-group(label="Event time" label-for="event-time" label-cols-sm="auto" label-class="event-edit-label")
    b-input#event-time(plaintext :value="`${event.start} to ${event.end}`")
  b-form-group(label="Event type" label-for="event-type" label-cols-sm="auto" label-class="event-edit-label")
    b-input#event-type(plaintext :value="eventTypes[event.type]")
  b-form-group(label="Roles" label-for="event-roles" label-cols-sm="auto" label-class="event-edit-label")
    b-textarea#event-roles(plaintext v-text="eventRoleList")
</template>

<script>
const eventTypes = {
  'Train': 'Training',
  'Drill': 'Drill',
  'Civic': 'Civic Event',
  'Incid': 'Incident',
  'CE': 'Continuing Ed',
  'Meeting': 'Meeting',
  'Class': 'Class',
}

export default {
  props: {
    event: Object,
    roles: Array,
  },
  data: () => ({
    eventTypes,
  }),
  computed: {
    eventRoleList() {
      const list = []
      this.event.roles.forEach(tid => {
        list.push(this.roles.find(t => t.id === tid).name)
      })
      return list.join('\n')
    },
  },
}
</script>

<style lang="stylus">
#event-view
  margin 1.5rem 0.75rem
</style>
