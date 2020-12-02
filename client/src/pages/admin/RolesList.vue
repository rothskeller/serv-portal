<!--
RolesList displays the list of roles.
-->

<template lang="pug">
#roles-list
  SSpinner(v-if='loading')
  template(v-else)
    table#roles-list-table(v-if='roles.length')
      tr(
        :style='dragOverStyle("TOP")',
        @dragenter='onDragEnter($event, "TOP")',
        @dragover='onDragOver($event, "TOP")',
        @dragleave='onDragLeave($event, "TOP")',
        @drop='onDrop($event, "TOP")'
      )
        td.roles-list-heading Org
        td.roles-list-heading Priv
        td.roles-list-heading Role
      tr(
        v-for='r in roles',
        draggable='true',
        :style='dragOverStyle(r)',
        @dragstart='onDragStart($event, r)',
        @dragenter='onDragEnter($event, r)',
        @dragover='onDragOver($event, r)',
        @dragleave='onDragLeave($event, r)',
        @drop='onDrop($event, r)',
        @dragend='onDragEnd()'
      )
        td(v-text='orgNames[r.org] || "—"')
        td(v-text='privLevelNames[r.privLevel] || "—"')
        td
          router-link(:to='`/admin/roles/${r.id}`', draggable='false', v-text='r.name')
          span.roles-list-people(v-text='` [${r.people}]`')
    div(v-else) No roles currently defined.
    #roles-list-buttons
      SButton#roles-list-saveOrder(
        v-if='orderChanged',
        variant='warning',
        :disabled='submitting',
        @click='onSaveOrder'
      ) Save Order
      SButton(variant='primary', :disabled='submitting', @click='onAdd') Add Role
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SButton, SSpinner } from '../../base'

interface GetRolesRole {
  id: number
  name: string
  org?: string
  privLevel?: string
  implicitOnly?: boolean
  people: number
}

const orgNames = {
  admin: 'Admin',
  'cert-d': 'CERT-D',
  'cert-t': 'CERT-T',
  listos: 'Listos',
  sares: 'SARES',
  snap: 'SNAP',
}
const privLevelNames = {
  student: 'Student',
  member: 'Member',
  leader: 'Leader',
}

export default defineComponent({
  components: { SButton, SSpinner },
  setup() {
    const router = useRouter()
    setPage({ title: 'Roles' })

    // Load page data.
    const loading = ref(true)
    const roles = ref([] as Array<GetRolesRole>)
    axios.get<Array<GetRolesRole>>('/api/roles').then((resp) => {
      roles.value = resp.data
      loading.value = false
    })

    function onAdd() {
      router.push(`/admin/roles/NEW`)
    }

    // Making order changes with drag and drop.
    const draggingOverRole = ref(null as null | 'TOP' | GetRolesRole)
    let draggingOverRoleCount = 0
    function onDragStart(evt: DragEvent, role: GetRolesRole) {
      evt.dataTransfer!.setData('x-serv-role', role.id.toString())
      evt.dataTransfer!.effectAllowed = 'move'
    }
    function onDragEnter(evt: DragEvent, role: 'TOP' | GetRolesRole) {
      if (!evt.dataTransfer!.types.includes('x-serv-role')) return
      evt.preventDefault()
      if (draggingOverRole.value === role) {
        draggingOverRoleCount++
      } else {
        draggingOverRole.value = role
        draggingOverRoleCount = 1
      }
    }
    function onDragOver(evt: DragEvent, role: 'TOP' | GetRolesRole) {
      if (!evt.dataTransfer!.types.includes('x-serv-role')) return
      evt.preventDefault()
    }
    function onDragLeave(evt: DragEvent, role: 'TOP' | GetRolesRole) {
      if (!evt.dataTransfer!.types.includes('x-serv-role')) return
      if (draggingOverRole.value === role && draggingOverRoleCount > 1) draggingOverRoleCount--
      else if (draggingOverRole.value === role) draggingOverRole.value = null
    }
    function onDrop(evt: DragEvent, role: 'TOP' | GetRolesRole) {
      const ridMoving = parseInt(evt.dataTransfer!.getData('x-serv-role'))
      const idxMoving = roles.value.findIndex((r) => r.id === ridMoving)
      let idxMoveTo = role === 'TOP' ? -1 : roles.value.findIndex((r) => r === role)
      if (idxMoving === idxMoveTo || idxMoving === idxMoveTo + 1) return
      const roleMoving = roles.value.splice(idxMoving, 1)
      if (idxMoveTo < idxMoving) idxMoveTo++
      roles.value.splice(idxMoveTo, 0, roleMoving[0])
      orderChanged.value = true
    }
    function onDragEnd() {
      draggingOverRole.value = null
    }
    function dragOverStyle(role: 'TOP' | GetRolesRole) {
      return role === draggingOverRole.value ? { borderBottom: '2px solid #888' } : null
    }

    // Saving order changes.
    const orderChanged = ref(false)
    const submitting = ref(false)
    async function onSaveOrder() {
      const body = new FormData()
      roles.value.forEach((r) => {
        body.append('role', r.id.toString())
      })
      submitting.value = true
      await axios.post<Array<GetRolesRole>>('/api/roles', body)
      submitting.value = false
      orderChanged.value = false
    }

    return {
      dragOverStyle,
      loading,
      onAdd,
      onDragEnd,
      onDragEnter,
      onDragLeave,
      onDragOver,
      onDragStart,
      onDrop,
      onSaveOrder,
      orderChanged,
      orgNames,
      privLevelNames,
      roles,
      submitting,
    }
  },
})
</script>

<style lang="postcss">
#roles-list {
  margin: 1.5rem 0.75rem;
}
#roles-list-table {
  td {
    padding-left: 1rem;
    vertical-align: middle;
    &:first-child {
      padding-left: 0;
    }
    .touch & {
      height: 40px;
    }
  }
}
.roles-list-heading {
  font-weight: bold;
}
.roles-list-people {
  color: #888;
}
#roles-list-buttons {
  margin-top: 1.5rem;
}
#roles-list-saveOrder {
  margin-right: 0.5rem;
}
</style>
