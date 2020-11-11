<!--
PersonHours allows editing the hours a person has spent on SERV activities.
-->

<template lang="pug">
#person-hours-unregistered(v-if='unregistered')
  | You are not currently registered as a City of Sunnyvale volunteer. We
  | appreciate your volunteer efforts, but we cannot record your hours until you
  | are registered. To do so, please fill out
  |
  a(href='https://www.volgistics.com/ex/portal.dll/ap?AP=929478828', target='_blank') this form
  | . In the “City employee status or referral” box, please enter
  pre
    | Rebecca Elizondo
    | Department of Public Safety
  div
    | and the names of the organizations you're volunteering for (CERT, Listos,
    | SNAP, and/or SARES). Come back a week or so later and we should have your
    | registration on file. If you have any difficulties with this, please
    | contact Rebecca at
    |
    a(href='mailto:RElizondo@sunnyvale.ca.gov') RElizondo@sunnyvale.ca.gov
    | .
#person-hours(v-else-if='notFound')
  div The link you used is not valid or has expired.
  div(style='margin-top: 1.5rem')
    SButton(to='/login', variant='primary') Go to Login Page
#person-hours(v-else-if='!months')
  SSpinner
#person-hours(v-else-if='saved')
  div Your volunteer hours have been saved. Thank you for volunteering!
  div(style='margin-top: 1.5rem')
    SButton(to='/login', variant='primary') Go to Login Page
form#person-hours(v-else, @submit.prevent='onSubmit')
  .person-hours(v-for='month in months')
    .person-hours-heading(v-text='`Volunteer Hours for ${month.month}`')
    table.person-hours-table
      PersonHoursEvent(v-for='event in month.events', :event='event', :key='event.id')
      tr
        td.person-hours-total-label TOTAL
        td
          .person-hours-total-time(v-text='totalHours(month)')
  #person-hours-buttons
    SButton(type='submit', variant='primary') Save Hours
    SButton(v-if='showCancel', @click='onCancel') Cancel
  table#person-hours-guide
    tr
      td Volunteer Hours
      td Not Volunteer Hours
    tr
      td In general, time you spend helping or preparing to help the community as part of SERV. For example:
      td In general, time you spend preparing yourself or your household; or time you spend becoming a SERV volunteer. For example:
    tr
      td Organizing or teaching CERT Basic, Listos, PEP, or SNAP events
      td Attending CERT Basic, Listos, PEP, or ham cram classes
    tr
      td Preparing and maintaining a CERT or SARES “go kit” for deployment
      td Preparing and maintaining a personal or household evacuation kit
    tr
      td SERV team meetings, radio nets, and drills; CERT continuing education seminars; SARES or county ARES training classes
      td SERV team social gatherings
    tr
      td Responding in an emergency when activated by the city
      td Responding in an emergency when not activated by the city
    tr
      td Travel to and from the above
      td
    tr
      td SERV administration activities
      td
</template>

<script lang="ts">
import { computed, defineComponent, inject, Ref, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import { SButton, SSpinner } from '../../base'
import PersonHoursEvent from './PersonHoursEvent.vue'
import { LoginData } from '../../plugins/login'
import setPage from '../../plugins/page'

export type GetPersonHoursEvent = {
  id: number
  date: string
  name: string
  minutes: number
  placeholder?: boolean
}
type GetPersonHoursMonth = {
  month: string
  events: Array<GetPersonHoursEvent>
}
type GetPersonHoursMonths = {
  name: string
  months: Array<GetPersonHoursMonth>
}
type GetPersonHours = false | GetPersonHoursMonths

export default defineComponent({
  components: { PersonHoursEvent, SButton, SSpinner },
  props: {
    onLoadPerson: Function, // not used
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const me = inject<Ref<LoginData>>('me')!
    setPage({ title: 'Volunteer Hours' })

    // Load the data.
    const months = ref(null as null | Array<GetPersonHoursMonth>)
    const unregistered = ref(false)
    const notFound = ref(false)
    axios
      .get<GetPersonHours>(`/api/people/${route.params.id}/hours`)
      .catch((e) => {
        if (e.response && e.response.status === 404) return null
        else throw e
      })
      .then((resp) => {
        if (!resp) notFound.value = true
        else if (resp.data) months.value = resp.data.months
        else unregistered.value = true
      })

    function totalHours(month: GetPersonHoursMonth) {
      const minutes = month.events.reduce((sum, e) => sum + e.minutes, 0)
      return Math.floor(minutes / 30) / 2
    }

    const saved = ref(false)
    const showCancel = computed(() => route.params.id.length <= 5)
    async function onSubmit() {
      const body = new FormData()
      months.value!.forEach((m) => {
        m.events.forEach((e) => {
          body.append(`e${e.id}`, e.minutes.toString())
        })
      })
      await axios.post(`/api/people/${route.params.id}/hours`, body)
      if (route.params.id.length <= 5) router.push(`/people/${route.params.id}`)
      else if (me.value) router.push(`/people/${me.value.id}`)
      else saved.value = true
    }
    function onCancel() {
      router.go(-1)
    }

    return {
      months,
      notFound,
      onCancel,
      onSubmit,
      saved,
      showCancel,
      totalHours,
      unregistered,
    }
  },
})
</script>

<style lang="postcss">
#person-hours-unregistered {
  margin: 1.5rem 0.75rem;
  max-width: 600px;
  & pre {
    margin-left: 1.5rem;
  }
}
#person-hours {
  margin: 1.5rem 0.75rem;
}
.person-hours-heading {
  font-weight: bold;
  font-size: 1.5rem;
}
.person-hours-total-label {
  padding-right: 1rem;
  text-align: right;
  font-weight: bold;
}
.person-hours-total-time {
  width: 3rem;
  text-align: right;
  font-weight: bold;
}
#person-hours-buttons {
  margin-top: 1rem;
  & .sbtn {
    margin-right: 0.5rem;
  }
}
#person-hours-guide {
  margin-top: 1.5rem;
  max-width: 800px;
  tr:nth-child(1) {
    background-color: #5b9bd5;
    td {
      color: white;
      text-align: center;
      font-weight: bold;
    }
  }
  tr:nth-child(2) {
    font-weight: bold;
  }
  tr:nth-child(even) {
    background-color: #deeaf6;
  }
  td {
    padding: 0.25rem;
    width: 50%;
    border: 1px solid #eee;
    vertical-align: top;
    line-height: 1.2;
  }
}
</style>
