<!--
People displays the list of people.
-->

<template lang="pug">
Page(title="Text Messages" menuItem="texts" noPadding)
  b-card#texts-card(no-body)
    b-tabs(card)
      b-tab.texts-tab-pane(title="Sent Messages" no-body)
        TextsList
      b-tab.person-tab-pane(title="New Message" no-body)
        | World
</template>

<script>
export default {
}
</script>

<style lang="stylus">
#texts-card
  height calc(100vh - 40px)
  border none
.texts-tab-pane
  overflow-y auto
  height calc(100vh - 3.25rem - 42px)
</style>
