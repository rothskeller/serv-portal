<!--
PersonEdit displays the editor page for a person.
-->

<template lang="pug">
form#person-edit(@submit.prevent="onSubmit")
  b-form-group(label="Name" label-for="person-informalName" label-cols-sm="auto" label-class="person-edit-label" :state="informalNameError ? false : null" :invalid-feedback="informalNameError")
    b-input#person-informalName.person-edit-input(autofocus :plaintext="!canEditDetails" :state="informalNameError ? false : null" trim v-model="person.informalName")
    b-form-text(v-if="canEditDetails") What you like to be called, e.g. “Joe Banks”
  b-form-group(label="Formal name" label-for="person-formalName" label-cols-sm="auto" label-class="person-edit-label" :state="formalNameError ? false : null" :invalid-feedback="formalNameError")
    b-input#person-formalName.person-edit-input(:plaintext="!canEditDetails" :state="formalNameError ? false : null" v-model="person.formalName")
    b-form-text(v-if="canEditDetails") For formal documents, e.g. “Joseph A. Banks, Jr.”
  b-form-group(label="Sort name" label-for="person-sortName" label-cols-sm="auto" label-class="person-edit-label" :state="sortNameError ? false : null" :invalid-feedback="sortNameError")
    b-input#person-sortName.person-edit-input(:plaintext="!canEditDetails" :state="sortNameError ? false : null" v-model="person.sortName")
    b-form-text(v-if="canEditDetails") For appearance in sorted lists, e.g. “Banks, Joe”
  b-form-group(v-if="canEditUsername" label="Username" label-for="person-username" label-cols-sm="auto" label-class="person-edit-label" :state="usernameError ? false : null" :invalid-feedback="usernameError")
    b-input#person-username.person-edit-input(:state="usernameError ? false : null" v-model="person.username")
  b-form-group(label="Call sign" label-for="person-callSign" label-cols-sm="auto" label-class="person-edit-label" :state="callSignError ? false : null" :invalid-feedback="callSignError")
    b-input#person-callSign.person-edit-input(:plaintext="!canEditDetails" :state="callSignError ? false : null" v-model="person.callSign")
  PasswordEntry(v-if="canEditDetails" label="Password" labelClass="person-edit-label" :deferValidation="!submitted" :allowBadPassword="allowBadPassword" :passwordHints="myPasswordHints" @change="onPasswordChange")
  b-form-group(v-for="(e, i) in person.emails" :key="`e${i}`" :label="person.emails.length === 1 ? 'Email' : `Email #${i+1}`" :label-for="`person-email${i}`" label-cols-sm="auto" label-class="person-edit-label" :state="emailErrors[i] ? false : null" :invalid-feedback="emailErrors[i]")
    b-input.person-edit-input.d-inline(:id="`person-email${i}`" :plaintext="!canEditDetails" :state="emailErrors[i] ? false : null" trim v-model="e.email")
    b-input.person-edit-label-input(v-if="person.emails.length > 1" :id="`person-email${i}-label`" placeholder="Label" :plaintext="!canEditDetails" trim v-model="e.label")
    b-button.mt-3.d-block(v-if="i === person.emails.length-1" size="sm" @click="addEmail") Add another email
  b-form-group(v-for="(p, i) in person.phones" :key="`p${i}`" :label="person.phones.length === 1 ? 'Phone' : `Phone #${i+1}`" :label-for="`person-phone${i}`" label-cols-sm="auto" label-class="person-edit-label" :state="phoneErrors[i] ? false : null" :invalid-feedback="phoneErrors[i]")
    b-input.person-edit-input.d-inline(:id="`person-phone${i}`" :plaintext="!canEditDetails" :state="phoneErrors[i] ? false : null" trim v-model="p.phone")
    b-input.person-edit-label-input(v-if="person.phones.length > 1" :id="`person-phone${i}-label`" placeholder="Label" :plaintext="!canEditDetails" trim v-model="p.label")
    b-checkbox(v-model="p.sms") Receive text messages at this number
    b-button.mt-3(v-if="i === person.phones.length-1" size="sm" @click="addPhone") Add another phone number
  PersonEditAddress(type="Home" v-model="person.homeAddress")
  PersonEditAddress(type="Work" v-model="person.workAddress" :hasHome="!!person.homeAddress.address")
  PersonEditAddress(type="Mail" v-model="person.mailAddress" :hasHome="!!person.homeAddress.address")
  b-form-group.mt-3(:label="rolesLabel" :state="rolesError ? false : null" :invalid-feedback="rolesError")
    b-checkbox(v-if="canEditRoles" v-for="role in person.roles" :key="role.id" v-model="role.held" :disabled="!role.canAssign") {{role.name}}
    template(v-else v-for="role in person.roles")
      div(v-if="role.held" v-text="role.name")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="submitLabel")
    b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
