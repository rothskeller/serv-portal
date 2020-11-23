<!--
PersonView displays the information about a person, in non-editable form.
-->

<template lang="pug">
#person-view-spinner(v-if='!person')
  SSpinner
#person-view(v-else)
  #person-view-name
    #person-view-informalName
      span(v-text='person.informalName')
      span#person-view-callSign(v-if='person.callSign', v-text='person.callSign')
    #person-view-formalName(
      v-if='person.formalName !== person.informalName',
      v-text='`(formally ${person.formalName})`'
    )
  #person-view-roles
    div(v-for='role in person.roles', v-text='role.name')
  #person-view-emails(v-if='person.email')
    div: a(:href='`mailto:${person.email}`', v-text='person.email')
    div(v-if='person.email2')
      a(:href='`mailto:${person.email2}`', v-text='person.email2')
  .person-view-phone(v-if='person.cellPhone')
    a(:href='`tel:${person.cellPhone}`', v-text='person.cellPhone')
    span.person-view-phone-label (Cell)
  .person-view-phone(v-if='person.homePhone')
    a(:href='`tel:${person.homePhone}`', v-text='person.homePhone')
    span.person-view-phone-label (Home)
  .person-view-phone(v-if='person.workPhone')
    a(:href='`tel:${person.workPhone}`', v-text='person.workPhone')
    span.person-view-phone-label (Work)
  .person-view-address(v-if='person.homeAddress && person.homeAddress.address')
    div
      span(v-if='person.workAddress && person.workAddress.sameAsHome') Home Address (all day):
      span(v-else) Home Address:
      a.person-view-map(
        target='_blank',
        :href='`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.homeAddress.address)}`'
      ) Map
    div(v-text='person.homeAddress.address.split(",")[0]')
    div(v-text='person.homeAddress.address.replace(/^[^,]*, */, "")')
    div(
      v-if='person.homeAddress.fireDistrict',
      v-text='`Sunnyvale Fire District ${person.homeAddress.fireDistrict}`'
    )
  .person-view-address(v-if='person.workAddress && person.workAddress.address')
    div
      span Work Address:
      a.person-view-map(
        target='_blank',
        :href='`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.workAddress.address)}`'
      ) Map
    div(v-text='person.workAddress.address.split(",")[0]')
    div(v-text='person.workAddress.address.replace(/^[^,]*, */, "")')
    div(
      v-if='person.workAddress.fireDistrict',
      v-text='`Sunnyvale Fire District ${person.workAddress.fireDistrict}`'
    )
  .person-view-address(v-if='person.mailAddress && person.mailAddress.address')
    div Mailing Address:
    div(v-text='person.mailAddress.address.split(",")[0]')
    div(v-text='person.mailAddress.address.replace(/^[^,]*, */, "")')
  .person-view-clearances(
    v-if='person.dsw || person.volgistics !== undefined || person.backgroundCheck !== undefined' || person.identification.length
  )
    div Volunteer status:
    div(v-if='person.volgistics === false') Sunnyvale volunteer: <span style="color:red">not registered.</span>
    div(
      v-else-if='person.volgistics && person.volgistics !== true',
      v-text='`Sunnyvale volunteer #${person.volgistics}`'
    )
    template(v-if='person.dsw')
      div(v-for='(dsw, cls) in person.dsw')
        | DSW for {{ cls }}:
        |
        span(v-if='dsw.needed && !dsw.registered', style='color:red') not registered.
        template(v-else)
          | registered {{ dsw.registered }},
          |
          span(v-if='dsw.expired', style='color:red') expired
          span(v-else) expires
          |
          | {{ dsw.expires }}.
    div(v-if='person.backgroundCheck === true') Background check cleared.
    div(v-else-if='person.backgroundCheck === false') Background check: <span style="color:red">not completed.</span>
    div(v-else-if='person.backgroundCheck') Background check cleared on {{ person.backgroundCheck }}.
    div(v-if='person.identification.length') IDs issued: {{ person.identification.join(', ') }}
  #person-view-attended(v-if='person.attended')
    div Events attended:
    .person-view-attended(v-for='e in person.attended')
      span.person-view-attended-date(v-text='e.date')
      span(v-text='e.name')
  #person-view-notes(v-if='person.notes && person.notes.length')
    div Notes:
    .person-view-note(v-for='n in person.notes')
      span.person-view-note-date(v-text='n.date')
      span(v-text='n.note')
  #person-view-unsub(v-if='person.noEmail || person.noText')
    div Unsubscribe:
    div(v-if='person.noEmail && person.noText') from all emails and text messages
    div(v-else-if='person.noEmail') from all emails
    div(v-else) from all text messages
  #person-view-roles2(v-if="me.webmaster")
    div Roles:
    div(v-for="r in person.roles2" v-text="r")
    div(v-if="!person.roles2.length") None
    div(v-if="person.canEditRoles2")
      SButton(variant='secondary', small, @click="onEditRoles") Edit
    PersonEditRoles(ref='editRolesModal', :pid='person.id')
  #person-view-lists(v-if="me.webmaster && person.lists")
    div Lists:
    div(v-for="l in person.lists" v-text="l")
    div(v-if="!person.lists.length") None
    div: SButton(variant='secondary', small, @click='onEditLists') Edit
    PersonEditLists(ref='editListsModal', :pid='person.id')
