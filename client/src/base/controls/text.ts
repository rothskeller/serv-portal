// Text input controls (oneline and area).

import { defineComponent, getCurrentInstance, h, onMounted, PropType, ref, watch } from 'vue'
import { propagateModel } from '../util'
import './controls.css'

function defineTextInput(name: string, tag: string) {
  return defineComponent({
    name,
    props: {
      modelValue: { type: String, required: true },
      trim: { type: Boolean, default: false },
      restrictFn: Function as PropType<(s: string) => string>,
      autofocus: { type: Boolean, default: false },
    },
    emits: ['blur', 'update:modelValue'],
    setup(props, { attrs, emit, expose }) {
      const instance = getCurrentInstance()!
      const input = propagateModel(props, emit)
      const inputRef = ref<HTMLInputElement | HTMLTextAreaElement>()

      // Handle autofocus and manual focus.
      if (props.autofocus)
        onMounted(() => {
          inputRef.value?.focus()
        })
      function focus() {
        inputRef.value?.focus()
      }
      expose({ focus })

      // Apply trim when losing focus.  (If we apply it more eagerly, backspacing
      // over the start of a word also removes the space before it, which is
      // disturbing.)
      function onBlur(evt: FocusEvent) {
        if (props.trim) input.value = input.value.trim()
        emit('blur', evt)
      }

      // Apply restrictions when receiving changes.
      function onInput(evt: InputEvent) {
        const tgt = evt.target! as HTMLInputElement
        const nv = props.restrictFn ? props.restrictFn(tgt.value) : tgt.value
        input.value = nv
        if (nv !== tgt.value) {
          // Forcible update required to rerender control if value changed.
          instance?.update()
        }
      }

      return () =>
        h(tag, {
          ...attrs,
          class: 'control',
          ref: (el) => (inputRef.value = el as HTMLInputElement | HTMLTextAreaElement),
          value: input.value,
          onBlur,
          onInput,
        })
    },
  })
}

const SInput = defineTextInput('SInput', 'input')
const STextArea = defineTextInput('STextArea', 'textarea')

export { SInput, STextArea }
