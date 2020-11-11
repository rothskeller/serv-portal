<!--
SFSelect is a select control in an SForm.
-->

<template lang="pug">
label.form-item-label.sfselect-label(:for='id', v-text='label')
select.form-item-input.form-control(
  :id='id',
  :class='{ "form-control-invalid": error }',
  v-bind='$attrs',
  v-model='input',
  @focus='onFocus',
  @blur='onBlur'
)
  option(
    v-for='o in options',
    :value='optionValue(o)',
    v-text='optionLabel(o)',
    :selected='optionValue(o) === input'
  )
.form-item-help.sfselect-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  toRefs,
  watch,
  PropType,
  Ref,
  computed,
  ComputedRef,
  watchEffect,
} from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'

export type { ErrorFunction }

export default defineComponent({
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

    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    // Watch for local changes and send them to the parent.
    watch(input, () => {
      emit('update:modelValue', input.value)
    })

    return { input, onFocus, onBlur, error, optionValue, optionLabel }
  },
})
</script>

<style lang="postcss">
.sfselect-label {
  padding-top: calc(0.375rem + 1px);
}
.form-l .sfselect-helpbox {
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
