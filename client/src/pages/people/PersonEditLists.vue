<!--
PersonEditLists is the dialog box for editing a person's list subscriptions.
-->

<template lang="pug">
Modal(ref='modal')
  SubscriptionsForm(inModal, :pid='pid', @done='onDone')
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { Modal } from '../../base'
import SubscriptionsForm from './SubscriptionsForm.vue'

export default defineComponent({
  components: { Modal, SubscriptionsForm },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    function show() {
      return modal.value.show()
    }
    function onDone(submitted: boolean) {
      modal.value.close(submitted)
    }
    return { modal, onDone, show }
  },
})
</script>
