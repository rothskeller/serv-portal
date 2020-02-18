<!--
Home page, shown when no other page is selected.
-->

<template lang="pug">
#home
  svg#home-closed-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 2048")
    path(fill="currentColor" d="M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V1992c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z")
  #home-welcome Welcome to SunnyvaleSERV.org!
  #home-closed-text Click here to open the menu.
  #home-helpers
    .home-helper
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to get details of upcoming classes and events.
    .home-helper(v-if="me.canViewRosters")
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to view team rosters and maps.
    .home-helper
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here for class materials and other documents.
    .home-helper(v-if="me.canViewReports")
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to generate reports.
    .home-helper(v-if="me.canSendTextMessages")
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to send group text messages.
    .home-helper(v-if="me.webmaster")
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div If you can see this, you shouldn't need handholding.
    .home-helper
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to change your password or contact info.
    .home-helper
      svg.home-helper-arrow(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M257.5 445.1l-22.2 22.2c-9.4 9.4-24.6 9.4-33.9 0L7 273c-9.4-9.4-9.4-24.6 0-33.9L201.4 44.7c9.4-9.4 24.6-9.4 33.9 0l22.2 22.2c9.5 9.5 9.3 25-.4 34.3L136.6 216H424c13.3 0 24 10.7 24 24v32c0 13.3-10.7 24-24 24H136.6l120.5 114.8c9.8 9.3 10 24.8.4 34.3z")
      div Click here to log out of the web site.
    .home-also-see
      div.mb-1 Also see:
      div.mb-1: b-link(to="/help/subscribe-calendar") Subscribe to the SERV calendar on your phone
      div.mb-1: b-link(to="/help/email-lists") Information about SERV email lists
</template>

<script>
export default {
  computed: {
    me() { return this.$store.state.me },
  },
  created() {
    this.$store.commit('setPage', { title: 'Sunnyvale SERV' })
  },
}
</script>

<style lang="stylus">
#home
  display grid
  grid-template 'arrow welcome' min-content 'text text' auto / 44px 1fr
  .page-menu-open &
    grid-template 'pad welcome' auto 'pad helpers' 1fr / 7rem 1fr
  @media (min-width: 576px)
    grid-template unquote("'welcome' auto 'helpers' 1fr / calc(576px - 7rem)")
#home-closed-arrow
  margin 0 6px
  color #f032e6
  grid-area arrow
  .page-menu-open &
    display none
  @media (min-width: 576px)
    display none
#home-welcome
  display flex
  justify-content center
  align-items center
  text-align center
  font-weight bold
  font-size 1.5rem
  grid-area welcome
  .page-menu-open &
    margin-bottom 0.5rem
    height calc(3.75rem + 1px)
    font-size 1.25rem
    line-height 1.2
  @media (min-width: 576px)
    margin-bottom 0.5rem
    height calc(3.75rem + 1px)
    font-size 1.25rem
    line-height 1.2
#home-closed-text
  margin-left 6px
  color #f032e6
  font-weight bold
  grid-area text
  .page-menu-open &
    display none
  @media (min-width: 576px)
    display none
#home-helpers
  display none
  .page-menu-open &
    display block
    grid-area helpers
  @media (min-width: 576px)
    display block
    grid-area helpers
.home-helper
  display flex
  align-items center
  height 2.125rem
  color #f032e6
  font-size 0.75rem
  line-height 1.2
  @media (min-width: 576px)
    font-size 1rem
.home-helper-arrow
  flex none
  margin-right 0.5rem
  width 1rem
  height 1rem
  color #f032e6
.home-also-see
  margin 1rem 0 0 1.5rem
  line-height 1.2
</style>
