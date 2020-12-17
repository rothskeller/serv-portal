// Form framework with validation support.
//
// Children of the form must have one of the following forms:
//   .form-item-label + .form-item-input + .form-item-help
//   .form-item-label + .form-item-input2
//   .form-item
//
// Three named slots are supported in addition to the default.
//   extraButtons allows adding extra buttons to the standard form button bar.
//   buttons replaces the standard form button bar entirely.
//   feedback provides feedback at the bottom of the form.
//
// Unless overridden by a "buttons" named slot, the standard form buttons are
// Save and Cancel.  The Save button is displayed according to the variant
// property; the Cancel button always uses the "secondary" variant.  The button
// labels can be changed with the submitLabel and cancelLabel properties.  The
// Cancel button is omitted if cancelLabel is set to an empty string.  Both
// buttons are disabled if the disabled property is true.
//
// If the dialog property is true, the form is displayed appropriately for a
// modal dialog box.  Different breakpoint sizes are used, and button placement
// is different.  When this property is true, the title property should also be
// set; it controls what appears in the title bar for the dialog.  The title bar
// color is determined by the variant property.  The title property should not
// be set if the dialog property is false.
//
// The form emits a 'submit' event if the user submits the form with no
// validation errors reported by the controls within it.  If an event handler is
// registered for the 'cancel' event, the form emits that when the user cancels
// the form; otherwise, the form forces a router "back" operation when the user
// cancels the form.
//
// Controls within the form can participate in validation by injecting the
// setValidity function provided by this form, and calling it to indicate when
// the validity of the control changes.  Its parameters are the ID of the
// control (which must be unique) and a Boolean indicating whether the control
// has valid data.
//
// Controls may want to know whether the user has attempted to submit the form,
// in order to apply more rigorous validation.  They can inject the
// formSubmitted Ref<boolean> for this purpose.  Form parents that need to reset
// the submitted state can do so by calling the resetSubmitted method of the
// form.

import { computed, defineComponent, getCurrentInstance, h, nextTick, provide, ref } from 'vue'
import { useRouter } from 'vue-router'
import provideSize from '../../plugins/size'
import SButton from '../controls/SButton'
import './form.css'

const SForm = defineComponent({
  name: 'SForm',
  props: {
    submitLabel: { type: String, default: 'Save' },
    cancelLabel: { type: String, default: 'Cancel' },
    disabled: { type: Boolean, default: false },
    variant: { type: String, default: 'primary' },
    dialog: { type: Boolean, default: false },
    title: { type: String },
  },
  emits: ['cancel', 'submit'],
  setup(props, { emit, expose, slots }) {
    const router = useRouter()
    const instance = getCurrentInstance()!

    // Set our form layout based on our container size.
    const size = provideSize()
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
    expose({ resetSubmitted })

    // On submit, we set the flag indicating that submit of the form has been
    // attempted, wait for the reactive value of the valid flag to be adjusted
    // as a result, and then, if the form is valid, emit the submit event to the
    // form parent.
    const onSubmit = async (evt: Event) => {
      evt.preventDefault()
      if (!submitted.value) {
        submitted.value = true
        await nextTick()
      }
      if (valid.value) emit('submit')
    }

    // On cancel, if our caller registered a cancel handler, we call it;
    // otherwise we navigate back to the previous page.
    function onCancel() {
      if (instance.vnode.props?.onCancel) emit('cancel')
      else router.go(-1)
    }

    // renderButtons renders the buttons for the form.
    function renderButtons() {
      if (slots.buttons) return slots.buttons()
      let children: Array<any> = [
        h(
          SButton,
          {
            type: 'submit',
            variant: props.variant,
            disabled: props.disabled,
          },
          () => props.submitLabel
        ),
      ]
      if (props.cancelLabel) {
        const cancel = h(
          SButton,
          {
            variant: 'secondary',
            disabled: props.disabled,
            onClick: onCancel,
          },
          () => props.cancelLabel
        )
        if (props.dialog) children.unshift(cancel)
        else children.push(cancel)
      }
      if (slots.extraButtons) {
        children = [h('span', { class: 'form-buttons-1' }, children)]
        if (props.dialog) children.unshift(slots.extraButtons())
        else children.push(slots.extraButtons())
      }
      return h(`div`, { class: 'form-buttons' }, children)
    }

    return () => {
      const children: Array<any> = []
      if (props.title)
        children.push(
          h(
            'div',
            {
              class: ['form-item', 'form-title', `form-title-${props.variant}`],
            },
            props.title
          )
        )
      children.push(slots.default!())
      children.push(h('div', { class: 'form-item' }, renderButtons()))
      if (slots.feedback) children.push(h('div', { class: 'form-item' }, slots.feedback()))
      return h('form', { class: ['form', layoutClass.value], onSubmit }, children)
    }
  },
})
export default SForm
