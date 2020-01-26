<!--
Login displays the login page.
-->

<template lang="pug">
PublicPage(title="SERV Portal")
  #login-top
    #login-banner Please log in.
    #login-forserv
      | This web site is for SERV volunteers only.
      | If you are interested in joining one of the SERV volunteer organizations,
      | please visit Sunnyvale’s <a href="https://sunnyvale.ca.gov/government/safety/emergency.htm">emergency response&nbsp;page</a>.
    form#login-form(@submit.prevent="onSubmit")
      b-form-group(label="Email address" label-for="login-email" label-cols-sm="4")
        b-input#login-email(autocorrect="off" autocapitalize="none" autofocus required trim v-model="email")
      b-form-group(label="Password" label-for="login-password" label-cols-sm="4")
        b-input#login-password(ref="password" type="password" v-model="password")
      #login-submit-row
        b-button(type="submit" variant="primary") Log in
      #login-failed(v-if="failed") Login incorrect. Please try again.
    #login-reset
      b-btn(to="/password-reset") Reset my password
</template>

<script>
import PublicPage from '@/base/PublicPage'

export default {
  components: { PublicPage },
  data: () => ({ email: '', password: '', failed: false }),
  methods: {
    async onSubmit() {
      if (!this.email || !this.password) return
      const body = new (FormData)
      body.append('username', this.email)
      body.append('password', this.password)
      try {
        const data = (await this.$axios.post('/api/login', body)).data
        this.$store.commit('login', data)
        this.$router.replace(this.$route.query.redirect || '/events')
      } catch (err) {
        console.error(err)
        this.failed = true
        this.password = ''
        this.$refs.password.focus()
      }
    },
  },
}
</script>

<style lang="stylus" scoped>
#login-top
  margin 0 auto
  padding 0 0.75rem
  max-width 0.75rem + 7rem + 20rem + 0.75rem
  // #login-top.padding-left .form-label.width #login-email.max-width #login-top.padding-right
#login-banner
  margin-top 1rem
  text-align center
  font-weight bold
  font-size 1.5rem
#login-forserv
  margin-top 0.5rem
  font-size 0.9rem
  line-height 1.2
#login-form
  margin-top 1.5rem
#login-submit-row
  margin 1rem 0 2rem
  text-align center
#login-failed
  color red
  text-align center
#login-reset
  margin-top 1rem
  text-align center
</style>
