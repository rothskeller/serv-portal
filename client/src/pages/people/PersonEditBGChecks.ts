// PersonEditBGChecks is the background-check editor part of the
// PersonEditStatus dialog box.

import moment from 'moment-mini'
import { defineComponent, Fragment, h, PropType, ref } from 'vue'
import { SFCheck, SFCheckGroup, SFInput } from '../../base'
import type { GetPersonStatusBGCheck } from './api'

const PersonEditBGChecks = defineComponent({
  name: 'PersonEditBGChecks',
  props: {
    checks: { type: Array as PropType<Array<GetPersonStatusBGCheck>>, required: true },
    types: { type: Array as PropType<Array<string>>, required: true },
  },
  setup(props, { expose }) {
    const before = ref([...props.checks])
    const after = ref([] as Array<GetPersonStatusBGCheck>)
    const edit = ref({ date: '', types: new Set<string>(), assumed: false })
    const editing = ref(props.checks.length === 0)

    function closeEdit(): boolean {
      if (!editing.value) return true
      if (edit.value.date && !edit.value.date.match(/^20\d\d-\d\d-\d\d$/)) return false
      if (edit.value.types.size)
        before.value.push({
          date: edit.value.date,
          types: props.types.filter((t) => edit.value.types.has(t)),
          assumed: edit.value.assumed,
        })
      before.value.splice(before.value.length, 0, ...after.value.splice(0))
      before.value.sort((a, b) => a.date.localeCompare(b.date))
      editing.value = false
      return true
    }

    function onEdit(item: GetPersonStatusBGCheck) {
      if (!closeEdit()) return
      edit.value.date = item.date
      edit.value.types = new Set(item.types)
      edit.value.assumed = item.assumed
      const index = before.value.findIndex((i) => i === item)
      after.value.splice(after.value.length, 0, ...before.value.splice(index + 1))
      before.value.splice(index)
      editing.value = true
    }

    function prepareForSave() {
      closeEdit()
      props.checks.splice(0, props.checks.length, ...before.value)
    }

    function startAdd() {
      console.log('startAdd called')
      if (!closeEdit()) return
      edit.value.date = moment().format('YYYY-MM-DD')
      edit.value.types = new Set<string>()
      edit.value.assumed = false
      editing.value = true
    }

    function render() {
      return [
        renderHeader(props.checks),
        renderList(before.value, onEdit),
        editing.value ? renderEditing(edit.value, props.types) : null,
        renderList(after.value, onEdit),
      ]
    }

    expose({ startAdd, prepareForSave })
    return render
  },
})
export default PersonEditBGChecks

function renderHeader(checks: Array<GetPersonStatusBGCheck>) {
  return [
    h('div', { id: 'person-edbg-header', class: 'form-item' }, 'Background Checks'),
    checks.length
      ? h(
          'div',
          { id: 'person-edbg-help', class: 'form-item' },
          'Click on a check to edit it, or click Add to add one.'
        )
      : null,
  ]
}

function renderList(
  list: Array<GetPersonStatusBGCheck>,
  onEdit: (item: GetPersonStatusBGCheck) => void
) {
  return list.map((item) =>
    h(
      'div',
      {
        class: 'form-item-input',
        onClick: () => {
          onEdit(item)
        },
      },
      `${item.date || 'Unknown date'}: ${item.types.join(', ')}${item.assumed ? ' (assumed)' : ''}`
    )
  )
}

function renderEditing(
  edit: { date: string; types: Set<string>; assumed: boolean },
  types: Array<string>
) {
  return [
    h(SFInput, {
      id: 'person-edbg-date',
      label: 'Date',
      help: 'Date on which the background check was cleared (empty if not known)',
      type: 'date',
      modelValue: edit.date,
      'onUpdate:modelValue': (v: string) => (edit.date = v),
      errorFn: (lostFocus: boolean): string => {
        if (!lostFocus || edit.date === '') return ''
        if (!edit.date.match(/^20\d\d-\d\d-\d\d$/)) return 'This is not a valid date.'
        return ''
      },
    }),
    h(SFCheckGroup, {
      id: 'person-edbg-types',
      label: 'Type(s)',
      options: types,
      inline: true,
      modelValue: edit.types,
      'onUpdate:modelValue': (v: Set<string>) => (edit.types = v),
    }),
    h(SFCheck, {
      id: 'person-edbg-assumed',
      label: 'Assumed (no paper record)',
      help: 'Checks are assumed for long-standing volunteers whose historical records are missing.',
      modelValue: edit.assumed,
      'onUpdate:modelValue': (v: boolean) => (edit.assumed = v),
    }),
  ]
}
