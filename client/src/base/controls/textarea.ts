// Text area control.

import { defineComponent, getCurrentInstance, h, onMounted, ref, watch } from 'vue'
import './controls.css'

export default defineComponent({
  props: {
    modelValue: { type: String, required: true },
    trim: { type: Boolean, default: false },
    autofocus: { type: Boolean, default: false },
  },
  emits: ['blur', 'update:modelValue'],
  setup(props, { attrs, emit, expose }) {
    const instance = getCurrentInstance()!

    const input = ref(props.modelValue)
    watch(
      () => props.modelValue,
      () => {
        input.value = props.modelValue
      }
    )
    watch(input, () => {
      emit('update:modelValue', input.value)
    })

    if (props.autofocus)
      onMounted(() => {
        ;(instance.refs['input'] as HTMLInputElement)?.focus()
      })
    function focus() {
      ;(instance.refs['input'] as HTMLInputElement)?.focus()
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
      input.value = tgt.value
    }

    return () =>
      h('textarea', {
        ...attrs,
        class: 'control',
        ref: 'input',
        value: input.value,
        onBlur,
        onInput,
      })
  },
})
