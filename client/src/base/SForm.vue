<!--
SForm is a container for a form.  It handles styling and validation of the form.
-->

<template lang="pug">
form.form(:class='layoutClass', @submit.prevent='onSubmitInner')
  .form-item.form-title(v-if='title', :class='`form-title-${variant}`', v-text='title')
  slot
  .form-item
    slot(name='buttons')
      .form-buttons
        slot(v-if='dialog', name='extraButtons')
        span(:class='{ "form-buttons-1": hasExtraButtons }')
          SButton(
            v-if='dialog',
            @click='onCancelInner',
            variant='secondary',
            :disabled='disabled',
            v-text='cancelLabel'
          )
          SButton(type='submit', :variant='variant', :disabled='disabled', v-text='submitLabel')
          SButton(
            v-if='!dialog && cancelLabel',
            @click='onCancelInner',
            variant='secondary',
            :disabled='disabled',
            v-text='cancelLabel'
          )
        slot(v-if='!dialog', name='extraButtons')
  .form-item(v-if='hasFeedback')
    slot(name='feedback')
</template>

<script lang="ts">
import {
  defineComponent,
  inject,
  computed,
  ref,
  provide,
  onUpdated,
  nextTick,
  getCurrentInstance,
  onMounted,
} from 'vue'
import { useRouter } from 'vue-router'
import provideSize, { Size } from '../plugins/size'
import SButton from './SButton.vue'
import './sfcontrol.css'

export default defineComponent({
  components: { SButton },
  props: {
    submitLabel: { type: String, default: 'Save' },
    cancelLabel: { type: String, default: 'Cancel' },
    dialog: { type: Boolean, default: false },
    title: String,
    variant: { type: String, default: 'primary' },
    disabled: { type: Boolean, default: false },
  },
  emits: ['cancel', 'submit'],
  setup(props, { attrs, emit, slots }) {
    const router = useRouter()

    // Were we given a cancel handler?
    let haveCancelHandler = false
    onMounted(() => {
      if (getCurrentInstance()!.vnode.props!.onCancel) haveCancelHandler = true
    })

    // Set our form layout based on our container size.
    const size: Size = provideSize()
    const layoutClass = computed(() => {
      if (props.dialog && size.w >= 30) return 'form-dialog-m'
      if (props.dialog) return 'form-dialog-xs'
      if (size.w >= 51) return 'form-l'
      if (size.w >= 30) return 'form-m'
      if (size.w >= 20) return 'form-s'
      return 'form-xs'
    })

    // invalidFields is a set of IDs of form controls with invalid contents.
    // valid is a reactive flag indicating whether there are any form controls
    // with invalid contents.  A setValidity function is provided so that form
    // controls can add and remove themselves from the invalidFields set.
    const invalidFields: Set<string> = new Set()
    const valid = ref(true)
    function setValidity(id: string, isValid: boolean) {
      if (!isValid) invalidFields.add(id)
      else invalidFields.delete(id)
      valid.value = invalidFields.size === 0
    }
    provide('setValidity', setValidity)

    // submitted is a reactive flag indicating that the user has attempted to
    // submit the form.  It is provided to form controls because some of them
    // don't perform full validation until submit is attempted.  A reset method
    // is provided for the form parent to use when desired.
    const submitted = ref(false)
    provide('formSubmitted', submitted)
    function resetSubmitted() {
      submitted.value = false
    }

    // On cancel, we route back to the previous page.
    function onCancelInner() {
      if (haveCancelHandler) emit('cancel')
      else router.go(-1)
    }

    // On submit, we set the flag indicating that submit of the form has been
    // attempted, wait for the reactive value of the valid flag to be adjusted
    // as a result, and then, if the form is valid, emit the submit event to the
    // form parent.
    const onSubmitInner = async () => {
      if (!submitted.value) {
        submitted.value = true
        await nextTick()
      }
      if (valid.value) emit('submit')
    }

    // It's better not to show the feedback form-input row if there is no form
    // feedback to display in it. Ditto the extra-buttons span.
    const hasFeedback = ref(!!slots.feedback)
    const hasExtraButtons = ref(!!slots.extraButtons)
    onUpdated(() => {
      hasFeedback.value = !!slots.feedback
      hasExtraButtons.value = !!slots.extraButtons
    })

    return {
      layoutClass,
      hasFeedback,
      hasExtraButtons,
      onCancelInner,
      onSubmitInner,
      valid,
      resetSubmitted,
    }
  },
})
</script>

