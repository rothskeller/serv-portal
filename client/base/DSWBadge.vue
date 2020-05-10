<!--
DSWBadge displays a DSW badge.
-->

<template lang="pug">
span(v-if="badge == 'valid'")
  |
  span.dsw-valid DSW
span(v-else-if="badge == 'invalid'")
  |
  span.dsw-invalid No DSW
</template>

<script>
export default {
  props: {
    badge: String,
  },
}
</script>

<style lang="stylus">
.dsw-valid
  align-self center
  padding 0 0.25rem
  border-radius 4px
  background-color green
  color white
  white-space nowrap
  font-size 0.75rem
.dsw-invalid
  align-self center
  padding 0 0.25rem
  border-radius 4px
  background-color red
  color white
  white-space nowrap
  font-size 0.75rem
</style>
