<!--
People displays the list of people.
-->

<template lang="pug">
Page(title="People" menuItem="people")
  #people-title
    | People
    select#people-role(v-if="roles && roles.length > 1" v-model="role")
      option(v-for="t in roles" :key="t.id" :value="t.id" v-text="t.name")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#people-table(v-else)
    tr
      th Person
      th Contact Info
      th Roles
    tr(v-for="p in people" :key="p.id")
      td: router-link(:to="`/people/${p.id}`" v-text="`${p.lastName}, ${p.nickname}`")
      td
        div: a(:href="`mailto:${p.email}`" v-text="p.email")
        div(v-if="p.phone"): a(:href="`tel:${p.phone}`" v-text="p.phone")
      td
        div(v-for="(r, i) in p.roles" :key="i" v-text="r")
        div(v-if="!p.roles.length") &mdash;
  div.mt-3(v-if="canAdd")
    b-btn(to="/people/NEW") Add Person
</template>

<script>
export default {
  data: () => ({ role: 0, roles: null, people: null, canAdd: false, loading: false }),
  created() {
    this.load()
  },
  watch: {
    role() {
      this.load()
    },
  },
  methods: {
    async load() {
      this.loading = true
      const data = (await this.$axios.get('/api/people', { params: { role: this.role } })).data
      this.people = data.people
      this.canAdd = data.canAdd
      if (data.viewableRoles.length > 1) {
        data.viewableRoles.unshift({ id: 0, name: '(all)' })
        this.roles = data.viewableRoles
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
#people-role
  margin-left 1rem
  font-size 1rem
#people-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
    line-height 1.2
</style>
