<!--
Roles displays the list of roles.
-->

<template lang="pug">
Page(title="Roles" subtitle="Roles" menuItem="roles")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#roles-table(v-else)
    tr(v-for="role in roles" :key="role.id")
      td: b-link(:to="`/roles/${role.id}`" v-text="role.name")
  div.mt-3
    b-btn(to="/roles/NEW") Add Role
</template>

<script>
export default {
  data: () => ({ loading: false, roles: null }),
  async created() {
    this.loading = true
    this.roles = (await this.$axios.get('/api/roles')).data
    this.loading = false
  },
}
</script>

<style lang="stylus">
#roles-table
  margin-top 1.5rem
  th, td
    padding 0.25rem 1em 0 0
</style>
