// Check box and radio button groups.

import { ComponentOptions, defineComponent, h, Ref, ref, watch } from 'vue'
import { propagateModel } from '../util'
import { SCheck, SRadio } from './rc'

function defineRCGroup(
  name: string,
  control: typeof SCheck,
  modelType: any,
  isChecked: (selected: any, option: any) => boolean,
  apply: (selected: Ref<any>, option: any, checked: boolean) => void
) {
  return defineComponent({
    name,
    props: {
      id: { type: String, required: true },
      options: { type: Array, required: true },
      valueKey: { type: String, default: 'value' },
      labelKey: { type: String, default: 'label' },
      disabledKey: { type: String },
      modelValue: { type: modelType, required: true },
      inline: { type: Boolean, default: false },
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
      const selected = propagateModel(props, emit)

      function optionValue(option: any) {
        return typeof option === 'object' ? option[props.valueKey] : option
      }
      function optionLabel(option: any) {
        return typeof option === 'object' ? option[props.labelKey] : option.toString()
      }
      function isDisabled(option: any) {
        if (typeof option !== 'object') return false
        if (props.disabledKey) return !!option[props.disabledKey]
        return false
      }

      return () =>
        props.options.map((o, i) =>
          h(control, {
            id: `${props.id}-${i}`,
            label: optionLabel(o),
            modelValue: isChecked(selected.value, optionValue(o)),
            inline: props.inline,
            disabled: isDisabled(o),
            'onUpdate:modelValue': (checked: boolean) => {
              apply(selected, optionValue(o), checked)
            },
          })
        )
    },
  })
}

const SCheckGroup = defineRCGroup(
  'SCheckGroup',
  SCheck,
  Set,
  (selected: Set<unknown>, option: unknown) => selected.has(option),
  (selected: Ref<Set<unknown>>, option: unknown, checked: boolean) => {
    if (checked) selected.value.add(option)
    else selected.value.delete(option)
    selected.value = new Set(selected.value) // to force reactivity
  }
)
const SRadioGroup = defineRCGroup(
  'SRadioGroup',
  SRadio,
  String,
  (selected: string, option: string) => selected === option,
  (selected: Ref<string>, option: string, checked: boolean) => {
    if (checked) selected.value = option
  }
)

export { SCheckGroup, SRadioGroup }
