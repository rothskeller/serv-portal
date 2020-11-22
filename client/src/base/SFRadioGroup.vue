<!--
SFRadioGroup is a group of radio buttons in an SForm.  The v-model value is the
value of the selected button.  The options prop contains an array listing the
radio buttons.  Each option is an object with a value property and a label
property, and possibly an enabled or disabled flag.
-->

<template lang="pug">
label.form-item-label(v-text='label')
.form-item-input
  div(v-for='(o, i) in options', :class='classes')
    input.sfradiogroup-radio-input(
      :id='`${id}-${i}`',
      type='radio',
      :value='o[valueKey]',
      :checked='isChecked(o)',
      :disabled='isDisabled(o)',
      @change='onChange($event, o)'
    )
    label.sfradiogroup-radio-label(:for='`${id}-${i}`', v-text='o[labelKey]')
.form-item-help.sfradiogroup-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, toRefs, watchEffect, PropType, computed } from 'vue'
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
    modelValue: { type: String, required: true },
    inline: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
  },
  setup(props, { emit }) {
    // Use the incoming modelValue as the initial value for selected, and update
    // selected whenever the parent changes the modelValue prop.
    const { modelValue } = toRefs(props)
    const selected = ref('')
    watch(
      modelValue,
      () => {
        selected.value = modelValue.value
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
      return selected.value === option[props.valueKey]
    }

    // Return whether a particular option is disabled.
    function isDisabled(option: any) {
      if (props.enabledKey) return !option[props.enabledKey]
      if (props.disabledKey) return !!option[props.disabledKey]
      return false
    }

    // When a button state changes, update our local copy and notify owner.
    function onChange({ target: { checked: nc } }: { target: HTMLInputElement }, option: any) {
      if (nc) selected.value = option[props.valueKey]
      else selected.value = ''
      emit('update:modelValue', selected.value)
    }

    // Display classes.
    const classes = computed(() => [
      props.inline ? 'sfradiogroup-hradio' : 'sfradiogroup-vradio',
      error.value ? 'sfradiogroup-invalid' : null,
    ])

    return { classes, error, isChecked, isDisabled, onChange }
  },
})
</script>

<style lang="postcss">
.sfradiogroup-vradio {
  position: relative;
  display: block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
}
.sfradiogroup-hradio {
  position: relative;
  display: inline-block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
  margin-right: 1rem;
}
.sfradiogroup-radio-input {
  position: absolute;
  left: 0;
  z-index: -1;
  width: 1rem;
  height: 1.25rem;
  opacity: 0;
  padding: 0;
}
.sfradiogroup-radio-label {
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
    border-radius: 50%;
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
.sfradiogroup-radio-input:not(:disabled):active ~ .sfradiogroup-radio-label:before {
  color: #fff;
  background-color: #b3d7ff;
  border-color: #b3d7ff;
}
.sfradiogroup-radio-input:focus:not(:checked) ~ .sfradiogroup-radio-label:before {
  border-color: #80bdff;
}
.sfradiogroup-radio-input:focus ~ .sfradiogroup-radio-label:before {
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}
.sfradiogroup-radio-input:checked ~ .sfradiogroup-radio-label:before {
  color: #fff;
  border-color: #007bff;
  background-color: #007bff;
}
.sfradiogroup-radio-input:checked ~ .sfradiogroup-radio-label:after {
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='-4 -4 8 8'%3E%3Ccircle r='3' fill='%23fff'/%3E%3C/svg%3E");
}
.sfradiogroup-invalid .sfradiogroup-radio-label:before {
  border-color: #dc3545;
}
.form-l .sfradiogroup-helpbox {
  display: block;
}
</style>
