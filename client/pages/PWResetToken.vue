<!--
PWResetToken displays the second password reset page (the one the email links to).
-->

<template lang="pug">
PublicPage(title="Sunnyvale SERV")
  #pwreset-token-top
    #pwreset-token-banner Password Reset
    template(v-if="invalid")
      #pwreset-token-intro
        | This password reset link is invalid or has expired.
      #pwreset-token-reset
        a.btn.btn-secondary(href="/login/start-reset") Try Again
    div.mt-3.text-center(v-else-if="!hints")
      b-spinner(small)
    template(v-else)
      #pwreset-token-intro
        | Please provide a new password.
      form#pwreset-token-form(@submit.prevent="onSubmit")
        PasswordEntry(label="Password" label-class="pwreset-token-label" :deferValidation="!submitted" :allowBadPassword="false" :passwordHints="hints" @change="onPasswordChange")
        #pwreset-token-submit-row
          b-button(type="submit" variant="primary") Reset Password
</template>

<script>
import PasswordEntry from '@/base/PasswordEntry'
import PublicPage from '@/base/PublicPage'

export default {
  components: { PasswordEntry, PublicPage },
  data: () => ({ hints: null, invalid: false, password: '', submitted: false }),
  async created() {
    try {
      this.hints = (await this.$axios.get(`/api/password-reset/${this.$route.params.token}`)).data
    } catch (e) {
      this.invalid = true
    }
  },
  methods: {
    onPasswordChange(p) {
      this.password = p
    },
    async onSubmit() {
      this.submitted = true
      if (!this.password) return
      const body = new (FormData)
      body.append('password', this.password)
      const data = (await this.$axios.post(`/api/password-reset/${this.$route.params.token}`, body)).data
      this.$store.commit('login', data)
      this.$router.replace('/events')
    },
  },
}
</script>

<style lang="stylus" scoped>
#pwreset-token-top
  margin 0 auto
  padding 0 0.75rem
  max-width 0.75rem + 7rem + 20rem + 0.75rem
  // #pwreset-token-top.padding-left .form-label.width #pwreset-token-email.max-width #pwreset-token-top.padding-right
#pwreset-token-banner
  margin-top 1rem
  text-align center
  font-weight bold
  font-size 1.5rem
#pwreset-token-intro
  margin-top 0.5rem
  text-align center
  font-size 0.9rem
  line-height 1.2
#pwreset-token-form
  margin-top 1.5rem
#pwreset-token-submit-row
  margin 1.5rem 0 2rem
  text-align center
#pwreset-token-reset
  margin-top 1rem
  text-align center
</style>
