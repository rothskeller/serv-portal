// Validated form radio button group control.

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import SRadioGroup from '../controls/rgroup'
import FormItem, { ErrorFunction, useLostFocus } from './item'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    disabledKey: { type: String },
    modelValue: { type: String, required: true },
    inline: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
  },
  emits: ['update:modelValue'],
  setup(props, { attrs, emit }) {
    // Set up for form validation.
    const { submitted } = useLostFocus()
    const error = ref('')
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(false, submitted?.value || false)
      })

    function renderInput() {
      return h(SRadioGroup, {
        id: props.id,
        options: props.options,
        valueKey: props.valueKey,
        labelKey: props.labelKey,
        disabledKey: props.disabledKey,
        modelValue: props.modelValue,
        'onUpdate:modelValue': (v: string) => emit('update:modelValue', v),
      })
    }

    // Render a FormItem, with an SRadioGroup inside it.
    return () =>
      h(
        FormItem,
        {
          id: props.id,
          label: props.label,
          help: props.help,
          error: error.value,
        },
        renderInput
      )
  },
})
