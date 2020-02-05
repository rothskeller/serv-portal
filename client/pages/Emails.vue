<!--
Emails displays and handles email messages.
-->

<template lang="pug">
b-card#emails-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/emails" exact exact-active-class="active") Emails
  #emails-scroll
    router-view(:onLoadEmail="onLoadEmail")
</template>

<script>
export default {
  data: () => ({ email: null }),
  watch: {
    $route() {
      if (!this.$route.params.id) {
        this.email = null
        this.$store.commit('setPage', { title: 'Emails' })
      }
    },
  },
  mounted() {
    this.$store.commit('setPage', { title: 'Emails' })
  },
  methods: {
    onLoadEmail(e) {
      this.email = e
      this.$store.commit('setPage', { title: e.name })
    },
  },
}
</script>

<style lang="stylus">
#emails-card
  height 100%
  border none
  .card-header
    @media print
      display none
#emails-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
