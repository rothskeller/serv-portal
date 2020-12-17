// Validated form text input controls (oneline and area).

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import { SInput, STextArea } from '../controls/text'
import FormItem, { ErrorFunction, useLostFocus } from './FormItem'

function defineTextInput(name: string, control: typeof SInput, xclass: string) {
  return defineComponent({
    name,
    props: {
      id: { type: String, required: true },
      label: String,
      help: String,
      modelValue: { type: String, required: true },
      trim: { type: Boolean, default: false },
      autofocus: { type: Boolean, default: false },
      errorFn: Function as PropType<ErrorFunction>,
      restrictFn: Function as PropType<(s: string) => string>,
    },
    emits: ['update:modelValue'],
    setup(props, { attrs, emit, expose }) {
      // Provide a focus() function so parent can direct focus to input field.
      const inputRef = ref<HTMLInputElement | HTMLTextAreaElement>()
      function focus() {
        inputRef.value?.focus()
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
        return h(control, {
          ...attrs,
          id: props.id,
          ref: (el) => (inputRef.value = el as HTMLInputElement | HTMLTextAreaElement),
          class: [attrs.class, { 'form-control-invalid': !!error.value }, 'w100'],
          modelValue: props.modelValue,
          'onUpdate:modelValue': (v: string) => emit('update:modelValue', v),
          trim: props.trim,
          restrictFn: props.restrictFn,
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
            xclass,
            label: props.label,
            help: props.help,
            error: error.value,
          },
          renderInput
        )
    },
  })
}

const SFInput = defineTextInput('SFInput', SInput, 'oneline')
const SFTextArea = defineTextInput('SFTextArea', STextArea, 'textarea')

export { SFInput, SFTextArea }
