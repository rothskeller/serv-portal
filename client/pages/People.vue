<!--
People displays the list of people.
-->

<template lang="pug">
b-card#people-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/people/list" exact exact-active-class="active") List
      b-nav-item(to="/people/map" exact exact-active-class="active") Map
      b-nav-item(v-if="canAdd" to="/people/NEW/edit") Add Person
      b-nav-item(v-if="canView" :to="`/people/${$route.params.id}`" exact exact-active-class="active") Details
      b-nav-item(v-if="canEdit" :to="`/people/${$route.params.id}/edit`" exact exact-active-class="active") {{editLabel}}
  #people-scroll
    router-view(:onLoadPerson="onLoadPerson")
</template>

<script>
export default {
  data: () => ({ person: null }),
  computed: {
    canAdd() { return !this.$route.params.id && this.$store.state.me.canAddPeople },
    canEdit() { return this.person && this.person.canEdit },
    canView() { return this.person && this.$route.params.id !== 'NEW' },
    editLabel() { return this.$route.params.id === 'NEW' ? 'Add Person' : 'Edit' },
  },
  watch: {
    $route() {
      if (!this.$route.params.id) {
        this.person = null
        this.$store.commit('setPage', { title: 'People' })
      }
    },
  },
  mounted() {
    this.$store.commit('setPage', { title: this.$route.params.id === 'NEW' ? 'New Person' : 'People' })
  },
  methods: {
    onLoadPerson(p) {
      this.person = p
      if (this.$route.params.id !== 'NEW') this.$store.commit('setPage', { title: p.informalName })
    },
  },
}
</script>

<style lang="stylus">
#people-card
  height 100%
  border none
  .card-header
    @media print
      display none
#people-scroll
  display flex
  flex auto
  flex-direction column
  overflow-x hidden
  overflow-y auto
</style>
