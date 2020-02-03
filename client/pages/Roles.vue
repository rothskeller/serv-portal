<!--
Roles displays and edits roles and groups.
-->

<template lang="pug">
b-card#roles-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/roles" exact exact-active-class="active") Roles
      b-nav-item(v-if="canAddRole" to="/roles/NEW") Add Role
      b-nav-item(v-if="canEditRole" :to="`/roles/${$route.params.id}`" exact exact-active-class="active") {{editLabel}}
  #roles-scroll
    router-view(:onLoadRole="onLoadRole")
</template>

<script>
export default {
  data: () => ({ role: null }),
  computed: {
    canAddRole() { return !this.$route.params.id },
    canEditRole() { return this.role },
    editLabel() { return this.$route.params.id === 'NEW' ? 'Add Role' : 'Edit Role' },
  },
  watch: {
    $route() {
      if (!this.$route.params.id) {
        this.role = null
        this.$store.commit('setPage', { title: 'Roles' })
      }
    },
  },
  mounted() {
    this.$store.commit('setPage', { title: this.$route.params.id === 'NEW' ? 'New Role' : 'Roles' })
  },
  methods: {
    onLoadRole(r) {
      this.role = r
      if (this.$route.params.id !== 'NEW') this.$store.commit('setPage', { title: r.name })
    },
  },
}
</script>

<style lang="stylus">
#roles-card
  height 100%
  border none
  .card-header
    @media print
      display none
#roles-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
