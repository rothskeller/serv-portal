<!--
Event displays the event viewing/editing page.
-->

<template lang="pug">
#event-view(v-if='!event')
  SSpinner
#event-view(v-else)
  #event-view-name(v-text='event.name')
  #event-view-orgtype(v-text='`${orgNames[event.org]} ${event.type}`')
  #event-view-date-time(v-text='dateTimeFmt')
  #event-view-venue
    #event-view-venue-name(v-text='event.venue ? event.venue.name : "Location TBD"')
    #event-view-venue-address(v-if='event.venue')
      span(v-text='event.venue.address')
      span(v-if='event.venue.city', v-text='`, ${event.venue.city}`')
      span#event-view-venue-map(v-if='event.venue.url')
        |
        | (
        a(target='_blank', :href='event.venue.url') map
        | )
  #event-view-details(v-if='event.details', v-html='event.details')
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import moment from 'moment-mini'
import axios, { AxiosResponse } from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SSpinner } from '../../base'

export type GetEventVenue = {
  id: number
  name: string
  // Included in GetEventEvent.venue but not in GetEvent.venues:
  address?: string
  city?: string
  url?: string
}
export type GetEventEvent = {
  id: number
  name: string
  date: string
  start: string
  end: string
  venue: GetEventVenue
  details: string
  renewsDSW: boolean
  coveredByDSW: boolean
  org: string
  type: string
  roles: Array<number>
  canEdit: boolean
  canAttendance: boolean
  canEditDSWFlags: boolean
}
type GetEvent = {
  event: GetEventEvent
}

const orgNames = {
  admin: 'Admin',
  'cert-d': 'CERT Deployment',
  'cert-t': 'CERT Training',
  listos: 'Listos',
  sares: 'SARES',
  snap: 'SNAP',
}

export default defineComponent({
  components: { SSpinner },
  props: {
    onLoadEvent: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    setPage({ title: 'Events' })

    // Get the data for the event.
    const event = ref(null as null | GetEventEvent)
    axios.get(`/api/events/${route.params.id}`).then((resp: AxiosResponse<GetEvent>) => {
      event.value = resp.data.event
      setPage({ title: `${event.value.date} ${event.value.name}`, browserTitle: event.value.date })
      props.onLoadEvent(event.value)
    })

    // Format the date and time of the event.
    const dateTimeFmt = computed(() => {
      if (!event.value) return ''
      const date = moment(event.value.date, 'YYYY-MM-DD')
      const start = moment(event.value.start, 'HH:mm')
      const end = moment(event.value.end, 'HH:mm')
      if (start.format('a') !== end.format('a'))
        return `${date.format('dddd, MMMM D, YYYY')}\n${start.format('h:mma')} to ${end.format(
          'h:mma'
        )}`
      else
        return `${date.format('dddd, MMMM D, YYYY')}\n${start.format('h:mm')} to ${end.format(
          'h:mma'
        )}`
    })

    return { event, dateTimeFmt, orgNames }
  },
})
</script>

<style lang="postcss">
#event-view {
  padding: 1.5rem 0.75rem;
}
#event-view-name {
  font-weight: bold;
  font-size: 1.25rem;
  line-height: 1.2;
}
#event-view-orgtype {
  color: #888;
}
#event-view-date-time {
  margin-top: 0.75rem;
  white-space: pre-line;
  line-height: 1.2;
}
#event-view-venue {
  margin-top: 0.75rem;
  line-height: 1.2;
}
#event-view-venue-address {
  font-size: 0.875rem;
}
#event-view-details {
  margin-top: 0.75rem;
  max-width: 40rem;
  white-space: pre-line;
  line-height: 1.2;
}
</style>
