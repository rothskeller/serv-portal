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
    PasswordEntry(v-if="canEditInfo" label="Change password" labelClass="person-label" :deferValidation="!submitted" :allowBadPassword="allowBadPassword" :passwordHints="myPasswordHints" @change="onPasswordChange")
    b-form-group(:label="rolesLabel" :state="rolesError ? false : null" :invalid-feedback="rolesError")
      div
        b-checkbox(v-for="role in roles" :key="role.id" v-model="role.held" :disabled="!role.enabled") {{role.memberLabel || role.name}}
    div.mt-3(v-if="canEditInfo")
      b-btn(type="submit" variant="primary" :disabled="!valid" v-text="submitLabel")
      b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
import zxcvbn from 'zxcvbn'

export default {
  data: () => ({
    loading: false,
    person: null,
    roles: null,
    password: '',
    canEditInfo: false,
    allowBadPassword: false,
    firstNameError: null,
    lastNameError: null,
    duplicateName: null,
    emailError: null,
    duplicateEmail: null,
    phoneError: null,
    passwordHints: null,
    rolesError: null,
    submitted: false,
  }),
  computed: {
    me() { return this.$route.params.id == this.$store.state.me.id },
    newp() { return this.$route.params.id === 'NEW' },
    myPasswordHints() {
      const hints = [...this.passwordHints]
      if (this.person.firstName) hints.push(this.person.firstName)
      if (this.person.lastName) hints.push(this.person.lastName)
      if (this.person.email) hints.push(this.person.email)
      if (this.person.phone) hints.push(this.person.phone)
      return hints
    },
    rolesLabel() {
      if (this.me) return 'You hold these roles:'
      if (this.new) return 'This person will hold these roles:'
      return 'This person holds these roles:'
    },
    submitLabel() {
      if (this.me) return 'Save Changes'
      return this.newp ? 'Create Person' : 'Save Person'
    },
    subtitle() {
      if (this.me) return 'Edit Profile'
      return this.newp ? 'New Person' : 'Edit Person'
    },
    title() {
      if (this.person && this.person.id) return `${this.person.firstName} ${this.person.lastName}`
      return this.newp ? 'New Person' : 'Edit Person'
    },
    valid() {
      return !this.firstNameError && !this.lastNameError && !this.emailError && !this.phoneError && !this.rolesError && this.password !== null
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}`)).data
    this.canEditInfo = data.canEditInfo
    this.allowBadPassword = data.allowBadPassword
    this.person = data.person
    this.passwordHints = data.passwordHints
    this.roles = data.roles
    this.loading = false
  },
  watch: {
    'person.firstName': 'validate',
    'person.lastName': 'validate',
    'person.email': 'validate',
    'person.phone': 'validate',
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    onPasswordChange(p) { this.password = p },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('firstName', this.person.firstName)
      body.append('lastName', this.person.lastName)
      body.append('email', this.person.email)
      body.append('phone', this.person.phone)
      if (this.password) body.append('password', this.password)
      this.roles.filter(role => role.held && role.enabled).forEach(role => { body.append('role', role.id) })
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
        this.emailError = null
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
      if (this.newp && !this.roles.some(role => role.held))
        this.rolesError = 'At least one role must be selected.'
      else
        this.rolesError = null
    },
  },
}
</script>

<style lang="stylus">
.person-label
  width 9em
#person-firstName, #person-lastName, #person-email, #person-phone
  max-width 20em
</style>
