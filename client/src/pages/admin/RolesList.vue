<!--
RolesList displays the list of roles.
-->

<template lang="pug">
#roles-list
  SSpinner(v-if='loading')
  #roles-list-table(v-else)
    .roles-list-name.roles-list-heading Role
    .roles-list-groups.roles-list-heading Member of Groups
    template(v-for='r in roles')
      .roles-list-name
        router-link(:to='`/admin/roles/${r.id}`', v-text='r.name')
      .roles-list-groups
        div(v-for='g in r.groups', v-text='g')
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import SSpinner from '../../base/SSpinner.vue'

type GetRole = {
  id: number
  name: string
  groups: Array<string>
}

export default defineComponent({
  components: { SSpinner },
  setup() {
    setPage({ title: 'Roles' })
    const loading = ref(true)
    const roles = ref([] as Array<GetRole>)
    axios.get<Array<GetRole>>(`/api/roles`).then((resp) => {
      roles.value = resp.data
      loading.value = false
    })
    return { roles, loading }
  },
})
</script>

<style lang="postcss">
#roles-list {
  padding: 1.5rem 0.75rem;
}
#roles-list-table {
  display: grid;
  grid: auto / 1fr 1fr;
  @media (min-width: 576px) {
    grid: auto / 16rem 1fr;
  }
}
.roles-list-heading {
  font-weight: bold;
}
.roles-list-name {
  flex: none;
  margin: 0.75rem 0.75rem 0 0;
  font-variant: tabular-nums;
}
.roles-list-groups {
  flex: none;
  overflow: hidden;
  margin-top: 0.75rem;
  white-space: nowrap;
  div {
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
