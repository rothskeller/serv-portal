<!--
SFInput is an input control in an SForm.
-->

<template lang="pug">
label.form-item-label.sfinput-label(:id='id + "_label"', :for='id', v-text='label')
input.form-item-input.form-control(
  ref='inputRef',
  :id='id',
  :class='{ "form-control-invalid": error }',
  v-bind='$attrs',
  :value='input',
  @input='onInput',
  @focus='onFocus',
  @blur='onBlur'
)
.form-item-help.sfinput-helpbox
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
  watchEffect,
  getCurrentInstance,
  Prop,
} from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'

export type { ErrorFunction }

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
    restrictFn: Function as PropType<(s: string) => string>,
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

    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    // Apply trim when losing focus.  If we apply it more eagerly, backspacing
    // over the start of a word also removes the space before it, which is
    // disturbing.
    function myOnBlur() {
      const trimmed = props.trim ? input.value.trim() : input.value
      if (input.value !== trimmed) {
        input.value = trimmed
        emit('update:modelValue', input.value)
      }
      onBlur() // the one from provideValidation
    }

    // Watch for local changes and send them to the parent.
    const inst = getCurrentInstance()
    function onInput({ target: { value } }: { target: HTMLInputElement }) {
      if (props.restrictFn) {
        input.value = props.restrictFn(value)
        if (input.value !== value) inst?.update()
      } else {
        input.value = value
      }
      emit('update:modelValue', input.value)
    }

    return { inputRef, input, onInput, onFocus, onBlur: myOnBlur, error, focus }
  },
})
</script>

<style lang="postcss">
.sfinput-label {
  padding-top: calc(0.375rem + 1px);
}
.form-l .sfinput-helpbox {
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
