<!--
PersonEditDSW displays the editor control for a single DSW class.
-->

<template lang="pug">
SFInput(
  :id='`person-dsw-${type}`',
  type='date',
  :label='`DSW ${type}`',
  v-model='value',
  :errorFn='errorFn',
  :help='`Date when ${type} DSW registration form was signed.`'
)
</template>

<script lang="ts">
import { defineComponent, ref, toRefs, watch } from 'vue'
import SFInput from '../../base/SFInput.vue'

export default defineComponent({
  components: { SFInput },
  props: {
    type: { type: String, required: true },
    modelValue: { type: String, required: true },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // Reflect changes of property in local value.
    const { modelValue } = toRefs(props)
    const value = ref(modelValue.value)
    watch(modelValue, () => {
      value.value = modelValue.value
    })

    // Propagate value changes to parent.
    watch(value, () => {
      emit('update:modelValue', value.value)
    })

    // Validate input.
    function errorFn(lostFocus: boolean) {
      if (!lostFocus || !value.value) return ''
      if (!value.value.match(/^20\d\d-\d\d-\d\d$/)) return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }

    return { errorFn, value }
  },
})
</script>
