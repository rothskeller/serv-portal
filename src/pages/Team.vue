<!--
Team displays and edits the details for a team.
-->

<template lang="pug">
Page(:title="title" :subtitle="subtitle" menuItem="teams")
  div.mt-3(v-if="loading")
    b-spinner(small)
  form(v-else @submit.prevent="onSubmit")
    b-form-group(label="Team name" label-for="team-name" label-cols-sm="auto" label-class="team-label" :state="nameError ? false : null" :invalid-feedback="nameError")
      b-input#team-name(autofocus :state="nameError ? false : null" trim v-model="team.name")
    b-form-group(label="Parent team" label-for="team-parent" label-cols-sm="auto" label-class="team-label")
      b-input#team-parent(plaintext :value="team.parent || '(none)'")
    b-form-group(label="Team email" label-for="team-email" label-cols-sm="auto" label-class="team-label" :state="emailError ? false : null" :invalid-feedback="emailError")
      b-input#team-email(:state="emailError ? false : null" trim v-model="team.email")
    b-form-group(label="This team has the following privileges on other teams:")
      table#team-aprivs
        tr
          th
          th Member
          th Access
        tr(v-for="other in team.privs" :key="other.id")
          td(:class="`indent-${other.indent}`" v-text="other.id ? other.name : '(new team)'")
          td: b-checkbox(v-model="other.actor.member")
          td: b-radio-group(:options="accessList" v-model="other.actor.access")
    b-form-group(label="These other teams have the following privileges on this team:")
      table#team-tprivs
        tr
          th
          th Member
          th Access
        tr(v-for="other in team.privs" :key="other.id")
          template(v-if="other.id !== team.id")
            td(:class="`indent-${other.indent}`" v-text="other.id ? other.name : '(new team)'")
            td: b-checkbox(v-model="other.target.member")
            td: b-radio-group(:options="accessList" v-model="other.target.access")
    div.mt-3
      b-btn(type="submit" variant="primary" :disabled="!valid" v-text="team.id ? 'Save Team' : 'Create Team'")
      b-btn.ml-2(@click="onCancel") Cancel
      b-btn.ml-5(v-if="canDelete" variant="danger" @click="onDelete") Delete Team
</template>

<script>
const accessList = [
  { value: 'none', text: 'None' },
  { value: 'view', text: 'View' },
  { value: 'admin', text: 'Admin' },
  { value: 'manage', text: 'Manage' },
]

export default {
  props: {
    parent: String,
  },
  data: () => ({
    accessList,
    loading: false,
    origName: null,
    team: null,
    canDelete: false,
    nameError: null,
    duplicateName: null,
    emailError: null,
    duplicateEmail: null,
    submitted: false,
    valid: true,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/teams/${this.$route.params.id}?parent=${this.parent || ''}`)).data
    this.team = { id: data.id, name: data.name, email: data.email, parent: data.parent, privs: data.privs }
    this.origName = data.name
    this.canDelete = data.canDelete
    this.loading = false
  },
  computed: {
    subtitle() { return this.$route.params.id === 'NEW' ? 'Create Team' : 'Edit Team' },
    title() {
      if (this.$route.params.id === 'NEW') return 'New Team'
      return this.origName ? `Team: ${this.origName}` : 'Edit Team'
    },
  },
  watch: {
    'team.name': 'validate',
    'team.email': 'validate',
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onDelete() {
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/teams/${this.$route.params.id}`, body)
      this.$router.push('/teams')
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('parent', this.parent)
      body.append('name', this.team.name)
      body.append('email', this.team.email)
      this.team.privs.forEach(t => {
        body.append(`a:member-${t.id}`, t.actor.member)
        body.append(`a:access-${t.id}`, t.actor.access)
        body.append(`t:member-${t.id}`, t.target.member)
        body.append(`t:access-${t.id}`, t.target.access)
      })
      const resp = (await this.$axios.post(`/api/teams/${this.$route.params.id}`, body)).data
      if (!resp)
        this.$router.push('/teams')
      if (resp.duplicateName)
        this.duplicateName = this.team.name
      if (resp.duplicateEmail)
        this.duplicateEmail = this.team.email
      this.validate()
    },
    validate() {
      if (!this.submitted) return
      if (!this.team.name)
        this.nameError = 'The team name is required.'
      else if (this.duplicateName === this.team.name)
        this.nameError = 'A different team has this name.'
      else
        this.nameError = null
      if (this.duplicateEmail && this.duplicateEmail === this.team.email)
        this.emailError = 'A different team has this email.'
      else
        this.emailError = null
      this.valid = !this.nameError && !this.emailError
    },
  },
}
</script>

<style lang="stylus">
.team-label
  width 7rem
#team-name, #team-email
  max-width 20rem
#team-aprivs, #team-tprivs
  margin-top 0.5rem
  th
    padding-right 1em
    font-weight normal
  td
    padding-right 1em
    vertical-align middle
    text-align center
    &:first-child
      text-align left
    &.indent-1
      padding-left 1em
    &.indent-2
      padding-left 2em
    &.indent-3
      padding-left 3em
    &.indent-4
      padding-left 4em
    &.indent-5
      padding-left 5em
</style>
