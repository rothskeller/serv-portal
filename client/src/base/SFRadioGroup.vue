<!--
SFRadioGroup is a group of radio buttons in an SForm.  The v-model value is the
value of the selected button.  The options prop contains an array listing the
radio buttons.  Each option is either a plain value, an object with a value
property and a label property, and possibly an enabled or disabled flag.
-->

<template lang="pug">
label.form-item-label(v-text='label')
.form-item-input
  SRadioGroup(
    :id='id',
    :options='options',
    :valueKey='valueKey',
    :labelKey='labelKey',
    :enabledKey='enabledKey',
    :disabledKey='disabledKey',
    :inline='inline',
    :class='classes',
    v-model='selected'
  )
.form-item-help.sfradiogroup-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, watch, watchEffect, PropType, computed } from 'vue'
import provideValidation, { ErrorFunction } from './sfvalidate'
import SRadioGroup from './SRadioGroup.vue'

export type { ErrorFunction }

export default defineComponent({
  components: { SRadioGroup },
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
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // Use the incoming modelValue as the initial value for selected, and update
    // selected whenever the parent changes the modelValue prop.
    const selected = ref(props.modelValue)
    watch(
      () => props.modelValue,
      () => {
        selected.value = props.modelValue
      }
    )

    // Set up for form control validation.
    const error = ref('')
    const { submitted, lostFocus, onFocus, onBlur } = provideValidation(props.id, error)
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value, submitted.value)
      })

    // When a button state changes, update our local copy and notify owner.
    watch(selected, () => {
      emit('update:modelValue', selected.value)
    })

    // Display an error border if needed.
    const classes = computed(() => [error.value ? 'sfradiogroup-invalid' : null])

    return { classes, error, selected }
  },
})
</script>

<style lang="postcss">
.sfradiogroup-invalid .sradiogroup-radio-label:before {
  border-color: #dc3545;
}
.form-l .sfradiogroup-helpbox {
  display: block;
}
</style>
