<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!event")
  b-spinner(small)
form#event-edit(v-else @submit.prevent="onSubmit")
  b-form-group(label="Event name" label-for="event-name" label-cols-sm="auto" label-class="event-edit-label" :state="nameError ? false : null" :invalid-feedback="nameError")
    b-input#event-name(autofocus :state="nameError ? false : null" trim v-model="event.name")
  b-form-group(label="Event type" label-for="event-type" label-cols-sm="auto" label-class="event-edit-label" :state="typeError ? false : null" :invalid-feedback="typeError")
    b-form-select#event-type(:options="types" v-model="event.type" :state="typeError ? false : null")
  b-form-group(label="Organization" label-for="event-organization" label-cols-sm="auto" label-class="event-edit-label" :state="organizationError ? false : null" :invalid-feedback="organizationError")
    b-form-select#event-organization(:options="organizations" v-model="event.organization" :state="organizationError ? false : null")
  b-form-group(label="Event date" label-for="event-date" label-cols-sm="auto" label-class="event-edit-label" :state="dateError ? false : null" :invalid-feedback="dateError")
    b-input#event-date(type="date" :state="dateError ? false : null" v-model="event.date")
  b-form-group(label="Event time" label-cols-sm="auto" label-class="event-edit-label" :state="timeError ? false : null" :invalid-feedback="timeError")
    #event-time
      b-input#event-start(type="time" :state="timeError ? false : null" v-model="event.start")
      span to
      b-input#event-end(type="time" :state="timeError ? false : null" v-model="event.end")
  b-form-group(label="Location" label-for="event-venue" label-cols-sm="auto" label-class="event-edit-label")
    b-select#event-venue(v-model="event.venue.id" :options="venues" text-field="name" value-field="id")
      template(v-slot:first)
        b-select-option(value="0") TBD
      b-select-option(value="NEW") (create a new venue)
  template(v-if="event.venue.id === 'NEW'")
    b-form-group(label="Venue name" label-for="venue-name" label-cols-sm="auto" label-class="event-edit-label" :state="venueNameError ? false : null" :invalid-feedback="venueNameError")
      b-input#venue-name(ref="venueName" :state="venueNameError ? false : null" trim v-model="event.venue.name")
    b-form-group(label="Address" label-for="venue-address" label-cols-sm="auto" label-class="event-edit-label" :state="venueAddressError ? false : null" :invalid-feedback="venueAddressError")
      b-input#venue-address(:state="venueAddressError ? false : null" trim v-model="event.venue.address")
    b-form-group(label="City" label-for="venue-city" label-cols-sm="auto" label-class="event-edit-label" :state="venueCityError ? false : null" :invalid-feedback="venueCityError")
      b-input#venue-city(:state="venueCityError ? false : null" trim v-model="event.venue.city")
    b-form-group(label="Map URL" label-for="venue-url" label-cols-sm="auto" label-class="event-edit-label" :state="venueURLError ? false : null" :invalid-feedback="venueURLError")
      b-input#venue-url(:state="venueURLError ? false : null" trim v-model="event.venue.url")
      b-form-text
        | Generate this by opening Google Maps and searching for the venue.
        | (Search by the venue’s name rather than its address, if it has a
        | name, so that Google’s sidebar of information about the venue
        | appears.)  Once the venue is found, zoom out a few steps until the
        | nearest freeways are shown — enough to give the viewer an idea of
        | where in the Bay Area the venue is.  Copy the URL from the browser
        | address bar and paste it here.
  b-form-group(label="Details" label-for="event-details" label-cols-sm="auto" label-class="event-edit-label")
    b-textarea#event-details(v-model="event.details" rows="3")
    b-form-text
      | This may contain HTML &lt;a&gt; tags for links, but no other tags.
  b-form-group(v-if="event.canEditDSWFlags" label="Flags" label-cols-sm="auto" label-class="event-edit-label pt-0")
    #event-edit-flags
      b-checkbox(v-model="event.renewsDSW") Attendance renews DSW registration
      b-checkbox(v-model="event.coveredByDSW") Event is covered by DSW insurance
  b-form-group(label="Event is for these groups:" :state="groupsError ? false : null" :invalid-feedback="groupsError")
    b-form-checkbox-group(stacked :options="filteredGroups" text-field="name" value-field="id" v-model="event.groups")
  b-form-group(label="Visibility" label-for="event-private" label-cols-sm="auto" label-class="event-edit-label pt-0")
    b-form-radio-group(stacked v-model="event.private")
      b-form-radio(:value="false") Visible to everyone
      b-form-radio(:value="true") Visible only to selected groups
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="event.id ? 'Save Event' : 'Create Event'")
    b-btn.ml-2(@click="onCancel") Cancel
    b-btn.ml-5(v-if="event.id" variant="danger" @click="onDelete") Delete Event
</template>

<script>
import Cookies from 'js-cookie'

