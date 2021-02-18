import { computed, inject, ref, Ref, watch, WatchSource } from 'vue'

// propagateModel takes the incoming modelValue to a component and uses it to
// create a local ref.  Updates to the incoming modelValue are propagated to the
// local ref; updates to the local ref are emitted as updates to the modelValue.
export function propagateModel<T>(
  props: { modelValue: T },
  emit: (event: 'update:modelValue', value: T) => void
): Ref<T> {
  const value = ref(props.modelValue) as Ref<T>
  watch(
    () => props.modelValue,
    () => {
      value.value = props.modelValue
    }
  )
  watch(value, () => {
    emit('update:modelValue', value.value)
  })
  return value
}

// propagateError watches an error message string and notifies the containing
// form when it is non-empty so that the form knows its contents aren't valid.
export function propagateError(id: string, error: WatchSource<string>) {
  const setValidity = inject<(id: string, isValid: boolean) => void>('setValidity')
  watch(error, (n) => {
    setValidity?.(id, !n)
  })
}

// useLostFocus is a helper for form control validation.  It returns an object
// containing:
//   - onFocus and onBlur handlers that should be attached to the input control
//   - lostFocus reactive boolean that is true when lost-focus validation
//     should occur (see below)
//   - submitted reactive boolean that is true when the user has attempted to
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

// ErrorFunction is a function that returns an error message string for a form
// control, given some flags describing the form state.
export type ErrorFunction = (lostFocus: boolean, submitted: boolean) => string
