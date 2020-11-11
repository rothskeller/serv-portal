<!--
SFFile is a file input control in an SForm.
-->

<template lang="pug">
label.form-item-label.sffile-label(:id='id + "_label"', :for='id', v-text='label')
input.form-item-input.form-control(
  :id='id',
  :class='{ "form-control-invalid": error }',
  type='file',
  v-bind='$attrs',
  @change='onChange',
  @focus='onFocus',
  @blur='onBlur'
)
.form-item-help.sffile-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, toRefs, watch, PropType, watchEffect } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'

export type { ErrorFunction }

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: FileList,
    errorFn: Function as PropType<ErrorFunction>,
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    // Watch for local changes and send them to the parent.
    function onChange({ target: { files } }: { target: HTMLInputElement }) {
      emit('update:modelValue', files)
    }

    return { onChange, onFocus, onBlur, error }
  },
})
</script>

<style lang="postcss">
.sffile-label {
  padding-top: calc(0.375rem + 1px);
}
.form-l .sffile-helpbox {
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
