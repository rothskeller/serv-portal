<!--
PersonEditPassword is the dialog box for changing a person's password.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Change Password',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
      SFInput#person-oldPassword(
        v-if='!me.webmaster',
        ref='oldPasswordRef',
        type='password',
        label='Old Password',
        v-model='oldPassword',
        :errorFn='oldPasswordError'
      )
      SFPassword#person-password(
        ref='passwordRef',
        label='New Password',
        v-model='password',
        :allowBadPassword='me.webmaster',
        :passwordHints='hints'
      )
</template>

<script lang="ts">
import { defineComponent, inject, nextTick, Ref, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { LoginData } from '../../plugins/login'
import { Modal, SForm, SFInput, SFPassword, SSpinner } from '../../base'

export default defineComponent({
  components: { Modal, SForm, SFInput, SFPassword, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const me = inject<Ref<LoginData>>('me')!
    const modal = ref(null as any)
    function show() {
      oldPassword.value = password.value = ''
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const hints = ref([] as Array<string>)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      hints.value = (await axios.get<Array<string>>(`/api/people/${props.pid}/password`)).data
      loading.value = false
      nextTick(() => {
        if (me.value.webmaster) passwordRef.value.focus()
        else oldPasswordRef.value.focus()
      })
    }

    // Field validation.
    const oldPassword = ref('')
    const oldPasswordRef = ref(null as any)
    const password = ref('')
    const passwordRef = ref(null as any)
    const wrongOldPassword = ref('')
    function oldPasswordError(lostFocus: boolean, submitted: boolean) {
      if (!submitted) return ''
      if (password.value && !oldPassword.value && !me.value.webmaster)
        return 'You must supply your old password in order to change your password.'
      if (oldPassword.value && oldPassword.value === wrongOldPassword.value)
        return 'This is not the correct old password.'
      return ''
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      body.append('oldPassword', oldPassword.value)
      body.append('password', password.value)
      submitting.value = true
      try {
        await axios.post(`/api/people/${props.pid}/password`, body)
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        wrongOldPassword.value = oldPassword.value
      } finally {
        submitting.value = false
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      hints,
      loading,
      me,
      modal,
      oldPassword,
      oldPasswordError,
      oldPasswordRef,
      onCancel,
      onSubmit,
      password,
      passwordRef,
      show,
      submitting,
    }
  },
})
</script>

<style lang="postcss">
</style>
