<!--
SFInput is an input control in an SForm.
-->

<template lang="pug">
label.form-item-label.sfinput-label(:id='id + "_label"', :for='id', v-text='label')
SInput.form-item-input(
  ref='inputRef',
  :trim='trim',
  :restrictFn='restrictFn',
  :id='id',
  :class='{ "form-control-invalid": error }',
  v-bind='$attrs',
  v-model='input',
  @focus='onFocus',
  @blur='onBlur'
)
.form-item-help.sfinput-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, PropType, watchEffect } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'
import SInput from './SInput.vue'
export type { ErrorFunction }

export default defineComponent({
  components: { SInput },
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
    restrictFn: Function as PropType<(s: string) => string>,
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // Reference to the input field, so we can pass focus() calls to it.
    const inputRef = ref<typeof SInput>()
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

    return { error, focus, input, inputRef, onBlur, onFocus }
  },
})
</script>

<style lang="postcss">
.sfinput-label {
  /* Align label with baseline inside control. */
  padding-top: calc(0.375rem + 1px);
}
.form-l .sfinput-helpbox {
  /* Help box must not be shorter than control, so that vertical centering works. */
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
