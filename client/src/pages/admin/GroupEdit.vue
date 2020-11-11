<!--
GroupEdit displays the group viewing/editing page.
-->

<template lang="pug">
#group-edit(v-if='!group')
  SSpinner
SForm#group-edit(v-else, @submit='onSubmit', :submitLabel='submitLabel')
  SFInput#group-edit-name(
    label='Group name',
    trim,
    autofocus,
    v-model='group.name',
    :errorFn='nameError'
  )
  SFInput#group-edit-email(
    label='Group email',
    trim,
    v-model='group.email',
    :errorFn='emailError',
    help='@sunnyvaleserv.org'
  )
  SFCheckGroup#group-edit-flags(label='Flags', :options='allFlags', v-model='flags')
  SFSelect#group-edit-organization(
    label='Organization',
    :options='organizations',
    v-model='group.organization',
    help='For volunteer hours tracking'
  )
  #group-edit-privs.form-item
    .group-edit-role.group-edit-heading Role
    .group-edit-privs.group-edit-heading
      | Privileges (
      a(href='#', @click='onShowKey') Key
      | )
    template(v-for='(r, i) in privs')
      .group-edit-role(v-text='r.name')
      PrivilegeMask.group-edit-privs(v-model='privs[i]')
  .group-edit-unsubscribed(v-if='group.noEmail && group.noEmail.length')
    div Unsubscribed from emails to this group:
    div(v-for='p in group.noEmail', v-text='p')
  .group-edit-unsubscribed(v-if='group.noText && group.noText.length')
    div Unsubscribed from text messages to this group:
    div(v-for='p in group.noText', v-text='p')
  template(v-if='canDelete || canClone', #extraButtons)
    SButton(v-if='canClone', @click='onClone') Clone Group
    SButton(v-if='canDelete', variant='danger', @click='onDelete') Delete Group
  Modal(ref='keyModal', v-slot='{ close }')
    SForm(dialog, title='Privileges Key')
      #group-edit-priv-key-body.form-item
        div M = Role conveys membership in group
        div R = Role allows viewing group roster
        div C = Role allows viewing contact info of group members
        div A = Role allows admin of group (adding/removing members)
        div E = Role allows management of events for group
        div F = Role allows management of files for group
        div T = Role allows sending text messages to group
        div @ = Role allows sending email messages to group
        div B = Role gets bcc'd on email messages to group
      template(#buttons)
        SButton#group-edit-priv-key-ok(variant='primary', @click='close(null)') OK
  MessageBox(
    ref='deleteModal',
    title='Delete Group',
    cancelLabel='Keep',
    okLabel='Delete',
    variant='danger'
  )
    | Are you sure you want to delete this group? All associated data,
    | including privileges and memberships, will be permanently lost.
</template>

<script lang='ts'>
import { computed, defineComponent, ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import {
  MessageBox,
  Modal,
  Privileges,
  PrivilegeMask,
  SButton,
  SForm,
  SFCheckGroup,
  SFInput,
  SFSelect,
  SSpinner,
} from '../../base'

interface GetGroupEditPrivilege extends Privileges {
  id: number
  name: string
}
type GetGroupEditGroup = {
  id: number
  name: string
  email: string
  getHours: boolean
  dswRequired: boolean
  backgroundCheckRequired: boolean
  organization: string
  noEmail?: Array<string>
  noText?: Array<string>
}
type GetGroupEdit = {
  group: GetGroupEditGroup
  canDelete: boolean
  privs: Array<GetGroupEditPrivilege>
  organizations: Array<string>
}
type PostGroup = {
  duplicateName?: boolean
  duplicateEmail?: boolean
}

export default defineComponent({
  components: {
    MessageBox,
    Modal,
    PrivilegeMask,
    SButton,
    SForm,
    SFCheckGroup,
    SFInput,
    SFSelect,
    SSpinner,
  },
  setup() {
    // Handle page load and change.
    const route = useRoute()
    const router = useRouter()
    const group = ref(null as null | GetGroupEditGroup)
    const privs = ref([] as Array<GetGroupEditPrivilege>)
    const canDelete = ref(false)
    const organizations = ref([] as Array<string>)
    watchEffect(() => {
      setPage({ title: route.params.gid === 'NEW' ? 'Create Group' : 'Edit Group' })
      const gid = route.params.clone || route.params.gid
      axios.get<GetGroupEdit>(`/api/groups/${gid}`).then((resp) => {
        group.value = resp.data.group
        privs.value = resp.data.privs
        canDelete.value = resp.data.canDelete
        organizations.value = resp.data.organizations
        organizations.value.unshift('(none)')
        flags.value.clear()
        allFlags.forEach((f) => {
          // @ts-ignore
          if (resp.data.group[f.value]) flags.value.add(f.value)
        })
      })
    })

    // Name field.
    let duplicateName = ref('')
    function nameError(lostFocus: boolean) {
      if (!lostFocus || !group.value) return ''
      if (!group.value.name) return 'The group name is required.'
      if (duplicateName.value === group.value.name) return 'Another group has this name.'
      return ''
    }

    // Email field.
    let duplicateEmail = ref('')
    function emailError(lostFocus: boolean) {
      if (!lostFocus || !group.value || !group.value.email) return ''
      if (!group.value.email.match(/^[a-zA-Z][-a-zA-Z0-9]*$/))
        return 'This is not a valid email list name.  Letters, digits, and hyphens only.'
      if (duplicateEmail.value === group.value.email)
        return 'Another group has this email list name.'
      return ''
    }

    // Flags field.
    const flags = ref(new Set() as Set<string>)
    const allFlags = [
      { value: 'getHours', label: 'Request volunteer hours from this group' },
      { value: 'dswRequired', label: 'DSW registration required for this group' },
      { value: 'backgroundCheckRequired', label: 'Background check required for this group' },
    ]

    // Clone.
    const canClone = computed(() => route.params.gid !== 'NEW')
    function onClone() {
      router.push({ name: 'groups-gid', params: { gid: 'NEW', clone: route.params.gid } })
    }

    // Delete.
    const deleteModal = ref(null as any)
    async function onDelete() {
      if (!(await deleteModal.value.show())) return
      const body = new FormData()
      body.append('delete', 'true')
      await axios.post(`/api/groups/${route.params.gid}`, body)
      router.push('/admin/groups')
    }

    // Submit.
    const submitLabel = computed(() => (route.params.gid === 'NEW' ? 'Create Group' : 'Save Group'))
    async function onSubmit() {
      if (!group.value) return
      const body = new FormData()
      body.append('name', group.value.name)
      body.append('email', group.value.email)
      allFlags.forEach((f) => {
        body.append(f.value, flags.value.has(f.value).toString())
      })
      body.append('organization', group.value.organization)
      privs.value.forEach((r) => {
        if (r.member) body.append(`member:${r.id}`, 'true')
        if (r.roster) body.append(`roster:${r.id}`, 'true')
        if (r.contact) body.append(`contact:${r.id}`, 'true')
        if (r.admin) body.append(`admin:${r.id}`, 'true')
        if (r.events) body.append(`events:${r.id}`, 'true')
        if (r.texts) body.append(`texts:${r.id}`, 'true')
        if (r.emails) body.append(`emails:${r.id}`, 'true')
        if (r.bcc) body.append(`bcc:${r.id}`, 'true')
        if (r.folders) body.append(`folders:${r.id}`, 'true')
      })
      const resp = (await axios.post<PostGroup>(`/api/groups/${route.params.gid}`, body)).data
      if (resp) {
        if (resp.duplicateName) duplicateName.value = group.value.name
        if (resp.duplicateEmail) duplicateEmail.value = group.value.email
      } else {
        router.push('/admin/groups')
      }
    }

    // Privileges Key.
    const keyModal = ref(null as any)
    function onShowKey() {
      keyModal.value.show()
    }

    return {
      allFlags,
      canClone,
      canDelete,
      deleteModal,
      emailError,
      flags,
      group,
      keyModal,
      nameError,
      onClone,
      onDelete,
      onShowKey,
      onSubmit,
      organizations,
      privs,
      submitLabel,
    }
  },
})
</script>

<style lang="postcss">
#group-edit {
  padding: 1.5rem 0.75rem;
}
#group-edit-privs {
  display: grid;
  grid: auto / 1fr;
  @media (min-width: 450px) {
    justify-content: start;
    grid: auto / auto min-content;
  }
}
.group-edit-heading {
  display: none;
  @media (min-width: 450px) {
    display: block;
    font-weight: bold;
  }
}
.group-edit-group {
  overflow: hidden;
  margin-top: 0.75rem;
  min-width: 0;
  text-overflow: ellipsis;
  white-space: nowrap;
  @media (min-width: 450px) {
    align-self: center;
    margin-top: 0;
    margin-right: 0.75rem;
  }
}
#group-edit-priv-key-body {
  line-height: 1.2;
  margin: 1.5rem 0.75rem;
}
#group-edit-priv-key-ok {
  margin: 0 0.75rem 0.75rem;
}
.group-edit-unsubscribed {
  margin-top: 1.5rem;
}
</style>
