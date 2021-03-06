// Validated form radio button and check box group controls.

import { defineComponent, h, PropType, ref, watchEffect } from 'vue'
import { SCheckGroup, SRadioGroup } from '../controls/rcgroup'
import FormItem, { ErrorFunction, useLostFocus } from './FormItem'

function defineRCG(name: string, modelType: any, Group: typeof SCheckGroup) {
  return defineComponent({
    name,
    props: {
      id: { type: String, required: true },
      label: String,
      help: String,
      options: { type: Array, required: true },
      valueKey: { type: String, default: 'value' },
      labelKey: { type: String, default: 'label' },
      disabledKey: { type: String },
      modelValue: { type: modelType, required: true },
      inline: { type: Boolean, default: false },
      errorFn: Function as PropType<ErrorFunction>,
    },
    emits: ['update:modelValue'],
    setup(props, { attrs, emit }) {
      // Set up for form validation.
      const { submitted } = useLostFocus()
      const error = ref('')
      if (props.errorFn)
        watchEffect(() => {
          error.value = props.errorFn!(false, submitted?.value || false)
        })

      // Render a FormItem, with a control group inside it.
      return () =>
        h(
          FormItem,
          {
            id: props.id,
            label: props.label,
            help: props.help,
            error: error.value,
          },
          () =>
            h(Group, {
              id: props.id,
              options: props.options,
              valueKey: props.valueKey,
              labelKey: props.labelKey,
              disabledKey: props.disabledKey,
              modelValue: props.modelValue,
              inline: props.inline,
              'onUpdate:modelValue': (v: any) => emit('update:modelValue', v),
            })
        )
    },
  })
}

const SFCheckGroup = defineRCG('SFCheckGroup', Set, SCheckGroup)
const SFRadioGroup = defineRCG('SFRadioGroup', String, SRadioGroup)

export { SFCheckGroup, SFRadioGroup }
