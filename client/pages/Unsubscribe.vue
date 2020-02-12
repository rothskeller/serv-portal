<!--
Unsubscribe displays the Unsubscribe page.
-->

<template lang="pug">
PublicPage(title="Sunnyvale SERV")
  div.mt-3.ml-2(v-if="loading")
    b-spinner(small)
  div.mt-3.mx-2(v-else-if="submitted")
    #unsub-head Unsubscribe
    #unsub-intro(v-if="noEmail").
      You have been removed from all of our email lists.
    #unsub-intro(v-else).
      You have been removed from the selected email lists.
    #unsub-warn.
      We would appreciate it if you’d drop us a note at
      <a href="mailto:admin@sunnyvaleserv.org">admin@sunnyvaleserv.org</a>
      and let us know why you unsubscribed.
    #unsub-warn.
      If you ever want to get back on the email lists, come back to this page
      and let us know.
  div.mt-3.mx-2(v-else)
    #unsub-head Unsubscribe
    #unsub-intro Which email list(s) do you want to unsubscribe from?
    b-checkbox(v-for="g in groups" :key="g.id" v-model="g.unsub") {{g.email}}@SunnyvaleSERV.org
    b-checkbox.mt-1(v-model="noEmail") All Sunnyvale SERV email lists
    #unsub-warn.
      Please note that, if you unsubscribe from a critical mailing list for one
      of our volunteer groups, you may no longer be able to participate in that
      group.
    #unsub-buttons
      b-btn(variant="primary" @click="onSubmit") Unsubscribe
</template>

<script>
import PublicPage from '@/base/PublicPage'

export default {
  components: { PublicPage },
  data: () => ({
    loading: false,
    noEmail: false,
    groups: null,
    submitted: false,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/unsubscribe/${this.$route.params.token}`)).data
    this.noEmail = data.noEmail
    this.groups = data.groups
    this.loading = false
    if (this.$route.params.email) {
      this.groups.forEach(g => {
        if (g.email === this.$route.params.email) g.unsub = true
      })
    }
  },
  methods: {
    async onSubmit() {
      const body = new FormData
      body.append('noEmail', this.noEmail)
      this.groups.forEach(g => {
        body.append(`unsub:${g.id}`, g.unsub)
      })
      await this.$axios.post(`/api/unsubscribe/${this.$route.params.token}`, body)
      this.submitted = true
    },
  },
}
</script>

<style lang="stylus">
#unsub-head
  font-weight bold
  font-size 1.25rem
#unsub-intro
  margin 0.5rem 0
#unsub-warn
  margin 0.5rem 0 0.75rem
  max-width 40rem
  line-height 1.2
</style>
