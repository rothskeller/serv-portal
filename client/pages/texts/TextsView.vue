<!--
TextsView displays the status of an outgoing text message, including its
deliveries and replies.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!message")
  b-spinner(small)
#texts-view(v-else)
  #texts-view-meta
    .texts-view-meta-label Message sent
    .texts-view-meta(v-text="message.timestamp")
    .texts-view-meta-label Sent by
    .texts-view-meta(v-text="message.sender")
    .texts-view-meta-label Sent to
    .texts-view-meta
      div(v-for="g in message.groups" v-text="g")
    .texts-view-meta-label Message text
    .texts-view-meta(v-text="message.message")
  #texts-view-deliveries
    TextsViewDelivery(v-for="d in message.deliveries" :key="d.id" :delivery="d")
</template>

<script>
import moment from 'moment-mini'
import TextsViewDelivery from './TextsViewDelivery'

export default {
  components: { TextsViewDelivery },
  data: () => ({ message: null }),
  created() {
    this.load()
  },
  methods: {
    async load() {
      this.message = (await this.$axios.get(`/api/sms/${this.$route.params.id}`)).data
      this.message.deliveries.sort((a, b) => (a.recipient < b.recipient ? -1 : a.recipient > b.recipient ? +1 : 0))
      window.setTimeout(this.load, 5000)
    },
  }
}
</script>

<style lang="stylus">
#texts-view
  margin 0.75rem
  @media (min-width: 752px)
    display flex
    flex-direction column
    overflow auto
    height calc(100vh - 4.75rem - 41px)
    // 40px title bar
    // 3.25rem + 1px tab bar
    // 1.5rem tab margin
#texts-view-meta
  display flex
  flex none
  flex-wrap wrap
  line-height 1.2
.texts-view-meta-label
  width 8rem
.texts-view-meta
  width calc(100vw - 9.5rem)
  @media (min-width: 576px)
    width calc(100vw - 16.5rem)
#texts-view-deliveries
  margin-top 1rem
  @media (min-width: 752px)
    display flex
    flex 1 1 0
    flex-direction column
    flex-wrap wrap
    align-content flex-start
    min-height 1rem
</style>
