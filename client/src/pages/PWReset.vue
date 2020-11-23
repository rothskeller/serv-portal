<!--
PWReset displays the password reset page.
-->

<template lang="pug">
#pwreset-top
  #pwreset-banner Password Reset
  template(v-if='!finished')
    #pwreset-intro
      | To reset your password, please enter your email address. If it's one we
      | have on file, we'll send a password reset link to it.
    SForm#pwreset-form(@submit='onSubmit')
      SFInput#pwreset-email(
        label='Email address',
        autocorrect='off',
        autocapitalize='none',
        autofocus,
        trim,
        v-model='email'
      )
      template(#buttons)
        #pwreset-submit-row
          SButton(type='submit', variant='primary') Reset Password
  #pwreset-intro(v-else)
    | We have sent a password reset link to the email address you provided.
    | It is valid for one hour. Please check you email and follow the link we
    | sent to reset your password.
    br
    br
    | If you do not receive an email with a password reset link, it may be that
    | the email address you provided is not the one we have on file for you.
    | Contact
    |
    a(href='mailto:admin@sunnyvaleserv.org') admin@SunnyvaleSERV.org
    |
    | for assistance.
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import axios from '../plugins/axios'
import setPage from '../plugins/page'
import { SForm, SFInput, SButton } from '../base'

export default defineComponent({
  components: { SForm, SFInput, SButton },
  setup() {
    setPage({ title: '', browserTitle: 'Password Reset' })

    const email = ref('')
    const finished = ref(false)

    async function onSubmit() {
      if (!email.value) return
      const body = new FormData()
      body.append('email', email.value)
      const data = (await axios.post('/api/password-reset', body)).data
      finished.value = true
    }

    return { email, finished, onSubmit }
  },
})
</script>

<style lang="postcss">
#pwreset-top {
  margin: 0 auto;
  padding: 0 0.75rem;
  max-width: calc(0.75rem + 10rem + 20rem + 0.75rem + 1.5rem);
}
#pwreset-banner {
  margin-top: 1rem;
  text-align: center;
  font-weight: bold;
  font-size: 1.5rem;
}
#pwreset-intro {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  line-height: 1.2;
}
#pwreset-form {
  margin-top: 1.5rem;
  justify-content: center;
}
#pwreset-submit-row {
  margin-bottom: 2rem;
  text-align: center;
}
</style>
