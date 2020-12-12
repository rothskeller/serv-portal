<!--
SFSelect is a select control in an SForm.
-->

<template lang="pug">
label.form-item-label.sfselect-label(:for='id', v-text='label')
SSelect.form-item-input(
  :id='id',
  :options='options',
  :valueKey='valueKey',
  :labelKey='labelKey',
  :class='{ "form-control-invalid": error }',
  v-bind='$attrs',
  v-model='input',
  @focus='onFocus',
  @blur='onBlur'
)
.form-item-help.sfselect-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, PropType, watchEffect } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'
import SSelect from './SSelect.vue'

export type { ErrorFunction }

export default defineComponent({
  components: { SSelect },
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
    // Get the initial value from the props.  Update our local value whenever
    // the props change.
    const input = ref(props.modelValue)
    watch(
      () => props.modelValue,
      () => {
        input.value = props.modelValue
      }
    )

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

    return { error, input, onBlur, onFocus }
  },
})
</script>

<style lang="postcss">
.sfselect-label {
  /* Align baseline of label with baseline in control. */
  padding-top: calc(0.375rem + 1px);
}
.form-l .sfselect-helpbox {
  /* Ensure help box is at least as tall as control so vertical centering works. */
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
