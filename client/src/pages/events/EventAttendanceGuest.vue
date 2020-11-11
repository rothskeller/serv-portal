<!--
EventAttendanceGuest
-->

<template lang="pug">
Modal(ref='modal', v-slot='{ close }')
  #event-attendance-guest-title Add Guest
  form#event-attendance-guest(@submit.prevent='onSubmit')
    SInput#event-attendance-guest-input(ref='input', placeholder='Guest Name', v-model='search')
    select(size='10', v-model='guest')
      option(v-for='o in options', :value='o.id', v-text='o.sortName')
    #event-attendance-guest-buttons
      SButton(@click='close(null)') Cancel
      SButton(type='submit', variant='primary') OK
</template>

<script lang="ts">
import { defineComponent, ref, watch, onMounted, nextTick } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SButton, SInput } from '../../base'

export type GuestOption = {
  id: number
  sortName: string
}

export default defineComponent({
  components: { Modal, SButton, SInput },
  setup(props, { emit }) {
    const modal = ref(null as any)
    const input = ref(null as any)
    function show() {
      function focus() {
        if (input.value) input.value.focus()
        else nextTick(focus)
      }
      focus()
      return modal.value.show()
    }

    // Handle searching.
    const search = ref('')
    const guest = ref(0)
    const options = ref([] as Array<GuestOption>)
    watch(search, async () => {
      if (!search.value.trim()) {
        options.value = []
        guest.value = 0
        return
      }
      options.value = (
        await axios.get('/api/people', { params: { search: search.value.trim() } })
      ).data
      guest.value = options.value.length === 1 ? options.value[0].id : 0
    })

    // Submit
    function onSubmit() {
      if (!guest.value) return
      modal.value.close(options.value.find((g) => g.id === guest.value))
    }

    return { guest, input, modal, onSubmit, options, search, show }
  },
})
</script>

<style lang="postcss">
#event-attendance-guest-title {
  font-size: 1.25rem;
  font-weight: 500;
  padding: 0.75rem;
  color: #fff;
  background-color: #007bff;
}
#event-attendance-guest {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  padding: 0.75rem;
}
#event-attendance-guest-input {
  margin-bottom: 0.25rem;
}
#event-attendance-guest-buttons {
  padding-top: 0.75rem;
  border-top: 1px solid rgba(0, 0, 0, 0.2);
  text-align: right;
  & .sbtn {
    margin-right: 0.5rem;
  }
}
</style>
