<!--
PrivilegeMaskButton displays a button inside a PrivilegeMask.
-->

<template lang="pug">
button.privileges-btn(:class='classes', v-text='label', @click.prevent='onClick')
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'

export default defineComponent({
  props: {
    label: { type: String, required: true },
    modelValue: { type: Boolean, required: true },
    variant: String,
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const classes = computed(() => {
      const classes = []
      if (props.modelValue) classes.push('privileges-btn-selected')
      if (props.variant) classes.push(`privileges-btn-${props.variant}`)
      return classes
    })
    function onClick() {
      emit('update:modelValue', !props.modelValue)
    }
    return { classes, onClick }
  },
})
</script>

<style lang="postcss">
.privileges-btn {
  display: inline-block;
  font-weight: 400;
  color: #212529;
  text-align: center;
  vertical-align: middle;
  cursor: pointer;
  user-select: none;
  background-color: transparent;
  border: 1px solid #007bff;
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
  line-height: 1.5;
  border-radius: 0.2rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out,
    border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}
.privileges-btn-selected {
  color: #fff;
  background-color: #007bff;
}
.privileges-btn-left {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}
.privileges-btn-inner {
  border-left: none;
  border-radius: 0;
}
.privileges-btn-right {
  border-left: none;
}
</style>
