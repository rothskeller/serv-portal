<!--
PWResetToken displays the second password reset page (the one the email links to).
-->

<template lang="pug">
#pwreset-token-top
  #pwreset-token-banner Password Reset
  template(v-if='invalid')
    #pwreset-token-intro
      | This password reset link is invalid or has expired.
    #pwreset-token-reset
      SButton(to='/password-reset', variant='primary') Try Again
  #pwreset-token-intro(v-else-if='!hints')
    SSpinner
  template(v-else)
    #pwreset-token-intro
      | Please provide a new password. Enter it twice.
    SForm#pwreset-token-form(@submit='onSubmit')
      SFPassword#pwreset-password(
        :allowBadPassword='false',
        :passwordHints='hints',
        required,
        v-model='password'
      )
      template(#buttons)
        #pwreset-token-submit-row
          SButton(type='submit', variant='primary') Reset Password
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../plugins/axios'
import { passwordReset } from '../plugins/login'
import setPage from '../plugins/page'
import { SFPassword, SForm, SButton, SSpinner } from '../base'

export default defineComponent({
  components: { SForm, SButton, SFPassword, SSpinner },
  setup() {
    setPage({ title: '', browserTitle: 'Password Reset' })

    const route = useRoute()
    const router = useRouter()

    const invalid = ref(false)
    const hints = ref(null as null | Array<string>)
    axios
      .get(`/api/password-reset/${route.params.token}`)
      .then((result) => {
        hints.value = result.data
      })
      .catch(() => {
        invalid.value = true
      })

    const password = ref('')
    async function onSubmit() {
      const token = Array.isArray(route.params.token) ? route.params.token[0] : route.params.token
      if (await passwordReset(token, password.value)) router.replace('/')
    }

    return { password, hints, invalid, onSubmit }
  },
})
</script>

<style lang="postcss">
#pwreset-token-top {
  margin: 0 auto;
  padding: 0 0.75rem;
  max-width: calc(0.75rem + 10rem + 20rem + 0.75rem + 1.5rem);
}
#pwreset-token-banner {
  margin-top: 1rem;
  text-align: center;
  font-weight: bold;
  font-size: 1.5rem;
}
#pwreset-token-intro {
  margin-top: 0.5rem;
  text-align: center;
  font-size: 0.9rem;
  line-height: 1.2;
}
#pwreset-token-form {
  margin-top: 1.5rem;
  justify-content: center;
}
#pwreset-token-submit-row {
  margin: 1.5rem 0 2rem;
  text-align: center;
}
#pwreset-token-reset {
  margin-top: 1rem;
  text-align: center;
}
</style>
