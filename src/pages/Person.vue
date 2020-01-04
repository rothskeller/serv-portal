<!--
Person displays and allows editing of the information about a person.
-->

<template lang="pug">
Page(:title="title" :subtitle="subtitle" menuItem="people")
  div.mt-3(v-if="loading")
    b-spinner(small)
  form.mt-3(v-else @submit.prevent="onSubmit")
    b-form-group(label="First name" label-for="person-firstName" label-cols-sm="auto" label-class="person-label" :state="firstNameError ? false : null" :invalid-feedback="firstNameError")
      b-input#person-firstName(autofocus :plaintext="!canEditInfo" :state="firstNameError ? false : null" trim v-model="person.firstName")
    b-form-group(label="Last name" label-for="person-lastName" label-cols-sm="auto" label-class="person-label" :state="lastNameError ? false : null" :invalid-feedback="lastNameError")
      b-input#person-lastName(:plaintext="!canEditInfo" :state="lastNameError ? false : null" trim v-model="person.lastName")
    b-form-group(label="Email address" label-for="person-email" label-cols-sm="auto" label-class="person-label" :state="emailError ? false : null" :invalid-feedback="emailError")
      b-input#person-email(:plaintext="!canEditInfo" :state="emailError ? false : null" trim v-model="person.email")
    b-form-group(label="Phone number" label-for="person-phone" label-cols-sm="auto" label-class="person-label" :state="phoneError ? false : null" :invalid-feedback="phoneError")
      b-input#person-phone(:plaintext="!canEditInfo" :state="phoneError ? false : null" trim v-model="person.phone")
    b-form-group#person-password(v-if="canEditInfo" label="Change password" label-cols-sm="auto" label-class="person-label" :state="passwordError ? false : passwordSuccess ? true : null" :invalid-feedback="passwordError" :valid-feedback="passwordSuccess")
      div
        b-input#person-password1(type="password" :state="passwordState" v-model="password1")
        b-input.mt-2#person-password2(type="password" :state="passwordState" v-model="password2")
        b-progress#person-password-score(v-if="passwordScore" height="0.5rem")
          b-progress-bar(:value="100" :variant="passwordVariant")
          b-progress-bar(:value="passwordScore > 0 ? 100 : 0" :variant="passwordVariant")
          b-progress-bar(:value="passwordScore > 1 ? 100 : 0" :variant="passwordVariant")
          b-progress-bar(:value="passwordScore > 2 ? 100 : 0" :variant="passwordVariant")
          b-progress-bar(:value="passwordScore > 3 ? 100 : 0" :variant="passwordVariant")
    div.mt-3(v-text="`${me ? 'You belong' : 'This person belongs'} to these teams:`")
    table#person-team-table
      PersonTeam(v-for="t in teams" :key="t.id" :manageAny="manageAny" :team="t" @change="onChangeTeam")
    div.mt-3(v-if="adminAny")
      b-btn(type="submit" variant="primary" :disabled="!valid" v-text="submitLabel")
      b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
import zxcvbn from 'zxcvbn'

