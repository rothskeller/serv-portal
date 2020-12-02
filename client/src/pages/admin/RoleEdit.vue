<!--
RoleEdit displays the page for creating and editing roles.
-->

<template lang="pug">
#role
  SSpinner(v-if='loading')
  template(v-else)
    SForm(:submitLabel='submitLabel', :disabled='submitting', @submit='onSubmit')
      SFInput#role-name(
        label='Name',
        help='Collective name for those who hold the role',
        trim,
        v-model='role.name',
        :errorFn='nameError'
      )
      SFInput#role-title(
        label='Title',
        help='Name for a single person who holds this role, or empty',
        trim,
        v-model='role.title',
        :errorFn='titleError'
      )
      SFRadioGroup#role-org(label='Org', :options='orgOptions', v-model='role.org')
      SFRadioGroup#role-privLevel(
        label='Priv Level',
        :options='role.org ? privLevelOptions : privLevelNoOption',
        v-model='role.privLevel'
      )
      SFCheckGroup#role-flags(label='Flags', :options='flagOptions', v-model='flags')
      SFCheckGroup#role-implies(
        v-if='role.impliable.length',
        label='Implies',
        :options='role.impliable',
        valueKey='id',
        labelKey='name',
        v-model='implies'
      )
      label#role-lists-label.form-item-label Lists
      .form-item-input2
        template(v-for='l in role.lists')
          div(v-if='l.subModel || l.sender')
            a(href='#', @click.prevent='editList(l)') {{ l.type === "sms" ? `SMS: ${l.name}` : `${l.name}@SunnyvaleSERV.org` }}
            | : {{ l.subModel ? subModelNames[l.subModel] : "" }}{{ l.subModel && l.sender ? ", " : "" }}{{ l.sender ? "can send" : "" }}
        div: a(href='#', @click.prevent='editList(null)') Add List
      template(v-if='role.id', #extraButtons)
        SButton(@click='onDelete', variant='danger', :disabled='submitting') Delete Role
      Modal(ref='listEditModal', v-slot='{ close }')
        SForm(
          dialog,
          variant='primary',
          title='Connect to List',
          submitLabel='OK',
          @submit='close(true)',
          @cancel='close(false)'
        )
          SFSelect#role-listEdit-list(
            label='List',
            :options='selectableLists',
            valueKey='id',
            labelKey='nameFmt',
            v-model='editingList.list'
          )
          SFRadioGroup#role-listEdit-subModel(
            label='Subscription',
            :options='subModelOptions',
            v-model='editingList.subModel'
          )
          SFCheck#role-listEdit-sender(label='Can Send', v-model='editingList.sender')
      MessageBox(
        ref='deleteModal',
        title='Delete Role',
        cancelLabel='Keep',
        okLabel='Delete',
        variant='danger'
      )
        | Are you sure you want to delete this role? All associated data,
        | including role assignments, list associations, etc. will be
        | permanently lost.
</template>

