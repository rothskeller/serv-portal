// Validated form text area control.

import { defineComponent, getCurrentInstance, h, PropType, ref, watchEffect } from 'vue'
import STextArea from '../controls/textarea'
import FormItem, { ErrorFunction, useLostFocus } from './item'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    autofocus: { type: Boolean, default: false },
    errorFn: Function as PropType<ErrorFunction>,
  },
  emits: ['update:modelValue'],
  setup(props, { attrs, emit, expose }) {
    // Provide a focus() function so parent can direct focus to text field.
    const instance = getCurrentInstance()!
    function focus() {
      ;(instance?.refs['input'] as typeof STextArea)?.focus()
    }
    expose({ focus })

    // Set up for form validation.
    const { lostFocus, submitted, onFocus, onBlur } = useLostFocus()
    const error = ref('')
    if (props.errorFn)
      watchEffect(() => {
        error.value = props.errorFn!(lostFocus.value!, submitted?.value || false)
      })

    function renderInput() {
      return h(STextArea, {
        ...attrs,
        id: props.id,
        ref: 'input',
        class: [attrs.class, { 'form-control-invalid': !!error.value }, 'w100'],
        modelValue: props.modelValue,
        'onUpdate:modelValue': (v: string) => emit('update:modelValue', v),
        trim: props.trim,
        autofocus: props.autofocus,
        onFocus,
        onBlur,
      })
    }

    // Render a FormItem, with an SInput inside it.
    return () =>
      h(
        FormItem,
        {
          id: props.id,
          xclass: 'textarea',
          label: props.label,
          help: props.help,
          error: error.value,
        },
        renderInput
      )
  },
})
