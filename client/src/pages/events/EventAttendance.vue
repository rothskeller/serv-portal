<!--
EventAttendance shows and allows changes to the attendance for an event.
-->

<template lang="pug">
#event-attendance-spinner(v-if='!event')
  SSpinner
form#event-attendance(v-else, @submit.prevent='onSave')
  #event-attend-settings
    label(for='event-attend-type') Set attendance for:
    SSelect#event-attend-type(:options='typeOptions', v-model='setType')
    label(for='event-attend-hours') Hours:
    SInput#event-attend-hours(v-model='setHoursS', type='number', min='0', max='24', step='0.5')
  #event-attend-group
    EventAttendancePerson(
      v-for='p in people',
      :key='p.id',
      :person='p',
      @toggle='onTogglePerson(p)'
    )
  #event-attend-submit
    SButton(type='submit', variant='primary') Save Attendance
    SButton#event-attend-cancel(@click='onCancel') Cancel
    SButton(@click='onAddGuest') Add Guest
    SButton(@click='onTimesheetView') Timesheet View
  EventAttendanceGuest(ref='guestDialog')
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import moment from 'moment-mini'
// import VSelect from 'vue-select'
// import 'vue-select/dist/vue-select.css'
import axios, { AxiosResponse } from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SButton, SInput, SSelect, SSpinner } from '../../base'
import EventAttendanceGuest, { GuestOption } from './EventAttendanceGuest.vue'
import EventAttendancePerson from './EventAttendancePerson.vue'

import type { GetEventEvent } from './EventView.vue'
type GetEventPersonAttendance = {
  type: string
  minutes: number
}
export type GetEventPerson = {
  id: number
  sortName: string
  attended: false | GetEventPersonAttendance
}
type GetEventAttendance = {
  event: GetEventEvent
  people: Array<GetEventPerson>
}

export default defineComponent({
  components: {
    EventAttendanceGuest,
    EventAttendancePerson,
    SButton,
    SInput,
    SSelect,
    SSpinner,
  },
  props: {
    onLoadEvent: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Events' })

    // Load the event attendance to be edited.
    const event = ref(null as null | GetEventEvent)
    const people = ref(null as null | Array<GetEventPerson>)
    axios
      .get(`/api/events/${route.params.id}?attendance=1`)
      .then((resp: AxiosResponse<GetEventAttendance>) => {
        event.value = resp.data.event
        people.value = resp.data.people
        setHours.value = moment(event.value.end, 'HH:mm').diff(
          moment(event.value.start, 'HH:mm'),
          'hours',
          true
        )
        if (route.params.id === 'NEW') setPage({ title: 'New Event' })
        else
          setPage({
            title: `${event.value.date} ${event.value.name}`,
            browserTitle: event.value.date,
          })
        props.onLoadEvent(event.value)
      })

    // Attendance type to set.
    const setType = ref('Volunteer')
    const typeOptions = [
      { value: 'Volunteer', label: 'Volunteer' },
      { value: 'Student', label: 'Student' },
      { value: 'Audit', label: 'Auditor' },
      { value: 'Absent', label: 'Absentee' },
    ]

    // Number of hours to set.
    const setHours = ref(1.0)
    const setHoursS = computed({
      get: () => setHours.value.toString(),
      set: (h) => {
        setHours.value = parseFloat(h)
      },
    })

    // Toggling the state of a person.
    function onTogglePerson(person: GetEventPerson) {
      if (person.attended && person.attended.type === setType.value) person.attended = false
      else if (person.attended) person.attended.type = setType.value
      else person.attended = { type: setType.value, minutes: 60 * setHours.value }
    }

    // Adding a guest.
    const guestDialog = ref(null as any)
    async function onAddGuest() {
      const guest: null | GuestOption = await guestDialog.value.show()
      if (!guest) return
      const existing = people.value!.find((p) => p.id === guest.id)
      if (existing) {
        existing.attended = { type: setType.value, minutes: 60 * setHours.value }
      } else {
        people.value!.push({
          id: guest.id,
          sortName: guest.sortName,
          attended: { type: setType.value, minutes: 60 * setHours.value },
        })
        people.value!.sort((a, b) =>
          a.sortName < b.sortName ? -1 : a.sortName > b.sortName ? +1 : 0
        )
      }
    }

    // Save and cancel.
    function onCancel() {
      router.go(-1)
    }
    async function onSave() {
      await save()
      router.push(`/events/${route.params.id}`)
    }
    async function onTimesheetView() {
      await save()
      router.replace(`/events/${route.params.id}/timesheet`)
    }
    async function save() {
      const body = new FormData()
      people.value!.forEach((p) => {
        if (p.attended) {
          body.append('person', p.id.toString())
          body.append('type', p.attended.type)
          body.append('minutes', p.attended.minutes.toString())
        }
      })
      await axios.post(`/api/events/${route.params.id}/attendance`, body)
    }

    return {
      event,
      guestDialog,
      onAddGuest,
      onCancel,
      onSave,
      onTimesheetView,
      onTogglePerson,
      people,
      setHoursS,
      setType,
      typeOptions,
    }
  },
})
</script>

<style lang="postcss">
#event-attendance-spinner {
  padding: 1.5rem 0.75rem;
}
#event-attendance {
  padding: 1.5rem 0.75rem;
  @media (min-width: 576px) {
    display: grid;
    height: 100%;
    grid: max-content 1fr max-content / 100%;
  }
}
#event-attend-settings {
  margin-bottom: 0.75rem;
  @media print {
    display: none;
  }
}
#event-attend-type {
  margin: 0 1.5rem 0 0.5rem;
}
#event-attend-hours {
  margin-left: 0.5rem;
  max-width: 5rem;
}
@media (min-width: 576px) {
  #event-attend-group {
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    align-content: flex-start;
    overflow-x: auto;
    overflow-y: hidden;
    min-height: 0;
    .custom-checkbox {
      margin-right: 1.5rem;
    }
  }
}
#event-attend-submit {
  display: flex;
  flex-wrap: wrap;
  margin: 0.75rem 0 0 -0.5rem;
  & .sbtn {
    margin-left: 0.5rem;
  }
  @media print {
    display: none;
  }
}
#event-attend-cancel {
  margin-right: 2rem;
}
</style>
