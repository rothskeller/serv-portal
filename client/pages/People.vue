<!--
People displays the list of people.
-->

<template lang="pug">
Page(title="People" menuItem="people" noPadding)
  b-card#people-card(no-body)
    b-tabs(card)
      b-tab.people-tab-pane(title="List" no-body)
        PeopleList
      b-tab.person-tab-pane(title="Map" no-body lazy)
        PeopleMap
</template>

<script>
const PeopleMap = () => import(/* webpackChunkName: "maps" */ './people/-PeopleMap.vue')

export default {
  components: { PeopleMap },
}
</script>

<style lang="stylus">
#people-card
  height calc(100vh - 40px)
  border none
.people-tab-pane
  overflow-y auto
  height calc(100vh - 3.25rem - 42px)
</style>
