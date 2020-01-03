<!--
People displays the list of people.
-->

<template lang="pug">
Page(title="People" menuItem="people")
  #people-title
    | People
    select#people-team(v-if="teams && teams.length > 1" v-model="team")
      option(v-for="t in teams" :key="t.id" :value="t.id" v-text="t.name")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#people-table(v-else)
    tr
      th Person
      th Contact Info
      th Roles
    tr(v-for="p in people" :key="p.id")
      td: router-link(:to="`/people/${p.id}`" v-text="`${p.lastName}, ${p.firstName}`")
      td
        div: a(:href="`mailto:${p.email}`" v-text="p.email")
        div(v-if="p.phone"): a(:href="`tel:${p.phone}`" v-text="p.phone")
      td
        div(v-for="r in p.roles" :key="r.team" v-text="r.role ? `${r.team}: ${r.role}` : r.team")
        div(v-if="!p.roles.length") &mdash;
  div.mt-3(v-if="canAdd")
    b-btn(to="/people/NEW") Add Person
</template>

<script>
export default {
  data: () => ({ team: 0, teams: null, people: null, canAdd: false, loading: false }),
  created() {
    this.load()
  },
  watch: {
    team() {
      this.load()
    },
  },
  methods: {
    async load() {
      this.loading = true
      const data = (await this.$axios.get('/api/people', { params: { team: this.team } })).data
      this.people = data.people
      this.canAdd = data.canAdd
      if (data.viewableTeams.length > 1) {
        data.viewableTeams.unshift({ id: 0, name: '(all)' })
        this.teams = data.viewableTeams
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
#people-team
  margin-left 1rem
  font-size 1rem
#people-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
    line-height 1.2
</style>
