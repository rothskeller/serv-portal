<!--
Events displays the list of events.
-->

<template lang="pug">
b-card#events-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/events/calendar" exact exact-active-class="active") Calendar
      b-nav-item(to="/events/list" exact exact-active-class="active") List
      b-nav-item(v-if="canAdd" to="/events/NEW") Add Event
      b-nav-item(v-if="canView" :to="`/events/${$route.params.id}`" exact exact-active-class="active") Details
      b-nav-item(v-if="canEdit" :to="`/events/${$route.params.id}/edit`" exact exact-active-class="active") {{editLabel}}
      b-nav-item(v-if="canAttendance" :to="`/events/${$route.params.id}/attendance`" exact exact-active-class="active") Attendance
  #events-scroll
    router-view(:onLoadEvent="onLoadEvent")
</template>

<script>
export default {
  data: () => ({ event: null }),
  computed: {
    canAdd() { return !this.$route.params.id && this.$store.state.me.canAddEvents },
    canAttendance() { return this.event && this.event.canAttendance },
    canEdit() { return this.event && this.event.canEdit },
    canView() { return this.event && this.$route.params.id !== 'NEW' },
    editLabel() { return this.$route.params.id === 'NEW' ? 'Add Event' : 'Edit' },
  },
  watch: {
    $route() {
      if (!this.$route.params.id) {
        this.event = null
        this.$store.commit('setPage', { title: 'Events' })
      }
    },
  },
  mounted() {
    this.$store.commit('setPage', { title: this.$route.params.id === 'NEW' ? 'New Event' : 'Events' })
  },
  methods: {
    onLoadEvent(e) {
      this.event = e
      if (this.$route.params.id !== 'NEW') this.$store.commit('setPage', { title: `${e.date} ${e.name}` })
    },
  },
}
</script>

<style lang="stylus">
#events-card
  height 100%
  border none
  .card-header
    @media print
      display none
#events-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
