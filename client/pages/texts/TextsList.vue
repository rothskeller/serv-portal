<!--
TextsList displays the list of previously sent text messages.
-->

<template lang="pug">
#texts-list
  #texts-list-spinner(v-if="loading")
    b-spinner(small)
  #texts-list-table(v-else)
    .texts-list-timestamp.texts-list-heading Time Sent
    .texts-list-groups.texts-list-heading Recipients
    .texts-list-message.texts-list-heading Message
    template(v-for="m in messages")
      .texts-list-timestamp
        b-link(:to="`/texts/${m.id}`" v-text="m.timestamp")
      .texts-list-groups
        .texts-list-group(v-for="g in m.groups" v-text="'CERT Team Alpha'")
      .texts-list-message(v-text="m.message")
</template>

<script>
export default {
  data: () => ({
    messages: null,
    loading: true,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/sms`)).data
    this.messages = data.messages
    this.loading = false
  },
}
</script>

<style lang="stylus">
#texts-list
  margin 0.75rem
#texts-list-table
  display flex
  flex-wrap wrap
  line-height 1.2
.texts-list-heading
  font-weight bold
  &:nth-child(3)
    display none
    @media (min-width: 700px)
      display block
.texts-list-timestamp
  flex none
  margin-top 0.75rem
  width 10rem
  font-variant tabular-nums
.texts-list-groups
  flex none
  margin-top 0.75rem
  width calc(100vw - 11.5rem)
  white-space nowrap
  @media (min-width: 576px)
    width calc(100vw - 18.5rem)
  @media (min-width: 700px)
    width 10rem
.texts-list-group
  overflow hidden
  text-overflow ellipsis
.texts-list-message
  padding-left 4em
  width calc(100vw - 1.5rem)
  font-size 0.75rem
  @media (min-width: 700px)
    padding-top 0.75rem
    padding-left 0
    width calc(100vw - 28.5rem)
    font-size 1rem
</style>
