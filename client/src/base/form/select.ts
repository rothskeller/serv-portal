// Validated form select dropdown control.

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import SSelect from '../controls/select'
import FormItem, { ErrorFunction, useLostFocus } from './item'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: [String, Number], required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    errorFn: Function as PropType<ErrorFunction>,
  },
  emits: ['update:modelValue'],
  setup(props, { attrs, emit }) {
    // Set up for form validation.
    const { lostFocus, submitted, onFocus, onBlur } = useLostFocus()
    const error = ref('')
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value!, submitted?.value || false)
      })

    function renderSelect() {
      return h(SSelect, {
        ...attrs,
        id: props.id,
        class: [attrs.class, { 'form-control-invalid': !!error.value }, 'w100'],
        modelValue: props.modelValue,
        'onUpdate:modelValue': (v: string) => emit('update:modelValue', v),
        options: props.options,
        valueKey: props.valueKey,
        labelKey: props.labelKey,
        onFocus,
        onBlur,
      })
    }

    // Render a FormItem, with an SSelect inside it.
    return () =>
      h(
        FormItem,
        {
          id: props.id,
          xclass: 'oneline',
          label: props.label,
          help: props.help,
          error: error.value,
        },
        renderSelect
      )
  },
})
