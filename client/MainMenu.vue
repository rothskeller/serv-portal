<!--
MainMenu displays the main menu and routes to the page being displayed within
it.
-->

<template lang="pug">
#page-top(:class="[$store.state.touch ? 'touch' : 'mouse', menuOpen ? 'page-menu-open' : null]")
  #page-heading
    #page-menu-trigger-box(@click="onMenuTrigger")
      svg#page-menu-trigger.bi.bi-list(viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg")
        path(fill-rule="evenodd" d="M4.5 13.5A.5.5 0 015 13h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 9h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 5h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5z" clip-rule="evenodd")
    #page-titlebox
      #page-title(v-text="$store.state.page.title")
    #page-menu-spacer
  #page-menu
    #page-menu-welcome
      | Welcome
      br
      b(v-text="$store.state.me.informalName")
    b-nav#page-nav(pills vertical)
      b-nav-item(to="/events" :active="menuItem === 'events'") Events
      b-nav-item(to="/people" :active="menuItem === 'people'") People
      b-nav-item(to="/reports" :active="menuItem === 'reports'") Reports
      b-nav-item(v-if="$store.state.me.canSendTextMessages" to="/texts" :active="menuItem === 'texts'") Texts
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
sidebarWidth = 7em
#page-top
  display grid
  height 100vh
  grid 'header header' titlebarHeight 'menu content' 1fr / sidebarWidth 1fr
  @media print
    display block
    height auto
#page-heading
  display flex
  justify-content space-between
  align-items center
  background-color #006600
  color #fff
  grid-area header
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
  flex auto
  flex-direction column
  text-align center
#page-title
  overflow hidden
  width 100%
  text-overflow ellipsis
  white-space nowrap
  font-size 1.5rem
#page-menu
  display none
  overflow visible
  border-right 1px solid #888
  background-color #ccc
  grid-area menu
  .page-menu-open &
    z-index 1
    display block
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
  overflow auto
  padding 1.5rem 0.75rem
  grid-area 2 / 1 / 3 / 3
  @media (min-width: 576px)
    grid-area content
  &.page-no-padding
    padding 0
#page-subtitle
  font-size 1.5rem
</style>
