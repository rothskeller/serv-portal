<!--
Texts displays the list of text messages we have sent.
-->

<template lang="pug">
b-card#texts-card(no-body)
  b-card-header(header-tag="nav")
    b-nav(card-header tabs)
      b-nav-item(to="/texts" exact exact-active-class="active") Sent Messages
      b-nav-item(to="/texts/send" exact exact-active-class="active") New Message
      b-nav-item(v-if="$route.params.id" :to="`/texts/${$route.params.id}`" exact exact-active-class="active") Delivery Status
  #texts-scroll
    router-view
</template>

<script>
export default {
  created() { this.$store.commit('setPage', { title: 'Text Messages' }) },
}
</script>

<style lang="stylus">
#texts-card
  height 100%
  border none
  .card-header
    @media print
      display none
#texts-scroll
  flex auto
  overflow-x hidden
  overflow-y auto
</style>