<style lang="postcss">
.form {
  display: grid;
  margin: 1.5rem 0.75rem;
  &.form-dialog-m,
  &.form-dialog-xs {
    margin: 0;
  }
}
.form-title {
  font-size: 1.25rem;
  font-weight: 500;
  padding: 0.75rem;
  color: #fff;
  margin-bottom: 1.5rem;
}
.form-title-primary {
  background-color: #007bff;
}
.form-title-danger {
  background-color: #dc3545;
}
.form-title-warning {
  background-color: #ffc107;
}
.form-item-label {
  grid-column: 1;
  margin-right: 1rem;
  color: #212529;
  white-space: nowrap;
  &:empty {
    margin-right: 0;
  }
}
.form-l {
  grid: auto-flow / max-content 20rem 21rem;
  & .form-item-label {
    margin-bottom: 1rem;
  }
  & .form-item-input {
    grid-column: 2;
    margin-bottom: 1rem;
  }
  & .form-item-input2 {
    grid-column: 2 / 4;
    margin-bottom: 1rem;
  }
  & .form-item-help {
    grid-column: 3;
    margin-left: 1rem;
    margin-bottom: 1rem;
  }
  & .form-item {
    grid-column: 1 / 4;
    max-width: 51rem;
  }
}
.form-m {
  grid: auto-flow / min-content 20rem;
  & .form-item-label {
    grid-row-end: span 2;
    margin-bottom: 1rem;
  }
  & .form-item-input {
    grid-column: 2;
  }
  & .form-item-input2 {
    grid-column: 2;
    margin-bottom: 1rem;
  }
  & .form-item-help {
    grid-column: 2;
    margin-top: 0.25rem;
    margin-bottom: 1rem;
  }
  & .form-item {
    grid-column: 1 / 3;
    max-width: 30rem;
  }
}
.form-dialog-m {
  grid: auto-flow / min-content 1fr;
  & .form-item-label {
    grid-row-end: span 2;
    margin-left: 0.75rem;
    margin-bottom: 1rem;
  }
  & .form-item-input {
    grid-column: 2;
    margin-right: 0.75rem;
  }
  & .form-item-input2 {
    grid-column: 2;
    margin-bottom: 1rem;
    margin-right: 0.75rem;
  }
  & .form-item-help {
    grid-column: 2;
    margin-top: 0.25rem;
    margin-bottom: 1rem;
    margin-right: 0.75rem;
  }
  & .form-item {
    grid-column: 1 / 3;
  }
}
.form-s {
  grid: auto-flow / 20rem;
  & .form-item-input2 {
    margin-bottom: 1rem;
  }
  & .form-item-help {
    margin-top: 0.25rem;
    margin-bottom: 1rem;
  }
}
.form-xs {
  grid: auto-flow / 100%;
  & .form-item-input2 {
    margin-bottom: 1rem;
  }
  & .form-item-help {
    margin-top: 0.25rem;
    margin-bottom: 1rem;
  }
}
.form-dialog-xs {
  grid: auto-flow / 100%;
  & .form-item-label {
    margin: 0 0.75rem;
  }
  & .form-item-input {
    margin: 0 0.75rem;
  }
  & .form-item-input2 {
    margin: 0 0.75rem;
  }
  & .form-item-help {
    margin: 0.25rem 0.75rem 1rem;
  }
}
.form-item-help {
  display: flex;
  flex-direction: column;
  justify-content: center;
  &:empty {
    margin-top: 0 !important;
  }
}
.form-buttons {
  display: flex;
  flex-wrap: wrap;
  margin: 1rem 0 0 -0.5rem;
  & .sbtn {
    margin: 0 0 0.5rem 0.5rem;
  }
  .form-dialog-m &,
  .form-dialog-xs & {
    justify-content: flex-end;
    padding: 0.75rem 0.25rem 0.25rem 0.75rem;
    border-top: 1px solid #dee2e6;
    & .sbtn {
      margin: 0 0.5rem 0.5rem 0;
    }
  }
}
.form-buttons-1 {
  margin: 0 2rem 0 0;
  .form-dialog & {
    margin: 0 0 0 2rem;
  }
}
.form-item-error-text {
  color: #dc3545;
  font-size: 80%;
  line-height: 1.2;
}
.form-item-help-text {
  color: #6c757d;
  font-size: 80%;
  line-height: 1.2;
}
</style>
