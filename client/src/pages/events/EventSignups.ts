// EventSignups displays the event signups page.

import moment from 'moment-mini'
import { defineComponent, h, inject, Ref, ref, VNode } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import axios from '../../plugins/axios'
import { SButton, SCheck, SIcon, SSpinner } from '../../base'
import './events.css'
import setPage from '../../plugins/page'
import { LoginData } from '../../plugins/login'

type GetEventSignups = {
  id: number
  sortName: string
  events: Array<GetEventSignupsEvent>
  // added locally
  saved: boolean
  loggedIn: boolean
}
type GetEventSignupsEvent = {
  id: number
  date: string
  name: string
  signupText: string
  shifts: Array<GetEventSignupsShift>
  // added locally
  declined: boolean
}
type GetEventSignupsShift = {
  start: string
  end: string
  task: string
  min: number
  max: number
  count: number
  list: Array<{
    id: number
    sortName: string
  }>
  signedUp: boolean
  // added locally
  baseCount: number
  expanded: boolean
}

const EventSignups = defineComponent({
  name: 'EventSignups',
  setup() {
    const route = useRoute()
    const data = ref<GetEventSignups>()
    const me = inject<Ref<LoginData>>('me')!
    setPage({ title: 'Event Signups' })

    // Load the data on startup.
    loadData()
    async function loadData() {
      data.value = (
        await axios.get<GetEventSignups>(
          `/api/events/signups${route.params.id ? `/${route.params.id}` : ''}`
        )
      ).data
      data.value.loggedIn = !!me.value
      data.value.events.forEach((e) => {
        if (!e.shifts.find((s) => s.signedUp)) e.declined = true
        e.shifts.forEach((s) => {
          s.baseCount = s.signedUp ? s.count - 1 : s.count
          s.expanded = false
        })
      })
    }

    // Send the signups to the server.
    async function onSubmit(evt: Event) {
      evt.preventDefault()
      const body = new FormData()
      data.value!.events.forEach((e) => {
        e.shifts.forEach((s) => {
          body.append(
            `${e.id}.${s.start}.${s.task}`,
            e.declined ? 'declined' : s.signedUp.toString()
          )
        })
      })
      await axios.post(`/api/events/signups${route.params.id ? `/${route.params.id}` : ''}`, body)
      data.value!.saved = true
    }

    function render() {
      if (!data.value) return h(SSpinner, { id: 'events-signup' })
      if (!data.value.events.length)
        return h(
          'div',
          { id: 'events-signup' },
          'There are no upcoming events for which you can sign up.'
        )
      return h('form', { id: 'events-signup', onSubmit }, [
        data.value.events.map((e) => renderEvent(data.value!, e)),
        h('div', { id: 'events-signup-buttons' }, [
          h(SButton, { variant: 'primary', type: 'submit' }, () => 'Save'),
          data.value.saved ? h('div', { id: 'events-signup-saved' }, 'Saved.') : null,
        ]),
      ])
    }
    return render
  },
})
export default EventSignups

function renderEvent(d: GetEventSignups, e: GetEventSignupsEvent) {
  function onDecline(checked: boolean) {
    if (!checked) return
    d.saved = false
    e.declined = true
    e.shifts.forEach((s) => {
      s.signedUp = false
    })
  }
  return [
    h(
      'div',
      { class: 'events-signup-date' },
      moment(e.date, 'YYYY-MM-DD').format('dddd, MMMM D, YYYY')
    ),
    h('div', { class: 'events-signup-event' }, e.name),
    e.signupText ? h('div', { class: 'events-signup-text', innerHTML: e.signupText }) : null,
    h('div', { class: 'events-signup-shifts' }, [
      e.shifts.map((s) => renderShift(d, e, s)),
      h(SCheck, {
        id: `${e.id}.declined`,
        label: 'Decline',
        modelValue: e.declined,
        disabled: e.declined,
        'onUpdate:modelValue': onDecline,
      }),
    ]),
  ]
}

function renderShift(d: GetEventSignups, e: GetEventSignupsEvent, s: GetEventSignupsShift) {
  function onUpdate(checked: boolean) {
    d.saved = false
    s.signedUp = checked
    if (!checked) {
      if (!e.shifts.find((s) => s.signedUp)) e.declined = true
    } else {
      e.declined = false
      e.shifts.forEach((o) => {
        if (s === o) return
        if (s.start < o.end && s.end > o.start) o.signedUp = false
      })
    }
  }
  return [
    h(SCheck, {
      id: `${e.id}.${s.start}.${s.task}`,
      class: 'events-signup-check',
      label: `${s.start}–${s.end} ${s.task}`,
      modelValue: s.signedUp,
      disabled: !s.signedUp && s.max > 0 && s.count >= s.max,
      'onUpdate:modelValue': onUpdate,
    }),
    renderStatus(e, s),
    s.expanded ? renderList(d, e, s) : null,
  ]
}

function renderStatus(e: GetEventSignupsEvent, s: GetEventSignupsShift) {
  function onClick() {
    s.expanded = !s.expanded
  }
  const count = s.signedUp ? s.baseCount + 1 : s.baseCount
  // If the count and the min are both 10 or less, we can use silhouettes.
  if (count <= 10 && s.min <= 10) {
    const silhouettes: Array<VNode> = []
    for (let i = 0; i < count; i++)
      silhouettes.push(h(SIcon, { class: 'events-signup-signedup', icon: 'user' }))
    for (let i = count; i < s.min; i++)
      silhouettes.push(h(SIcon, { class: 'events-signup-needed', icon: 'user-outline' }))
    // If the max is also <= 10, we can use silhouettes for that too.
    if (s.max !== 0 && s.max <= 10)
      for (let i = Math.max(count, s.min); i < s.max; i++)
        silhouettes.push(h(SIcon, { class: 'events-signup-allowed', icon: 'user-outline' }))
    else if (s.max !== 0)
      silhouettes.push(h('span', { class: 'events-signup-max' }, `(${s.max} allowed)`))
    else silhouettes.push(h('span', { class: 'events-signup-max' }, '(no maximum)'))
    return h('div', { class: 'events-signup-status', onClick }, silhouettes)
  }

  // If the silhouettes won't fit, use text.
  if (count >= s.min)
    if (s.max > 0)
      return h(
        'div',
        { class: 'events-signup-status', onClick },
        `${count} signed up, ${s.max} allowed`
      )
    else
      return h('div', { class: 'events-signup-status', onClick }, `${count} signed up (no maximum)`)
  return h('div', { class: 'events-signup-status', onClick }, `${count} signed up, need ${s.min}`)
}

function renderList(d: GetEventSignups, e: GetEventSignupsEvent, s: GetEventSignupsShift) {
  return h('div', { class: 'events-signup-list' }, [
    s.list.map((p) =>
      h(
        'div',
        d.loggedIn
          ? h(RouterLink, { to: `/api/person/${p.id}` }, p.sortName)
          : h('span', p.sortName)
      )
    ),
    s.signedUp
      ? h(
          'div',
          d.loggedIn
            ? h(RouterLink, { to: `/api/person/${d.id}` }, d.sortName)
            : h('span', d.sortName)
        )
      : null,
  ])
}