<script lang="ts">
import { computed, defineComponent, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import {
  MessageBox,
  Modal,
  SButton,
  SForm,
  SFCheck,
  SFCheckGroup,
  SFInput,
  SFRadioGroup,
  SFSelect,
  SSpinner,
} from '../../base'

interface GetRoleImpliableRole {
  id: number
  name: string
}
interface GetRoleList {
  id: number
  type: string
  name: string
  subModel: string
  sender: boolean
  nameFmt: string // added locally
}
interface GetRole {
  id: number
  name: string
  title: string
  org: string
  privLevel: string
  showRoster: boolean
  implicitOnly: boolean
  implies: Array<number>
  impliable: Array<GetRoleImpliableRole>
  lists: Array<GetRoleList>
}
interface PostRole {
  duplicateName?: true
  duplicateTitle?: true
}
interface EditingList {
  list: number
  subModel: string
  sender: boolean
}

const orgOptions = [
  { value: 'admin', label: 'Admin' },
  { value: 'cert-d', label: 'CERT-D' },
  { value: 'cert-t', label: 'CERT-T' },
  { value: 'listos', label: 'Listos' },
  { value: 'sares', label: 'SARES' },
  { value: 'snap', label: 'SNAP' },
]
const privLevelOptions = [
  { value: '', label: '(none)' },
  { value: 'student', label: 'Student' },
  { value: 'member', label: 'Member' },
  { value: 'leader', label: 'Leader' },
]
const privLevelNoOption = [{ value: '', label: '(none)' }]
const flagOptions = [
  { value: 'showRoster', label: 'Available choice on People list page' },
  { value: 'implicitOnly', label: 'Role can only be implied, not assigned' },
]
const subModelNames = {
  allow: 'manual subscr.',
  auto: 'auto subscr.',
  warn: 'warn on unsub.',
}
const subModelOptions = [
  { value: '', label: 'No Subscription' },
  { value: 'allow', label: 'Manual Subscription' },
  { value: 'auto', label: 'Automatic Subscription' },
  { value: 'warn', label: 'Warn on Unsubscribe' },
]

export default defineComponent({
  components: {
    MessageBox,
    Modal,
    SButton,
    SForm,
    SFCheck,
    SFCheckGroup,
    SFInput,
    SFRadioGroup,
    SFSelect,
    SSpinner,
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: route.params.lid === 'NEW' ? 'Add Role' : 'Edit Role' })
    const submitLabel = route.params.lid === 'NEW' ? 'Add Role' : 'Save Role'

    // Load the page data.
    const loading = ref(true)
    const role = ref({} as GetRole)
    const flags = ref(new Set<string>())
    const implies = ref(new Set<number>())
    axios.get<GetRole>(`/api/roles/${route.params.rid}`).then((resp) => {
      role.value = resp.data
      flags.value.clear()
      if (role.value.showRoster) flags.value.add('showRoster')
      if (role.value.implicitOnly) flags.value.add('implicitOnly')
      implies.value = new Set(role.value.implies)
      loading.value = false
    })

    // The name field.
    const duplicateName = ref('')
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!role.value.name) return 'The role name is required.'
      if (duplicateName.value === role.value.name) return 'Another role has this name.'
      return ''
    }

    // The title field.
    const duplicateTitle = ref('')
    function titleError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (duplicateTitle.value && duplicateTitle.value === role.value.title)
        return 'Another role has this title.'
      return ''
    }
    watch(
      computed(() => role.value.name),
      (n, o) => {
        if (role.value.title === o) role.value.title = n
      }
    )

    // Submit the changes.
    const submitting = ref(false)
    async function onSubmit() {
      const body = new FormData()
      body.append('name', role.value.name)
      body.append('title', role.value.title)
      body.append('org', role.value.org)
      body.append('privLevel', role.value.org ? role.value.privLevel : '')
      body.append('showRoster', flags.value.has('showRoster').toString())
      body.append('implicitOnly', flags.value.has('implicitOnly').toString())
      implies.value.forEach((v) => {
        body.append('implies', v.toString())
      })
      role.value.lists.forEach((l) => {
        if (l.subModel || l.sender) body.append('lists', `${l.id}:${l.subModel}:${l.sender}`)
      })
      submitting.value = true
      const resp = (await axios.post<PostRole>(`/api/roles/${route.params.rid}`, body)).data
      submitting.value = false
      if (resp) {
        if (resp.duplicateName) duplicateName.value = role.value.name
        if (resp.duplicateTitle) duplicateTitle.value = role.value.title
      } else {
        router.push('/admin/roles')
      }
    }

    // List editing.
    const listEditModal = ref(null as any)
    const editingList = ref({} as EditingList)
    const selectableLists = ref([] as Array<GetRoleList>)
    async function editList(list: null | GetRoleList) {
      selectableLists.value = role.value.lists.filter(
        (l) => l === list || (!l.subModel && !l.sender)
      )
      selectableLists.value.forEach((l) => {
        l.nameFmt = l.type === 'sms' ? `SMS: ${l.name}` : `${l.name}@SunnyvaleSERV.org`
      })
      if (list) {
        editingList.value.list = list.id
        editingList.value.subModel = list.subModel
        editingList.value.sender = list.sender
      } else {
        if (!selectableLists.value.length) return
        editingList.value.list = selectableLists.value[0].id
        editingList.value.subModel = ''
        editingList.value.sender = false
      }
      if (!(await listEditModal.value.show())) return
      if (list && list.id !== editingList.value.list) {
        list.subModel = ''
        list.sender = false
      }
      list = role.value.lists.find((l) => l.id === editingList.value.list)!
      list.subModel = editingList.value.subModel
      list.sender = editingList.value.sender
    }

    // Handle deletions of roles.
    const deleteModal = ref(null as any)
    async function onDelete() {
      if (deleteModal.value) {
        const confirmed: boolean = await deleteModal.value.show()
        if (confirmed) {
          submitting.value = true
          await axios.delete(`/api/roles/${route.params.rid}`)
          submitting.value = false
          router.push('/admin/roles')
        }
      }
    }

    return {
      deleteModal,
      editList,
      editingList,
      flags,
      flagOptions,
      implies,
      listEditModal,
      loading,
      nameError,
      onDelete,
      onSubmit,
      orgOptions,
      privLevelNoOption,
      privLevelOptions,
      role,
      selectableLists,
      submitLabel,
      submitting,
      subModelNames,
      subModelOptions,
      titleError,
    }
  },
})
</script>

<style lang="postcss">
#role {
  margin: 1.5rem 0.75rem;
}
</style>
