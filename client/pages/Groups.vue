<!--
Groups displays and edits groups.
-->

<template lang="pug">
b-card#groups-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/groups" exact exact-active-class="active") Groups
      b-nav-item(v-if="canAddGroup" to="/groups/NEW") Add Group
      b-nav-item(v-if="canEditGroup" :to="`/groups/${$route.params.id}`" exact exact-active-class="active") {{editLabel}}
  #groups-scroll
    router-view(:onLoadGroup="onLoadGroup")
</template>

<script>
export default {
  data: () => ({ group: null }),
  computed: {
    canAddGroup() { return !this.$route.params.id },
    canEditGroup() { return this.group },
    editLabel() { return this.$route.params.id === 'NEW' ? 'Add Group' : 'Edit Group' },
  },
  watch: {
    $route() {
      if (!this.$route.params.id) {
        this.group = null
        this.$store.commit('setPage', { title: 'Groups' })
      }
    },
  },
  mounted() {
    this.$store.commit('setPage', { title: this.$route.params.id === 'NEW' ? 'New Group' : 'Groups' })
  },
  methods: {
    onLoadGroup(g) {
      this.group = g
      if (this.$route.params.id !== 'NEW') this.$store.commit('setPage', { title: g.name })
    },
  },
}
</script>

<style lang="stylus">
#groups-card
  height 100%
  border none
  .card-header
    @media print
      display none
#groups-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
