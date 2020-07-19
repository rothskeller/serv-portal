<!--
Admin displays the administrative UI.
-->

<template lang="pug">
b-card#admin-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/admin/groups" exact exact-active-class="active") Groups
      b-nav-item(v-if="canEditGroup" :to="`/admin/groups/${$route.params.gid}`" exact exact-active-class="active") {{editGroupLabel}}
      b-nav-item(to="/admin/roles" exact exact-active-class="active") Roles
      b-nav-item(v-if="canEditRole" :to="`/admin/roles/${$route.params.rid}`" exact exact-active-class="active") {{editRoleLabel}}
  #admin-scroll
    router-view(:onLoadGroup="onLoadGroup" :onLoadRole="onLoadRole")
</template>

<script>
export default {
  data: () => ({ group: null, role: null }),
  computed: {
    canEditGroup() { return this.group },
    canEditRole() { return this.role },
    editGroupLabel() { return this.$route.params.gid === 'NEW' ? 'Add Group' : 'Edit Group' },
    editRoleLabel() { return this.$route.params.rid === 'NEW' ? 'Add Role' : 'Edit Role' },
  },
  watch: {
    $route() {
      if (!this.$route.params.gid) this.group = null
      if (!this.$route.params.rid) this.role = null
      if (this.group && this.group.id) this.$store.commit('setPage', { title: this.group.name })
      else if (this.group) this.$store.commit('setPage', { title: 'New Group' })
      else if (this.role && this.role.id) this.$store.commit('setPage', { title: this.role.name })
      else if (this.role) this.$store.commit('setPage', { title: 'New Role' })
      else this.$store.commit('setPage', { title: 'Administration' })
    },
  },
  methods: {
    onLoadGroup(g) {
      this.group = g
      this.$store.commit('setPage', { title: this.group.id ? this.group.name : 'New Group' })
    },
    onLoadRole(r) {
      this.role = r
      this.$store.commit('setPage', { title: this.role.id ? this.role.name : 'New Group' })
    },
  },
}
</script>

<style lang="stylus">
#admin-card
  height 100%
  border none
  .card-header
    @media print
      display none
#admin-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
