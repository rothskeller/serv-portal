<!--
EventTypeDots displays a set of appropriately colored dots for a set of event
types.
-->

<template lang="pug">
span
  span.dot(v-for="type in types" :class="`dot-${typeToClass(type)}`" :title="type")
</template>

<script>
export default {
  props: {
    types: Array,
  },
  methods: {
    typeToClass(type) {
      return type.toLowerCase().replace(/ /g, '-')
    },
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
.dot-incident
  background-color #e6194b // red
.dot-civic-event
  background-color #f032e6 // magenta
.dot-drill
  background-color #f58231 // orange
.dot-continuing-education
  border 2px solid #4363d8 // blue circle
.dot-training
  background-color #4363d8 // blue
.dot-class
  background-color #3cb44b // green
.dot-work-event
  border 2px solid #a9a9a9 // gray circle
.dot-meeting
  background-color #a9a9a9 // gray
.dot-party
  background-color #e6beff // lavender
  // colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/
</style>
