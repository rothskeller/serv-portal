import { defineComponent } from 'vue'
import { propagateError } from '../util'

export { propagateError, useLostFocus } from '../util'
export type { ErrorFunction } from '../util'

export default defineComponent({
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
    return () => [
      <label
        class={['form-item-label', props.xclass]}
        for={props.id}
      >{props.label}</label>,
      <div class={['form-item-input', props.xclass]}>{slots.default()}</div>,
      <div class={['form-item-help', props.xclass]}>
        {props.error ? <div class="form-item-error-text">{props.error}</div> : null}
        {props.help ? <div class="form-item-help-text">{props.help}</div> : null}
      </div>
    ]
  },
})
