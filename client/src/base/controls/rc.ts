// Radio buttons and check boxes.

import { defineComponent, h, ref, watch } from "vue"
import './controls.css'

function defineRC(radio: boolean) {
  return defineComponent({
    props: {
      id: { type: String, required: true },
      label: String,
      modelValue: { type: Boolean, required: true },
      inline: { type: Boolean, default: false },
      disabled: { type: Boolean, default: false },
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
      const value = ref(props.modelValue)
      watch(() => props.modelValue, () => { value.value = props.modelValue })
      watch(value, () => { emit(`update:modelValue`, value.value) })
      return () => h('div', {
        class: [
          'rc',
          props.inline ? 'rc-inline' : 'rc-stacked',
          radio ? 'radio' : 'check',
        ],
      }, [
        h('input', {
          id: props.id,
          class: 'rc-input',
          type: radio ? 'radio' : 'checkbox',
          checked: value.value,
          disabled: props.disabled,
          onChange: (evt: Event) => {
            value.value = (evt.target as HTMLInputElement).checked
          },
        }),
        h('label', {
          class: 'rc-label',
          for: props.id,
        }, props.label)
      ])
    }
  })
}

const SRadio = defineRC(true)
const SCheck = defineRC(false)

export { SRadio, SCheck }
