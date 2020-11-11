<!--
SButton displays a styled button.
-->

<template lang="pug">
component.sbtn(
  :is='tag',
  :class='classes',
  :to='to',
  :type='!to ? type : null',
  :disabled='disabled'
)
  slot
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'

export default defineComponent({
  props: {
    disabled: { type: Boolean, default: false },
    variant: { type: String, default: 'secondary' },
    type: { type: String, default: 'button' },
    to: [String, Object],
  },
  setup(props) {
    const tag = computed(() => (props.to ? 'router-link' : 'button'))
    const classes = computed(() => [`sbtn-${props.variant}`, { 'sbtn-disabled': props.disabled }])
    return { tag, classes }
  },
})
</script>

<style lang="postcss">
.sbtn {
  display: inline-block;
  font-weight: 400;
  color: #212529;
  text-align: center;
  vertical-align: middle;
  cursor: pointer;
  user-select: none;
  background-color: transparent;
  border: 1px solid transparent;
  padding: 0.375rem 0.75rem;
  font-size: 1rem;
  line-height: 1.5;
  border-radius: 0.25rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out,
    border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}
.sbtn-primary {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}
.sbtn-secondary {
  color: #fff;
  background-color: #6c757d;
  border-color: #6c757d;
}
.sbtn-danger {
  color: #fff;
  background-color: #dc3545;
  border-color: #dc3545;
}
.sbtn-warning {
  color: #fff;
  background-color: #ffc107;
  border-color: #ffc107;
}
.sbtn-disabled {
  opacity: 0.65;
}
</style>
