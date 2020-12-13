// Framework for a basic form item.

import { computed, defineComponent, h, inject, Ref, ref, Slot, watch } from 'vue'

function renderLabel(id: string, label: string, xclass?: string) {
  return h(
    'label',
    {
      id: `${id}_label`,
      class: ['form-item-label', xclass],
      for: id,
    },
    label
  )
}

function renderInput(slot: Slot, xclass?: string) {
  return h('div', { class: ['form-item-input', xclass] }, slot())
}

function renderHelp(help: string, error: string, xclass?: string) {
  const children = []
  if (error) children.push(h('div', { class: 'form-item-error-text' }, error))
  if (help) children.push(h('div', { class: 'form-item-help-text' }, help))
  return h(
    'div',
    {
      class: ['form-item-help', xclass],
    },
    children
  )
}

export default defineComponent({
  props: {
    id: { type: String, required: true },
    xclass: String,
    label: String,
    help: String,
    error: String,
  },
  setup(props, { slots }) {
    // Notify the form when the control validity changes.
    const setValidity = inject<(id: string, isValid: boolean) => void>('setValidity')
    watch(
      () => props.error,
      () => {
        setValidity?.(props.id, !props.error)
      }
    )

    return () => [
      renderLabel(props.id, props.label || '', props.xclass),
      renderInput(slots.default!, props.xclass),
      renderHelp(props.help || '', props.error || '', props.xclass),
    ]
  },
})

// This module also provides a useLostFocus function which callers can use to
// track focus status for validation.  It returns an object containing:
//   - onFocus and onBlur handlers that should be attached to the input control
//   - lostFocus reactive boolean that is true when lost-focus validation
//     should occur
//   - submitted reactive boolean  that is true when the user has attempted to
//     submit the form (and the form has not been subsequently reset)
//
// Lost-focus validation should happen, and lostFocus is true, when:
//   (control does not currently have focus)
//   and (
//     (control has had focus since form was reset) or
//     (user has attempted to submit since form was reset)
//   )
export function useLostFocus() {
  const hasHadFocus = ref(false)
  const hasFocus = ref(false)
  const submitted = inject<Ref<boolean>>('formSubmitted')
  const lostFocus = computed(() => !hasFocus.value && (hasHadFocus.value || submitted?.value))
  if (submitted)
    watch(submitted, () => {
      hasHadFocus.value = submitted.value
      // change from true to false means the form was reset
    })
  function onFocus() {
    hasFocus.value = hasHadFocus.value = true
  }
  function onBlur() {
    hasFocus.value = false
  }
  return { submitted, lostFocus, onFocus, onBlur }
}

// Many form items use this definition; it's here for convenience since they all
// import this module anyway.
export type ErrorFunction = (lostFocus: boolean, submitted: boolean) => string
