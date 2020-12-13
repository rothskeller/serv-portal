// Radio button group.

import { defineComponent, h, ref, watch } from "vue"
import { SRadio } from './rc'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    disabledKey: { type: String },
    modelValue: { type: String, required: true },
    inline: { type: Boolean, default: false },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const selected = ref(props.modelValue)
    watch(() => props.modelValue, () => { selected.value = props.modelValue })
    watch(selected, () => { emit('update:modelValue', selected.value) })

    function optionValue(option: any) {
      return typeof option === 'object' ? option[props.valueKey] : option
    }
    function optionLabel(option: any) {
      return typeof option === 'object' ? option[props.labelKey] : option.toString()
    }
    function isChecked(option: any) {
      return selected.value === optionValue(option)
    }
    function isDisabled(option: any) {
      if (typeof option !== 'object') return false
      if (props.disabledKey) return !!option[props.disabledKey]
      return false
    }

    return () => props.options.map((o, i) => h(SRadio, {
      id: `${props.id}-${i}`,
      label: optionLabel(o),
      modelValue: isChecked(o),
      inline: props.inline,
      disabled: isDisabled(o),
      'onUpdate:modelValue': (checked: boolean) => {
        if (checked) selected.value = optionValue(o)
      },
    }))
  },
})
