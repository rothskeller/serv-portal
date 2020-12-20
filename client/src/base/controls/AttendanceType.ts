// AttendanceType displays an attendance type badge, and optionally allows it to
// be changed.

import { defineComponent, h } from 'vue'

const typeLabels: Record<string, string> = {
  '': '',
  Volunteer: 'Volunteer',
  Student: 'Student',
  Audit: 'Auditor',
  Absent: 'Absentee',
}

const AttendanceType = defineComponent({
  name: 'AttendanceType',
  props: {
    value: { type: String, required: true },
    allowEmpty: { type: Boolean, default: false },
    disabled: { type: Boolean, default: false },
  },
  emits: ['update'],
  setup(props, { emit }) {
    function onClick() {
      if (props.disabled) return
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
          class: ['attendance-type', `attendance-type-${props.value}`],
          onClick,
        },
        typeLabels[props.value]
      )
  },
})
export default AttendanceType
