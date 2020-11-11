<!--
SInput is an input control, not in an SForm.
-->

<template lang="pug">
input.form-control(ref='inputRef', v-bind='$attrs', :value='input', @input='onInput')
</template>

<script lang="ts">
import { defineComponent, ref, toRefs, watch } from 'vue'
import './sfcontrol.css'

export default defineComponent({
  props: {
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
  },
  setup(props, { emit }) {
    // Reference to the input field, so we can pass focus() calls to it.
    const inputRef = ref(null as HTMLInputElement | null)
    function focus() {
      if (inputRef.value) inputRef.value.focus()
    }

    // Get the initial value from the props.  Update our local value whenever
    // the props change.
    const { modelValue } = toRefs(props)
    const input = ref(modelValue.value as string)
    watch(modelValue, () => {
      input.value = modelValue.value
    })

    // Watch for local changes and send them to the parent.
    function onInput({ target: { value } }: { target: HTMLInputElement }) {
      input.value = props.trim ? value.trim() : value
      emit('update:modelValue', input.value)
    }

    return { inputRef, input, onInput, focus }
  },
})
</script>
