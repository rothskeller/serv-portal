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
        .texts-list-group(v-for="g in m.groups" v-text="g")
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
  padding 0.75rem
#texts-list-table
  display grid
  line-height 1.2
  grid-auto-columns 10rem 1fr
  @media (min-width: 700px)
    grid-auto-columns 10rem 10rem 1fr
.texts-list-heading
  font-weight bold
  &:nth-child(3)
    display none
    @media (min-width: 700px)
      display block
.texts-list-timestamp
  margin-top 0.75rem
  font-variant tabular-nums
.texts-list-groups
  margin-top 0.75rem
  white-space nowrap
.texts-list-group
  overflow hidden
  text-overflow ellipsis
.texts-list-message
  padding-left 4em
  font-size 0.75rem
  grid-column 1 / 3
  @media (min-width: 700px)
    padding-top 0.75rem
    padding-left 0
    font-size 1rem
    grid-column 3 / 4
</style>