export default {
  data: () => ({
    loading: false,
    person: null,
    password1: null,
    password2: null,
    teams: null,
    canEditInfo: false,
    allowBadPassword: false,
    manageAny: false,
    adminAny: false,
    firstNameError: null,
    lastNameError: null,
    duplicateName: null,
    emailError: null,
    duplicateEmail: null,
    phoneError: null,
    passwordHints: null,
    passwordError: null,
    passwordSuccess: null,
    passwordScore: 0,
    submitted: false,
  }),
  computed: {
    me() { return this.$route.params.id == this.$store.state.me.id },
    passwordState() {
      return this.password1 && this.passwordScore < 2 && !this.allowBadPassword ? false : null
    },
    passwordVariant() {
      if (this.passwordScore >= 3) return 'success'
      if (this.passwordScore === 2) return 'warning'
      return 'danger'
    },
    submitLabel() {
      if (this.me) return 'Save Changes'
      return this.$route.params.id === 'NEW' ? 'Create Person' : 'Save Person'
    },
    subtitle() {
      if (this.me) return 'Edit Profile'
      return this.$route.params.id === 'NEW' ? 'New Person' : 'Edit Person'
    },
    title() {
      if (this.person && this.person.id) return `${this.person.firstName} ${this.person.lastName}`
      return this.$route.params.id === 'NEW' ? 'New Person' : 'Edit Person'
    },
    valid() {
      return !this.firstNameError && !this.lastNameError && !this.emailError && !this.phoneError &&
        ((!this.password1 && !this.password2) || (this.passwordScore >= (this.allowBadPassword ? 1 : 2)))
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}`)).data
    this.canEditInfo = data.canEditInfo
    this.allowBadPassword = data.allowBadPassword
    this.person = data.person
    this.teams = data.teams
    this.teams.forEach(t => {
      if (t.canManage) this.manageAny = true
      if (t.canAdmin) this.adminAny = true
    })
    this.passwordHints = data.passwordHints
    this.loading = false
  },
  watch: {
    'person.firstName': 'validate',
    'person.lastName': 'validate',
    'person.email': 'validate',
    'person.phone': 'validate',
    password1: 'validate',
    password2: 'validate',
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    onChangeTeam({ team, role }) { team.role = role ? role.id : 0 },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('firstName', this.person.firstName)
      body.append('lastName', this.person.lastName)
      body.append('email', this.person.email)
      body.append('phone', this.person.phone)
      if (this.password1) body.append('password', this.password1)
      this.teams.forEach(t => {
        if (t.role) body.append('role', t.role)
      })
      const resp = (await this.$axios.post(`/api/people/${this.$route.params.id}`, body)).data
      if (resp) {
        if (resp.duplicateName) this.duplicateName = { firstName: this.person.firstName, lastName: this.person.lastName }
        if (resp.duplicateEmail) this.duplicateEmail = this.person.email
        // disregarding resp.weakPassword since we catch that locally
        this.validate()
      } else {
        this.$router.push('/people')
      }
    },
    validate() {
      if ((this.submitted || this.password2) && this.password1 !== this.password2) {
        this.passwordScore = 0
        this.passwordError = 'These two password entries do not match.'
        this.passwordSuccess = null
      } else if (this.password1) {
        const hints = [...this.passwordHints]
        if (this.person.firstName) hints.push(this.person.firstName)
        if (this.person.lastName) hints.push(this.person.lastName)
        if (this.person.email) hints.push(this.person.email)
        if (this.person.phone) hints.push(this.person.phone)
        const result = zxcvbn(this.password1, hints)
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
        this.passwordScore = 0
        this.passwordError = this.passwordSuccess = null
      }
      if (!this.submitted) return
      if (!this.person.firstName)
        this.firstNameError = 'A first name is required.'
      else
        this.firstNameError = null
      if (!this.person.lastName)
        this.lastNameError = 'A last name is required.'
      else if (this.duplicateName && this.person.firstName == this.duplicateName.firstName && this.person.lastName == this.duplicateName.lastName)
        this.lastNameError = 'A different person has this name.'
      else
        this.lastNameError = null
      if (!this.person.email)
        this.emailError = 'An email address is required.'
      else if (!this.person.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        this.emailError = 'This is not a valid email address.'
      else if (this.person.email === this.duplicateEmail)
        this.emailError = 'This email address is in use by another person.'
      else
        this.emailError = null
      if (this.person.phone && this.person.phone.replace(/[^0-9]/g, '').length !== 10)
        this.phoneError = 'A valid phone number must have 10 digits.'
      else
        this.phoneError = null
    },
  },
}
</script>

<style lang="stylus">
.person-label
  width 9em
#person-firstName, #person-lastName, #person-email, #person-phone, #person-password1, #person-password2
  max-width 20em
#person-password .invalid-feedback
  white-space pre
#person-password-score
  margin-top 0.5rem
  width 11rem
  .progress-bar
    margin-left 0.25rem
    max-width 2rem
    &:first-child
      margin-left 0
</style>
