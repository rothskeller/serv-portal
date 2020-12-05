<!--
PersonViewStatus displays the volunteer status part of the person view page.
-->

<template lang="pug">
PersonViewSection(
  v-if='person.status',
  title='Volunteer Status',
  :editable='person.status.canEdit',
  @edit='onEditStatus'
)
  #person-view-status-grid
    template(v-if='person.status.level === "admin"')
      div Volgistics
      div(v-if='person.status.volgistics.id', v-text='`#${person.status.volgistics.id}`')
      div(v-else, :style='{ color: person.status.volgistics.needed ? "red" : null }') Not registered
    template(v-else-if='person.status.volgistics.needed && !person.status.volgistics.id')
      div City volunteer
      div(v-if='me.id === person.id')
        a(href='https://www.volgistics.com/ex/portal.dll/ap?AP=929478828', target='_blank') Please register
      div(v-else, style='color: red') Not registered
    template(v-else-if='nothingElseToShow')
      div(v-if='person.status.volgistics.id') Registered
      div(v-else) Not registered
    template(v-if='person.status.dswCERT.registered')
      div DSW CERT
      template(v-if='person.status.dswCERT.expired')
        div(:style='{ color: person.status.dswCERT.needed ? "red" : null }') Expired on {{ person.status.dswCERT.expires }}
      template(v-else-if='person.status.level == "admin"')
        div Registered {{ person.status.dswCERT.registered }}, expires&nbsp;{{ person.status.dswCERT.expires.replace(/-/g, "\u2011") }}
      template(v-else)
        div Expires on {{ person.status.dswCERT.expires }}
    template(v-else-if='person.status.dswCERT.needed')
      div DSW CERT
      div(style='color: red') Not registered
    template(v-if='person.status.dswComm.registered')
      div DSW SARES
      template(v-if='person.status.dswComm.expired')
        div(:style='{ color: person.status.dswComm.needed ? "red" : null }') Expired on {{ person.status.dswComm.expires }}
      template(v-else-if='person.status.level == "admin"')
        div Registered {{ person.status.dswComm.registered }}, expires&nbsp;{{ person.status.dswComm.expires.replace(/-/g, "\u2011") }}
      template(v-else)
        div Expires on {{ person.status.dswComm.expires }}
    template(v-else-if='person.status.dswComm.needed')
      div DSW SARES
      div(style='color: red') Not registered
    template(v-if='me.webmaster')
      template(v-if='person.status.backgroundCheck.cleared === "true"')
        div Background check
        div Cleared
      template(v-else-if='person.status.backgroundCheck.cleared')
        div Background check
        div Cleared on {{ person.status.backgroundCheck.cleared }}
      template(v-else-if='person.status.backgroundCheck.needed')
        div Background check
        div(style='color: red') Not cleared
    template(v-if='person.status.identification.length')
      div IDs issued
      #person-view-status-identification
        div(v-for='id in person.status.identification', v-text='id')
  PersonEditStatus(v-if='person.canEdit', ref='editStatusModal', :pid='person.id')
</template>

<script lang="ts">
import { computed, defineComponent, inject, PropType, Ref, ref } from 'vue'
import { LoginData } from '../../plugins/login'
import type { GetPerson } from './PersonView.vue'
import PersonEditStatus from './PersonEditStatus.vue'
import PersonViewSection from './PersonViewSection.vue'

export default defineComponent({
  components: { PersonEditStatus, PersonViewSection },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const me = inject<Ref<LoginData>>('me')!
    const nothingElseToShow = computed(
      () =>
        props.person.status &&
        !props.person.status.dswCERT.registered &&
        !props.person.status.dswCERT.needed &&
        !props.person.status.dswComm.registered &&
        !props.person.status.dswComm.needed &&
        (!me.value.webmaster ||
          (!props.person.status.backgroundCheck.cleared &&
            !props.person.status.backgroundCheck.needed)) &&
        !props.person.status.identification.length
    )
    const editStatusModal = ref(null as any)
    async function onEditStatus() {
      if (!(await editStatusModal.value.show())) return
      emit('reload')
    }
    return {
      editStatusModal,
      onEditStatus,
      me,
      nothingElseToShow,
    }
  },
})
</script>

<style lang="postcss">
#person-view-status-grid {
  margin-top: 0.75rem;
  display: grid;
  grid: auto-flow / 1fr;
  & > div:nth-child(2n + 1) {
    white-space: nowrap;
  }
  & > div:nth-child(2n) {
    margin-left: 2rem;
  }
  @media (min-width: 400px) {
    grid: auto-flow / min-content 1fr;
    & > div:nth-child(2n) {
      margin-left: 1rem;
    }
  }
  @media (min-width: 740px) {
    grid: auto-flow / 1fr;
    & > div:nth-child(2n) {
      margin-left: 2rem;
    }
  }
}
#person-view-status-identification {
  display: flex;
  flex-wrap: wrap;
  & div {
    white-space: nowrap;
    &::after {
      content: ', '; /* non-breaking space */
    }
    &:last-child::after {
      content: '';
    }
  }
}
</style>
