// EventTimesheetEntry displays and updates the grid row for a single person.

import moment from 'moment-mini'
import { defineComponent, Fragment, h, PropType, ref, VNode, watch, watchEffect } from 'vue'
import { AttendanceType, SInput, SSelect } from '../../base'
import type { TimesheetEntry } from './EventTimesheet'

function fmtMinutes(m: number): string {
  if (!m) return ''
  const hours = Math.floor(m / 60).toString()
  return m === 30
    ? '½ hour'
    : m === 60
    ? '1 hour'
    : m % 60 !== 0
    ? `${hours}½ hours`
    : `${hours} hours`
}

const EventTimesheetEntry = defineComponent({
  name: 'EventTimesheetEntry',
  props: {
    entry: { type: Object as PropType<TimesheetEntry>, required: true },
    defaultType: { type: String, required: true },
  },
  emits: ['newType'],
  setup(props, { emit }) {
    // Keep track of when one of the two input fields has focus.
    const startRef = ref<VNode>()
    const endRef = ref<VNode>()
    const focused = ref(false)
    function onFocus() {
      focused.value = true
    }
    function onBlur(evt: FocusEvent) {
      if (startRef.value?.el && evt.relatedTarget === startRef.value?.el) return
      if (endRef.value?.el && evt.relatedTarget === endRef.value?.el) return
      focused.value = false
    }

    // Update the type, minutes, and invalid fields when the start and end
    // fields change.
    const showInvalid = ref(props.entry.invalid)
    watchEffect(() => {
      if (!props.entry.start && !props.entry.end) {
        props.entry.minutes = props.entry.origMinutes
        if (!props.entry.origType) props.entry.type = ''
        props.entry.invalid = false
      } else {
        const start = moment(props.entry.start, 'H:mm', true)
        const end = moment(props.entry.end, 'H:mm', true)
        if (!start.isValid() || !end.isValid() || end.isBefore(start)) {
          props.entry.minutes = 0
          props.entry.invalid = true
        } else {
          let minutes = end.diff(start, 'minutes')
          props.entry.minutes = Math.floor((minutes + 20) / 30) * 30
          props.entry.invalid = false
        }
      }
      if (props.entry.minutes && !props.entry.type) props.entry.type = props.defaultType
      showInvalid.value = props.entry.invalid && !focused.value
    })
    watch(
      () => props.entry.type,
      () => {
        if (props.entry.type && props.entry.type !== props.defaultType)
          emit('newType', props.entry.type)
      }
    )

    return () =>
      h(Fragment, [
        h('div', { class: 'event-timesheet-name' }, props.entry.sortName),
        h(AttendanceType, {
          value: props.entry.type,
          allowEmpty: !props.entry.minutes,
          onUpdate: (v: string) => (props.entry.type = v),
        }),
        h(SInput, {
          ref: (r: any) => (startRef.value = r),
          class: 'event-timesheet-time',
          type: 'time',
          modelValue: props.entry.start,
          onBlur,
          onFocus,
          'onUpdate:modelValue': (v: string) => (props.entry.start = v),
        }),
        h(SInput, {
          ref: (r: any) => (endRef.value = r),
          class: 'event-timesheet-time',
          type: 'time',
          modelValue: props.entry.end,
          onBlur,
          onFocus,
          'onUpdate:modelValue': (v: string) => (props.entry.end = v),
        }),
        showInvalid.value
          ? h('div', { class: 'event-timesheet-invalid' }, 'Invalid')
          : h('div', { class: 'event-timesheet-minutes' }, fmtMinutes(props.entry.minutes)),
      ])
  },
})
export default EventTimesheetEntry
