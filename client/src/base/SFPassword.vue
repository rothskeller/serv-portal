<!--
SFPassword is a control for getting a new password in an SForm.  It has two
input fields, for cross check, and a meter for password quality.
-->

<template lang="pug">
label.form-item-label.sfpassword-label(:id='id + "_label"', :for='id', v-text='label')
.form-item-input
  input.form-control.sfpassword-input(
    ref='input1Ref',
    :id='id',
    :class='{ "form-control-invalid": error }',
    type='password',
    v-model='input1'
  )
  input.form-control.sfpassword-input(
    :id='id + "2"',
    :class='{ "form-control-invalid": error }',
    type='password',
    v-model='input2'
  )
.form-item-help.sfpassword-helpbox
  .sfpassword-meter(v-if='score >= 0')
    .sfpassword-meter-step(v-for='n in score + 1', :class='meterColor')
  .form-item-error-text(v-if='error', v-text='error')
  .sfpassword-success(v-if='success', v-text='success')
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
import provideValidation from './sfvalidate'
import type { ZXCVBNResult } from 'zxcvbn'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    modelValue: { type: String, required: true },
    required: { type: Boolean, default: false },
    allowBadPassword: { type: Boolean, default: false },
    passwordHints: Array as PropType<Array<string>>,
  },
  setup(props, { emit }) {
    // Get the initial value from the props.
    const { modelValue } = toRefs(props)
    const input1 = ref(modelValue.value as string)
    const input2 = ref(modelValue.value as string)

    // Set up for form control and password quality validation.
    const zxcvbn = import('zxcvbn')
    const score = ref(-1)
    const meterColor = ref('')
    const error = ref('')
    const success = ref('')
    const { submitted } = provideValidation(props.id, error)
    function crackMessage(result: ZXCVBNResult) {
      return `This password would take ${result.crack_times_display.offline_slow_hashing_1e4_per_second} to crack.`
    }
    async function validate() {
      if ((submitted.value || input2.value) && input1.value !== input2.value) {
        score.value = -1
        error.value = 'These two password entries do not match.'
        success.value = ''
        return
      }
      if (!input1.value) {
        score.value = -1
        error.value =
          submitted.value && props.required ? 'You must enter a new password, twice.' : ''
        success.value = ''
        return
      }
      const result = (await zxcvbn).default(input1.value, props.passwordHints)
      score.value = result.score
      meterColor.value = [
        'sfpassword-meter-bad',
        'sfpassword-meter-bad',
        'sfpassword-meter-warn',
        'sfpassword-meter-good',
        'sfpassword-meter-good',
      ][score.value]
      error.value = success.value = ''
      if (result.feedback) {
        const acceptable = result.score > 2 || props.allowBadPassword
        const message = [
          result.feedback.warning,
          ...result.feedback.suggestions,
          crackMessage(result),
        ]
          .filter((s) => !!s)
          .join('\n')
        if (acceptable) success.value = message
        else error.value = message
      }
    }
    watchEffect(validate)
    watch(input1, () => {
      emit('update:modelValue', input1.value)
    })

    // Pass on focus requests to first input field.
    const input1Ref = ref(null as any)
    function focus() {
      input1Ref.value.focus()
    }

    return { focus, input1, input1Ref, input2, score, meterColor, error, success }
  },
})
</script>

<style lang="postcss">
.sfpassword-label {
  padding-top: calc(0.375rem + 1px);
}
.sfpassword-input {
  display: block;
  width: 100%;
  &:first-child {
    margin-bottom: 0.375rem;
  }
}
.form-l .sfpassword-helpbox {
  min-height: calc(1.5rem + 0.75rem + 2px);
}
.sfpassword-meter {
  margin: 0.5rem 0;
  height: 0.5rem;
  width: 11rem;
  background-color: #e9ecef;
}
.sfpassword-meter-step {
  display: inline-block;
  vertical-align: top;
  width: 2rem;
  height: 0.5rem;
  margin: 0 0 0 0.25rem;
  &:first-child {
    margin-left: 0;
  }
}
.sfpassword-meter-bad {
  background-color: #dc3545;
}
.sfpassword-meter-warn {
  background-color: #ffc107;
}
.sfpassword-meter-good {
  background-color: #28a745;
}
.sfpassword-success {
  color: #28a745;
  font-size: 80%;
  line-height: 1.2;
}
</style>
