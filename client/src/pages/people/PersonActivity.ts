// PersonActivity displays the Activity tab for a person.

import moment from 'moment-mini'
import { defineComponent, Fragment, h, inject, Ref, ref, render, watchEffect } from 'vue'
import { RouteLocationNormalizedLoaded, Router, useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import { MonthSelect, SButton, SSpinner } from '../../base'
import { fmtMinutes, GetPersonHours, GetPersonHoursEvent } from './api'
import { LoginData } from '../../plugins/login'
import PersonActivityEvent from './PersonActivityEvent'
import './people.css'
import setPage from '../../plugins/page'

const PersonActivity = defineComponent({
  name: 'PersonActivity',
  props: {
    onLoadPerson: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    const me = inject<Ref<LoginData>>('me')!
    setPage({ title: 'Volunteer Activity', browserTitle: 'Activity' })
    const person = ref<GetPersonHours>()
    const notFound = ref(false)
    const saved = ref(false)
    watchEffect(async () => {
      if (!route.params.month) return // probably leaving page
      if (route.params.month === 'current') {
        const now = moment()
        const month = now.date() <= 10 ? now.subtract(1, 'month') : now
        router.replace(`/people/${route.params.id}/activity/${month.format('YYYY-MM')}`)
        return
      }
      try {
        person.value = (
          await axios.get<GetPersonHours>(
            `/api/people/${route.params.id}/hours/${route.params.month}`
          )
        ).data
        setPage({ title: person.value.name, browserTitle: 'Activity' })
        person.value.canHours = true
        props.onLoadPerson(person.value)
      } catch (e) {
        if (e.response && e.response.status === 403 && !me.value)
          router.replace({ path: '/login', query: { redirect: route.fullPath } })
        if (e.response && e.response.status === 404 && !me.value) notFound.value = true
        else throw e
      }
    })
    async function onSubmit(evt: Event) {
      evt.preventDefault()
      const body = new FormData()
      person.value!.events.forEach((e) => {
        if (!e.canEdit) return
        body.append(`e${e.id}`, `${e.minutes}:${e.type}`)
      })
      await axios.post(`/api/people/${route.params.id}/hours/${route.params.month}`, body)
      if (me.value) router.push(`/people/${me.value.id}`)
      else saved.value = true
    }
    function render() {
      if (notFound.value) return renderNotFoundPage()
      if (saved.value) return renderSavedPage()
      if (!person.value) return h('div', { id: 'person-activity' }, h(SSpinner))
      if (person.value.needsVolgistics && (!me.value || me.value.id === person.value.id))
        return renderVolgisticsWarningPage()
      const canEdit = person.value.events.some((e) => e.canEdit)
      return h(canEdit ? 'form' : 'div', { id: 'person-activity', onSubmit }, [
        renderMonthSelector(route, router),
        renderEvents(person.value.events),
        canEdit ? renderButtons(!!me.value, router) : null,
        canEdit ? renderGuidanceTable() : null,
      ])
    }
    return render
  },
})
export default PersonActivity

function renderButtons(showCancel: boolean, router: Router) {
  function cancel() {
    router.go(-1)
  }
  return h('div', { id: 'person-activity-buttons' }, [
    h(SButton, { type: 'submit', variant: 'primary' }, () => 'Save Activity'),
    showCancel ? h(SButton, { onClick: cancel }, () => 'Cancel') : null,
  ])
}

function renderTotalRow(events: Array<GetPersonHoursEvent>) {
  const total = events.reduce((sum, e) => sum + e.minutes, 0)
  const count = events.reduce((count, e) => count + (e.minutes ? 1 : 0), 0)
  const canViewTypes = events.some((e) => e.canViewType)
  if (count > 1) {
    return [
      h('div', { id: 'person-activity-total-label' }, 'TOTAL'),
      canViewTypes ? h('div') : null,
      h('div', { style: 'person-activity-total' }, fmtMinutes(total)),
    ]
  } else {
    return null
  }
}

function renderEvents(events: Array<GetPersonHoursEvent>) {
  if (!events.length) return h('div', 'No activity has been recorded for this month.')
  const canEdit = events.some((e) => e.canEdit)
  const canViewTypes = events.some((e) => e.canViewType)
  const total = events.reduce((sum, e) => sum + e.minutes, 0)
  return h(
    'div',
    {
      id: 'person-activity-grid',
      class: { 'person-activity-editable': canEdit, 'person-activity-with-types': canViewTypes },
    },
    [
      h('div', { style: 'font-weight: bold' }, 'Event'),
      canViewTypes ? h('div', { style: 'font-weight: bold' }, 'Type') : null,
      h('div', { style: 'font-weight: bold' }, 'Hours'),
      events.map((e) => h(PersonActivityEvent, { event: e, typesColumn: canViewTypes })),
      renderTotalRow(events),
    ]
  )
}

function renderGuidanceTable() {
  return h('table', { id: 'person-activity-guide' }, [
    h('tr', [h('td', 'Volunteer Hours'), h('td', 'Not Volunteer Hours')]),
    h('tr', [
      h(
        'td',
        'In general, time you spend helping or preparing to help the community as part of SERV. For example:'
      ),
      h(
        'td',
        'In general, time you spend preparing yourself or your household; or time you spend becoming a SERV volunteer. For example:'
      ),
    ]),
    h('tr', [
      h('td', 'Organizing or teaching CERT Basic, Listos, PEP, or SNAP events'),
      h('td', 'Attending CERT Basic, Listos, PEP, or ham cram classes'),
    ]),
    h('tr', [
      h('td', 'Preparing and maintaining a CERT or SARES “go kit” for deployment'),
      h('td', 'Preparing and maintaining a personal or household evacuation kit'),
    ]),
    h('tr', [
      h(
        'td',
        'SERV team meetings, radio nets, and drills; CERT continuing education seminars; SARES or county ARES training classes'
      ),
      h('td', 'SERV team social gatherings'),
    ]),
    h('tr', [
      h('td', 'Responding in an emergency when activated by the city'),
      h('td', 'Responding in an emergency when not activated by the city'),
    ]),
    h('tr', [h('td', 'Travel to and from the above'), h('td')]),
    h('tr', [h('td', 'SERV administration activities'), h('td')]),
  ])
}

function renderMonthSelector(route: RouteLocationNormalizedLoaded, router: Router) {
  return h(
    'div',
    { id: 'person-activity-month' },
    h(MonthSelect, {
      modelValue: route.params.month as string,
      'onUpdate:modelValue': (v: string) => {
        router.replace(`/people/${route.params.id}/activity/${v}`)
      },
    })
  )
}

function renderNotFoundPage() {
  return h('div', { id: 'person-activity' }, [
    h('div', 'The link you used is not valid or has expired.'),
    h('div', { style: 'margin-top: 1.5rem' }, [
      h(SButton, { to: '/login', variant: 'primary' }, 'Go to Login Page'),
    ]),
  ])
}

function renderSavedPage() {
  return h('div', { id: 'person-activity' }, [
    h('div', 'Your volunteer hours have been saved. Thank you for volunteering!'),
    h('div', { style: 'margin-top: 1.5rem' }, [
      h(SButton, { to: '/login', variant: 'primary' }, 'Go to Login Page'),
    ]),
  ])
}

function renderVolgisticsWarningPage() {
  return h('div', { id: 'person-activity-unregistered' }, [
    'You are not currently registered as a City of Sunnyvale volunteer. We appreciate your volunteer efforts, but we cannot record your hours until you are registered. To do so, please fill out ',
    h(
      'a',
      { href: 'https://www.volgistics.com/ex/portal.dll/ap?AP=929478828', target: '_blank' },
      'this form'
    ),
    '. In the “City employee status or referral” box, please enter',
    h('pre', 'Rebecca Elizondo\nDepartment of Public Safety'),
    'and the names of the organizations you’re volunteering for (CERT, Listos, SARES, and/or SNAP). Come back a week or so later and we should have your registration on file. If you have any difficulties with this, please contact Rebecca at ',
    h('a', { href: 'mailto:RElizondo@sunnyvale.ca.gov' }, 'RElizondo@sunnyvale.ca.gov'),
    '.',
  ])
}
