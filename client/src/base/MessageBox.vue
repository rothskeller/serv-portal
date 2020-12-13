<!--
MessageBox displays a message box as a modal dialog.  Pop it up by calling its
show() method through a reference to it.  The show() method returns a Promise
that resolves to true if the user clicks OK or false if the user clicks Cancel.
Use the title prop to set the message box title, the variant prop to set the
color of the title bar and the OK button, and the okLabel and cancelLabel props
to set the labels on the OK and Cancel buttons.
-->

<template lang="pug">
Modal(ref='modal', v-slot='{ close }')
  SForm(
    dialog,
    :variant='variant',
    :title='title',
    :submitLabel='okLabel',
    :cancelLabel='cancelLabel',
    @submit='close(true)',
    @cancel='close(false)'
  )
    .form-item.mbox-message
      slot
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import Modal from './Modal.vue'
import SForm from './form/form'

export default defineComponent({
  components: { Modal, SForm },
  props: {
    title: { type: String },
    variant: { type: String, default: 'primary' },
    okLabel: { type: String, default: 'OK' },
    cancelLabel: { type: String, default: 'Cancel' },
  },
  setup() {
    const modal = ref(null as any)
    function show() {
      if (!modal.value) throw new Error('showing MessageBox2 with no reference to modal')
      return modal.value.show()
    }
    return { modal, show }
  },
})
</script>

<style lang="postcss">
.mbox-message {
  padding: 0.75rem;
}
</style>
