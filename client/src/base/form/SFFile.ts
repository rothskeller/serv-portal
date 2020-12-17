// Validated form file upload control.

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import { SInput } from '../controls/text'
import FormItem, { ErrorFunction, useLostFocus } from './FormItem'

const SFFile = defineComponent({
  name: 'SFFile',
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    modelValue: FileList,
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

    // Updates come on the files attribute of the input.
    function onChange(evt: Event) {
      const target = evt.target! as HTMLInputElement
      emit('update:modelValue', target.files)
    }

    function renderInput() {
      return h(SInput as any, {
        ...attrs,
        id: props.id,
        class: [attrs.class, { 'form-control-invalid': !!error.value }, 'w100'],
        type: 'file',
        onChange,
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
          xclass: 'oneline',
          label: props.label,
          help: props.help,
          error: error.value,
        },
        renderInput
      )
  },
})
export default SFFile
