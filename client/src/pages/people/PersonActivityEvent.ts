// PersonActivityEvent is a single event on a person's activity page.

import { defineComponent, Fragment, h, PropType, watch } from 'vue'
import { SInput } from '../../base'
import AttendanceType from '../../base/controls/AttendanceType'
import { fmtMinutes, GetPersonHoursEvent } from './api'

const PersonActivityEvent = defineComponent({
  name: 'PersonActivityEvent',
  props: {
    event: { type: Object as PropType<GetPersonHoursEvent>, required: true },
    typesColumn: { type: Boolean, required: true },
  },
  setup(props) {
    watch(
      () => props.event,
      () => {
        if (!props.event.minutes && props.event.canEdit && !props.event.canViewType)
          props.event.type = props.event.placeholder ? 'Volunteer' : 'Absent'
      },
      { immediate: true }
    )
    function render() {
      const date = props.event.placeholder ? props.event.date.substr(0, 7) : props.event.date
      const label = `${date} ${props.event.name}`
      return h(Fragment, [
        h('div', { class: 'person-activity-label' }, label),
        props.typesColumn && props.event.minutes
          ? h(AttendanceType, {
              class: 'person-activity-type',
              value: props.event.placeholder ? 'Volunteer' : props.event.type,
              disabled: !props.event.canViewType || !props.event.canEdit || props.event.placeholder,
              onUpdate: (v: string) => (props.event.type = v),
            })
          : props.typesColumn
          ? h('div')
          : null,
        props.event.canEdit
          ? h(SInput, {
              class: 'person-activity-hours',
              modelValue: (props.event.minutes / 60.0).toString(),
              type: 'number',
              min: '0',
              step: '0.5',
              'onUpdate:modelValue': (v: string) => {
                props.event.minutes = parseFloat(v) * 60.0
                if (isNaN(props.event.minutes)) props.event.minutes = 0
              },
            })
          : h('div', { class: 'person-activity-hours' }, fmtMinutes(props.event.minutes)),
      ])
    }
    return render
  },
})
export default PersonActivityEvent
