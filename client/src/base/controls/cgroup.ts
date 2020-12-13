// Check box group.

import { defineComponent, h, ref, watch } from "vue"
import { SCheck } from './rc'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    options: { type: Array, required: true },
    valueKey: { type: String, default: 'value' },
    labelKey: { type: String, default: 'label' },
    disabledKey: { type: String },
    modelValue: { type: Set, required: true },
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
      return selected.value.has(optionValue(option))
    }
    function isDisabled(option: any) {
      if (typeof option !== 'object') return false
      if (props.disabledKey) return !!option[props.disabledKey]
      return false
    }

    return () => props.options.map((o, i) => h(SCheck, {
      id: `${props.id}-${i}`,
      label: optionLabel(o),
      modelValue: isChecked(o),
      inline: props.inline,
      disabled: isDisabled(o),
      'onUpdate:modelValue': (checked: boolean) => {
        if (checked) selected.value.add(optionValue(o))
        else selected.value.delete(optionValue(o))
      },
    }))
  },
})
