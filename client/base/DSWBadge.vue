<!--
DSWBadge displays a DSW badge.
-->

<template lang="pug">
span.dsw-valid(v-if="badge == 'valid'") DSW
span.dsw-invalid(v-else-if="badge == 'invalid'") No DSW
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
  border-radius 4px
  background-color green
  color white
  font-size 0.75rem
.dsw-invalid
  border-radius 4px
  background-color red
  color white
  font-size 0.75rem
</style>
