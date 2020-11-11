// sfvalidate.ts is a mixin for form controls allowing them to plug into the
// form validation mechanism.  From their setup() function, controls should call
// its default export function, passing the id for the control and a reactive
// reference to the error message string (which is assumed to be truthy when the
// control value is invalid and falsey when it is valid).
//
// The default export function returns an object with three keys: onBlur,
// lostFocus, and submitted.  onBlur is a function that should be called by the
// form control when it loses focus.  lostFocus and submitted are both reactive
// references to Boolean flags; they indicate when the control has lost focus
// (i.e., has had focus at some time but doesn't have it now), and when the user
// has attempted to submit the form.  These flags can be used by the form parent
// while computing the error message string.
//
// This module also exports an ErrorFunction type, which is used by most form
// controls that use this mixin.  It's here only for convenience.

import { inject, ref, Ref, watch } from "vue"

export type ErrorFunction = (lostFocus: boolean, submitted: boolean) => string

export default function provideValidation(id: string, error: Ref<string>) {
  // submitted is a reference to a flag indicating whether the user has ever
  // attempted to submit the form.
  const submitted = inject<Ref<boolean>>('formSubmitted')!

  // lostFocus is a reference to a flag indicating whether the control has lost
  // focus, i.e., it had focus at one time and no longer has it.
  const lostFocus = ref(false)
  const onFocus = () => { lostFocus.value = false }
  const onBlur = () => { lostFocus.value = true }

  // Update lostFocus when submitted flag changes.
  watch(submitted, () => { lostFocus.value = submitted.value })

  // Notify form when validity changes.
  const setValidity = inject<(id: string, isValid: boolean) => void>('setValidity')!
  watch(error, () => { setValidity(id, !error.value) })

  return { onFocus, onBlur, submitted, lostFocus }
}