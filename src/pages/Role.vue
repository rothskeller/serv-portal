<!--
Role displays and edits the details for a role.
-->

<template lang="pug">
Page(:title="title" :subtitle="subtitle" menuItem="roles")
  div.mt-3(v-if="loading")
    b-spinner(small)
  form(v-else @submit.prevent="onSubmit")
    b-form-group(label="Role name" label-for="role-name" label-cols-sm="auto" label-class="role-label" :state="nameError ? false : null" :invalid-feedback="nameError")
      b-input#role-name(autofocus :state="nameError ? false : null" trim v-model="role.name")
    b-form-group(label="Holder label" label-for="role-holder-label" label-cols-sm="auto" label-class="role-label")
      b-input#role-holder-label(trim v-model="role.memberLabel")
    b-form-group(label="Attributes" label-cols-sm="auto" label-class="role-label pt-0")
      div
        b-checkbox(v-model="role.implyOnly") Role can only be implied, not assigned
        b-checkbox(v-model="role.individual") Role can only be held by one person
    b-form-group(label="This role has the following privileges on other roles:")
      table#role-aprivs
        tr
          th
          th Holds
          th View
          th Assign
          th Events
        tr(v-for="other in privs" :key="other.id")
          td(v-text="other.id === role.id && newr ? '(new role)' : other.name")
          td: b-checkbox(:checked="other.actor.holdsRole" :disabled="!other.actor.holdsRoleEnabled" @change="onChangePriv(other, 'actor', 'holdsRole')")
          td: b-checkbox(:checked="other.actor.viewHolders" :disabled="!other.actor.viewHoldersEnabled" @change="onChangePriv(other, 'actor', 'viewHolders')")
          td: b-checkbox(:checked="other.actor.assignRole" :disabled="!other.actor.assignRoleEnabled" @change="onChangePriv(other, 'actor', 'assignRole')")
          td: b-checkbox(:checked="other.actor.manageEvents" :disabled="!other.actor.manageEventsEnabled" @change="onChangePriv(other, 'actor', 'manageEvents')")
    b-form-group(label="These other roles have the following privileges on this role:")
      table#role-tprivs
        tr
          th
          th Holds
          th View
          th Assign
          th Events
        template(v-for="other in privs")
          tr(v-if="other.id !== role.id" :key="other.id")
            td(v-text="other.name")
            td: b-checkbox(:checked="other.target.holdsRole" :disabled="!other.target.holdsRoleEnabled" @change="onChangePriv(other, 'target', 'holdsRole')")
            td: b-checkbox(:checked="other.target.viewHolders" :disabled="!other.target.viewHoldersEnabled" @change="onChangePriv(other, 'target', 'viewHolders')")
            td: b-checkbox(:checked="other.target.assignRole" :disabled="!other.target.assignRoleEnabled" @change="onChangePriv(other, 'target', 'assignRole')")
            td: b-checkbox(:checked="other.target.manageEvents" :disabled="!other.target.manageEventsEnabled" @change="onChangePriv(other, 'target', 'manageEvents')")
    div.mt-3
      b-btn(type="submit" variant="primary" :disabled="!valid" v-text="newr ? 'Create Role' : 'Save Role'")
      b-btn.ml-2(@click="onCancel") Cancel
      b-btn.ml-5(v-if="role.id" variant="danger" @click="onDelete") Delete Role
</template>

<script>
export default {
  data: () => ({
    loading: false,
    origName: null,
    role: null,
    privs: null,
    nameError: null,
    duplicateName: null,
    submitted: false,
    valid: true,
  }),
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/roles/${this.$route.params.id}`)).data
    this.role = data.role
    this.privs = data.privs
    this.origName = this.role.name
    this.loading = false
  },
  computed: {
    newr() { return this.$route.params.id === 'NEW' },
    subtitle() { return this.newr ? 'Create Role' : 'Edit Role' },
    title() {
      if (this.newr) return 'New Role'
      return this.origName ? `Role: ${this.origName}` : 'Edit Role'
    },
  },
  watch: {
    'role.name'(n, o) {
      if (this.role.memberLabel === o) this.role.memberLabel = n
      this.validate()
    },
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onChangePriv(role, dir, priv) {
      role[dir][priv] = !role[dir][priv]
      const body = new FormData
      this.privs.forEach(r => {
        body.append(`a:holdsRole-${r.id}`, r.actor.holdsRole)
        body.append(`a:viewHolders-${r.id}`, r.actor.viewHolders)
        body.append(`a:assignRole-${r.id}`, r.actor.assignRole)
        body.append(`a:manageEvents-${r.id}`, r.actor.manageEvents)
        body.append(`t:holdsRole-${r.id}`, r.target.holdsRole)
        body.append(`t:viewHolders-${r.id}`, r.target.viewHolders)
        body.append(`t:assignRole-${r.id}`, r.target.assignRole)
        body.append(`t:manageEvents-${r.id}`, r.target.manageEvents)
      })
      this.privs = (await this.$axios.post(`/api/roles/${this.$route.params.id}/reloadPrivs`, body)).data
    },
    async onDelete() {
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/roles/${this.$route.params.id}`, body)
      this.$router.push('/roles')
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('name', this.role.name)
      body.append('memberLabel', this.role.memberLabel)
      body.append('implyOnly', this.role.implyOnly)
      body.append('individual', this.role.individual)
      this.privs.forEach(r => {
        body.append(`a:holdsRole-${r.id}`, r.actor.holdsRole)
        body.append(`a:viewHolders-${r.id}`, r.actor.viewHolders)
        body.append(`a:assignRole-${r.id}`, r.actor.assignRole)
        body.append(`a:manageEvents-${r.id}`, r.actor.manageEvents)
        body.append(`t:holdsRole-${r.id}`, r.target.holdsRole)
        body.append(`t:viewHolders-${r.id}`, r.target.viewHolders)
        body.append(`t:assignRole-${r.id}`, r.target.assignRole)
        body.append(`t:manageEvents-${r.id}`, r.target.manageEvents)
      })
      const resp = (await this.$axios.post(`/api/roles/${this.$route.params.id}`, body)).data
      if (!resp)
        this.$router.push('/roles')
      if (resp.duplicateName)
        this.duplicateName = this.role.name
      this.validate()
    },
    validate() {
      if (!this.submitted) return
      if (!this.role.name)
        this.nameError = 'The role name is required.'
      else if (this.duplicateName === this.role.name)
        this.nameError = 'A different team has this name.'
      else
        this.nameError = null
      this.valid = !this.nameError
    },
  },
}
</script>

<style lang="stylus">
.role-label
  width 7rem
#role-name, #role-holder-label
  max-width 20rem
#role-aprivs, #role-tprivs
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
