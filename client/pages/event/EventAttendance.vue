<!--
EventAttendance shows and allows changes to the attendance for an event.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!event")
  b-spinner(small)
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
  #event-attend-submit
    b-btn(type="submit" variant="primary") Save Attendance
    b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
import moment from 'moment-mini'
import EventAttendancePerson from './EventAttendancePerson'

export default {
  components: { EventAttendancePerson },
  props: {
    onLoadEvent: Function,
  },
  data: () => ({
    event: null,
    people: null,
    setType: 'Volunteer',
    setHours: 1.0,
  }),
  async created() {
    const data = (await this.$axios.get(`/api/events/${this.$route.params.id}?attendance=1`)).data
    this.event = data.event
    this.people = data.people
    this.setHours = moment(this.event.end, 'HH:mm').diff(moment(this.event.start, 'HH:mm'), 'hours', true)
    this.onLoadEvent(this.event)
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
      this.$router.push(`/events/${this.$route.params.id}`)
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
  padding 1.5rem 0.75rem
  @media (min-width: 576px)
    display grid
    height 100%
    grid max-content 1fr max-content / 100%
#event-attend-settings
  margin-bottom 0.75rem
  @media print
    display none
#event-attend-hours
  max-width 5rem
@media (min-width: 576px)
  #event-attend-group
    display flex
    flex-direction column
    flex-wrap wrap
    align-content flex-start
    min-height 0
    .custom-checkbox
      margin-right 1.5rem
#event-attend-submit
  margin-top 0.75rem
  @media print
    display none
</style>
