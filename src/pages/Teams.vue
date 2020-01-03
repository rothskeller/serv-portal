<!--
Teams displays the list of teams and associated roles.
-->

<template lang="pug">
Page(title="Teams and Roles" subtitle="Teams and Roles" menuItem="teams")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#teams-table(v-else)
    tr
      th Team
      th Email
      th Roles
      th
    tr(v-for="team in teams" :key="team.id")
      td(:class="`indent-${team.indent}`")
        b-link(:to="`/teams/${team.id}`" v-text="team.name")
      td(v-text="team.email || '\u2014'")
      td
        div(v-for="role in team.roles" :key="role.id")
          b-link(:to="`/teams/${team.id}/roles/${role.id}`" v-text="role.name || '(member)'")
      td
        b-link(v-if="team.type === 'normal'" :to="`/teams/${team.id}/roles/NEW`") Add Role
        b-link(v-if="team.type === 'ancestor'" :to="`/teams/NEW?parent=${team.id}`") Add Child Team
</template>

<script>
export default {
  data: () => ({ loading: false, teams: null }),
  async created() {
    this.loading = true
    this.teams = (await this.$axios.get('/api/teams')).data
    this.loading = false
  },
}
</script>

<style lang="stylus">
#teams-table
  margin-top 1.5rem
  th, td
    padding 0.75rem 1em 0 0
    vertical-align top
    &.indent-1
      padding-left 1em
    &.indent-2
      padding-left 2em
    &.indent-3
      padding-left 3em
    &.indent-4
      padding-left 4em
    &.indent-5
      padding-left 5em
</style>
