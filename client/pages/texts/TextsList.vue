<!--
TextsList displays the list of previously sent text messages.
-->

<template lang="pug">
#texts-list
  #texts-list-spinner(v-if="loading")
    b-spinner(small)
  #texts-list-table(v-else)
    .texts-list-date.texts-list-heading Date
    .texts-list-event.texts-list-heading Recipients
    .texts-list-location.texts-list-heading Message
    template(v-for="m in messages")
      .texts-list-timestamp
        b-link(:to="`/texts/${m.id}`" v-text="m.timestamp")
      .texts-list-groups
        div(v-for="g in m.groups" v-text="g.name")
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
  margin 1.5rem 0.75rem
#texts-list-table
  display flex
  flex-wrap wrap
.texts-list-heading
  display none
  @media (min-width: 576px)
    display block
    font-weight bold
.texts-list-timestamp
  flex none
  margin-top 0.25rem
  width 10rem
  font-variant tabular-nums
.texts-list-groups
  flex none
  margin-top 0.25rem
.texts-list-message
  display none
</style>
