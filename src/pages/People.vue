<!--
People displays the list of people.
-->

<template lang="pug">
Page(title="People" menuItem="people")
  #people-title
    | People
    select#people-group(v-if="groups && groups.length > 1" v-model="group")
      option(v-for="t in groups" :key="t.id" :value="t.id" v-text="t.name")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#people-table(v-else)
    tr
      th Person
      th Contact Info
      th Roles
    tr(v-for="p in people" :key="p.id")
      td: router-link(:to="`/people/${p.id}`" v-text="p.callSign ? `${p.sortName} (${p.callSign})` : p.sortName")
      td
        div(v-for="email in p.emails")
          a(:href="`mailto:${email.email}`" v-text="email.email" :class="email.bad ? '.people-bad-email' : null")
          span(v-if="email.label" v-text="` (${email.label})`")
        div(v-for="phone in p.phones")
          a(:href="`tel:${phone.phone}`" v-text="phone.phone")
          span(v-if="phone.label" v-text="` (${phone.label})`")
      td
        div(v-for="(r, i) in p.roles" :key="i" v-text="r")
        div(v-if="!p.roles.length") &mdash;
  //-div.mt-3(v-if="canAdd")
    b-btn(to="/people/NEW") Add Person
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
#people-title
  display flex
  align-items center
  font-size 1.5rem
#people-group
  margin-left 1rem
  font-size 1rem
#people-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
    line-height 1.2
.people-bad-email
  text-decoration line-through
</style>
