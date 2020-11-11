<!--
RoleEdit displays the role viewing/editing page.
-->

<template lang="pug">
#role-edit(v-if='!role')
  SSpinner
SForm#role-edit(v-else, @submit='onSubmit', :submitLabel='submitLabel')
  SFInput#role-edit-name(
    label='Role name',
    trim,
    autofocus,
    v-model='role.name',
    :errorFn='nameError'
  )
  SFCheckGroup#role-edit-flags(label='Flags', :options='allFlags', v-model='flags')
  #role-edit-privs.form-item
    .role-edit-group.role-edit-heading Group
    .role-edit-privs.role-edit-heading
      | Privileges (
      a(href='#', @click='onShowKey') Key
      | )
    template(v-for='(g, i) in privs')
      .role-edit-group(v-text='g.name')
      PrivilegeMask.role-edit-privs(v-model='privs[i]')
  template(v-if='canDelete || canClone', #extraButtons)
    SButton(v-if='canClone', @click='onClone') Clone Role
    SButton(v-if='canDelete', variant='danger', @click='onDelete') Delete Role
  Modal(ref='keyModal', v-slot='{ close }')
    SForm(dialog, title='Privileges Key')
      #role-edit-priv-key-body.form-item
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
        SButton#role-edit-priv-key-ok(variant='primary', @click='close(null)') OK
  MessageBox(
    ref='deleteModal',
    title='Delete Role',
    cancelLabel='Keep',
    okLabel='Delete',
    variant='danger'
  )
    | Are you sure you want to delete this role? All associated data,
    | including privileges and role assignments, will be permanently lost.
</template>

<script lang="ts">
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
  SSpinner,
} from '../../base'

interface GetRoleEditPrivilege extends Privileges {
  id: number
  name: string
}
type GetRoleEditRole = {
  id: number
  name: string
  individual: boolean
  detail: boolean
  permViewClearances: boolean
  permEditClearances: boolean
}
type GetRoleEdit = {
  role: GetRoleEditRole
  canDelete: boolean
  privs: Array<GetRoleEditPrivilege>
}
type PostRole = {
  duplicateName?: boolean
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
    SSpinner,
  },
  setup() {
    // Handle page load and change.
    const route = useRoute()
    const router = useRouter()
    const role = ref(null as null | GetRoleEditRole)
    const privs = ref([] as Array<GetRoleEditPrivilege>)
    const canDelete = ref(false)
    watchEffect(() => {
      setPage({ title: route.params.rid === 'NEW' ? 'Create Role' : 'Edit Role' })
      const rid = route.params.clone || route.params.rid
      axios.get<GetRoleEdit>(`/api/roles/${rid}`).then((resp) => {
        role.value = resp.data.role
        privs.value = resp.data.privs
        canDelete.value = resp.data.canDelete
        flags.value.clear()
        allFlags.forEach((f) => {
          // @ts-ignore
          if (resp.data.role[f.value]) flags.value.add(f.value)
        })
      })
    })

    // Name field.
    let duplicateName = ref('')
    function nameError(lostFocus: boolean) {
      if (!lostFocus || !role.value) return ''
      if (!role.value.name) return 'The role name is required.'
      if (duplicateName.value === role.value.name) return 'Another role has this name.'
      return ''
    }

    // Flags field.
    const flags = ref(new Set() as Set<string>)
    const allFlags = [
      { value: 'individual', label: 'Individual (one person only)' },
      { value: 'detail', label: 'Hide in person list' },
      { value: 'permViewClearances', label: 'Can view clearances' },
      { value: 'permEditClearances', label: 'Can edit clearances' },
    ]

    // Clone.
    const canClone = computed(() => route.params.rid !== 'NEW')
    function onClone() {
      router.push({ name: 'roles-rid', params: { rid: 'NEW', clone: route.params.rid } })
    }

    // Delete.
    const deleteModal = ref(null as any)
    async function onDelete() {
      if (!(await deleteModal.value.show())) return
      const body = new FormData()
      body.append('delete', 'true')
      await axios.post(`/api/roles/${route.params.rid}`, body)
      router.push('/admin/roles')
    }

    // Submit.
    const submitLabel = computed(() => (route.params.rid === 'NEW' ? 'Create Role' : 'Save Role'))
    async function onSubmit() {
      if (!role.value) return
      const body = new FormData()
      body.append('name', role.value.name)
      allFlags.forEach((f) => {
        body.append(f.value, flags.value.has(f.value).toString())
      })
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
      const resp = (await axios.post<PostRole>(`/api/roles/${route.params.rid}`, body)).data
      if (resp) {
        if (resp.duplicateName) duplicateName.value = role.value.name
      } else {
        router.push('/admin/roles')
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
      flags,
      keyModal,
      nameError,
      onClone,
      onDelete,
      onShowKey,
      onSubmit,
      privs,
      role,
      submitLabel,
    }
  },
})
</script>

<style lang="postcss">
#role-edit {
  padding: 1.5rem 0.75rem;
}
#role-edit-privs {
  display: grid;
  grid: auto / 1fr;
  @media (min-width: 450px) {
    justify-content: start;
    grid: auto / auto min-content;
  }
}
.role-edit-heading {
  display: none;
  @media (min-width: 450px) {
    display: block;
    font-weight: bold;
  }
}
.role-edit-group {
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
#role-edit-priv-key-body {
  line-height: 1.2;
  margin: 1.5rem 0.75rem;
}
#role-edit-priv-key-ok {
  margin: 0 0.75rem 0.75rem;
}
</style>
