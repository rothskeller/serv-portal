<!--
GroupEdit displays the group viewing/editing page.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!group")
  b-spinner(small)
form#group-edit(v-else @submit.prevent="onSubmit")
  b-form-group(label="Group name" label-for="group-name" label-cols-sm="auto" label-class="group-edit-label" :state="nameError ? false : null" :invalid-feedback="nameError")
    b-input#group-edit-name(autofocus :state="nameError ? false : null" trim v-model="group.name")
  b-form-group(label="Group email" label-for="group-email" label-cols-sm="auto" label-class="group-edit-label" :state="emailError ? false : null" :invalid-feedback="emailError")
    b-input#group-edit-email(:state="emailError ? false : null" trim v-model="group.email")
    b-form-text @sunnyvaleserv.org
  #group-edit-privs
    .group-edit-role.group-edit-heading Role
    .group-edit-privs.group-edit-heading
      | Privileges (
      b-link(href="#" v-b-modal.group-edit-priv-key) Key
      | )
      b-modal#group-edit-priv-key(title="Privileges Key" ok-only)
        #group-edit-priv-key-body
          div.mb-2 M = Role conveys membership in group
          div R = Role allows viewing group roster
          div C = Role allows viewing contact info of group members
          div A = Role allows administration of group (adding/removing members)
          div E = Role allows management of events for group
          div F = Role allows management of files for group
          div T = Role allows sending text messages to group
          div @ = Role allows sending email messages to group
          div B = Role gets bcc'd on email messages to group
    template(v-for="(r, i) in privs")
      .group-edit-role(v-text="r.name")
      Privileges.group-edit-privs(v-model="privs[i]")
  div.mt-3(v-if="group.noEmail && group.noEmail.length")
    div Unsubscribed from emails to this group:
    div(v-for="p in group.noEmail" v-text="p")
  div.mt-3(v-if="group.noText && group.noText.length")
    div Unsubscribed from text messages to this group:
    div(v-for="p in group.noText" v-text="p")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!valid" v-text="newg ? 'Create Group' : 'Save Group'")
    b-btn.ml-2(@click="onCancel") Cancel
    b-btn.ml-5(v-if="!newg" @click="onClone") Clone Group
    b-btn.ml-2(v-if="canDelete" variant="danger" @click="onDelete") Delete Group
</template>

<script>
import Privileges from '@/base/Privileges'

export default {
  components: { Privileges },
  props: {
    onLoadGroup: Function,
  },
  data: () => ({
    group: null,
    privs: null,
    canDelete: false,
    submitted: false,
    nameError: null,
    duplicateName: null,
    emailError: null,
    duplicateEmail: null,
  }),
  computed: {
    newg() { return this.$route.params.gid === 'NEW' },
    valid() { return !this.nameError },
  },
  watch: {
    $route: 'onChangeRoute',
    'group.name': 'validate',
    'group.email': 'validate',
  },
  created() {
    this.onChangeRoute()
  },
  methods: {
    onCancel() { this.$router.go(-1) },
    async onChangeRoute() {
      const gid = this.$route.params.clone || this.$route.params.gid
      const data = (await this.$axios.get(`/api/groups/${gid}`)).data
      this.group = data.group
      this.privs = data.privs
      this.canDelete = data.canDelete && this.$route.params.gid !== 'NEW'
      this.onLoadGroup(this.group)
    },
    onClone() {
      this.$router.push({ name: 'groups-gid', params: { gid: 'NEW', clone: this.$route.params.gid } })
    },
    async onDelete() {
      const resp = await this.$bvModal.msgBoxConfirm(
        'Are you sure you want to delete this group?  All associated data, including privileges and memberships, will be permanently lost.', {
        title: 'Delete Group', headerBgVariant: 'danger', headerTextVariant: 'white',
        okTitle: 'Delete', okVariant: 'danger', cancelTitle: 'Keep',
      }).catch(err => { })
      if (!resp) return
      const body = new FormData
      body.append('delete', 'true')
      await this.$axios.post(`/api/groups/${this.$route.params.gid}`, body)
      this.$router.push('/admin/groups')
    },
    async onSubmit() {
      this.submitted = true
      this.validate()
      if (!this.valid) return
      const body = new FormData
      body.append('name', this.group.name)
      body.append('email', this.group.email)
      this.privs.forEach(r => {
        if (r.member) body.append(`member:${r.id}`, true)
        if (r.roster) body.append(`roster:${r.id}`, true)
        if (r.contact) body.append(`contact:${r.id}`, true)
        if (r.admin) body.append(`admin:${r.id}`, true)
        if (r.events) body.append(`events:${r.id}`, true)
        if (r.texts) body.append(`texts:${r.id}`, true)
        if (r.emails) body.append(`emails:${r.id}`, true)
        if (r.bcc) body.append(`bcc:${r.id}`, true)
        if (r.folders) body.append(`folders:${r.id}`, true)
      })
      const resp = (await this.$axios.post(`/api/groups/${this.$route.params.gid}`, body)).data
      if (resp) {
        if (resp.duplicateName) this.duplicateName = this.group.name
        if (resp.duplicateEmail) this.duplicateEmail = this.group.email
      } else {
        this.$router.push('/admin/groups')
      }
    },
    validate() {
      if (!this.submitted) return
      if (!this.group.name)
        this.nameError = 'The group name is required.'
      else if (this.duplicateName && this.duplicateName === this.group.name)
        this.nameError = 'Another group has this name.'
      else
        this.nameError = null
      if (!this.group.email)
        this.emailError = null
      else if (!this.group.email.match(/^[a-zA-Z][-a-zA-Z0-9]*$/))
        this.emailError = 'This is not a valid email list name.  Letters, digits, and hyphens only.'
      else if (this.duplicateEmail && this.duplicateEmail === this.group.email)
        this.emailError = 'Another group has this email list name.'
      else
        this.emailError = null
    },
  },
}
</script>

<style lang="stylus">
#group-edit
  padding 1.5rem 0.75rem
.group-edit-label
  width 7rem
#group-edit-name, #group-edit-email, #group-edit-flags
  min-width 14rem
  max-width 20rem
#group-edit-privs
  display grid
  grid auto / 1fr
  @media (min-width: 450px)
    justify-content start
    grid auto / auto min-content
.group-edit-heading
  display none
  @media (min-width: 450px)
    display block
    font-weight bold
.group-edit-group
  overflow hidden
  margin-top 0.75rem
  min-width 0
  text-overflow ellipsis
  white-space nowrap
  @media (min-width: 450px)
    align-self center
    margin-top 0
    margin-right 0.75rem
#group-edit-priv-key-body
  line-height 1.2
</style>
