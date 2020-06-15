<!--
TextsSend sends a new text message.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!groups")
  b-spinner(small)
b-form#texts-send(v-else @submit.prevent="onSubmit")
  b-form-group(label="Message" label-for="texts-send-message" label-cols-sm="auto" label-class="texts-send-label" :state="messageError ? false: null" :invalidFeedback="messageError")
    b-textarea#texts-send-message(v-model="message" rows="5" autofocus)
    b-form-text(v-if="countMessage" :class="countClass") {{countMessage}}
  b-form-group(label="Recipients" label-for="texts-send-groups" label-cols-sm="auto" label-class="texts-send-label pt-0" :state="groupsError ? false : null" :invalidFeedback="groupsError")
    b-form-checkbox-group#texts-send-groups(v-model="recipients" :options="groups" stacked text-field="name" value-field="id")
  div.mt-3
    b-button(type="submit" variant="primary" :disabled="sending || !valid")
      b-spinner(v-if="sending" small)
      span(v-else) Send Message
</template>

<script>
export default {
  data: () => ({
    groups: null,
    message: '',
    recipients: [],
    countClass: '',
    countMessage: '0/160',
    messageError: null,
    groupsError: null,
    submitted: false,
    sending: false,
  }),
  async created() {
    const data = (await this.$axios.get('/api/sms/NEW')).data
    this.groups = data.groups
  },
  watch: {
    message: 'validate',
    recipients: 'validate',
  },
  computed: {
    valid() { return !this.messageError && !this.groupsError }
  },
  methods: {
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('message', this.message)
      this.recipients.forEach(r => body.append('group', r))
      this.sending = true
      const resp = (await this.$axios.post('/api/sms', body)).data
      this.sending = false
      this.$router.push(`/texts/${resp.id}`)
    },
    validate() {
      let unicode = false
      let runes = 0
      let chars = 0
      for (const ch of this.message) {
        runes++
        if ('£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !"#¤%&\'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà'.includes(ch))
          chars++
        else if ('\f\n^{}\\[~]|€'.includes(ch))
          chars += 2
        else
          unicode = true
      }
      if (this.submitted && !runes) {
        this.messageError = 'Please enter the text of your message.'
        this.countMessage = ''
        this.countClass = ''
      } else if (unicode) {
        this.messageError = ''
        this.countMessage = `${runes}/70`
        this.countClass = runes > 70 ? 'texts-send-long' : ''
      } else {
        this.messageError = ''
        this.countMessage = `${chars}/160`
        this.countClass = chars > 160 ? 'texts-send-long' : ''
      }
      if (this.submitted && !this.recipients.length)
        this.groupsError = 'Please select the recipients of your message.'
      else
        this.groupsError = null
    },
  },
}
</script>

<style lang="stylus">
#texts-send
  padding 1.5rem 0.75rem
.texts-send-label
  width 7rem
#texts-send-message
  min-width 12rem
  max-width 20rem
.texts-send-long
  color red
</style>
