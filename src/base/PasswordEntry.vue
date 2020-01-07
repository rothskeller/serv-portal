<!--
PasswordEntry displays a form group for password entry, with validation.
-->

<template lang="pug">
b-form-group#pwentry-password(:label="label" label-cols-sm="auto" :label-class="labelClass" :state="passwordError ? false : passwordSuccess ? true : null" :invalid-feedback="passwordError" :valid-feedback="passwordSuccess")
  div
    b-input#pwentry-password1(type="password" :state="passwordState" v-model="password1")
    b-input.mt-2#pwentry-password2(type="password" :state="passwordState" v-model="password2")
    b-progress#pwentry-password-score(v-if="passwordScore>=0" height="0.5rem")
      b-progress-bar(:value="100" :variant="passwordVariant")
      b-progress-bar(:value="passwordScore > 0 ? 100 : 0" :variant="passwordVariant")
      b-progress-bar(:value="passwordScore > 1 ? 100 : 0" :variant="passwordVariant")
      b-progress-bar(:value="passwordScore > 2 ? 100 : 0" :variant="passwordVariant")
      b-progress-bar(:value="passwordScore > 3 ? 100 : 0" :variant="passwordVariant")
</template>

<script>
import zxcvbn from 'zxcvbn'

export default {
  props: {
    label: String,
    labelClass: String,
    deferValidation: Boolean,
    allowBadPassword: Boolean,
    passwordHints: Array,
  },
  data: () => ({
    password1: null,
    password2: null,
    passwordHints: null,
    passwordError: null,
    passwordSuccess: null,
    passwordScore: -1,
  }),
  computed: {
    passwordState() {
      return this.password1 && this.passwordScore < 3 && !this.allowBadPassword ? false : null
    },
    passwordVariant() {
      if (this.passwordScore >= 3) return 'success'
      if (this.passwordScore === 2) return 'warning'
      return 'danger'
    },
  },
  watch: {
    password1: 'validate',
    password2: 'validate',
  },
  methods: {
    validate() {
      if ((!this.deferValidation || this.password2) && this.password1 !== this.password2) {
        this.passwordScore = -1
        this.passwordError = 'These two password entries do not match.'
        this.passwordSuccess = null
      } else if (this.password1) {
        const result = zxcvbn(this.password1, this.passwordHints)
        this.passwordScore = result.score
        if (result.feedback) {
          this.passwordError = [
            result.feedback.warning,
            ...result.feedback.suggestions,
            this.passwordScore < 3
              ? `This password would take ${result.crack_times_display.offline_slow_hashing_1e4_per_second} to crack.`
              : null,
          ].filter(s => !!s).join('\n')
          this.passwordSuccess = this.passwordScore > 2
            ? `This password would take ${result.crack_times_display.offline_slow_hashing_1e4_per_second} to crack.`
            : null
        } else {
          this.passwordError = this.passwordSuccess = null
        }
      } else {
        this.passwordScore = -1
        this.passwordError = this.passwordSuccess = null
      }
      if (!this.password1 && !this.password2) this.$emit('change', '')
      else if (this.password1 !== this.password2) this.$emit('change', null)
      else if (this.allowBadPassword || this.passwordScore >= 3) this.$emit('change', this.password1)
      else this.emit('change', null)
    },
  },
}
</script>

<style lang="stylus">
#pwentry-password1, #pwentry-password2
  max-width 20em
#pwentry-password .invalid-feedback
  white-space pre-line
#pwentry-password-score
  margin-top 0.5rem
  width 11rem
  .progress-bar
    margin-left 0.25rem
    max-width 2rem
    &:first-child
      margin-left 0
</style>
