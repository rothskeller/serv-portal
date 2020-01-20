<!--
PWReset displays the password reset page.
-->

<template lang="pug">
Page(title="SERV Portal")
  #pwreset-top
    #pwreset-banner Password Reset
    template(v-if="!finished")
      #pwreset-intro
        | To reset your password, please enter your email address.  If it's one we
        | have on file, we'll send a password reset link to it.
      form#pwreset-form(@submit.prevent="onSubmit")
        b-form-group(label="Email address" label-for="pwreset-username" label-cols-sm="4")
          b-input#pwreset-username(autocorrect="off" autocapitalize="none" autofocus required trim v-model="username")
        #pwreset-submit-row
          b-button(type="submit" variant="primary") Reset Password
    #pwreset-intro(v-else)
      | We have sent a password reset link to the email address you provided.
      | It is valid for one hour.  Please check you email and follow the link we
      | sent to reset your password.
</template>

<script>
export default {
  data: () => ({ username: '', finished: false }),
  watch: {
    email: 'validate',
  },
  methods: {
    async onSubmit() {
      if (!this.username) return
      const body = new (FormData)
      body.append('username', this.username)
      const data = (await this.$axios.post('/api/password-reset', body)).data
      this.finished = true
    },
  },
}
</script>

<style lang="stylus" scoped>
#pwreset-top
  margin 0 auto
  padding 0 0.75rem
  max-width 0.75rem + 7rem + 20rem + 0.75rem
  // #pwreset-top.padding-left .form-label.width #pwreset-email.max-width #pwreset-top.padding-right
#pwreset-banner
  margin-top 1rem
  text-align center
  font-weight bold
  font-size 1.5rem
#pwreset-intro
  margin-top 0.5rem
  font-size 0.9rem
  line-height 1.2
#pwreset-form
  margin-top 1.5rem
#pwreset-submit-row
  margin 1.5rem 0 2rem
  text-align center
</style>
