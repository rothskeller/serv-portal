<!--
EventOrgDot displays an appropriately colored dot for an event organization.
-->

<template lang="pug">
span.dot(:class="`dot-${organization.replace(' ', '')}`" :title="organization")
</template>

<script>
export default {
  props: {
    organization: String,
  },
}
</script>

<style lang="stylus">
.dot
  display inline-block
  margin 0 0 0 2px
  width 10px
  height 10px
  border-radius 50%
  &:first-child
    margin 0
.dot-Admin, .dot-
  background-color #a9a9a9 // gray
.dot-CERTTraining
  background-color #3cb44b // green
.dot-CERTDeployment
  border 2px solid #3cb44b // green
  background-color white
.dot-LISTOS
  background-color #f58231 // orange
.dot-Outreach
  background-color #000000 // black
.dot-PEP
  background-color #e6194b // red
.dot-SARES
  background-color #ffe119 // yellow
.dot-SNAP
  background-color #4363d8 // blue
  // colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/
</style>
