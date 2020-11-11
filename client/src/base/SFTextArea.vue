<!--
SFTextArea is a textarea control in an SForm.
-->

<template lang="pug">
label.form-item-label.sftextarea-label(:for='id', v-text='label')
textarea.form-item-input.form-control.sftextarea(
  ref='inputRef',
  :id='id',
  :class='{ "form-control-invalid": error }',
  v-bind='$attrs',
  :value='input',
  @input='onInput',
  @focus='onFocus',
  @blur='onBlur'
)
.form-item-help.sftextarea-helpbox
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
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
  },
  setup(props, { emit }) {
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

    // Watch for local changes and send them to the parent.
    function onInput({ target: { value } }: { target: HTMLTextAreaElement }) {
      input.value = props.trim ? value.trim() : value
      emit('update:modelValue', input.value)
    }

    return { input, onInput, onFocus, onBlur, error, focus }
  },
})
</script>

<style lang="postcss">
.sftextarea-label {
  padding-top: calc(0.375rem + 1px);
}
.sftextarea {
  height: calc(4.5rem + 0.75rem + 1px);
}
.form-l .sftextarea-helpbox {
  display: block;
}
</style>
