<!--
ClearanceBadge displays a clearance badge.
-->

<template lang="pug">
span.clearance(v-if="req === true && (!v || !d || !b)")
  | {{ ' ' }}
  span.clear-v(:class="v ? 'clear-g' : 'clear-r'") V
  span.clear-d(:class="d ? 'clear-g' : 'clear-r'") D
  span.clear-b(:class="b ? 'clear-g' : 'clear-r'") B
span.clearance(v-else-if="req === false && (v || d || b)")
  | {{ ' ' }}
  span.clear-v(:class="v ? 'clear-g' : 'clear-x'") V
  span.clear-d(:class="d ? 'clear-g' : 'clear-x'") D
  span.clear-b(:class="b ? 'clear-g' : 'clear-x'") B
</template>

<script>
export default {
  props: {
    req: Boolean,
    v: Boolean,
    d: Boolean,
    b: Boolean,
  },
}
</script>

<style lang="stylus">
.clearance
  align-self center
  font-size 0.75rem
.clear-v, .clear-d, .clear-b
  padding 0 0.25rem
  color white
.clear-v
  border-top-left-radius 4px
  border-bottom-left-radius 4px
.clear-b
  border-top-right-radius 4px
  border-bottom-right-radius 4px
.clear-g
  background-color green
.clear-r
  background-color red
.clear-x
  background-color gray
</style>
