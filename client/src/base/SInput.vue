<!--
SInput is a text input control.
-->

<template lang="pug">
input.form-control(ref='inputRef', v-bind='$attrs', v-model='input', @blur='onBlur')
</template>

<script lang="ts">
import { defineComponent, getCurrentInstance, PropType, ref, watch } from 'vue'
import './sfcontrol.css'

export default defineComponent({
  props: {
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    restrictFn: Function as PropType<(s: string) => string>,
  },
  emits: ['blur', 'update:modelValue'],
  setup(props, { emit }) {
    const instance = getCurrentInstance()

    // Reference to the input field, so we can pass focus() calls to it.
    const inputRef = ref<HTMLInputElement>()
    function focus() {
      inputRef.value?.focus()
    }

    // Get the initial value from the props.  Update our local value whenever
    // the props change.
    const input = ref(props.modelValue)
    watch(
      () => props.modelValue,
      () => {
        input.value = props.modelValue
      }
    )

    // Apply trim when losing focus.  (If we apply it more eagerly, backspacing
    // over the start of a word also removes the space before it, which is
    // disturbing.)
    function onBlur() {
      const trimmed = props.trim ? input.value.trim() : input.value
      if (input.value !== trimmed) {
        input.value = trimmed
        emit('update:modelValue', input.value)
      }
      emit('blur')
    }

    // Watch for local changes and send them to the parent.
    watch(input, () => {
      let nv = props.restrictFn ? props.restrictFn(input.value) : input.value
      if (nv !== input.value) {
        input.value = nv
        // Forcible update required to rerender control if value changed.
        instance?.update()
      }
      emit('update:modelValue', input.value)
    })

    return { focus, input, inputRef, onBlur }
  },
})
</script>
