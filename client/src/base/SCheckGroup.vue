<!--
SCheckGroup is a group of checkbox controls in an SForm.  The v-model value is
a set containing the values associated with the boxes that are selected.  The
options prop contains an array listing the checkboxes.  Each option is an object
with a value property and a label property, and possibly an enabled or disabled
flag.
-->

<template lang="pug">
div(v-for='(o, i) in options', :class='classes')
  input.scheckgroup-check-input(
    :id='`${id}-${i}`',
    type='checkbox',
    :value='o[valueKey]',
    :checked='isChecked(o)',
    :disabled='isDisabled(o)',
    @change='onChange($event, o)'
  )
  label.scheckgroup-check-label(:for='`${id}-${i}`', v-text='o[labelKey]')
</template>

<script lang="ts">
import { defineComponent, ref, watch, toRefs, watchEffect, PropType, computed } from 'vue'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    enabledKey: { type: String },
    disabledKey: { type: String },
    modelValue: { type: Set, required: true },
    inline: { type: Boolean, default: false },
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

    // Display class.
    const classes = computed(() => (props.inline ? 'scheckgroup-hcheck' : 'scheckgroup-vcheck'))

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
      emit('update:modelValue', checked.value)
    }

    return { classes, isChecked, isDisabled, onChange }
  },
})
</script>

<style lang="postcss">
.scheckgroup-vcheck {
  position: relative;
  display: block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
}
.scheckgroup-hcheck {
  position: relative;
  display: inline-block;
  min-height: 1.5rem;
  padding-left: 1.5rem;
  margin-right: 1rem;
}
.scheckgroup-check-input {
  position: absolute;
  left: 0;
  z-index: -1;
  width: 1rem;
  height: 1.25rem;
  opacity: 0;
  padding: 0;
}
.scheckgroup-check-label {
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
.scheckgroup-check-input:not(:disabled):active ~ .scheckgroup-check-label:before {
  color: #fff;
  background-color: #b3d7ff;
  border-color: #b3d7ff;
}
.scheckgroup-check-input:focus:not(:checked) ~ .scheckgroup-check-label:before {
  border-color: #80bdff;
}
.scheckgroup-check-input:focus ~ .scheckgroup-check-label:before {
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}
.scheckgroup-check-input:checked ~ .scheckgroup-check-label:before {
  color: #fff;
  border-color: #007bff;
  background-color: #007bff;
}
.scheckgroup-check-input:checked ~ .scheckgroup-check-label:after {
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='8' height='8'%3E%3Cpath fill='%23fff' d='M6.564.75l-3.59 3.612-1.538-1.55L0 4.26l2.974 2.99L8 2.193z'/%3E%3C/svg%3E");
}
.scheckgroup-invalid .scheckgroup-check-label:before {
  border-color: #dc3545;
}
</style>
