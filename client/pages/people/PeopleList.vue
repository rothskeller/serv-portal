<!--
PeopleList displays the list of people.
-->

<template lang="pug">
#people-list
  select#people-group(v-if="groups && groups.length > 1" v-model="group")
    option(v-for="t in groups" :key="t.id" :value="t.id" v-text="t.name")
  div.mt-3(v-if="loading")
    b-spinner(small)
  #people-table(v-else)
    .people-person.people-heading Person
    .people-contact.people-heading Contact Info
    .people-roles.people-heading Roles
    template(v-for="p in people")
      .people-person
        td: router-link(:to="`/people/${p.id}`" v-text="p.callSign ? `${p.sortName} (${p.callSign})` : p.sortName")
      .people-contact
        div(v-for="email in p.emails")
          a(:href="`mailto:${email.email}`" v-text="email.email" :class="email.bad ? '.people-bad-email' : null")
          span(v-if="email.label" v-text="` (${email.label})`")
        .people-phone(v-if="p.cellPhone")
          a(:href="`tel:${p.cellPhone}`" v-text="p.cellPhone")
          |
          | (Cell)
        .people-phone(v-if="p.homePhone")
          a(:href="`tel:${p.homePhone}`" v-text="p.homePhone")
          |
          | (Home)
        .people-phone(v-if="p.workPhone")
          a(:href="`tel:${p.workPhone}`" v-text="p.workPhone")
          |
          | (Work)
      .people-roles
        div(v-for="(r, i) in p.roles" :key="i" v-text="r")
        div(v-if="!p.roles.length") &mdash;
</template>

<script>
import Cookies from 'js-cookie'

export default {
  data: () => ({ group: 0, groups: null, people: null, canAdd: false, loading: false }),
  created() {
    this.group = Cookies.get('serv-people-group') || 0
    this.load()
  },
  watch: {
    group() {
      Cookies.set('serv-people-group', this.group, { expires: 3650 })
      this.load()
    },
  },
  methods: {
    async load() {
      this.loading = true
      const data = (await this.$axios.get('/api/people', { params: { group: this.group } })).data
      this.people = data.people
      this.canAdd = data.canAdd
      if (data.viewableGroups.length > 1) {
        data.viewableGroups.unshift({ id: 0, name: '(all)' })
        this.groups = data.viewableGroups
      }
      this.loading = false
    },
  },
}
</script>

<style lang="stylus">
#people-list
  margin 1.5rem 0.75rem 0.75rem
#people-table
  display flex
  flex-wrap wrap
  margin-top 0.755rem
.people-heading
  display none
  @media (min-width: 576px)
    display block
    font-weight bold
.people-person
  overflow hidden
  margin-top 0.25rem
  width calc(100vw - 1.5rem)
  text-overflow ellipsis
  white-space nowrap
  @media (min-width: 576px)
    margin-top 0.75rem
    width 10rem
    white-space normal
.people-contact
  padding-left 8rem
  width calc(100vw - 1.5rem)
  div
    overflow hidden
    text-overflow ellipsis
    white-space nowrap
  @media (min-width: 576px)
    margin-top 0.75rem
    padding-left 0.25rem
    width calc(100vw - 18.5rem)
  @media (min-width: 800px)
    width calc(50vw - 9.25rem)
  @media (min-width: 960px)
    width 20.75rem
.people-phone
  font-variant tabular-nums
.people-roles
  display none
  @media (min-width: 800px)
    display block
    margin-top 0.75rem
    padding-left 0.25rem
    width calc(50vw - 9.25rem)
    div
      overflow hidden
      text-overflow ellipsis
      white-space nowrap
  @media (min-width: 960px)
    width calc(100vw - 39.25rem)
.people-bad-email
  text-decoration line-through
</style>
