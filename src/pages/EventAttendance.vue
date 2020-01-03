<!--
EventAttendance shows and allows changes to the attendance for an event.
-->

<template lang="pug">
Page(:title="title" subtitle="Event Attendance" menuItem="events")
  div.mt-3(v-if="loading")
    b-spinner(small)
  form.mt-3(v-else @submit.prevent="onSave")
    b-form-group(label="This event was attended by:")
      b-form-checkbox-group(stacked :options="people" text-field="name" value-field="id" v-model="attendees")
    div.mt-3
      b-btn(type="submit" variant="primary") Save Attendance
      b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
export default {
  data: () => ({ loading: false, event: null, people: null, attendees: null }),
  computed: {
    title() {
      if (this.loading) return 'Event Attendance'
      return `${this.event.date} ${this.event.name}`
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/events/${this.$route.params.id}/attendance`)).data
    this.event = data.event
    const attendees = []
    data.people.forEach(p => {
      p.name = `${p.lastName}, ${p.firstName}`
      if (p.attended) attendees.push(p.id)
    })
    this.people = data.people
    this.attendees = attendees
    this.loading = false
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
