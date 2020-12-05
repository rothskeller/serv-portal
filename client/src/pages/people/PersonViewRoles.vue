<!--
PersonViewRoles displays a person's roles.
-->

<template lang="pug">
PersonViewSection(
  v-if='badges.length || person.canEditRoles',
  :title='`SERV ${person.roles.length === 1 ? "Role" : "Roles"}`',
  :editable='person.canEditRoles',
  @edit='onEditRoles'
)
  #person-view-roles-badges
    template(v-for='org in badges')
      OrgBadge.person-view-roles-badge(:org='org.org')
      .person-view-roles-titles
        .person-view-roles-title(v-for='title in org.titles', v-text='title')
    .person-view-roles-title(v-if='!badges.length') No current role in any SERV org.
  PersonEditRoles(v-if='person.canEditRoles', ref='editRolesModal', :pid='person.id')
</template>

<script lang="ts">
import { defineComponent, PropType, ref, watch, watchEffect } from 'vue'
import { OrgBadge } from '../../base'
import { GetPerson, GetPersonRole } from './PersonView.vue'
import PersonEditRoles from './PersonEditRoles.vue'
import PersonViewSection from './PersonViewSection.vue'

function badgeForRole(role: GetPersonRole): string {
  if (role.org === 'admin' && role.title.startsWith('OES')) return 'dps'
  if (role.org === 'admin') return 'serv'
  return role.org.replace(/-.*/, '')
}

export default defineComponent({
  components: { OrgBadge, PersonEditRoles, PersonViewSection },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const badges = ref([] as Array<{ org: string; titles: Array<string> }>)
    watchEffect(() => {
      badges.value = []
      props.person.roles.forEach((r) => {
        const idx = badges.value.findIndex((o) => o.org === badgeForRole(r))
        if (idx < 0) badges.value.push({ org: badgeForRole(r), titles: [r.title] })
        else badges.value[idx].titles.push(r.title)
      })
    })
    const editRolesModal = ref(null as any)
    async function onEditRoles() {
      if (!(await editRolesModal.value.show())) return
      emit('reload')
    }
    return {
      editRolesModal,
      onEditRoles,
      badges,
    }
  },
})
</script>

<style lang="postcss">
#person-view-roles-badges {
  display: grid;
  grid: auto-flow / min-content 1fr;
  min-width: 0;
}
.person-view-roles-badge {
  margin-top: 1rem;
}
.person-view-roles-titles {
  margin-top: 1rem;
  line-height: 1.2;
  min-width: 0;
}
.person-view-roles-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