export default {
  props: {
    person: Object,
    allowBadPassword: Boolean,
    canEditDetails: Boolean,
    canEditRoles: Boolean,
    canEditUsername: Boolean,
    passwordHints: Array,
  },
  data: () => ({
    informalNameError: null,
    formalNameError: null,
    sortNameError: null,
    duplicateSortName: null,
    usernameError: null,
    duplicateUsername: null,
    callSignError: null,
    duplicateCallSign: null,
    password: '',
    emailErrors: [],
    phoneErrors: [],
    rolesError: null,
    submitted: false,
    suggestions: null,
  }),
  computed: {
    me() { return this.$route.params.id == this.$store.state.me.id },
    newp() { return this.$route.params.id === 'NEW' },
    rolesLabel() {
      if (this.me) return 'You hold these roles:'
      if (this.newp) return 'This person will hold these roles:'
      return 'This person holds these roles:'
    },
    myPasswordHints() {
      const hints = [...this.passwordHints]
      if (this.person.informalName) hints.push(this.person.informalName)
      if (this.person.formalName) hints.push(this.person.formalName)
      if (this.person.callSign) hints.push(this.person.callSign)
      if (this.person.username) hints.push(this.person.username)
      if (this.person.homeAddress && this.person.homeAddress.address) hints.push(this.person.homeAddress.address)
      if (this.person.mailAddress && this.person.mailAddress.address) hints.push(this.person.mailAddress.address)
      if (this.person.workAddress && this.person.workAddress.address) hints.push(this.person.workAddress.address)
      this.person.emails.forEach(e => { hints.push(e.email) })
      this.person.phones.forEach(p => { hints.push(p.phone) })
      return hints
    },
    submitLabel() {
      if (this.me) return 'Save Changes'
      return this.newp ? 'Create Person' : 'Save Person'
    },
    valid() {
      return !this.informalNameError && !this.formalNameError && !this.sortNameError && !this.usernameError && !this.callSignError && !this.rolesError && this.password !== null && this.person.homeAddress && this.person.mailAddress && this.person.workAddress && !this.emailErrors.some(e => e) && !this.phoneErrors.some(p => p)
    },
  },
  mounted() {
    this.person.emails.forEach(e => {
      this.$watch((() => e.email), this.validate)
      this.emailErrors.push(null)
    })
    if (!this.person.emails.length) this.addEmails()
    this.person.phones.forEach(p => { this.phoneErrors.push(null) })
    if (!this.person.phones.length) this.addPhones()
    if (this.canEditRoles && this.newp)
      this.person.roles.forEach(r => {
        if (r.canAssign) this.$watch((() => r.held), this.validate)
      })
  },
  watch: {
    'person.informalName'(n, o) {
      if (this.person.formalName === o) this.person.formalName = n
      if (this.person.sortName === this.informalToSort(o)) this.person.sortName = this.informalToSort(n)
      this.validate()
    },
    'person.formalName': 'validate',
    'person.sortName': 'validate',
    'person.username': 'validate',
    'person.callSign': 'validate',
  },
  methods: {
    addEmail() {
      const e = { email: '', label: '' }
      this.person.emails.push(e)
      this.emailErrors.push(null)
      this.$watch((() => e.email), this.validate)
    },
    addPhone() {
      const p = { phone: '', label: '', sms: true }
      this.person.phones.push(p)
      this.phoneErrors.push(null)
      this.$watch((() => p.phone), this.validate)
    },
    onCancel() { this.$router.go(-1) },
    onPasswordChange() { this.password = p },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      // TODO fill in body
      body.append('informalName', this.person.informalName)
      body.append('formalName', this.person.formalName)
      body.append('sortName', this.person.sortName)
      body.append('username', this.person.username)
      body.append('callSign', this.person.callSign)
      if (this.password) body.append('password', this.password)
      this.person.emails.forEach(e => {
        if (e.email) {
          body.append('email', e.email)
          body.append('emailLabel', e.label || '')
        }
      })
      this.person.phones.forEach(p => {
        if (p.phone) {
          body.append('phone', p.phone)
          body.append('phoneLabel', p.label || '')
          body.append('phoneSMS', p.sms)
        }
      })
      if (this.person.homeAddress.address) {
        body.append('homeAddress', this.person.homeAddress.address)
        body.append('homeAddressLatitude', this.person.homeAddress.latitude)
        body.append('homeAddressLongitude', this.person.homeAddress.longitude)
      }
      if (this.person.workAddress.address) {
        body.append('workAddress', this.person.workAddress.address)
        body.append('workAddressLatitude', this.person.workAddress.latitude)
        body.append('workAddressLongitude', this.person.workAddress.longitude)
      } else {
        body.append('workAddressSameAsHome', this.person.workAddress.sameAsHome)
      }
      if (this.person.mailAddress.address) {
        body.append('mailAddress', this.person.mailAddress.address)
      } else {
        body.append('mailAddressSameAsHome', this.person.mailAddress.sameAsHome)
      }
      this.person.roles.filter(role => role.held && role.canAssign).forEach(role => { body.append('role', role.id) })
      const resp = (await this.$axios.post(`/api/people/${this.$route.params.id}`, body)).data
      if (resp) {
        if (resp.duplicateSortName) this.duplicateSortName = this.person.sortName
        if (resp.duplicateUsername) this.duplicateUsername = this.person.username
        if (resp.duplicateCallSign) this.duplicateCallSign = this.person.callSign
        // disregarding resp.weakPassword since we catch that locally
        this.validate()
      } else {
        this.$router.push('/people')
      }
    },
    validate() {
      if (!this.submitted) return
      if (!this.person.informalName)
        this.informalNameError = 'A name is required.'
      else
        this.informalNameError = null
      if (!this.person.formalName)
        this.formalNameError = 'A name is required.'
      else
        this.formalNameError = null
      if (!this.person.sortName)
        this.sortNameError = 'A name is required.'
      else if (this.duplicateSortName === this.person.sortName)
        this.sortNameError = 'A different person has this name.'
      else
        this.sortNameError = null
      if (this.duplicateUsername && this.person.username === this.duplicateUsername)
        this.usernameError = 'A different person has this username.'
      else
        this.usernameError = null
      if (this.person.callSign && !this.person.callSign.match(/^[AKNW][A-Z]?[0-9][A-Z]{1,3}$/))
        this.callSignError = 'This is not a valid call sign.'
      else if (this.duplicateCallSign === this.person.callSign)
        this.callSignError = 'A different person has this call sign.'
      else
        this.callSignError = null
      this.person.emails.forEach((e, i) => {
        if (!e.email)
          this.emailErrors[i] = null
        else if (!e.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
          this.emailErrors[i] = 'This is not a valid email address.'
        else
          this.emailErrors[i] = null
      })
      this.person.phones.forEach((p, i) => {
        if (p.phone && p.phone.replace(/[^0-9]/g, '').length !== 10)
          this.phoneErrors[i] = 'A valid phone number must have 10 digits.'
        else
          this.phoneErrors[i] = null
      })
      if (this.newp && !this.person.roles.some(role => role.held))
        this.rolesError = 'At least one role must be selected.'
      else
        this.rolesError = null
    },
  },
}
</script>

<style lang="stylus">
#person-edit
  margin 1.5rem 0.75rem
.person-edit-label
  width 8rem
.person-edit-input
  min-width 14rem
  max-width 20rem
.person-edit-label-input
  margin-top 0.25rem
  min-width 14rem
  max-width 20rem
  @media (min-width: 41.75rem)
    display inline
    margin-top 0
    margin-left 0.25rem
    min-width 6rem
    width 6rem
</style>
