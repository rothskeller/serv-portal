<!--
SSelect is a select control, not in an SForm.
-->

<template lang="pug">
select.sselect.form-control(v-bind='$attrs', v-model='input')
  option(
    v-for='o in options',
    :value='optionValue(o)',
    v-text='optionLabel(o)',
    :selected='optionValue(o) === input'
  )
</template>

<script lang="ts">
import { defineComponent, ref, toRefs, watch } from 'vue'
import './sfcontrol.css'

export default defineComponent({
  props: {
    modelValue: { type: [String, Number], required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
  },
  setup(props, { emit }) {
    // Utility functions.
    function optionValue(o: any) {
      return typeof o === 'object' ? o[props.valueKey] : o
    }
    function optionLabel(o: any) {
      return typeof o === 'object' ? o[props.labelKey] : o.toString()
    }

    // Get the initial value from the props.  Update our local value whenever
    // the props change.
    const { modelValue } = toRefs(props)
    const input = ref(modelValue.value as string | number)
    if (!props.options.find((o) => input.value === optionValue(o))) {
      console.warn('Initial value for select is not one of the allowed options.')
      input.value = optionValue(props.options[0])
    }
    watch(modelValue, () => {
      if (props.options.find((o) => modelValue.value === optionValue(o)))
        input.value = modelValue.value
      else console.warn('Updated value for select is not one of the allowed options.')
    })

    // Watch for local changes and send them to the parent.
    watch(input, () => {
      emit('update:modelValue', input.value)
    })

    return { input, optionValue, optionLabel }
  },
})
</script>
