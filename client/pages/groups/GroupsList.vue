<!--
GroupsList displays the list of groups.
-->

<template lang="pug">
#groups-list
  #groups-list-spinner(v-if="loading")
    b-spinner(small)
  #groups-list-table(v-else)
    .groups-list-name.groups-list-heading Group
    .groups-list-roles.groups-list-heading Included in Roles
    template(v-for="g in groups")
      .groups-list-name
        router-link(:to="`/groups/${g.id}`" v-text="g.name")
      .groups-list-roles
        div(v-for="r in g.roles" v-text="r")
</template>

<script>
export default {
  data: () => ({
    groups: null,
    loading: true,
  }),
  async created() {
    this.loading = true
    this.groups = (await this.$axios.get(`/api/groups`)).data
    this.loading = false
  },
}
</script>

<style lang="stylus">
#groups-list
  padding 1.5rem 0.75rem
#groups-list-spinner
  margin-top 1.5rem
#groups-list-table
  display grid
  grid auto / 1fr 1fr
  @media (min-width: 576px)
    grid auto / 16rem 1fr
.groups-list-heading
  font-weight bold
.groups-list-name
  flex none
  margin 0.75rem 0.75rem 0 0
  font-variant tabular-nums
.groups-list-roles
  flex none
  overflow hidden
  margin-top 0.75rem
  white-space nowrap
  div
    overflow hidden
    text-overflow ellipsis
</style>
