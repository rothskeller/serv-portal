<!--
PersonEditAddress displays an address entry on the PersonEdit form, and handles
autocomplete and validation of addresses.
-->

<template lang="pug">
label.form-item-label(:class='`person-address-label-${type}`', v-text='label')
.form-item-input
  .person-address-same-box(v-if='type !== "Home"')
    input.person-address-same(
      :id='`person-address-${type}-same`',
      type='checkbox',
      autocomplete='off',
      v-model='sameAsHome',
      :disabled='!hasHome'
    )
    label.person-address-same-label(:for='`person-address-${type}-same`') Same as home address
  input.form-control.person-address-input(
    v-if='!sameAsHome',
    ref='line1input',
    :id='`person-address-${type}`',
    :class='{ "form-control-invalid": error }',
    v-model='line1',
    @blur='onBlur'
  )
  input.form-control.person-address-input(
    v-if='!sameAsHome',
    ref='line2input',
    :class='{ "form-control-invalid": error }',
    v-model='line2',
    @focus='onFocusLine2',
    @blur='onBlur'
  )
.form-item-help.person-address-helpbox
  .form-item-help-text(v-if='help', v-text='help')
  .form-item-error-text(v-if='error', v-text='error')
</template>

<script lang="ts">
import { defineComponent, inject, nextTick, PropType, ref, toRefs, watch } from 'vue'
import { useLostFocus } from '../../base/form/FormItem'
import type { GetPersonAddress } from './PersonView.vue'

export default defineComponent({
  props: {
    type: String as PropType<'Home' | 'Work' | 'Mail'>,
    modelValue: { type: Object as PropType<GetPersonAddress>, required: true },
    hasHome: { type: Boolean, default: false },
    help: String,
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // Get the initial value from the props.
    const { modelValue } = toRefs(props)
    const line1 = ref('')
    const line2 = ref('')
    const sameAsHome = ref(false)
    let lastChecked = ''
    let latitude = 0
    let longitude = 0
    if (props.modelValue?.address) {
      lastChecked = props.modelValue.address
      line1.value = props.modelValue.address.split(',')[0]
      line2.value = props.modelValue.address.replace(/^[^,]*, */, '')
      latitude = props.modelValue.latitude || 0
      longitude = props.modelValue.longitude || 0
    }
    if (!props.modelValue) sameAsHome.value = props.type === 'Mail'
    else sameAsHome.value = props.modelValue.sameAsHome || false

    // Update line2 automatically based on line1 entries.
    watch(line1, () => {
      if (line1.value === '') line2.value = ''
      else if (line2.value === '') line2.value = 'Sunnyvale, CA'
    })

    // Update the result when the sameAsHome checkbox changes.
    const line1input = ref(null as null | HTMLInputElement)
    watch(sameAsHome, () => {
      line1.value = line2.value = ''
      emit('update:modelValue', {
        address: '',
        latitude: 0,
        longitude: 0,
        sameAsHome: sameAsHome.value,
      })
      if (!sameAsHome.value)
        nextTick(() => {
          line1input.value?.focus()
        })
    })

    // Select the contents of line 2 when it receives focus.
    const line2input = ref(null as null | HTMLInputElement)
    function onFocusLine2() {
      line2input.value?.select()
    }

    // Set up for form control and address validation.
    const error = ref('')
    const setValidity = inject<(id: string, isValid: boolean) => void>('setValidity')
    watch(
      error,
      () => {
        setValidity?.(`person-address-${props.type}`, !error.value)
      }
    )
    const { submitted } = useLostFocus()
    async function validate() {
      if (!line1.value && !line2.value) {
        emit('update:modelValue', {
          address: '',
          latitude: 0,
          longitude: 0,
          sameAsHome: sameAsHome.value,
        })
        return
      }
      let check = `${line1.value}, ${line2.value}`
      if (!check.match(/\W[A-Za-z][A-Za-z]\W/)) check += ', CA'
      if (check === lastChecked) return
      const body = {address: {regionCode: 'US', addressLines: [check]}}
      const response = await fetch(
        `https://addressvalidation.googleapis.com/v1:validateAddress?key=AIzaSyCgi3GzjWG35S89-tnkxHgi5TJVD2eUe2o`,
        {method: 'POST', body: JSON.stringify(body)})
      if (!response || !response.ok) {
        error.value =
          'The address could not be verified because the address verification service is not available.'
        return
      }
      const result = await response.json()
      if (result.verdict.validationGranularity !== 'SUB_PREMISE' && result.verdict.validationGranularity !== 'PREMISE') {
        lastChecked = check
        error.value = "We couldn't locate this address.  Please provide a valid address."
      } else {
        let fmt = result.address.formattedAddress
        line1.value = fmt.split(',')[0]
        line2.value = fmt.replace(/^[^,]*, */, '')
        line2.value = line2.value.replace(/-[0-8]{4}/, '') // remove ZIP+4
        lastChecked = line1.value + ', ' + line2.value
        latitude = longitude = 0
        if (result.geocode && result.geocode.location) {
          latitude = result.geocode.location.latitude
          longitude = result.geocode.location.longitude
        }
        emit('update:modelValue', {
          address: lastChecked,
          latitude: latitude,
          longitude: longitude,
        })
      }
    }
    watch(submitted!, validate)
    function onBlur(evt: FocusEvent) {
      if (evt.relatedTarget === line1input.value || evt.relatedTarget === line2input.value) return // focus is still in one of the two inputs
      validate()
    }

    // The label to display depends on the address type.
    const label =
      props.type === 'Home'
        ? 'Home Address'
        : props.type === 'Work'
          ? 'Business Hours Address'
          : 'Mailing Address'

    return { error, label, line1, line1input, line2, line2input, onFocusLine2, onBlur, sameAsHome }
  },
})
</script>

<style lang="postcss">
.person-address-label-Home {
}
.person-address-label-Work,
.person-address-label-Mail {
  padding-top: 0;
}
.person-address-same-box {
  white-space: nowrap;
}
.person-address-same-label {
  margin: 0 0 0 0.5rem;
}
.person-address-input {
  display: block;
  width: 100%;
  margin-top: 0.5rem;
  &:first-child {
    margin-top: 0;
  }
}
</style>
