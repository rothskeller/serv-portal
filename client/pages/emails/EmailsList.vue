<!--
EmailsList displays the list of emails.
-->

<template lang="pug">
#emails-list
  #emails-list-spinner(v-if="loading")
    b-spinner(small)
  .emails-list-email(v-for="email in emails" :key="email.id")
    template(v-if="email.from")
      .emails-list-heading From:
      div(v-text="email.from")
    template(v-if="email.subject")
      .emails-list-heading Subj:
      div(v-text="email.subject")
    template(v-if="email.to && email.to.length")
      div.emails-list-heading To:
      .emails-list-to
        div(v-for="group in email.to" v-text="group")
    template(v-if="email.timestamp")
      .emails-list-heading Date:
      div(v-text="email.timestamp.substr(0, 10) + ' ' + email.timestamp.substr(11, 8)")
    template(v-if="email.error")
      .emails-list-heading Error:
      .emails-list-error(v-text="email.error")
    template(v-else-if="email.type === 'bounce'")
      .emails-list-heading Error:
      .emails-list-error Bounce message
    template(v-else-if="email.type === 'send_failed'")
      .emails-list-heading Error:
      .emails-list-error Send failed
    template(v-else-if="email.type === 'moderated'")
      .emails-list-heading Error:
      .emails-list-error Message requires moderation
    template(v-else-if="email.type === 'unrecognized'")
      .emails-list-heading Error:
      .emails-list-error Message not recognized
    .emails-list-buttons
      b-btn(size="sm" variant="primary") View
      b-btn.ml-2(v-if="email.type === 'moderated'" size="sm" variant="success" @click="onAccept(email)") Accept
      b-btn.ml-2(v-if="email.type !== 'sent'" size="sm" variant="danger" @click="onDiscard(email)") Discard
      b-btn.ml-2(size="sm" variant="info" @click="onSendToMe(email)") Send to Me
</template>

<script>
export default {
  data: () => ({
    emails: null,
    loading: true,
  }),
  async created() {
    this.loading = true
    this.emails = (await this.$axios.get(`/api/emails`)).data
    this.loading = false
  },
  methods: {
    async onAccept(email) {
      await this.$axios.post(`/api/emails/${email.id}?action=accept`)
      this.emails = (await this.$axios.get(`/api/emails`)).data
    },
    async onDiscard(email) {
      await this.$axios.post(`/api/emails/${email.id}?action=discard`)
      this.emails = (await this.$axios.get(`/api/emails`)).data
    },
    async onSendToMe(email) {
      await this.$axios.post(`/api/emails/${email.id}?action=sendToMe`)
      this.$bvModal.msgBoxOk("This message has been sent to your primary email address.  Because it is being sent from a different mailbox than its From: line, it may be delivered to your Spam folder.")
    },
  },
}
</script>

<style lang="stylus">
#emails-list
  padding 1.5rem 0.75rem
#emails-list-spinner
  margin-top 1.5rem
.emails-list-email
  display grid
  margin-top 0.5rem
  padding-top 0.5rem
  border-top 1px solid #ccc
  grid auto / min-content 1fr
  &:first-child
    margin-top 0
    padding-top 0
    border-top none
.emails-list-heading
  padding-right 0.5rem
  color #888
.emails-list-to
  padding-top 0.15rem // compensate for line-height 1.5 on heading
  line-height 1.2
.emails-list-error
  color red
.emails-list-buttons
  margin-top 0.25rem
  grid-column 1 / 3
</style>