export default {
  props: {
    onLoadEvent: Function,
  },
  data: () => ({
    event: null,
    groups: null,
    organizations: null,
    venues: null,
    types: null,
    submitted: false,
    dateError: null,
    timeError: null,
    nameError: null,
    duplicateName: null,
    venueNameError: null,
    venueAddressError: null,
    venueCityError: null,
    venueURLError: null,
    organizationError: null,
    typeError: null,
    groupsError: null,
    valid: true,
  }),
  computed: {
    filteredGroups() { return this.groups.filter(f => f.organization === this.event.organization || !f.organization) }
  },
  watch: {
    'event.name': 'validate',
    'event.date': 'validate',
    'event.start': 'validate',
    'event.end': 'validate',
    'event.venue.name': 'validate',
    'event.venue.address': 'validate',
    'event.venue.city': 'validate',
    'event.venue.url': 'validate',
    'event.organization': 'validate',
    'event.type': 'validate',
    'event.groups': 'validate',
    'event.venue.id'() {
      if (this.event.venue.id === 'NEW') {
        this.event.venue.name = this.event.venue.address = this.event.venue.city = this.event.venue.url = ''
        this.validate()
        this.$nextTick(() => { this.$refs.venueName.focus() })
      }
    },
  },
  async created() {
    const data = (await this.$axios.get(`/api/events/${this.$route.params.id}?edit=1`)).data
    this.event = data.event
    this.groups = data.groups
    this.venues = data.venues
    this.types = data.types
    this.organizations = data.organizations
    this.onLoadEvent(this.event)
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
      body.append('name', this.event.name)
      body.append('type', this.event.type)
      body.append('organization', this.event.organization)
      body.append('date', this.event.date)
      body.append('start', this.event.start)
      body.append('end', this.event.end)
      body.append('private', this.event.private)
      body.append('venue', this.event.venue.id)
      if (this.event.venue.id === 'NEW') {
        body.append('venueName', this.event.venue.name)
        body.append('venueAddress', this.event.venue.address)
        body.append('venueCity', this.event.venue.city)
        body.append('venueURL', this.event.venue.url)
      }
      body.append('details', this.event.details)
      if (this.event.canEditDSWFlags) {
        body.append('renewsDSW', this.event.renewsDSW)
        body.append('coveredByDSW', this.event.coveredByDSW)
      }
      this.event.groups.forEach(r => { body.append('group', r) })
      const resp = (await this.$axios.post(`/api/events/${this.$route.params.id}`, body)).data
      if (resp && resp.nameError)
        this.duplicateName = { date: this.event.date, name: this.event.name }
      else if (resp) {
        Cookies.set('serv-events-month', this.event.date.substr(0, 7))
        this.$router.push(`/events/${resp.id}`)
      }
    },
    validate() {
      if (!this.submitted) return
      if (!this.event.name)
        this.nameError = 'The event name is required.'
      else if (this.duplicateName && this.duplicateName.date === this.event.date && this.duplicateName.name === this.event.name)
        this.nameError = 'Another event on this date has this name.'
      else
        this.nameError = this.duplicateName = null
      if (!this.event.type)
        this.typeError = 'The event type is required.'
      else
        this.typeError = null
      if (!this.event.organization)
        this.organizationError = 'The organization is required.'
      else
        this.organizationError = null
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
      if (this.event.venue.id === 'NEW') {
        if (!this.event.venue.name)
          this.venueNameError = 'The venue name is required.'
        else
          this.venueNameError = null
        if (!this.event.venue.address)
          this.venueAddressError = 'The venue address is required.'
        else
          this.venueAddressError = null
        if (!this.event.venue.city)
          this.venueCityError = 'The venue city is required.'
        else
          this.venueCityError = null
        if (this.event.venue.url && !this.event.venue.url.startsWith('https://www.google.com/maps/'))
          this.venueURLError = 'The venue map URL must start with https://www.google.com/maps/.'
        else
          this.venueURLError = null
      } else
        this.venueNameError = this.venueAddressError = this.venueCityError = this.venueURLError = null
      if (!this.event.groups.length)
        this.groupsError = 'At least one group must be selected.'
      else
        this.groupsError = null
      this.valid = !this.nameError && !this.dateError && !this.timeError && !this.venueNameError && !this.venueAddressError &&
        !this.venueCityError && !this.venueURLError && !this.organizationError && !this.typeError && !this.groupsError
    },
  },
}
</script>

<style lang="stylus">
#event-edit
  padding 1.5rem 0.75rem
.event-edit-label
  width 7rem
#event-date, #event-name, #event-type, #event-groups, #venue-name, #venue-address, #venue-city, #venue-url, #event-details, #event-organization, #event-edit-flags
  min-width 14rem
  max-width 20rem
#event-venue
  min-width 14rem
  width auto
#event-time
  display flex
  align-items baseline
  min-width 14rem
  max-width 20rem
  span
    padding 0 0.5rem
</style>
