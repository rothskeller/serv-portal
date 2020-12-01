<!--
TextsSend sends a new text message.
-->

<template lang="pug">
#texts-send(v-if='!lists.length')
  SSpinner
SForm#texts-send(v-else, @submit='onSubmit', :submitLabel='submitLabel', :disabled='disabled')
  SFTextArea#texts-send-message(
    label='Message',
    autofocus,
    rows='5',
    trim,
    v-model='message',
    :help='countMessage',
    :errorFn='messageError'
  )
  SFCheckGroup#texts-send-lists(
    label='Recipients',
    :options='lists',
    valueKey='id',
    labelKey='name',
    v-model='recipients',
    :errorFn='listsError'
  )
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SForm, SFCheckGroup, SFTextArea, SSpinner } from '../../base'

type GetSMSNewList = {
  id: number
  name: string
}
type GetSMSNew = {
  lists: Array<GetSMSNewList>
}
type PostSMS = {
  id: number
}

export default defineComponent({
  components: { SForm, SFCheckGroup, SFTextArea, SSpinner },
  setup() {
    setPage({ title: 'New Text Message' })

    // Get the list of allowed SMS lists.
    const lists = ref([] as Array<GetSMSNewList>)
    axios.get<GetSMSNew>(`/api/sms/NEW`).then((resp) => {
      lists.value = resp.data.lists
    })

    // Message field.
    const message = ref('')
    const countMessage = ref('0/160')
    function messageError(lostFocus: boolean) {
      if (lostFocus && !message.value) return 'Please enter the text of your message.'
      let unicode = false
      let runes = 0
      let chars = 0
      for (const ch of message.value) {
        runes++
        if (
          '£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !"#¤%&\'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà'.includes(
            ch
          )
        )
          chars++
        else if ('\f\n^{}\\[~]|€'.includes(ch)) chars += 2
        else unicode = true
      }
      let max = unicode ? 70 : 160
      let nmsg = Math.ceil(runes / max)
      if (nmsg > 1) countMessage.value = `${runes}/${max} (${nmsg} messages)`
      else countMessage.value = `${runes}/${max}`
      return ''
    }

    // Recipients.
    const recipients = ref(new Set() as Set<number>)
    function listsError(lostFocus: boolean, submitted: boolean) {
      if (submitted && !recipients.value.size)
        return 'Please select the recipients of your message.'
      return ''
    }

    // Send the message.
    const sending = ref(false)
    const submitLabel = computed(() => (sending.value ? 'Sending...' : 'Send Message'))
    const disabled = computed(() => sending.value)
    const router = useRouter()
    async function onSubmit() {
      sending.value = true
      const body = new FormData()
      body.append('message', message.value)
      recipients.value.forEach((r) => body.append('list', r.toString()))
      const resp = await axios.post<PostSMS>(`/api/sms`, body)
      router.push(`/texts/${resp.data.id}`)
    }

    return {
      countMessage,
      disabled,
      lists,
      listsError,
      message,
      messageError,
      onSubmit,
      recipients,
      submitLabel,
    }
  },
})
</script>

<style lang="postcss">
#texts-send {
  padding: 1.5rem 0.75rem;
}
.texts-send-label {
  width: 7rem;
}
#texts-send-message {
  min-width: 12rem;
  max-width: 20rem;
}
.texts-send-long {
  color: red !important;
}
</style>
