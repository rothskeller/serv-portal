<!--
TextsView displays the status of an outgoing text message, including its
deliveries and replies.
-->

<template lang="pug">
#texts-view-spinner(v-if='!message')
  SSpinner
#texts-view(v-else)
  #texts-view-meta
    .texts-view-meta-label Message sent
    .texts-view-meta(v-text='message.timestamp')
    .texts-view-meta-label Sent by
    .texts-view-meta(v-text='message.sender')
    .texts-view-meta-label Sent to
    .texts-view-meta
      div(v-for='l in message.lists', v-text='l')
    .texts-view-meta-label Message text
    .texts-view-meta(v-text='message.message')
  #texts-view-deliveries
    .texts-view-name-num.texts-view-heading Recipient
    .texts-view-status-time.texts-view-heading Status
    .texts-view-responses.texts-view-heading Reply
    template(v-for='d in message.deliveries')
      .texts-view-name-num(@click='onClick(d)')
        .texts-view-recipient(v-text='d.recipient')
        .texts-view-number(v-text='formatNumber(d)')
      .texts-view-status-time
        .texts-view-status(:class='statusColor(d)', v-text='formatStatus(d)')
        .texts-view-time(v-text='formatTimestamp(d)')
      .texts-view-responses
        div(v-for='r in d.responses')
          span.texts-view-response-time(v-text='r.timestamp.substr(11, 8)')
          span(v-text='r.response')
</template>

<script lang="ts">
import { defineComponent, ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../../plugins/axios'
import { SSpinner } from '../../base'
import setPage from '../../plugins/page'

type GetSMS1Response = {
  response: string
  timestamp: string
}
type GetSMS1Delivery = {
  id: number
  recipient: string
  number: string
  status: string
  timestamp: string
  responses: Array<GetSMS1Response>
}
type GetSMS1 = {
  id: number
  sender: string
  lists: Array<string>
  timestamp: string
  deliveries: Array<GetSMS1Delivery>
}

export default defineComponent({
  components: { SSpinner },
  setup() {
    setPage({ title: 'Text Message' })

    // Load the message, and reload it every 5 seconds.
    const route = useRoute()
    const message = ref(null as null | GetSMS1)
    async function load() {
      message.value = (await axios.get<GetSMS1>(`/api/sms/${route.params.id}`)).data
      message.value.deliveries.sort((a, b) =>
        a.recipient < b.recipient ? -1 : a.recipient > b.recipient ? +1 : 0
      )
      window.setTimeout(load, 5000)
    }
    watchEffect(() => {
      load()
    })

    // When someone clicks on a delivery, on a phone, start a text to that
    // recipient's number.
    function onClick(d: GetSMS1Delivery) {
      if (!navigator.userAgent.match(/Android/i) && !navigator.userAgent.match(/iPhone/i)) return
      if (!d.number) return
      window.location.href = `sms:${d.number.substr(2)}`
    }

    // Formatting functions.
    function formatNumber(d: GetSMS1Delivery) {
      if (!d.number) return ''
      return `${d.number.substr(2, 3)}-${d.number.substr(5, 3)}-${d.number.substr(8, 4)}`
    }
    function formatStatus(d: GetSMS1Delivery) {
      if (d.responses && d.responses.length) return 'Replied'
      switch (d.status) {
        case 'queued':
          return 'Queued'
        case 'sending':
          return 'Sending'
        case 'sent':
          return 'Sent'
        case 'delivered':
          return 'Delivered'
        case 'undelivered':
          return 'Not Delivered'
        case 'failed':
          return 'Failed'
        case 'No Cell Phone':
          return 'No Cell Phone'
        default:
          return 'Pending'
      }
    }
    function formatTimestamp(d: GetSMS1Delivery) {
      if (d.responses && d.responses.length)
        return d.responses[d.responses.length - 1].timestamp.substr(11, 8)
      if (!d.status || d.status === 'No Cell Phone') return ''
      return d.timestamp.substr(11, 8)
    }
    function statusColor(d: GetSMS1Delivery) {
      if (d.responses && d.responses.length) return 'texts-view-delivery-replied'
      switch (d.status) {
        case 'queued':
          return 'texts-view-delivery-pending'
        case 'sent':
          return 'texts-view-delivery-sent'
        case 'sending':
          return 'texts-view-delivery-pending'
        case 'delivered':
          return 'texts-view-delivery-delivered'
        case 'undelivered':
          return 'texts-view-delivery-failed'
        case 'failed':
          return 'texts-view-delivery-failed'
        case 'No Cell Phone':
          return 'texts-view-delivery-failed'
        default:
          return 'texts-view-delivery-pending'
      }
    }

    return { formatNumber, formatStatus, formatTimestamp, message, onClick, statusColor }
  },
})
</script>

<style lang="postcss">
#texts-view-spinner {
  margin: 1.5rem 0.75rem;
}
#texts-view {
  display: flex;
  flex-direction: column;
  padding: 0.75rem;
  height: 100%;
}
#texts-view-meta {
  display: grid;
  flex: none;
  line-height: 1.2;
  grid: auto / 8rem 1fr;
}
#texts-view-deliveries {
  display: grid;
  justify-content: start;
  align-items: start;
  margin-top: 1rem;
  min-width: 0;
  width: 100%;
  line-height: 1.2;
  grid: auto / 1fr 7rem;
  @media (min-width: 328px) {
    grid: auto / 12rem 7rem;
  }
  @media (min-width: 576px) {
    grid: auto / 12rem 7rem 1fr;
  }
}
.texts-view-heading {
  font-weight: bold;
  &:nth-child(3) {
    display: none;
    @media (min-width: 576px) {
      display: block;
    }
  }
}
.texts-view-delivery-pending {
  color: #888;
}
.texts-view-delivery-failed {
  background-color: #e6194b;
  color: white;
}
.texts-view-delivery-delivered {
  background-color: #808000;
  color: white;
}
.texts-view-delivery-replied {
  background-color: #3cb44b;
  color: white;
}
.texts-view-delivery-sent {
  background-color: #9a6324;
  color: white;
}
.texts-view-name-num {
  overflow: hidden;
  margin-top: 0.75rem;
  padding-right: 0.5rem;
}
.texts-view-recipient {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.texts-view-number {
  color: #888;
  font-size: 0.75rem;
}
.texts-view-status-time {
  margin-top: 0.75rem;
}
.texts-view-status {
  padding-left: 2px;
}
.texts-view-time {
  color: #888;
  font-size: 0.75rem;
}
.texts-view-responses {
  margin-left: 4rem;
  text-indent: -2rem;
  grid-column: 1 / 3;
  @media (min-width: 576px) {
    margin-top: 0.75rem;
    margin-left: 1.5rem;
    text-indent: -1rem;
    grid-column: 3 / 4;
  }
}
.texts-view-response-time {
  display: none;
  color: #888;
  font-variant: tabular-nums;
  font-size: 0.75rem;
  &::before {
    content: '[';
  }
  &::after {
    content: '] ';
  }
  @media (min-width: 800px) {
    display: inline;
  }
}
</style>
