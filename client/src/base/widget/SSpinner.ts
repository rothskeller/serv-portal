import { defineComponent, h } from 'vue'
import './widget.css'

const SSpinner = defineComponent({
  name: 'SSpinner',
  render() {
    return h('span', { class: 'sspinner', 'aria-hidden': true })
  },
})
export default SSpinner
