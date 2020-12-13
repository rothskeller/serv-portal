// Drop-down select control.

import { defineComponent, h, ref, watch } from 'vue'

export default defineComponent({
  props: {
    id: String,
    modelValue: { type: [String, Number], required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
  },
  emits: ['update:modelValue'],
  setup(props, { attrs, emit }) {
    const input = ref(props.modelValue)
    if (!props.options.find((o) => input.value === optionValue(o))) {
      if (props.modelValue)
        console.warn(`Initial value for SSelect#${props.id} is not one of the allowed options.`)
      input.value = optionValue(props.options[0])
    }
    watch(
      () => props.modelValue,
      () => {
        if (props.options.find((o) => props.modelValue === optionValue(o)))
          input.value = props.modelValue
        else
          console.warn(`Updated value for SSelect#${props.id} is not one of the allowed options.`)
      }
    )
    watch(input, () => {
      emit('update:modelValue', input.value)
    })

    function optionValue(o: any) {
      return typeof o === 'object' ? o[props.valueKey] : o
    }
    function optionLabel(o: any) {
      return typeof o === 'object' ? o[props.labelKey] : o.toString()
    }

    return () =>
      h(
        'select',
        {
          ...attrs,
          id: props.id,
          class: 'control',
          value: input.value,
          onChange: (evt: Event) => {
            input.value = (evt.target as HTMLSelectElement).value
          },
        },
        props.options.map((o) =>
          h(
            'option',
            {
              value: optionValue(o),
              selected: optionValue(o) === input.value,
            },
            optionLabel(o)
          )
        )
      )
  },
})
