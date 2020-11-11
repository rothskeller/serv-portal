<!--
SProgress displays a progress bar.
-->

<template lang="pug">
.sprogress
  .sprogressbar(
    :style='{ width: `${100 * value}%` }',
    role='progressbar',
    aria-valuemin='0',
    aria-valuemax='1',
    :aria-valuenow='value.toFixed(2)'
  )
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    value: { type: Number, required: true },
  },
})
</script>

<style lang="postcss">
.sprogress {
  display: flex;
  overflow: hidden;
  height: 1rem;
  background-color: #e9ecef;
  border-radius: 0.25rem;
}
.sprogressbar {
  display: flex;
  overflow: hidden;
  background-color: #007bff;
  transition: width 0.6s ease;
}
</style>
