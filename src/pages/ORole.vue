<!--
Role displays and edits the details for a role.
-->

<template lang="pug">
Page(:title="title" :subtitle="subtitle" menuItem="teams")
  div.mt-3(v-if="loading")
    b-spinner(small)
  form(v-else @submit.prevent="onSubmit")
    b-form-group(label="Role name" label-for="role-name" label-cols-sm="auto" label-class="role-label" :state="nameError ? false : null" :invalid-feedback="nameError")
      b-input#role-name(autofocus :state="nameError ? false : null" trim v-model="role.name")
    b-form-group(label="Team" label-for="role-team" label-cols-sm="auto" label-class="role-label")
      b-input#role-team(plaintext :value="team.name")
    b-form-group(label="This role has the following privileges on teams:")
      table#role-privs
        tr
          th
          th Member
          th Access
        tr(v-for="team in privs" :key="team.id")
          td(:class="`indent-${team.indent}`" v-text="team.name")
          td: b-checkbox(v-model="team.member")
          td: b-radio-group(:options="accessList" v-model="team.access")
    div.mt-3
      b-btn(type="submit" variant="primary" :disabled="!!nameError" v-text="role.id ? 'Save Role' : 'Create Role'")
      b-btn.ml-2(@click="onCancel") Cancel
      b-btn.ml-5(v-if="canDelete" variant="danger" @click="onDelete") Delete Role
</template>

<script>
const accessList = [
  { value: 'none', text: 'None' },
  { value: 'view', text: 'View' },
  { value: 'admin', text: 'Admin' },
  { value: 'manage', text: 'Manage' },
]

export default {
  data: () => ({
    accessList,
    loading: false,
    origName: null,
    team: null,
    role: null,
    canDelete: false,
    privs: null,
    nameError: null,
    duplicateName: null,
    submitted: false,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/teams/${this.$route.params.tid}/roles/${this.$route.params.rid}`)).data
    this.team = data.team
    this.role = data.role
    this.origName = this.role.name
    this.canDelete = data.canDelete
    this.privs = data.privs
    this.loading = false
  },
  computed: {
    subtitle() { return this.$route.params.rid === 'NEW' ? 'Create Role' : 'Edit Role' },
    title() {
      if (this.$route.params.rid === 'NEW') return 'New Role'
      return this.origName ? `${this.team.name}: ${this.origName || '(member)'}` : 'Edit Role'
    },
  },
  watch: {
    'role.name': 'validate',
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onDelete() {
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/teams/${this.$route.params.tid}/roles/${this.$route.params.rid}`, body)
      this.$router.push('/teams')
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (this.nameError) return
      const body = new FormData
      body.append('name', this.role.name)
      this.privs.forEach(t => {
        body.append(`member-${t.id}`, t.member)
        body.append(`access-${t.id}`, t.access)
      })
      const resp = (await this.$axios.post(`/api/teams/${this.$route.params.tid}/roles/${this.$route.params.rid}`, body)).data
      if (!resp)
        this.$router.push('/teams')
      if (resp.duplicateName)
        this.duplicateName = this.role.name
      this.validate()
    },
    validate() {
      if (!this.submitted) return
      if (this.duplicateName === this.role.name)
        this.nameError = 'A different role on this team has this name.'
      else
        this.nameError = null
    },
  },
}
</script>

<style lang="stylus">
.role-label
  width 7rem
#role-name
  max-width 20rem
#role-privs
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
