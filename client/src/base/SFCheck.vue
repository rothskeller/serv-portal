<!--
SFCheck is a checkbox control in an SForm.  It is initially checked if the value
passed to it via v-model is truthy.  When checked by the user, the v-model value
is set to the value of the 'value' prop.  When unchecked by the user, the
v-model value is set to Boolean false.
-->

<template lang="pug">
.sfcheck.form-item-input(:class='{ "sfcheck-invalid": error }')
  input.sfcheck-input(
    :id='id',
    type='checkbox',
    autocomplete='off',
    :checked='checked',
    @change='onChange',
    @focus='onFocus',
    @blur='onBlur'
  )
  label.sfcheck-label(:for='id', v-text='label')
.form-item-help.sfcheck-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, toRefs, PropType, watchEffect } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    value: { type: [String, Number, Boolean, Object], default: true },
    modelValue: { type: [String, Number, Boolean, Object], default: null },
    errorFn: Function as PropType<ErrorFunction>,
  },
  setup(props, { emit }) {
    // Use the incoming modelValue as the initial value for checked, and update
    // checked whenever the parent changes the modelValue prop.
    const { modelValue } = toRefs(props)
    const checked = ref(!!modelValue.value)
    watch(modelValue, () => {
      checked.value = !!modelValue.value
    })

    // When the button state changes, update our local copy and notify owner.
    function onChange({ target: { checked: nc } }: { target: HTMLInputElement }) {
      checked.value = nc ? props.value : false
      emit('update:modelValue', checked.value)
    }

    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    return { checked, onFocus, onBlur, error, onChange }
  },
})
</script>

<style lang="postcss">
.sfcheck {
  position: relative;
  display: block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
}
.sfcheck-input {
  position: absolute;
  left: 0;
  z-index: -1;
  width: 1rem;
  height: 1.25rem;
  opacity: 0;
  padding: 0;
}
.sfcheck-label {
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
.sfcheck-input:not(:disabled):active ~ .sfcheck-label:before {
  color: #fff;
  background-color: #b3d7ff;
  border-color: #b3d7ff;
}
.sfcheck-input:focus:not(:checked) ~ .sfcheck-label:before {
  border-color: #80bdff;
}
.sfcheck-input:focus ~ .sfcheck-label:before {
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}
.sfcheck-input:checked ~ .sfcheck-label:before {
  color: #fff;
  border-color: #007bff;
  background-color: #007bff;
}
.sfcheck-invalid > .sfcheck-label:before {
  border-color: #dc3545;
}
.sfcheck-input:checked ~ .sfcheck-label:after {
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='8' height='8'%3E%3Cpath fill='%23fff' d='M6.564.75l-3.59 3.612-1.538-1.55L0 4.26l2.974 2.99L8 2.193z'/%3E%3C/svg%3E");
}
.form-l .sfcheck-helpbox {
  min-height: 1.5rem;
}
</style>
