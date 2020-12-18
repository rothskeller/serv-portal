// MonthSelect displays a control that allows the user to select a month.

import moment from 'moment-mini'
import { defineComponent, h, ref, watch } from 'vue'
import { SIcon } from '../../base'
import MonthSelectDD from './MonthSelectDD'

const MonthSelect = defineComponent({
  name: 'MonthSelect',
  props: {
    modelValue: { type: String, required: true },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const showingDropdown = ref(false)
    function onToggleDropdown() {
      showingDropdown.value = !showingDropdown.value
    }
    function onMonthArrow(dir: number) {
      showingDropdown.value = false
      emit(
        'update:modelValue',
        moment(props.modelValue, 'YYYY-MM', true).add(dir, 'month').format('YYYY-MM')
      )
    }
    function onDropdownResult(m: string) {
      showingDropdown.value = false
      emit('update:modelValue', m)
    }
    function renderArrow(left: boolean) {
      return h(
        'div',
        { class: 'mselect-arrow', onClick: () => onMonthArrow(left ? -1 : +1) },
        h(SIcon, { class: 'mselect-icon', icon: left ? 'chevron-left' : 'chevron-right' })
      )
    }
    function renderMonth() {
      return h(
        'div',
        { class: 'mselect-month', onClick: onToggleDropdown },
        moment(props.modelValue, 'YYYY-MM', true).format('MMMM YYYY')
      )
    }
    function renderDropdown() {
      return h(MonthSelectDD, { month: props.modelValue, onUpdate: onDropdownResult })
    }
    return () =>
      h('div', { class: 'mselect' }, [
        renderArrow(true),
        renderMonth(),
        renderArrow(false),
        showingDropdown.value ? renderDropdown() : null,
      ])
  },
})
export default MonthSelect
