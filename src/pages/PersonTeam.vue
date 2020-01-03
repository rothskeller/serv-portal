<!--
PersonTeam displays the row for a single team on the Person page.
-->

<template lang="pug">
tr.person-team-row
  td(v-if="!team.canAdmin && team.role" v-text="`${team.name}: ${roleName(team.role)}`")
  td(v-else-if="!team.canAdmin" v-text="team.name")
  td(v-else)
    b-checkbox(v-if="manageAny" :checked="!!team.role" @change="onMemberChange") {{team.name}}
    span(v-else v-text="team.name")
  td(v-if="team.canAdmin && team.roles.length>1 && !!team.role")
    select(v-model="team.role")
      option(v-for="r in team.roles" :key="r.id" :value="r.id" v-text="r.name || '(member)'")
</template>

<script>
export default {
  props: {
    manageAny: Boolean,
    team: Object,
  },
  methods: {
    onMemberChange() { this.team.role = this.team.role ? 0 : this.team.roles[0].id },
    roleName(id) { return this.team.roles.find(r => r.id === id).name },
  },
}
</script>

<style lang="stylus"></style>
