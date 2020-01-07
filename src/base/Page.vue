<!--
Page is the basic framework of all pages on the site.
-->

<template lang="pug">
#page-top
  #page-heading
    #page-menu-trigger-box(@click="onMenuTrigger")
      svg#page-menu-trigger.bi.bi-list(v-if="$store.state.me" viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg")
        path(fill-rule="evenodd" d="M4.5 13.5A.5.5 0 015 13h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 9h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 5h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5z" clip-rule="evenodd")
    #page-titlebox
      #page-title(v-text="title")
    #page-menu-spacer
  #page-main(:class="menuOpen ? 'page-menu-open' : null")
    #page-menu(v-if="$store.state.me")
      #page-menu-welcome
        | Welcome
        br
        b(v-text="`${$store.state.me.firstName} ${$store.state.me.lastName}`")
      router-link.page-menu-item(:class="{'page-menu-item-active': menuItem === 'events'}" to="/events") Events
      router-link.page-menu-item(:class="{'page-menu-item-active': menuItem === 'people'}" to="/people") People
      router-link.page-menu-item(v-if="$store.state.me.webmaster" :class="{'page-menu-item-active': menuItem === 'roles'}" to="/roles") Roles
      router-link.page-menu-item(:class="{'page-menu-item-active': menuItem === 'reports'}" to="/reports") Reports
      router-link.page-menu-item(:to="`/people/${$store.state.me.id}`") Profile
      router-link.page-menu-item(to="/logout") Logout
    #page-content
      #page-subtitle(v-if="subtitle" v-text="subtitle")
      slot
</template>

<script>
export default {
  props: {
    title: String,
    subtitle: String,
    menuItem: String,
  },
  data: () => ({ menuOpen: false }),
  methods: {
    onMenuTrigger() { this.menuOpen = !this.menuOpen },
  },
}
</script>

<style lang="stylus">
titlebarHeight = 40px
sidebarWidth = 6rem
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
  width sidebarWidth
  border-right 1px solid #888
  background-color #ccc
  .page-menu-open &
    position absolute
    display block
    height 'calc(100vh - %s)' % titlebarHeight
  @media (min-width: 576px)
    display block
  @media print
    display none
#page-menu-welcome
  margin-bottom 0.5rem
  padding 0.75rem
  width 'calc(%s - 1px)' % sidebarWidth
  border-bottom 1px solid white
  text-align center
  white-space nowrap
  font-size 0.75rem
.page-menu-item
  display block
  padding 0.5rem 0 0.5rem 0.75rem
  width sidebarWidth
  color black
  font-size 1.25rem
  &:hover
    color black
    text-decoration none
.page-menu-item-active
  .page-menu-open &
    color #00f
  @media (min-width: 576px)
    border-top 1px solid #888
    border-bottom 1px solid #888
    background-color white
#page-content
  flex auto
  overflow auto
  padding 1.5rem 0.75rem
#page-subtitle
  font-size 1.5rem
</style>
