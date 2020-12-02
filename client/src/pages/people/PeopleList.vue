<!--
PeopleList displays the list of people.
-->

<template lang="pug">
#people-list
  select#people-role(v-if='roles && roles.length > 1', v-model='role')
    option(v-for='t in roles', :key='t.id', :value='t.id', v-text='t.name')
  #people-list-spinner(v-if='loading')
    SSpinner
  template(v-else)
    #people-table
      .people-person.people-heading Person
      .people-contact.people-heading Contact Info
      .people-roles.people-heading Roles
      template(v-for='p in people')
        .people-person
          router-link(
            :to='`/people/${p.id}`',
            v-text='p.callSign ? `${p.sortName} (${p.callSign})` : p.sortName'
          )
          .people-badge(
            v-for='b in p.badges',
            :class='b.startsWith("No") ? "people-badge-no" : "people-badge-yes"',
            v-text='b'
          )
        .people-contact
          div(v-if='p.email')
            a(:href='`mailto:${p.email}`', v-text='p.email')
          div(v-if='p.email2')
            a(:href='`mailto:${p.email2}`', v-text='p.email2')
          .people-phone(v-if='p.cellPhone')
            a(:href='`tel:${p.cellPhone}`', v-text='p.cellPhone')
            |
            | (Cell)
          .people-phone(v-if='p.homePhone')
            a(:href='`tel:${p.homePhone}`', v-text='p.homePhone')
            |
            | (Home)
          .people-phone(v-if='p.workPhone')
            a(:href='`tel:${p.workPhone}`', v-text='p.workPhone')
            |
            | (Work)
        .people-roles
          div(v-for='(r, i) in p.roles', :key='i', v-text='r')
          div(v-if='!p.roles.length') &mdash;
    div(v-if='people.length === 1') 1 person listed.
    div(v-else, v-text='`${people.length} people listed.`')
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Cookies from 'js-cookie'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SSpinner } from '../../base'

export type GetPeopleAddress = {
  address?: string
  sameAsHome: boolean
  latitude?: number
  longitude?: number
  fireDistrict: number
}
export type GetPeoplePerson = {
  id: number
  informalName: string
  sortName: string
  callSign: string
  email?: string
  email2?: string
  homeAddress?: GetPeopleAddress
  mailAddress?: GetPeopleAddress
  workAddress?: GetPeopleAddress
  cellPhone?: string
  homePhone?: string
  workPhone?: string
  badges?: Array<string>
  roles?: Array<string>
  // temp
  canEdit: boolean
  canHours: boolean
}
export type GetPeopleViewableRole = {
  id: number
  name: string
}
export type GetPeople = {
  people: Array<GetPeoplePerson>
  viewableRoles: Array<GetPeopleViewableRole>
  canAdd: boolean
}

export default defineComponent({
  components: { SSpinner },
  setup() {
    const route = useRoute()
    setPage({ title: 'People' })

    // The role being viewed.
    const role = ref(
      parseInt((route.query.role as string) || Cookies.get('serv-people-role') || '0')
    )
    const roles = ref([] as Array<GetPeopleViewableRole>)
    const loading = ref(true)
    const people = ref([] as Array<GetPeoplePerson>)
    watch(
      role,
      async () => {
        Cookies.set('serv-people-role', role.value.toString(), { expires: 3650 })
        loading.value = true
        const data = (await axios.get<GetPeople>('/api/people', { params: { role: role.value } }))
          .data
        people.value = data.people
        people.value.forEach((p) => {
          if (p.roles) p.roles.sort()
        })
        if (data.viewableRoles.length > 1) {
          data.viewableRoles.sort((a, b) => a.name.localeCompare(b.name))
          data.viewableRoles.unshift({ id: 0, name: '(all)' })
          roles.value = data.viewableRoles
        }
        loading.value = false
      },
      { immediate: true }
    )

    return { role, roles, loading, people }
  },
})
</script>

<style lang="postcss">
#people-list {
  padding: 1.5rem 0.75rem 0.75rem;
}
#people-list-spinner {
  margin-top: 1.5rem;
}
#people-table {
  display: flex;
  flex-direction: column;
  margin-top: 0.75rem;
  @media (min-width: 576px) {
    display: grid;
    grid: auto / 10rem 1fr;
  }
  @media (min-width: 800px) {
    grid: auto / 10rem 1fr 1fr;
  }
  @media (min-width: 960px) {
    grid: auto / 10rem 21rem 1fr;
  }
}
.people-heading {
  display: none;
  @media (min-width: 576px) {
    display: block;
    font-weight: bold;
  }
}
.people-person {
  overflow: hidden;
  margin-top: 0.25rem;
  padding-left: 1.5rem;
  text-indent: -1.5rem;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
  @media (min-width: 576px) {
    margin-top: 0.75rem;
    text-overflow: clip;
    white-space: normal;
  }
}
.people-badge {
  display: inline-block;
  margin-left: 1.5rem;
  padding: 0 0.25rem;
  border-radius: 4px;
  color: white;
  text-indent: 0;
  font-size: 0.75rem;
}
.people-badge-no {
  background-color: red;
}
.people-badge-yes {
  background-color: green;
}
.people-contact {
  margin-left: 6rem;
  div {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  @media (min-width: 576px) {
    margin-top: 0.75rem;
    margin-left: 0.25rem;
  }
}
.people-phone {
  font-variant: tabular-nums;
}
.people-roles {
  display: none;
  @media (min-width: 800px) {
    display: block;
    margin-top: 0.75rem;
    margin-left: 0.25rem;
    div {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}
</style>
