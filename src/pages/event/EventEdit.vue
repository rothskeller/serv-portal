<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
form#event-edit(@submit.prevent="onSubmit")
  b-form-group(label="Event name" label-for="event-name" label-cols-sm="auto" label-class="event-edit-label" :state="nameError ? false : null" :invalid-feedback="nameError")
    b-input#event-name(:state="nameError ? false : null" trim v-model="event.name")
  b-form-group(label="Event date" label-for="event-date" label-cols-sm="auto" label-class="event-edit-label" :state="dateError ? false : null" :invalid-feedback="dateError")
    b-input#event-date(type="date" autofocus :state="dateError ? false : null" v-model="event.date")
  b-form-group(label="Event time" label-cols-sm="auto" label-class="event-edit-label" :state="timeError ? false : null" :invalid-feedback="timeError")
    #event-time
      b-input#event-start(type="time" :state="timeError ? false : null" v-model="event.start")
      span to
      b-input#event-end(type="time" :state="timeError ? false : null" v-model="event.end")
  b-form-group(label="Event type:" :state="typeError ? false : null" :invalid-feedback="typeError")
    b-form-radio-group(stacked :options="eventTypes" v-model="event.type")
  b-form-group(label="Event is for these roles:" :state="rolesError ? false : null" :invalid-feedback="rolesError")
    b-form-checkbox-group(stacked :options="roles" text-field="name" value-field="id" v-model="event.roles")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="event.id ? 'Save Event' : 'Create Event'")
    b-btn.ml-2(@click="onCancel") Cancel
    b-btn.ml-5(v-if="event.id" variant="danger" @click="onDelete") Delete Event
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
    event: null,
    roles: null,
  },
  data: () => ({
    eventTypes,
    submitted: false,
    dateError: null,
    timeError: null,
    nameError: null,
    duplicateName: null,
    typeError: null,
    rolesError: null,
    valid: true,
  }),
  watch: {
    'event.date': 'validate',
    'event.start': 'validate',
    'event.end': 'validate',
    'event.name': 'validate',
    'event.type': 'validate',
    'event.roles': 'validate',
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onDelete() {
      const resp = await this.$bvModal.msgBoxConfirm(
        'Are you sure you want to delete this event?  All associated data, including attendance records, will be permanently lost.', {
        title: 'Delete Event', headerBgVariant: 'danger', headerTextVariant: 'white',
        okTitle: 'Delete', okVariant: 'danger', cancelTitle: 'Keep',
      }).catch(err => { })
      if (!resp) return
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/events/${this.$route.params.id}`, body)
      this.$router.push({ path: '/events', params: { year: this.event.date.substr(0, 4) } })
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('date', this.event.date)
      body.append('start', this.event.start)
      body.append('end', this.event.end)
      body.append('name', this.event.name)
      body.append('type', this.event.type)
      this.event.roles.forEach(r => { body.append('role', r) })
      const resp = (await this.$axios.post(`/api/events/${this.$route.params.id}`, body)).data
      if (resp && resp.nameError)
        this.duplicateName = { date: this.event.date, name: this.event.name }
      else
        this.$router.push({ path: '/events', params: { year: this.event.date.substr(0, 4) } })
    },
    validate() {
      if (!this.submitted) return
      if (!this.event.date)
        this.dateError = 'The event date is required.'
      else if (!this.event.date.match(/^20\d\d-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$/))
        this.dateError = 'This is not a valid date.'
      else
        this.dateError = null
      if (!this.event.start || !this.event.end)
        this.timeError = 'The event times are required.'
      else if (!this.event.start.match(/^(?:[01][0-9]|2[0-3]):[0-5][0-9]$/))
        this.timeError = 'The start time is not valid.'
      else if (!this.event.end.match(/^(?:[01][0-9]|2[0-3]):[0-5][0-9]$/))
        this.timeError = 'The end time is not valid.'
      else if (this.event.end < this.event.start)
        this.timeError = 'The end time must come after the start time.'
      else
        this.timeError = null
      if (!this.event.name)
        this.nameError = 'The event name is required.'
      else if (this.duplicateName && this.duplicateName.date === this.event.date && this.duplicateName.name === this.event.name)
        this.nameError = 'Another event on this date has this name.'
      else
        this.nameError = this.duplicateName = null
      if (!eventTypes[this.event.type])
        this.typeError = 'The event type is required.'
      else
        this.typeError = null
      if (!this.event.roles.length)
        this.rolesError = 'At least one role must be selected.'
      else
        this.rolesError = null
      this.valid = !this.dateError && !this.timeError && !this.nameError && !this.typeError && !this.rolesError
    },
  },
}
</script>

<style lang="stylus">
#event-edit
  margin 1.5rem 0.75rem
.event-edit-label
  width 7rem
#event-date, #event-name, #event-type, #event-roles
  min-width 14rem
  max-width 20rem
#event-time
  display flex
  align-items baseline
  min-width 14rem
  max-width 20rem
  span
    padding 0 0.5rem
</style>
