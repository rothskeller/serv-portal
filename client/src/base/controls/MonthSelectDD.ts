// MonthSelectDD is the dropdown inside the MonthSelect control.

import { defineComponent, h, ref, watch } from 'vue'
import { SButton, SIcon } from '../../base'

const MonthSelectDD = defineComponent({
  name: 'MonthSelectDD',
  props: {
    month: { type: String, required: true },
  },
  emits: ['update'],
  setup(props, { emit }) {
    const year = ref(parseInt(props.month.substr(0, 4)))
    watch(
      () => props.month,
      () => (year.value = parseInt(props.month.substr(0, 4)))
    )
    function onYearArrow(dir: number) {
      year.value += dir
    }
    function onMonthButton(m: number) {
      emit('update', `${year.value}-${m < 10 ? '0' : ''}${m}`)
    }
    function renderArrow(left: boolean) {
      return h(
        'div',
        { class: 'mselect-arrow', onClick: () => onYearArrow(left ? -1 : +1) },
        h(SIcon, { class: 'mselect-icon', icon: left ? 'chevron-left' : 'chevron-right' })
      )
    }
    function renderYear() {
      return h('div', { class: 'mselectdd-year' }, year.value.toString())
    }
    function renderTopBar() {
      return h('div', { class: 'mselectdd-top' }, [
        renderArrow(true),
        renderYear(),
        renderArrow(false),
      ])
    }
    function renderMonthButtons() {
      return [
        'Jan',
        'Feb',
        'Mar',
        'Apr',
        'May',
        'Jun',
        'Jul',
        'Aug',
        'Sep',
        'Oct',
        'Nov',
        'Dec',
      ].map((m, i) =>
        h(
          SButton,
          {
            class: 'mselectdd-month',
            onClick: () => onMonthButton(i + 1),
          },
          () => m
        )
      )
    }
    return () => h('div', { class: 'mselectdd' }, [renderTopBar(), renderMonthButtons()])
  },
})
export default MonthSelectDD
