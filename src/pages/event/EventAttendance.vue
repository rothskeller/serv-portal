<!--
EventAttendance shows and allows changes to the attendance for an event.
-->

<template lang="pug">
form#event-attendance(v-else @submit.prevent="onSave")
  b-form-group(label="This event was attended by:")
    b-form-checkbox-group#event-attend-group(stacked :options="people" text-field="sortName" value-field="id" v-model="attendees")
  div.mt-3
    b-btn(type="submit" variant="primary") Save Attendance
    b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
export default {
  props: {
    event: null,
    people: null,
  },
  data: () => ({ attendees: null }),
  mounted() {
    const attendees = []
    this.people.forEach(p => {
      if (p.attended) attendees.push(p.id)
    })
    this.attendees = attendees
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onSave() {
      const body = new FormData
      this.attendees.forEach(pid => { body.append('person', pid) })
      await this.$axios.post(`/api/events/${this.$route.params.id}/attendance`, body)
      this.$router.push({ path: '/events', params: { year: this.event.date.substr(0, 4) } })
    },
  },
}
</script>

<style lang="stylus">
#event-attendance
  margin 1.5rem 0.75rem
@media (min-width: 576px)
  #event-attend-group
    display flex
    flex-direction column
    flex-wrap wrap
    align-content flex-start
    height calc(100vh - 40px - (3.25rem + 2px) - 3rem - (1.875rem + 1px) - 1rem - (2.25rem + 2px))
    // title bar, tab bar, tab margin, label, button margin, button
    .custom-checkbox
      margin-right 1.5rem
</style>
