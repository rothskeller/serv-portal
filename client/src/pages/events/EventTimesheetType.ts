// EventTimesheetType displays the attendance type for an entry and allows it to
// be toggled.

import { defineComponent, h } from 'vue'

const typeLabels: Record<string, string> = {
  '': '',
  Volunteer: 'Volunteer',
  Student: 'Student',
  Audit: 'Auditor',
  Absent: 'Absentee',
}

const EventTimesheetType = defineComponent({
  name: 'EventTimesheetType',
  props: {
    value: { type: String, required: true },
    allowEmpty: { type: Boolean, required: true },
  },
  emits: ['update'],
  setup(props, { emit }) {
    function onClick() {
      switch (props.value) {
        case '':
          emit('update', 'Volunteer')
          break
        case 'Volunteer':
          emit('update', 'Student')
          break
        case 'Student':
          emit('update', 'Audit')
          break
        case 'Audit':
          emit('update', 'Absent')
          break
        case 'Absent':
          if (props.allowEmpty) emit('update', '')
          else emit('update', 'Volunteer')
          break
      }
    }
    return () =>
      h(
        'div',
        {
          class: ['event-timesheet-type', `event-timesheet-type-${props.value}`],
          onClick,
        },
        typeLabels[props.value]
      )
  },
})
export default EventTimesheetType
