// SButton displays a styled button.

import { defineComponent, h } from 'vue'
import { RouterLink } from 'vue-router'

const SButton = defineComponent({
  name: 'SButton',
  props: {
    disabled: { type: Boolean, default: false },
    variant: { type: String, default: 'secondary' },
    type: { type: String, default: 'button' },
    to: [String, Object],
    small: { type: Boolean, default: false },
  },
  setup(props, { attrs, slots }) {
    function render() {
      return h(
        // @ts-ignore - it can't handle the polymorphism of the first argument
        props.to ? RouterLink : 'button',
        {
          ...attrs,
          class: [
            'sbtn',
            `sbtn-${props.variant}`,
            { 'sbtn-disabled': props.disabled, 'sbtn-small': props.small },
            attrs.class,
          ],
          to: props.to,
          type: props.to ? null : props.type,
          disabled: props.disabled,
        },
        slots.default?.()
      )
    }
    return render
  },
})
export default SButton
