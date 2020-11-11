<!--
SFTimeRange is a pair of input controls in an SForm that collect a starting and
ending time of a range.
-->

<template lang="pug">
label.form-item-label.sftimerange-label(:for='id', v-text='label')
.sftimerange.form-item-input
  input.form-control(
    :id='id',
    :class='{ "form-control-invalid": error }',
    type='time',
    v-model='start',
    @focus='onFocus1',
    @blur='onBlur1'
  )
  span.sftimerange-span to
  input.form-control(
    :class='{ "form-control-invalid": error }',
    type='time',
    v-model='end',
    @focus='onFocus1',
    @blur='onBlur1'
  )
.form-item-help.sftimerange-helpbox
  .form-item-error-text(v-if='error', v-text='error')
  .form-item-help-text(v-if='help', v-text='help')
</template>

<script lang="ts">
import { defineComponent, ref, toRefs, watch } from 'vue'
import provideValidation from './sfvalidate'

export default defineComponent({
  props: {
    id: { type: String, required: true },
    label: String,
    help: String,
    start: { type: String, required: true },
    end: { type: String, required: true },
  },
  emits: ['update:start', 'update:end'],
  setup(props, { emit }) {
    // Get the initial value from the props.  Update our local value whenever
    // the props change.
    const { start: propStart, end: propEnd } = toRefs(props)
    const start = ref(propStart.value)
    const end = ref(propEnd.value)
    watch(propStart, () => {
      end.value = propEnd.value
    })
    watch(propEnd, () => {
      end.value = propEnd.value
    })

    // Set up for form control validation.
    const error = ref('')
    provideValidation(props.id, error)

    // Keep track of when either of the inputs has focus.
    const everFocused = ref(false)
    const focused = ref(null as null | EventTarget)
    function onFocus1(evt: FocusEvent) {
      focused.value = evt.target
      everFocused.value = true
    }
    function onBlur1(evt: FocusEvent) {
      if (evt.target === focused.value) focused.value = null
    }

    // Watch for local changes and send them to the parent.
    watch(start, () => {
      emit('update:start', start.value)
    })
    watch(end, () => {
      emit('update:end', end.value)
    })

    // Validate the input.
    watch([start, end, focused, everFocused], () => {
      if (!everFocused.value || focused.value) error.value = ''
      else if (!start.value || !end.value) error.value = 'The start and end times are required.'
      else if (!start.value.match(/^(?:[01][0-9]|2[0-3]):[0-5][0-9]$/))
        error.value = 'The start time is not valid.'
      else if (!end.value.match(/^(?:[01][0-9]|2[0-3]):[0-5][0-9]$/))
        error.value = 'The end time is not valid.'
      else if (end.value < start.value) error.value = 'The end time must come after the start time.'
      else error.value = ''
    })

    return { start, end, onFocus1, onBlur1, error }
  },
})
</script>

<style lang="postcss">
.sftimerange-label {
  padding-top: calc(0.375rem + 1px);
}
.sftimerange {
  display: flex;
  align-items: baseline;
}
.sftimerange-span {
  padding: 0 0.5rem;
}
.form-l .sftimerange-helpbox {
  min-height: calc(1.5rem + 0.75rem + 2px);
}
</style>
