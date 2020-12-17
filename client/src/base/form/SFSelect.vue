<!--
SFSelect is a validated select drop-down control.
-->

<template lang="pug">
FormItem(:id='id', xclass='oneline', :label='label', :help='help', :error='error')
  SSelect.w100(
    :id='id',
    :class='classes',
    :options='options',
    :valueKey='valueKey',
    :labelKey='labelKey',
    v-bind='$attrs',
    v-model='value',
    @focus='onFocus',
    @blur='onBlur'
  )
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, watchEffect } from 'vue'
import { propagateModel } from '../util'
import SSelect from '../controls/SSelect'
import FormItem, { ErrorFunction, useLostFocus } from './FormItem'

export default defineComponent({
  components: { FormItem, SSelect },
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: [String, Number], required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    errorFn: Function as PropType<ErrorFunction>,
  },
  emits: ['update:modelValue'],
  setup(props, { attrs, emit }) {
    const value = propagateModel(props, emit)
    // Set up for form validation.
    const { lostFocus, submitted, onFocus, onBlur } = useLostFocus()
    const error = ref('')
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value!, submitted?.value || false)
      })
    const classes = computed(() => error.value ? 'form-control-invalid' : null)
    return { classes, error, onBlur, onFocus, value }
  }
})
</script>
