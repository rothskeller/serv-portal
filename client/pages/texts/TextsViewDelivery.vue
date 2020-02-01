<!--
TextsViewDelivery displays the delivery record for one recipient of an outgoing
message.
-->

<template lang="pug">
.texts-view-delivery(:class="classes")
  .texts-view-delivery-meta
    .texts-view-delivery-name-num
      .texts-view-delivery-name(v-text="delivery.recipient")
      .texts-view-delivery-num(v-text="delivery.number")
    .texts-view-delivery-status-time
      .texts-view-delivery-status(:class="classes" v-text="formatStatus")
      .texts-view-delivery-time(v-text="formatTimestamp")
  .texts-view-delivery-responses(v-if="delivery.responses.length")
    .texts-view-delivery-response(v-for="r in delivery.responses")
      span.texts-view-delivery-response-time(v-text="r.timestamp.substr(11, 8)")
      span(v-text="r.response")
</template>

<script>
export default {
  props: {
    delivery: Object,
  },
  computed: {
    classes() {
      switch (this.delivery.status) {
        case 'queued': return 'texts-view-delivery-pending'
        case 'sent': return 'texts-view-delivery-sent'
        case 'sending': return 'texts-view-delivery-pending'
        case 'delivered': return 'texts-view-delivery-delivered'
        case 'undelivered': return 'texts-view-delivery-failed'
        case 'failed': return 'texts-view-delivery-failed'
        case 'No Cell Phone': return 'texts-view-delivery-failed'
        default: return 'texts-view-delivery-pending'
      }
    },
    formatStatus() {
      switch (this.delivery.status) {
        case 'queued': return 'Queued'
        case 'sending': return 'Sending'
        case 'sent': return 'Sent'
        case 'delivered': return 'Delivered'
        case 'undelivered': return 'Not Delivered'
        case 'failed': return 'Failed'
        case 'No Cell Phone': return 'No Cell Phone'
        default: return 'Pending'
      }
    },
    formatTimestamp() {
      if (!this.delivery.status || this.delivery.status === 'No Cell Phone') return ''
      return `as of ${this.delivery.timestamp.substr(11, 8)}`
    },
  },
}
</script>

<style lang="stylus">
.texts-view-delivery
  margin-top 0.5rem
  padding 0.25rem
  width 296px
  border 2px solid #ccc
  border-radius 8px
.texts-view-delivery-meta
  display flex
  justify-content space-between
  align-items flex-start
  line-height 1.2
.texts-view-delivery-name-num
  flex 1 1 0
  min-width 1px
.texts-view-delivery-name
  color black
  font-weight bold
.texts-view-delivery-num
  color #888
  font-size 0.75rem
.texts-view-delivery-status-time
  flex 1 1 0
  min-width 1px
  text-align right
.texts-view-delivery-status
  font-weight bold
.texts-view-delivery-time
  color #888
  font-size 0.75rem
.texts-view-delivery-pending
  border-color #888
  color #888
.texts-view-delivery-failed
  border-color #e6194b
  color #e6194b
.texts-view-delivery-delivered
  border-color #3cb44b
  color #3cb44b
.texts-view-delivery-sent
  border-color #9a6324
  color #9a6324
.texts-view-delivery-responses
  margin-top 0.25rem
  color black
.texts-view-delivery-response
  padding-left 2rem
  text-indent -2rem
.texts-view-delivery-response-time
  color #888
  &::before
    content '['
  &::after
    content '] '
</style>
