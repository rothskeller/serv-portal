<!--
Login displays the login page.
-->

<template lang="pug">
#login-top
  SForm#login-form(ref='formRef', @submit='onSubmit')
    #login-banner.form-item Please log in.
    #login-forserv.form-item
      | This web site is for SERV volunteers only.
      | If you are interested in joining one of the SERV volunteer organizations,
      | please visit Sunnyvale’s <a href="https://sunnyvale.ca.gov/government/safety/emergency.htm">emergency response&nbsp;page</a>.
    #login-browserwarn.form-item
      | Your browser is out of date and lacks features needed by this web site.
      | The site may not look or behave correctly.
    #login-spacer.form-item
    SFInput#login-email(
      label='Email address',
      autocorrect='off',
      autocapitalize='none',
      autofocus,
      trim,
      v-model='email',
      :errorFn='emailError'
    )
    SFInput#login-password(
      ref='passwordRef',
      label='Password',
      type='password',
      v-model='password',
      :errorFn='passwordError'
    )
    SFCheck#login-remember(label='Remember me', v-model='remember')
    template(#buttons)
      #login-submit-row
        SButton(type='submit', variant='primary') Log in
    template(v-if='failed', #feedback)
      #login-failed Login incorrect. Please try again.
  #login-reset
    router-link(to='/password-reset') Reset my password
  #login-policies
    router-link(to='/policies') Site Policies / Legal Stuff
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { login } from '../plugins/login'
import setPage from '../plugins/page'
import { SForm, SFInput, SFCheck, SButton } from '../base'

export default defineComponent({
  components: { SForm, SFInput, SFCheck, SButton },
  setup() {
    setPage({ title: 'Sunnyvale SERV', browserTitle: 'Login' })

    // References to template elements.
    const formRef = ref(null as any)
    const passwordRef = ref(null as any)

    // Email input field.
    const email = ref('')
    function emailError(lostFocus: boolean) {
      if (email.value || !lostFocus) return ''
      return 'To log in, you must enter your email address.'
    }

    // Password input field.
    const password = ref('')
    function passwordError(lostFocus: boolean) {
      if (password.value || !lostFocus) return ''
      return 'To log in, you must enter your password.  If you do not know your password, click the “Reset my password” link below.'
    }

    // Remember me input field.
    const remember = ref(false)

    // Form submission.
    const failed = ref(false)
    const router = useRouter()
    const route = useRoute()
    async function onSubmit() {
      if (!email.value || !password.value) return
      if (await login(email.value, password.value, remember.value)) {
        const redir = route.query.redirect
        if (typeof redir === 'string') router.replace(redir)
        else router.replace('/home')
      } else {
        formRef.value.resetSubmitted()
        password.value = ''
        failed.value = true
        passwordRef.value.focus()
      }
    }

    return {
      formRef,
      email,
      emailError,
      passwordRef,
      password,
      passwordError,
      remember,
      onSubmit,
      failed,
    }
  },
})
</script>

<style lang="postcss">
#login-top {
  display: flex;
  flex-direction: column;
  margin: 0 auto;
  padding: 0 0.75rem;
  min-height: calc(100vh - 3rem - 40px);
  max-width: calc(0.75rem + 10rem + 20rem + 0.75rem + 1.5rem);
}
#login-form {
  justify-content: center;
}
#login-banner {
  text-align: center;
  font-weight: bold;
  font-size: 1.5rem;
}
#login-forserv {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  line-height: 1.2;
}
#login-browserwarn {
  margin-top: 0.5rem;
  padding: 2px 4px;
  background-color: red;
  color: white;
  line-height: 1.2;
  @supports (display: grid) {
    display: none;
  }
}
#login-spacer,
#login-form {
  margin-top: 1.5rem;
}
#login-submit-row {
  margin-bottom: 1rem;
  text-align: center;
}
#login-failed {
  color: red;
  text-align: center;
}
#login-reset {
  text-align: center;
}
#login-policies {
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  justify-content: flex-end;
  align-self: center;
  margin-top: 3rem;
  text-decoration: none;
}
</style>
