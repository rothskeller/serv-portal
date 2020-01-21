<!--
EventAttendance shows and allows changes to the attendance for an event.
-->

<template lang="pug">
form#event-attendance(v-else @submit.prevent="onSave")
  b-form#event-attend-settings(inline)
    label.mr-2(for="event-attend-type") Set attendance for:
    b-form-select#event-attend-type.mr-4(v-model="setType")
      option(value="Volunteer") Volunteer
      option(value="Student") Student
      option(value="Audit") Auditor
    label.mr-2(for="event-attend-hours") Hours:
    b-form-input#event-attend-hours(v-model="setHours" type="number" min="0" max="24" step="0.5")
  #event-attend-group
    EventAttendancePerson(v-for="p in people" :key="p.id" :person="p" @toggle="onTogglePerson")
  div.mt-3
    b-btn(type="submit" variant="primary") Save Attendance
    b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
import moment from 'moment'

export default {
  props: {
    event: null,
    people: null,
  },
  data: () => ({ setType: 'Volunteer', setHours: 1.0 }),
  mounted() {
    this.setHours = moment(this.event.end, 'HH:mm').diff(moment(this.event.start, 'HH:mm'), 'hours', true)
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onSave() {
      const body = new FormData
      this.people.forEach(p => {
        if (p.attended) {
          body.append('person', p.id)
          body.append('type', p.attended.type)
          body.append('minutes', p.attended.minutes)
        }
      })
      await this.$axios.post(`/api/events/${this.$route.params.id}/attendance`, body)
      this.$router.push({ path: '/events', params: { year: this.event.date.substr(0, 4) } })
    },
    onTogglePerson(person) {
      if (person.attended && person.attended.type === this.setType) person.attended = false
      else person.attended = { type: this.setType, minutes: 60 * this.setHours }
    },
  },
}
</script>

<style lang="stylus">
#event-attendance
  margin 1.5rem 0.75rem
#event-attend-settings
  margin-bottom 0.75rem
#event-attend-hours
  max-width 5rem
@media (min-width: 576px)
  #event-attend-group
    display flex
    flex-direction column
    flex-wrap wrap
    align-content flex-start
    height calc(100vh - 40px - (3.25rem + 2px) - 3rem - (3rem + 2px) - 1rem - (2.25rem + 2px))
    // title bar, tab bar, tab margin, settings form and margin, button margin, button
    .custom-checkbox
      margin-right 1.5rem
</style>