</template>

<script lang="ts">
import { defineComponent, inject, Ref, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../../plugins/axios'
import { LoginData } from '../../plugins/login'
import { SButton, SSpinner } from '../../base'
import PersonEditLists from './PersonEditLists.vue'
import PersonEditRoles from './PersonEditRoles.vue'

export type GetPersonAddress = {
  address?: string
  sameAsHome: boolean
  latitude?: number
  longitude?: number
  fireDistrict: number
}
export type GetPersonRole = {
  id: number
  name: string
}
interface GetPersonDSW {
  needed?: true
  registered?: string
  expires?: string
  expired?: true
}
interface GetPersonAttended {
  id: number
  date: string
  name: string
  type: string
  minutes: number
}
interface GetPersonNote {
  date: string
  note: string
}
export interface GetPersonPersonBase {
  id: number
  informalName: string
  formalName: string
  sortName: string
  callSign: string
  email?: string
  email2?: string
  homeAddress?: GetPersonAddress
  mailAddress?: GetPersonAddress
  workAddress?: GetPersonAddress
  cellPhone?: string
  homePhone?: string
  workPhone?: string
  roles: Array<GetPersonRole>
  roles2: Array<string>
  lists?: Array<string>
  canEdit: boolean
  canEditRoles2: boolean
  canHours: boolean
  noEmail: boolean
  noText: boolean
}
interface GetPersonPerson extends GetPersonPersonBase {
  volgistics?: number | boolean
  backgroundCheck?: string | boolean
  dsw?: Record<string, GetPersonDSW>
  attended?: Array<GetPersonAttended>
  notes: Array<GetPersonNote>
}
interface GetPerson {
  person: GetPersonPerson
}

export default defineComponent({
  components: { PersonEditLists, PersonEditRoles, SButton, SSpinner },
  props: {
    onLoadPerson: { type: Function, required: true },
  },
  setup(props) {
    const me = inject<Ref<LoginData>>('me')!
    const route = useRoute()
    const person = ref(null as null | GetPersonPerson)
    axios.get<GetPerson>(`/api/people/${route.params.id}`).then((resp) => {
      person.value = resp.data.person
      props.onLoadPerson(person.value)
    })
    const editRolesModal = ref(null as any)
    async function onEditRoles() {
      if (!(await editRolesModal.value.show())) return
      person.value = (await axios.get<GetPerson>(`/api/people/${route.params.id}`)).data.person
    }
    const editListsModal = ref(null as any)
    async function onEditLists() {
      if (!(await editListsModal.value.show())) return
      person.value = (await axios.get<GetPerson>(`/api/people/${route.params.id}`)).data.person
    }
    return { editListsModal, editRolesModal, me, onEditLists, onEditRoles, person }
  },
})
</script>

<style lang="postcss">
#person-view {
  margin: 1.5rem 0.75rem;
}
#person-view-name {
  display: flex;
  flex-direction: column;
  @media (min-width: 576px) {
    flex-direction: row;
  }
}
#person-view-informalName {
  font-weight: bold;
  font-size: 1.25rem;
  line-height: 1.2;
}
#person-view-callSign {
  margin-left: 0.5rem;
  font-weight: normal;
}
#person-view-formalName {
  color: #888;
  @media (min-width: 576px) {
    margin-left: 1rem;
  }
}
#person-view-cbadge {
  margin-left: 0.5rem;
}
#person-view-roles {
  line-height: 1.2;
}
#person-view-emails {
  margin-top: 0.75rem;
}
.person-view-email-label {
  color: #888;
}
.person-view-phone {
  font-variant: tabular-nums;
}
.person-view-phone-label {
  margin-left: 0.25rem;
  color: #888;
}
.person-view-address,
.person-view-clearances {
  margin-top: 0.75rem;
  line-height: 1.2;
}
.person-view-map {
  margin-left: 1rem;
  &::before {
    content: '[';
  }
  &::after {
    content: ']';
  }
}
.person-view-address-flag {
  color: #888;
  font-size: 0.875rem;
}
#person-view-attended {
  margin-top: 0.75rem;
  line-height: 1.2;
}
.person-view-attended {
  margin-left: 2rem;
  text-indent: -2rem;
}
.person-view-attended-date {
  margin-right: 0.75rem;
  color: #888;
  font-variant: tabular-nums;
}
#person-view-notes,
#person-view-unsub {
  margin-top: 0.75rem;
  line-height: 1.2;
}
.person-view-note {
  margin-left: 2rem;
  text-indent: -2rem;
}
.person-view-note-date {
  margin-right: 0.75rem;
  color: #888;
  font-variant: tabular-nums;
}
</style>
