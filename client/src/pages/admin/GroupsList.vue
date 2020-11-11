<!--
GroupsList displays the list of groups.
-->

<template lang="pug">
#groups-list
  SSpinner(v-if='loading')
  #groups-list-table(v-else)
    .groups-list-name.groups-list-heading Group
    .groups-list-roles.groups-list-heading Included in Roles
    template(v-for='g in groups')
      .groups-list-name
        router-link(:to='`/admin/groups/${g.id}`', v-text='g.name || "(none)"')
      .groups-list-roles
        div(v-for='r in g.roles', v-text='r')
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import SSpinner from '../../base/SSpinner.vue'

type GetGroup = {
  id: number
  name: string
  roles: Array<string>
}

export default defineComponent({
  components: { SSpinner },
  setup() {
    setPage({ title: 'Groups' })
    const loading = ref(true)
    const groups = ref([] as Array<GetGroup>)
    axios.get<Array<GetGroup>>(`/api/groups`).then((resp) => {
      groups.value = resp.data
      loading.value = false
    })
    return { groups, loading }
  },
})
</script>

<style lang="postcss">
#groups-list {
  padding: 1.5rem 0.75rem;
}
#groups-list-table {
  display: grid;
  grid: auto / 1fr 1fr;
  @media (min-width: 576px) {
    grid: auto / 16rem 1fr;
  }
}
.groups-list-heading {
  font-weight: bold;
}
.groups-list-name {
  flex: none;
  margin: 0.75rem 0.75rem 0 0;
  font-variant: tabular-nums;
}
.groups-list-roles {
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
