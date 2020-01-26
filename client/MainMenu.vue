<!--
MainMenu displays the main menu and routes to the page being displayed within
it.
-->

<template lang="pug">
#page-top(:class="$store.state.touch ? 'touch' : 'mouse'")
  #page-heading
    #page-menu-trigger-box(@click="onMenuTrigger")
      svg#page-menu-trigger.bi.bi-list(viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg")
        path(fill-rule="evenodd" d="M4.5 13.5A.5.5 0 015 13h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 9h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 5h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5z" clip-rule="evenodd")
    #page-titlebox
      #page-title(v-text="$store.state.page.title")
    #page-menu-spacer
  #page-main(:class="menuOpen ? 'page-menu-open' : null")
    #page-menu
      #page-menu-welcome
        | Welcome
        br
        b(v-text="$store.state.me.informalName")
      b-nav#page-nav(pills vertical)
        b-nav-item(to="/events" :active="menuItem === 'events'") Events
        b-nav-item(to="/people" :active="menuItem === 'people'") People
        b-nav-item(to="/reports" :active="menuItem === 'reports'") Reports
        b-nav-item(v-if="$store.state.me.webmaster" to="/texts" :active="menuItem === 'texts'") Texts
        //-b-nav-item(:to="`/people/${$store.state.me.id}`") Profile
        b-nav-item(to="/logout") Logout
    #page-content(:class="tabbed ? 'page-no-padding': null")
      #page-subtitle(v-if="$store.state.page.subtitle" v-text="$store.state.page.subtitle")
      router-view
</template>

<script>
export default {
  data: () => ({ menuOpen: false }),
  computed: {
    menuItem() {
      const record = this.$route.matched.find(rec => rec.meta.menuItem)
      return record ? record.meta.menuItem : null
    },
    tabbed() { return this.$route.matched.some(rec => rec.meta.tabbed) },
  },
  methods: {
    onMenuTrigger() { this.menuOpen = !this.menuOpen },
  },
}
</script>

<style lang="stylus">
titlebarHeight = 40px
#page-top
  display flex
  flex-direction column
  height 100vh
  @media print
    height auto
#page-heading
  z-index 1
  display flex
  flex none
  justify-content space-between
  align-items center
  height titlebarHeight
  background-color #006600
  color #fff
  @media print
    display none
#page-menu-trigger-box, #page-menu-spacer
  flex none
  margin 0 6px
  width 2rem
  cursor pointer
  user-select none
#page-menu-trigger
  @media (min-width: 576px)
    display none
#page-titlebox
  display flex
  flex none
  flex-direction column
  width calc(100vw - 5.5rem)
  text-align center
#page-title
  overflow hidden
  width 100%
  text-overflow ellipsis
  white-space nowrap
  font-size 1.5rem
#page-main
  position relative
  display flex
  flex none
  overflow-y auto
  width 100vw
  height 'calc(100vh - %s)' % titlebarHeight
  @media print
    height auto
#page-menu
  display none
  flex none
  overflow visible
  width 7rem
  border-right 1px solid #888
  background-color #ccc
  .page-menu-open &
    position absolute
    z-index 1
    display block
    height 'calc(100vh - %s)' % titlebarHeight
  @media (min-width: 576px)
    display block
  @media print
    display none
#page-menu-welcome
  margin-bottom 0.5rem
  padding 0.75rem
  border-bottom 1px solid white
  text-align center
  white-space nowrap
  font-size 0.75rem
#page-nav
  padding 0 0.5rem
  font-size 1.25rem
  .nav-link
    padding 0.125rem 0.5rem
    color black
    &.active
      color white
#page-content
  flex auto
  overflow auto
  padding 1.5rem 0.75rem
  &.page-no-padding
    padding 0
#page-subtitle
  font-size 1.5rem
</style>
