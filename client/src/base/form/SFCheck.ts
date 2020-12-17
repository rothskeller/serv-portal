// Validated form check box control (individual).

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import { SCheck } from '../controls/rc'
import FormItem, { ErrorFunction, useLostFocus } from './FormItem'

const SFCheck = defineComponent({
  name: 'SFCheck',
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: Boolean, required: true },
    disabled: { type: Boolean, default: false },
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

    function renderInput() {
      return h(SCheck, {
        ...attrs,
        id: props.id,
        label: props.label,
        modelValue: props.modelValue,
        'onUpdate:modelValue': (v: string) => emit('update:modelValue', v),
        disabled: props.disabled,
        onFocus,
        onBlur,
      })
    }

    // Render a FormItem, with an SCheck inside it.
    return () =>
      h(
        FormItem,
        {
          id: props.id,
          xclass: 'check',
          label: '',
          help: props.help,
          error: error.value,
        },
        renderInput
      )
  },
})
export default SFCheck
