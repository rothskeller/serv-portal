<!--
PersonEdit displays the editor page for a person.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!person")
  b-spinner(small)
form#person-edit(v-else @submit.prevent="onSubmit")
  template(v-if="canEditDetails")
    .person-edit-block-head Identification
    b-form-group(label="Name" label-for="person-informalName" label-cols-sm="auto" label-class="person-edit-label" :state="informalNameError ? false : null" :invalid-feedback="informalNameError")
      b-input#person-informalName.person-edit-input(autofocus :state="informalNameError ? false : null" trim v-model="person.informalName")
      b-form-text What you like to be called, e.g. “Joe Banks”
    b-form-group(label="Formal name" label-for="person-formalName" label-cols-sm="auto" label-class="person-edit-label" :state="formalNameError ? false : null" :invalid-feedback="formalNameError")
      b-input#person-formalName.person-edit-input(:state="formalNameError ? false : null" v-model="person.formalName")
      b-form-text For formal documents, e.g. “Joseph A. Banks, Jr.”
    b-form-group(label="Sort name" label-for="person-sortName" label-cols-sm="auto" label-class="person-edit-label" :state="sortNameError ? false : null" :invalid-feedback="sortNameError")
      b-input#person-sortName.person-edit-input(:state="sortNameError ? false : null" v-model="person.sortName")
      b-form-text For appearance in sorted lists, e.g. “Banks, Joe”
    b-form-group(v-if="canEditUsername" label="Username" label-for="person-username" label-cols-sm="auto" label-class="person-edit-label" :state="usernameError ? false : null" :invalid-feedback="usernameError")
      b-input#person-username.person-edit-input(:state="usernameError ? false : null" v-model="person.username")
    b-form-group(label="Call sign" label-for="person-callSign" label-cols-sm="auto" label-class="person-edit-label" :state="callSignError ? false : null" :invalid-feedback="callSignError")
      b-input#person-callSign.person-edit-input(:state="callSignError ? false : null" v-model="person.callSign")
    .person-edit-block-head Change Password
    b-form-group(v-if="!allowBadPassword" label="Old Password" label-for="person-oldPassword" label-cols-sm="auto" label-class="person-edit-label" :state="oldPasswordError ? false : null" :invalid-feedback="oldPasswordError")
      b-input#person-oldPassword.person-edit-input(type="password" :state="oldPasswordError ? false : null" v-model="oldPassword")
    PasswordEntry(label="New Password" labelClass="person-edit-label" :deferValidation="!submitted" :allowBadPassword="allowBadPassword" :passwordHints="myPasswordHints" @change="onPasswordChange")
    .person-edit-block-head Contact Information
    b-form-group(label="Email" label-for="person-email" label-cols-sm="auto" label-class="person-edit-label" :state="emailError ? false : null" :invalid-feedback="emailError")
      b-input#person-email.person-edit-input(:state="emailError ? false : null" trim v-model="person.email")
    b-form-group(label="Alt. Email" label-for="person-email2" label-cols-sm="auto" label-class="person-edit-label" :state="email2Error ? false : null" :invalid-feedback="email2Error")
      b-input#person-email2.person-edit-input(:state="email2Error ? false : null" trim v-model="person.email2")
    b-form-group(label="Cell Phone" label-for="person-cellPhone" label-cols-sm="auto" label-class="person-edit-label" :state="cellPhoneError ? false : null" :invalid-feedback="cellPhoneError")
      b-input#person-cellPhone.person-edit-input(:state="cellPhoneError ? false : null" trim v-model="person.cellPhone")
    b-form-group(label="Home Phone" label-for="person-homePhone" label-cols-sm="auto" label-class="person-edit-label" :state="homePhoneError ? false : null" :invalid-feedback="homePhoneError")
      b-input#person-homePhone.person-edit-input(:state="homePhoneError ? false : null" trim v-model="person.homePhone")
    b-form-group(label="Work Phone" label-for="person-workPhone" label-cols-sm="auto" label-class="person-edit-label" :state="workPhoneError ? false : null" :invalid-feedback="workPhoneError")
      b-input#person-workPhone.person-edit-input(:state="workPhoneError ? false : null" trim v-model="person.workPhone")
    PersonEditAddress(type="Home" v-model="person.homeAddress")
    PersonEditAddress(type="Work" v-model="person.workAddress" :hasHome="!!person.homeAddress.address")
    PersonEditAddress(type="Mail" v-model="person.mailAddress" :hasHome="!!person.homeAddress.address")
    .person-edit-block-head Roles
  b-form-group.mt-3(:label="rolesLabel" :state="rolesError ? false : null" :invalid-feedback="rolesError")
    b-checkbox(v-if="canEditRoles" v-for="role in person.roles" :key="role.id" v-model="role.held" :disabled="!role.canAssign") {{role.name}}
    template(v-else v-for="role in person.roles")
      div(v-if="role.held" v-text="role.name")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="submitLabel")
    b-btn.ml-2(@click="onCancel") Cancel
</template>

<script>
import PasswordEntry from '@/base/PasswordEntry'
import PersonEditAddress from './PersonEditAddress'

