import { defineComponent, h } from 'vue'
import { propagateError } from '../util'

export { propagateError, useLostFocus } from '../util'
export type { ErrorFunction } from '../util'

const FormItem = defineComponent({
  name: 'FormItem',
  props: {
    id: { type: String, required: true },
    xclass: String,
    label: String,
    help: String,
    error: String,
  },
  setup(props, { slots }) {
    propagateError(props.id, () => props.error || '')

    const renderLabel = () =>
      h('label', { class: ['form-item-label', props.xclass], for: props.id }, props.label)
    const renderInput = () =>
      h('div', { class: ['form-item-input', props.xclass] }, slots.default?.())
    const renderHelp = () =>
      h('div', { class: ['form-item-help', props.xclass] }, [
        props.error ? h('div', { class: 'form-item-error-text' }, props.error) : null,
        props.help ? h('div', { class: 'form-item-help-text' }, props.help) : null,
      ])

    return () => [renderLabel(), renderInput(), renderHelp()]
  },
})
export default FormItem
