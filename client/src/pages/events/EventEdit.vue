<!--
EventEdit is the event editing page.  (It also handles creation of new events.)
-->

<template lang="pug">
#event-edit-spinner(v-if='!event')
  SSpinner
SForm(v-else, @submit='onSubmit', submitLabel='Save Event')
  SFInput#event-name(
    label='Event name',
    autofocus,
    trim,
    v-model='event.name',
    :errorFn='nameError'
  )
  SFSelect#event-type(
    label='Event type',
    :options='types',
    v-model='event.type',
    :errorFn='typeError'
  )
  SFSelect#event-type(
    label='Organization',
    :options='orgs',
    v-model='event.org',
    :errorFn='orgError'
  )
  SFInput#event-date(label='Event date', type='date', v-model='event.date', :errorFn='dateError')
  SFTimeRange#event-times(label='Event time', v-model:start='event.start', v-model:end='event.end')
  SFSelect#event-venue(
    label='Location',
    :options='allVenues',
    valueKey='id',
    labelKey='name',
    v-model='event.venue.id'
  )
  template(v-if='event.venue.id === -1')
    SFInput#venue-name(
      label='Venue name',
      ref='venueName',
      trim,
      v-model='event.venue.name',
      :errorFn='venueNameError'
    )
    SFInput#venue-address(
      label='Address',
      trim,
      v-model='event.venue.address',
      :errorFn='venueAddressError'
    )
    SFInput#venue-city(label='City', trim, v-model='event.venue.city', :errorFn='venueCityError')
    SFInput#venue-url(
      label='Map URL',
      trim,
      v-model='event.venue.url',
      :errorFn='venueURLError',
      help='Generate this by opening Google Maps and searching for the venue. (Search by the venue’s name rather than its address, if it has a name, so that Google’s sidebar of information about the venue appears.)  Once the venue is found, zoom out a few steps until the nearest freeways are shown — enough to give the viewer an idea of where in the Bay Area the venue is.  Copy the URL from the browser address bar and paste it here.'
    )
  SFTextArea#event-details(
    label='Details',
    trim,
    rows=3,
    wrap='soft',
    v-model='event.details',
    help='This may contain HTML <a> tags for links, but no other tags.'
  )
  SFCheckGroup#event-flags(label='Flags', :options='flagOptions', v-model='flags')
  SFCheckGroup#event-roles(
    label='Invited roles',
    :options='filteredRoles',
    valueKey='id',
    labelKey='name',
    v-model='roles'
  )
  template(v-if='event.id', #extraButtons)
    SButton(@click='onDelete', variant='danger') Delete Event
  MessageBox(
    ref='deleteModal',
    title='Delete Event',
    cancelLabel='Keep',
    okLabel='Delete',
    variant='danger'
  )
    | Are you sure you want to delete this event? All associated data,
    | including attendance records, will be permanently lost.
</template>

<script lang="ts">
import { defineComponent, ref, watch, nextTick, computed } from 'vue'
import Cookies from 'js-cookie'
import { useRoute, useRouter } from 'vue-router'
import axios, { AxiosResponse } from '../../plugins/axios'
import setPage from '../../plugins/page'
import {
  MessageBox,
  SButton,
  SForm,
  SFCheckGroup,
  SFInput,
  SFSelect,
  SFTextArea,
  SFTimeRange,
  SSpinner,
} from '../../base'
import type { GetEventEvent, GetEventVenue } from './EventView.vue'

type GetEventRole = {
  id: number
  name: string
  org: string
}
type GetEventEdit = {
  event: GetEventEvent
  types: Array<string>
  roles: Array<GetEventRole>
  venues: Array<GetEventVenue>
  orgs: Array<string>
}
type PostEvent = {
  // It will be one or the other of:
  id?: number
  nameError?: true
}
type OrgLabel = {
  value: string
  label: string
}

const orgNames: Record<string, string> = {
  admin: 'Admin',
  'cert-d': 'CERT Deployment',
  'cert-t': 'CERT Training',
  listos: 'Listos',
  sares: 'SARES',
  snap: 'SNAP',
}

export default defineComponent({
  components: {
    MessageBox,
    SButton,
    SForm,
    SFCheckGroup,
    SFInput,
    SFSelect,
    SFTextArea,
    SFTimeRange,
    SSpinner,
  },
  props: {
    onLoadEvent: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Events' })

    // Load the event to be edited.
    const event = ref(null as null | GetEventEvent)
    const allRoles = ref([] as Array<GetEventRole>)
    const allVenues = ref([] as Array<GetEventVenue>)
    const types = ref([] as Array<string>)
    const orgs = ref([] as Array<OrgLabel>)
    axios.get(`/api/events/${route.params.id}?edit=1`).then((resp: AxiosResponse<GetEventEdit>) => {
      event.value = resp.data.event
      allRoles.value = resp.data.roles
      allVenues.value = resp.data.venues
      allVenues.value.unshift({ id: 0, name: 'TBD' })
      allVenues.value.push({ id: -1, name: '(create a new venue)' })
      types.value = resp.data.types
      if (!event.value.type) types.value.unshift('(select type)')
      orgs.value = resp.data.orgs.map((o) => ({ value: o, label: orgNames[o] }))
      if (!event.value.org) orgs.value.unshift({ value: '', label: '(select organization)' })
      roles.value = new Set(event.value.roles)
      if (route.params.id === 'NEW') setPage({ title: 'New Event' })
      else
        setPage({
          title: `${event.value.date} ${event.value.name}`,
          browserTitle: event.value.date,
        })
      props.onLoadEvent(event.value)
    })

    // Event name field.
    const duplicateName = ref(null as null | { name: string; date: string })
    function nameError(lostFocus: boolean) {
      if (!lostFocus || !event.value) return ''
      if (!event.value.name) return 'The event name is required.'
      if (
        duplicateName.value &&
        duplicateName.value.date === event.value.date &&
        duplicateName.value.name === event.value.name
      )
        return 'Another event on this date has this name.'
      return ''
    }

    // Event type field.
    function typeError(lostFocus: boolean) {
      if (!lostFocus || !event.value) return ''
      if (!event.value.type || event.value.type === '(select type)')
        return 'The event type is required.'
      return ''
    }

    // Organization field.
    function orgError(lostFocus: boolean) {
      if (!lostFocus || !event.value) return ''
      if (!event.value.org || event.value.org === '(select organization)')
        return 'The organization is required.'
      return ''
    }

    // Event date field.
    function dateError(lostFocus: boolean) {
      if (!lostFocus || !event.value) return ''
      if (!event.value.date) return 'The event date is required.'
      if (!event.value.date.match(/^20\d\d-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$/))
        return 'This is not a valid date.'
      return ''
    }

    // Venue selection.
    const venueName = ref(null as any)
    watch(
      () => event.value?.venue.id,
      () => {
        if (event.value && event.value.venue.id === -1) {
          event.value.venue.name = event.value.venue.address = ''
          event.value.venue.city = event.value.venue.url = ''
          nextTick(() => {
            venueName.value.focus()
          })
        }
      }
    )

    // Venue name field.
    function venueNameError(lostFocus: boolean) {
      if (!lostFocus || !event.value || event.value.venue.id !== -1) return ''
      if (!event.value.venue.name) return 'The venue name is required.'
      return ''
    }

    // Venue address field.
    function venueAddressError(lostFocus: boolean) {
      if (!lostFocus || !event.value || event.value.venue.id !== -1) return ''
      if (!event.value.venue.address) return 'The venue address is required.'
      return ''
    }

    // Venue city field.
    function venueCityError(lostFocus: boolean) {
      if (!lostFocus || !event.value || event.value.venue.id !== -1) return ''
      if (!event.value.venue.city) return 'The venue city is required.'
      return ''
    }

    // Venue map URL field.
    function venueURLError(lostFocus: boolean) {
      if (!lostFocus || !event.value || event.value.venue.id !== -1) return ''
      if (
        event.value.venue.url &&
        !event.value.venue.url.startsWith('https://www.google.com/maps/')
      )
        return 'The venue map URL must start with https://www.google.com/maps/.'
      return ''
    }

    // Flags.
    const flags = computed({
      get: () => {
        const flags = new Set()
        if (event.value!.coveredByDSW) flags.add('coveredByDSW')
        return flags
      },
      set: (flags) => {
        event.value!.coveredByDSW = flags.has('coveredByDSW')
      },
    })
    const flagOptions = [
      { value: 'coveredByDSW', label: 'Event is covered by DSW insurance' },
    ]

    // Invited roles.
    const filteredRoles = computed(() => allRoles.value.filter((f) => f.org === event.value!.org))
    const roles = ref<Set<number>>()

    async function onSubmit() {
      if (!event.value) return
      const body = new FormData()
      body.append('name', event.value.name)
      body.append('type', event.value.type)
      body.append('org', event.value.org)
      body.append('date', event.value.date)
      body.append('start', event.value.start)
      body.append('end', event.value.end)
      if (event.value.venue.id === -1) {
        body.append('venue', 'NEW')
        body.append('venueName', event.value.venue.name)
        body.append('venueAddress', event.value.venue.address!)
        body.append('venueCity', event.value.venue.city!)
        body.append('venueURL', event.value.venue.url!)
      } else body.append('venue', event.value.venue.id.toString())
      body.append('coveredByDSW', event.value.coveredByDSW.toString())
      body.append('details', event.value.details)
      roles.value!.forEach((r) => {
        body.append('role', r.toString())
      })
      const resp: PostEvent = (await axios.post(`/api/events/${route.params.id}`, body)).data
      if (resp.nameError) {
        duplicateName.value = { name: event.value.name, date: event.value.date }
        return
      }
      Cookies.set('serv-events-month', event.value!.date.substr(0, 7))
      router.push(`/events/${resp.id}`)
    }

    // Handle deleting an event.
    const showDeleteModal = ref(false)
    const deleteModal = ref(null as any)
    async function onDelete() {
      if (deleteModal.value) {
        const confirmed: boolean = await deleteModal.value.show()
        if (confirmed) {
          const body = new FormData()
          body.append('delete', 'true')
          await axios.post(`/api/events/${route.params.id}`, body)
          router.push({ path: '/events', params: { year: event.value!.date.substr(0, 4) } })
        }
      }
    }

    return {
      event,
      nameError,
      types,
      typeError,
      orgs,
      orgError,
      dateError,
      allVenues,
      venueName,
      venueNameError,
      venueAddressError,
      venueCityError,
      venueURLError,
      flags,
      flagOptions,
      filteredRoles,
      roles,
      onSubmit,
      onDelete,
      deleteModal,
    }
  },
})
</script>

<style lang="postcss">
#event-edit-spinner {
  margin: 1.5rem 0.75rem;
}
</style>
