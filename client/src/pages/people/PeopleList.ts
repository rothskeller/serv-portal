// PeopleList shows a roster of people belonging to a specified role.

import Cookies from 'js-cookie'
import { defineComponent, h, inject, nextTick, Ref, ref, watch } from 'vue'
import { Router, RouterLink, useRoute, useRouter } from 'vue-router'
import { SButton, SIcon } from '../../base'
import SSelect from '../../base/controls/SSelect'
import SSpinner from '../../base/widget/SSpinner'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import provideSize from '../../plugins/size'
import type { GetPeople, GetPeoplePerson, GetPeopleViewableRole } from './api'
import PersonDetailsPopup from './PersonDetailsPopup'
import './people.css'

const PeopleList = defineComponent({
  name: 'PeopleList',
  props: {
    onLoadPerson: Function, // not used
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const touch = inject<Ref<boolean>>('touch')!
    setPage({ title: 'People' })

    // Sort function for roles.
    const sort = ref('priority')
    function peopleSort(a: GetPeoplePerson, b: GetPeoplePerson): number {
      switch (sort.value) {
        case 'priority':
          if (a.priority !== b.priority) return a.priority - b.priority
        case 'name':
          return a.sortName.localeCompare(b.sortName)
        case 'callSignSuffix':
          const as = callSignSuffix(a.callSign)
          const bs = callSignSuffix(b.callSign)
          if (as !== bs) return as.localeCompare(bs)
        case 'callSign':
          if (a.callSign !== b.callSign) return a.callSign.localeCompare(b.callSign)
          return a.sortName.localeCompare(b.sortName)
        default:
          return 0
      }
    }
    function toggleSort() {
      switch (sort.value) {
        case 'priority':
          sort.value = 'name'
          break
        case 'name':
          sort.value = showCallSign.value > 0 ? 'callSignSuffix' : 'priority'
          showCallSign.value = showCallSign.value > 0 ? 2 : showCallSign.value
          break
        case 'callSignSuffix':
          sort.value = 'callSign'
          showCallSign.value = 1
          break
        case 'callSign':
          sort.value = 'priority'
          break
      }
      people.value.sort(peopleSort)
    }

    // The role being viewed.
    const role = ref(
      parseInt((route.query.role as string) || Cookies.get('serv-people-role') || '0')
    )
    const roles = ref([] as Array<GetPeopleViewableRole>)
    const loading = ref(true)
    const people = ref([] as Array<GetPeoplePerson>)
    const showCallSign = ref(0)
    const canAdd = ref(false)
    const compact = ref(false)
    watch(
      role,
      async () => {
        Cookies.set('serv-people-role', role.value.toString(), { expires: 3650 })
        loading.value = true
        const data = (await axios.get<GetPeople>('/api/people', { params: { role: role.value } }))
          .data
        showCallSign.value = !data.showCallSign ? 0 : sort.value === 'callSignSuffix' ? 2 : 1
        if (!data.showCallSign && (sort.value === 'callSign' || sort.value === 'callSignSuffix')) {
          sort.value = 'priority'
        }
        people.value = data.people.sort(peopleSort)
        canAdd.value = data.canAdd
        if (data.viewableRoles.length > 1) {
          data.viewableRoles.sort((a, b) => a.name.localeCompare(b.name))
          data.viewableRoles.unshift({ id: 0, name: '(all)' })
          roles.value = data.viewableRoles
        }
        compact.value = false
        desiredWidth.value = 0
        loading.value = false
        nextTick(() => nextTick(setCompact))
      },
      { immediate: true }
    )

    // Handle determination of compact vs. non-compact.
    const gridRef = ref<HTMLDivElement>()
    const mySize = provideSize()
    const desiredWidth = ref(0)
    watch([gridRef, () => mySize.pw], setCompact)
    function setCompact() {
      console.log('setCompact actual', mySize.pw, 'desired', desiredWidth.value)
      if (!desiredWidth.value) {
        if (!gridRef.value) {
          console.log('no grid ref')
          return
        }
        desiredWidth.value = gridRef.value.scrollWidth
      }
      compact.value = mySize.pw < desiredWidth.value
    }

    // The person whose details are being shown.
    const details = ref(null as null | GetPeoplePerson)

    function render() {
      return h('div', { id: 'people' }, [
        renderControls(roles.value, role, !loading.value && !!people.value.length, toggleSort),
        loading.value
          ? h(SSpinner)
          : renderPeople(
              people.value,
              showCallSign.value,
              touch.value,
              compact.value,
              !desiredWidth.value,
              (r: HTMLDivElement) => {
                nextTick(() => (gridRef.value = r))
              },
              router,
              details
            ),
        canAdd.value ? renderAddButton() : null,
      ])
    }

    return render
  },
})
export default PeopleList

