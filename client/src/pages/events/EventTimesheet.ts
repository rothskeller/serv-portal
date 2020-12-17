// EventTimesheet allows easy attendance entry from a timesheet.

import { defineComponent, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios, { AxiosResponse } from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SButton, SSpinner } from '../../base'
import type { GetEventAttendance, GetEventEvent, GetEventPerson, GuestOption } from './api'
import EventAttendanceGuest from './EventAttendanceGuest.vue'
import EventTimesheetEntry from './EventTimesheetEntry'
import './events.css'

export type TimesheetEntry = {
  id: number
  sortName: string
  type: string
  origType: string
  start: string
  end: string
  minutes: number
  origMinutes: number
  invalid: boolean
}

const EventTimesheet = defineComponent({
  props: {
    onLoadEvent: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Events' })

    // Load the event attendance to be edited.
    const event = ref<GetEventEvent>()
    const people = ref<Array<TimesheetEntry>>()
    axios
      .get(`/api/events/${route.params.id}?attendance=1`)
      .then((resp: AxiosResponse<GetEventAttendance>) => {
        event.value = resp.data.event
        people.value = resp.data.people.map((p) => ({
          id: p.id,
          sortName: p.sortName,
          type: p.attended ? p.attended.type : '',
          origType: p.attended ? p.attended.type : '',
          minutes: p.attended ? p.attended.minutes : 0,
          origMinutes: p.attended ? p.attended.minutes : 0,
          start: '',
          end: '',
          invalid: false,
        }))
        if (route.params.id === 'NEW') setPage({ title: 'New Event' })
        else
          setPage({
            title: `${event.value.date} ${event.value.name}`,
            browserTitle: event.value.date,
          })
        props.onLoadEvent(event.value)
      })

    // Attendance type to set.
    const defaultType = ref('Volunteer')

    // Adding a guest.
    const guestDialog = ref(null as any)
    async function onAddGuest() {
      const guest: null | GuestOption = await guestDialog.value.show()
      if (!guest) return
      const existing = people.value!.find((p) => p.id === guest.id)
      if (existing) return
      people.value!.push({
        id: guest.id,
        sortName: guest.sortName,
        start: '',
        end: '',
        minutes: 0,
        type: '',
        origType: '',
        origMinutes: 0,
        invalid: false,
      })
      people.value!.sort((a, b) =>
        a.sortName < b.sortName ? -1 : a.sortName > b.sortName ? +1 : 0
      )
    }

    // Save.
    async function save() {
      if (!people.value) return
      if (people.value.find((p) => p.invalid)) return
      const body = new FormData()
      people.value!.forEach((p) => {
        if (p.minutes || p.type) {
          body.append('person', p.id.toString())
          body.append('type', p.type)
          body.append('minutes', p.minutes.toString())
        }
      })
      await axios.post(`/api/events/${route.params.id}/attendance`, body)
    }
    async function onSave(evt: Event) {
      evt.preventDefault()
      await save()
      router.push(`/events/${route.params.id}`)
    }
    function onCancel(evt: Event) {
      evt.preventDefault()
      router.go(-1)
    }
    async function onHoursView(evt: Event) {
      evt.preventDefault()
      await save()
      router.replace(`/events/${route.params.id}/attendance`)
    }

    function renderGrid() {
      return h(
        'div',
        { id: 'event-timesheet-grid' },
        people.value!.map((p) =>
          h(EventTimesheetEntry, {
            entry: p,
            defaultType: defaultType.value,
            onNewType: (t: string) => (defaultType.value = t),
          })
        )
      )
    }

    function renderButtons() {
      return h('div', { id: 'event-timesheet-submit' }, [
        h(SButton, { type: 'submit', variant: 'primary' }, () => 'Save Timesheet'),
        h(SButton, { id: 'event-timesheet-cancel', onClick: onCancel }, () => 'Cancel'),
        h(SButton, { onClick: onAddGuest }, () => 'Add Guest'),
        h(SButton, { onClick: onHoursView }, () => 'Hours View'),
      ])
    }

    function renderGuestDialog() {
      return h(EventAttendanceGuest, { ref: (r) => (guestDialog.value = r) })
    }

    return () => {
      if (!event.value) return h(SSpinner)
      return h(
        'form',
        {
          id: 'event-timesheet',
          onSubmit: onSave,
        },
        [renderGrid(), renderButtons(), renderGuestDialog()]
      )
    }
  },
})
export default EventTimesheet