export default {
  components: { PasswordEntry, PersonEditAddress },
  props: {
    onLoadPerson: Function,
  },
  data: () => ({
    person: null,
    allowBadPassword: false,
    canEditDetails: false,
    canEditRoles: false,
    canEditUsername: false,
    passwordHints: [],
    informalNameError: null,
    formalNameError: null,
    sortNameError: null,
    duplicateSortName: null,
    usernameError: null,
    duplicateUsername: null,
    callSignError: null,
    duplicateCallSign: null,
    emailError: null,
    email2Error: null,
    cellPhoneError: null,
    duplicateCellPhone: null,
    homePhoneError: null,
    workPhoneError: null,
    oldPassword: null,
    oldPasswordError: null,
    wrongOldPassword: null,
    password: '',
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
      if (this.person.email) hints.push(this.person.email)
      if (this.person.email2) hints.push(this.person.email2)
      if (this.person.homeAddress && this.person.homeAddress.address) hints.push(this.person.homeAddress.address)
      if (this.person.mailAddress && this.person.mailAddress.address) hints.push(this.person.mailAddress.address)
      if (this.person.workAddress && this.person.workAddress.address) hints.push(this.person.workAddress.address)
      if (this.person.cellPhone) hints.push(this.person.cellPhone)
      if (this.person.homePhone) hints.push(this.person.homePhone)
      if (this.person.workPhone) hints.push(this.person.workPhone)
      return hints
    },
    submitLabel() {
      if (this.me) return 'Save Changes'
      return this.newp ? 'Create Person' : 'Save Person'
    },
    valid() {
      return !this.informalNameError && !this.formalNameError && !this.sortNameError && !this.usernameError && !this.callSignError && !this.emailError && !this.email2Error && !this.cellPhoneError && !this.homePhoneError && !this.workPhoneError && !this.rolesError && !this.oldPasswordError && this.password !== null && this.person.homeAddress && this.person.mailAddress && this.person.workAddress
    },
  },
  async created() {
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}?edit=1`)).data
    this.allowBadPassword = data.allowBadPassword
    this.canEditDetails = data.canEditDetails
    this.canEditRoles = data.canEditRoles
    this.canEditUsername = data.canEditUsername
    this.passwordHints = data.passwordHints
    this.person = data.person
    this.onLoadPerson(this.person)
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
    oldPassword: 'validate',
    'person.email': 'validate',
    'person.email2': 'validate',
    'person.cellPhone': 'validate',
    'person.homePhone': 'validate',
    'person.workPhone': 'validate',
  },
  methods: {
    informalToSort(n) {
      if (!n) return n
      const parts = n.split(/\s+/, 2)
      return parts.length > 1 ? `${parts[1]}, ${parts[0]}` : n
    },
    onCancel() { this.$router.go(-1) },
    onPasswordChange(p) {
      this.password = p
      this.validate()
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('informalName', this.person.informalName)
      body.append('formalName', this.person.formalName)
      body.append('sortName', this.person.sortName)
      body.append('username', this.person.username)
      body.append('callSign', this.person.callSign)
      body.append('email', this.person.email || this.person.email2)
      body.append('email2', this.person.email ? this.person.email2 : '')
      body.append('cellPhone', this.person.cellPhone)
      body.append('homePhone', this.person.homePhone)
      body.append('workPhone', this.person.workPhone)
      if (this.oldPassword) body.append('oldPassword', this.oldPassword)
      if (this.password) body.append('password', this.password)
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
        if (resp.duplicateCellPhone) this.duplicateCellPhone = this.person.cellPhone
        if (resp.wrongOldPassword) this.wrongOldPassword = this.oldPassword
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
      if (this.password && !this.oldPassword && !this.allowBadPassword)
        this.oldPasswordError = 'You must supply your old password in order to change your password.'
      else if (this.oldPassword && this.oldPassword === this.wrongOldPassword)
        this.oldPasswordError = 'This is not the correct old password.'
      else
        this.oldPasswordError = null
      if (!this.person.email)
        this.emailError = null
      else if (!this.person.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        this.emailError = 'This is not a valid email address.'
      else
        this.emailError = null
      if (!this.person.email2)
        this.email2Error = null
      else if (!this.person.email2.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        this.email2Error = 'This is not a valid email address.'
      else if (this.person.email2 == this.person.email)
        this.email2Error = 'The two email addresses should not be the same.  (Leave this field empty if you only have one.)'
      else
        this.email2Error = null
      if (this.person.cellPhone && this.person.cellPhone.replace(/[^0-9]/g, '').length !== 10)
        this.cellPhoneError = 'A valid phone number must have 10 digits.'
      else if (this.duplicateCellPhone === this.person.cellPhone)
        this.cellPhoneError = 'A different person has this cell phone number.'
      else
        this.cellPhoneError = null
      if (this.person.homePhone && this.person.homePhone.replace(/[^0-9]/g, '').length !== 10)
        this.homePhoneError = 'A valid phone number must have 10 digits.'
      else
        this.homePhoneError = null
      if (this.person.workPhone && this.person.workPhone.replace(/[^0-9]/g, '').length !== 10)
        this.workPhoneError = 'A valid phone number must have 10 digits.'
      else
        this.workPhoneError = null
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
.person-edit-block-head
  overflow hidden
  margin-top 2rem
  margin-bottom 1rem
  padding-left 2rem
  max-width 28rem
  color #888
  &:first-child
    margin-top 0
  &::before
    display inline-block
    margin-right 0.5rem
    margin-left -100%
    width 100%
    border-top 1px solid #888
    content ''
    vertical-align middle
  &::after
    display inline-block
    margin-right -100%
    margin-left 0.5rem
    width 100%
    border-top 1px solid #888
    content ''
    vertical-align middle
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
