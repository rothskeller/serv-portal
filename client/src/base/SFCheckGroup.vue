<!--
SFCheckGroup is a group of checkbox controls in an SForm.  The v-model value is
a set containing the values associated with the boxes that are selected.  The
options prop contains an array listing the checkboxes.  Each option is an object
with a value property and a label property, and possibly an enabled or disabled
flag.
-->

<template lang="pug">
label.form-item-label(v-text='label')
.form-item-input
  .sfcheckgroup-check(v-for='(o, i) in options', :class='{ "sfcheckgroup-invalid": error }')
    input.sfcheckgroup-check-input(
      :id='`${id}-${i}`',
      type='checkbox',
      autocomplete='off',
      :value='o[valueKey]',
      :checked='isChecked(o)',
      :disabled='isDisabled(o)',
      @change='onChange($event, o)'
    )
    label.sfcheckgroup-check-label(
      :for='`${id}-${i}`',
      :class='{ "sfcheckgroup-check-label-disabled": isDisabled(o) }',
      v-text='o[labelKey]'
    )
.form-item-help.sfcheckgroup-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, toRefs, watchEffect, PropType } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'

export type { ErrorFunction }

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    enabledKey: { type: String },
    disabledKey: { type: String },
    modelValue: { type: Set, required: true },
    errorFn: Function as PropType<ErrorFunction>,
  },
  setup(props, { emit }) {
    // Use the incoming modelValue as the initial value for checked, and update
    // checked whenever the parent changes the modelValue prop.
    const { modelValue } = toRefs(props)
    const checked = ref(new Set())
    watch(
      modelValue,
      () => {
        if (modelValue.value === checked.value) return
        checked.value.clear()
        modelValue.value.forEach((v) => checked.value.add(v))
      },
      { immediate: true }
    )

    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    // Return whether a particular option is checked.
    function isChecked(option: any) {
      return checked.value.has(option[props.valueKey])
    }

    // Return whether a particular option is disabled.
    function isDisabled(option: any) {
      if (props.enabledKey) return !option[props.enabledKey]
      if (props.disabledKey) return !!option[props.disabledKey]
      return false
    }

    // When a button state changes, update our local copy and notify owner.
    function onChange({ target: { checked: nc } }: { target: HTMLInputElement }, option: any) {
      if (nc) checked.value.add(option[props.valueKey])
      else checked.value.delete(option[props.valueKey])
      console.log('onChange', checked.value)
      emit('update:modelValue', checked.value)
    }

    return { isChecked, isDisabled, onChange, error }
  },
})
</script>

<style lang="postcss">
.sfcheckgroup-check {
  position: relative;
  display: block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
}
.sfcheckgroup-check-input {
  position: absolute;
  left: 0;
  z-index: -1;
  width: 1rem;
  height: 1.25rem;
  opacity: 0;
  padding: 0;
}
.sfcheckgroup-check-label {
  display: inline-block;
  position: relative;
  margin-bottom: 0;
  vertical-align: top;
  &:before {
    pointer-events: none;
    background-color: #fff;
    border: 1px solid #adb5bd;
    position: absolute;
    top: 0.25rem;
    left: -1.5rem;
    display: block;
    width: 1rem;
    height: 1rem;
    content: '';
    border-radius: 0.25rem;
    transition: background-color 0.15s ease-in-out, border-color 0.15s ease-in-out,
      box-shadow 0.15s ease-in-out;
  }
  &:after {
    position: absolute;
    top: 0.25rem;
    left: -1.5rem;
    display: block;
    width: 1rem;
    height: 1rem;
    content: '';
    background: no-repeat 50%/50% 50%;
  }
}
.sfcheckgroup-check-label-disabled {
  color: #888;
}
.sfcheckgroup-check-input:not(:disabled):active ~ .sfcheckgroup-check-label:before {
  color: #fff;
  background-color: #b3d7ff;
  border-color: #b3d7ff;
}
.sfcheckgroup-check-input:focus:not(:checked) ~ .sfcheckgroup-check-label:before {
  border-color: #80bdff;
}
.sfcheckgroup-check-input:focus ~ .sfcheckgroup-check-label:before {
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}
.sfcheckgroup-check-input:checked ~ .sfcheckgroup-check-label:before {
  color: #fff;
  border-color: #007bff;
  background-color: #007bff;
}
.sfcheckgroup-check-input:checked ~ .sfcheckgroup-check-label:after {
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='8' height='8'%3E%3Cpath fill='%23fff' d='M6.564.75l-3.59 3.612-1.538-1.55L0 4.26l2.974 2.99L8 2.193z'/%3E%3C/svg%3E");
}
.sfcheckgroup-invalid .sfcheckgroup-check-label:before {
  border-color: #dc3545;
}
.form-l .sfcheckgroup-helpbox {
  display: block;
}
</style>