function renderControls(
  roles: Array<GetPeopleViewableRole>,
  role: Ref<number>,
  showSort: boolean,
  toggleSort: Function
) {
  return h('div', { id: 'people-controls' }, [
    roles.length > 1
      ? h(SSelect, {
          options: roles,
          valueKey: 'id',
          labelKey: 'name',
          modelValue: role.value,
          'onUpdate:modelValue': (v: number) => (role.value = v),
        })
      : null,
    showSort ? h(SButton, { small: true, onClick: toggleSort }, () => 'Sort') : null,
  ])
}

function renderPeople(
  people: Array<GetPeoplePerson>,
  showCallSign: number,
  touch: boolean,
  compact: boolean,
  measure: boolean,
  setGridRef: (r: HTMLDivElement) => void,
  router: Router,
  details: Ref<null | GetPeoplePerson>
) {
  return h(
    //@ts-ignore doesn't like the ref
    'div',
    {
      id: 'people-grid',
      class: {
        'people-compact': compact,
        'people-grid-measure': measure,
      },
      ref: setGridRef,
    },
    people.map((p) => renderPerson(p, showCallSign, touch, router, details))
  )
}

function renderPerson(
  person: GetPeoplePerson,
  showCallSign: number,
  touch: boolean,
  router: Router,
  details: Ref<null | GetPeoplePerson>
) {
  const email = person.email || person.email2 || ''
  const phone = person.cellPhone || person.homePhone || person.workPhone || ''
  const roles = trimRoles(person.roles)
  function onClick() {
    router.push(`/people/${person.id}`)
  }
  return [
    h(
      'div',
      { class: 'people-callSign-prefix', onClick },
      showCallSign === 2 ? callSignPrefix(person.callSign) : ''
    ),
    h(
      'div',
      { class: 'people-callSign-suffix', onClick },
      showCallSign === 2
        ? callSignSuffix(person.callSign)
        : showCallSign === 1
        ? person.callSign
        : ''
    ),
    h('div', { class: 'people-nameRoles', onClick }, [
      h(RouterLink, { to: `/people/${person.id}`, class: 'people-name' }, () => person.sortName),
      h('div', { class: 'people-roles-n' }, roles),
    ]),
    h('div', { class: 'people-roles', onClick }, roles),
    h(
      //@ts-ignore optional content
      'div',
      { class: 'people-email' },
      email
        ? h('a', { href: `mailto:${encodeURIComponent(email)}`, target: '_blank' }, email)
        : null
    ),
    h('div', { class: 'people-phone' }, phone),
    h(
      'div',
      {
        class: 'people-details',
        onClick: () => (details.value = details.value === person ? null : person),
      },
      [
        h(SIcon, { icon: 'info' }),
        details.value === person ? h(PersonDetailsPopup, { person }) : null,
      ]
    ),
  ]
}

function renderAddButton() {}

function callSignPrefix(callSign: string): string {
  return callSign.replace(/[A-Z]*$/, '')
}
function callSignSuffix(callSign: string): string {
  return callSign.replace(/^[A-Z]*\d/, '')
}
function trimRoles(roles: Array<string>): string {
  if (!roles.length) return ''
  if (roles.length === 1) return roles[0]
  return `${roles[0]}, ...`
}
