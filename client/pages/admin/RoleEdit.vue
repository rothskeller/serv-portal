<!--
RoleEdit displays the role viewing/editing page.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!role")
  b-spinner(small)
form#role-edit(v-else @submit.prevent="onSubmit")
  b-form-group(label="Role name" label-for="role-name" label-cols-sm="auto" label-class="role-edit-label" :state="nameError ? false : null" :invalid-feedback="nameError")
    b-input#role-edit-name(autofocus :state="nameError ? false : null" trim v-model="role.name")
  b-form-group(label="Flags" label-cols-sm="auto" label-class="role-edit-label pt-0")
    #role-edit-flags
      b-checkbox(v-model="role.individual") Individual (one person only)
  #role-edit-privs
    .role-edit-group.role-edit-heading Group
    .role-edit-privs.role-edit-heading
      | Privileges (
      b-link(href="#" v-b-modal.role-edit-priv-key) Key
      | )
      b-modal#role-edit-priv-key(title="Privileges Key" ok-only)
        #role-edit-priv-key-body
          div.mb-2 M = Role conveys membership in group
          div R = Role allows viewing group roster
          div C = Role allows viewing contact info of group members
          div A = Role allows administration of group (adding/removing members)
          div E = Role allows management of events for group
          div T = Role allows sending text messages to group
          div @ = Role allows sending email messages to group
          div B = Role gets bcc'd on email messages to group
    template(v-for="(g, i) in privs")
      .role-edit-group(v-text="g.name")
      Privileges.role-edit-privs(v-model="privs[i]")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="newr ? 'Create Role' : 'Save Role'")
    b-btn.ml-2(@click="onCancel") Cancel
    b-btn.ml-5(v-if="!newr" @click="onClone") Clone Role
    b-btn.ml-2(v-if="canDelete" variant="danger" @click="onDelete") Delete Role
</template>

<script>
import Privileges from '@/base/Privileges'

export default {
  components: { Privileges },
  props: {
    onLoadRole: Function,
  },
  data: () => ({
    role: null,
    privs: null,
    canDelete: false,
    submitted: false,
    nameError: null,
    duplicateName: null,
  }),
  computed: {
    newr() { return this.$route.params.rid === 'NEW' },
    valid() { return !this.nameError },
  },
  watch: {
    $route: 'onChangeRoute',
    'role.name': 'validate',
  },
  created() {
    this.onChangeRoute()
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onChangeRoute() {
      const rid = this.$route.params.clone || this.$route.params.rid
      const data = (await this.$axios.get(`/api/roles/${rid}`)).data
      this.role = data.role
      this.privs = data.privs
      this.canDelete = data.canDelete && this.$route.params.rid !== 'NEW'
      this.onLoadRole(this.role)
    },
    onClone() {
      this.$router.push({ name: 'roles-rid', params: { rid: 'NEW', clone: this.$route.params.rid } })
    },
    async onDelete() {
      const resp = await this.$bvModal.msgBoxConfirm(
        'Are you sure you want to delete this role?  All associated data, including privileges and memberships, will be permanently lost.', {
        title: 'Delete Role', headerBgVariant: 'danger', headerTextVariant: 'white',
        okTitle: 'Delete', okVariant: 'danger', cancelTitle: 'Keep',
      }).catch(err => { })
      if (!resp) return
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/roles/${this.$route.params.rid}`, body)
      this.$router.push('/roles')
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('name', this.role.name)
      body.append('individual', this.role.individual)
      this.privs.forEach(g => {
        if (g.member) body.append(`member:${g.id}`, true)
        if (g.roster) body.append(`roster:${g.id}`, true)
        if (g.contact) body.append(`contact:${g.id}`, true)
        if (g.admin) body.append(`admin:${g.id}`, true)
        if (g.events) body.append(`events:${g.id}`, true)
        if (g.texts) body.append(`texts:${g.id}`, true)
        if (g.emails) body.append(`emails:${g.id}`, true)
        if (g.bcc) body.append(`bcc:${g.id}`, true)
      })
      const resp = (await this.$axios.post(`/api/roles/${this.$route.params.rid}`, body)).data
      if (resp && resp.duplicateName)
        this.duplicateName = this.role.name
      else
        this.$router.push('/roles')
    },
    validate() {
      if (!this.submitted) return
      if (!this.role.name)
        this.nameError = 'The role name is required.'
      else if (this.duplicateName && this.duplicateName === this.role.name)
        this.nameError = 'Another role has this name.'
      else
        this.nameError = null
    },
  },
}
</script>

<style lang="stylus">
#role-edit
  padding 1.5rem 0.75rem
.role-edit-label
  width 7rem
#role-edit-name, #role-edit-flags
  min-width 14rem
  max-width 20rem
#role-edit-privs
  display grid
  grid auto / 1fr
  @media (min-width: 450px)
    justify-content start
    grid auto / auto min-content
.role-edit-heading
  display none
  @media (min-width: 450px)
    display block
    font-weight bold
.role-edit-group
  overflow hidden
  margin-top 0.75rem
  min-width 0
  text-overflow ellipsis
  white-space nowrap
  @media (min-width: 450px)
    align-self center
    margin-top 0
    margin-right 0.75rem
#role-edit-priv-key-body
  line-height 1.2
</style>
