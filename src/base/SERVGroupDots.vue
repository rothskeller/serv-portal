<!--
SERVGroupDots displays a set of appropriately colored dots for a set of SERV
groups.
-->

<template lang="pug">
span
  span.dot(v-for="group in groups" :class="`dot-${group.toLowerCase()}`" :title="group")
</template>

<script>
export default {
  props: {
    groups: Array,
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
.dot-admin
  background-color #a9a9a9
.dot-cert-d
  background-color #469990
.dot-cert-t
  border 2px solid #469990
  background-color white
.dot-listos
  background-color #f032e6
.dot-outreach
  background-color #800000
.dot-pep
  border 2px solid #f032e6
  background-color white
.dot-sares
  background-color #f58231
.dot-snap
  background-color #4363d8
  // colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/
</style>
