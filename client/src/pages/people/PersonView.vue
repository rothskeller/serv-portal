<!--
PersonView displays the information about a person, in non-editable form.
-->

<template lang="pug">
#person-view-spinner(v-if='!person')
  SSpinner
#person-view(v-else)
  #person-view-grid
    PersonViewNames(:person='person', @reload='onReload')
    .person-view-spacer
    PersonViewRoles(:person='person', @reload='onReload', style='grid-row: 2/3')
    PersonViewContact(:person='person', @reload='onReload')
    PersonViewStatus(:person='person', @reload='onReload')
    PersonViewNotes(:person='person', @reload='onReload')
    PersonViewSubscriptions(:person='person', @reload='onReload')
    PersonViewSection(
      v-if='person.canChangePassword || person.canResetPassword',
      title='Password',
      :editable='false'
    )
      #person-view-chpw
        SButton(
          v-if='person.canChangePassword',
          variant='primary',
          small,
          @click='onChangePassword'
        ) Change Password
        SButton(v-if='person.canResetPassword', variant='primary', small, @click='onResetPassword') Reset Password
  PersonEditPassword(v-if='person.canChangePassword', ref='changePasswordModal', :pid='person.id')
  MessageBox(
    ref='resetPasswordModal',
    title='Reset Password',
    variant='warning',
    okLabel='Reset Password'
  )
    | This will replace {{ person.informalName }}’s current password with a new,
    | randomly generated one, and send it to them by email. You should do this
    | only if they have explicitly asked you to. Are you sure you want to reset
    | {{ person.informalName }}’s password?
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../../plugins/axios'
import { MessageBox, SButton, SSpinner } from '../../base'
import PersonEditPassword from './PersonEditPassword.vue'
import PersonViewContact from './PersonViewContact.vue'
import PersonViewNames from './PersonViewNames.vue'
import PersonViewNotes from './PersonViewNotes.vue'
import PersonViewRoles from './PersonViewRoles.vue'
import PersonViewSection from './PersonViewSection.vue'
import PersonViewStatus from './PersonViewStatus.vue'
import PersonViewSubscriptions from './PersonViewSubscriptions.vue'

export type GetPersonAddress = {
  address?: string
  sameAsHome: boolean
  latitude?: number
  longitude?: number
  fireDistrict: number
}
export type GetPersonRole = {
  title: string
  org: string
}
interface GetPersonNote {
  date: string
  note: string
}
interface GetPersonBGCheckAdmin {
  admin: true
  needed: string
  checks: Array<{
    type: Array<string>
    date?: string
    assumed: boolean
  }>
}
interface GetPersonBGCheckNonAdmin {
  admin: false
  needed: boolean
  cleared?: string
}
interface GetPersonStatus {
  canEdit?: true
  level: string
  volgistics: {
    needed: boolean
    id: number
  }
  dswCERT: {
    needed: boolean
    registered?: string
    expires?: string
    expired?: true
  }
  dswComm: {
    needed: boolean
    registered?: string
    expires?: string
    expired?: true
  }
  backgroundCheck: GetPersonBGCheckAdmin | GetPersonBGCheckNonAdmin
  identification: Array<string>
}
export interface GetPerson {
  id: number
  informalName: string
  formalName: string
  sortName: string
  callSign: string
  contact?: {
    email: string
    email2: string
    homeAddress: GetPersonAddress
    mailAddress: GetPersonAddress
    workAddress: GetPersonAddress
    cellPhone: string
    homePhone: string
    workPhone: string
  }
  roles: Array<GetPersonRole>
  lists?: Array<string>
  status?: GetPersonStatus
  notes: Array<GetPersonNote>
  canEdit: boolean
  canEditRoles: boolean
  canEditNotes: boolean
  canEditLists: boolean
  canChangePassword: boolean
  canResetPassword: boolean
  canHours: boolean
  noEmail: boolean
  noText: boolean
}

export default defineComponent({
  components: {
    MessageBox,
    PersonEditPassword,
    PersonViewContact,
    PersonViewNames,
    PersonViewNotes,
    PersonViewRoles,
    PersonViewSection,
    PersonViewStatus,
    PersonViewSubscriptions,
    SButton,
    SSpinner,
  },
  props: {
    onLoadPerson: { type: Function, required: true },
  },
  setup(props) {
    const route = useRoute()
    const person = ref(null as null | GetPerson)
    axios.get<GetPerson>(`/api/people/${route.params.id}`).then((resp) => {
      person.value = resp.data
      props.onLoadPerson(person.value)
    })
    async function onReload() {
      person.value = (await axios.get<GetPerson>(`/api/people/${route.params.id}`)).data
    }
    const changePasswordModal = ref(null as any)
    function onChangePassword() {
      changePasswordModal.value.show()
    }
    const resetPasswordModal = ref(null as any)
    function onResetPassword() {
      if (!resetPasswordModal.value.show()) return
      axios.post(`/api/people/${route.params.id}/pwreset`)
    }
    return {
      changePasswordModal,
      onChangePassword,
      onResetPassword,
      onReload,
      person,
      resetPasswordModal,
    }
  },
})
</script>

<style lang="postcss">
#person-view {
  margin: 1.5rem 0.75rem;
}
#person-view-grid {
  display: grid;
  grid: auto-flow / 1fr;
  @media (min-width: 740px) {
    grid: auto-flow / 1fr 1fr;
    column-gap: 0.75rem;
  }
  @media (min-width: 1048px) {
    grid: auto-flow / 1fr 1fr 1fr;
  }
}
.person-view-spacer {
  display: none;
  @media (min-width: 740px) {
    display: block;
  }
  @media (min-width: 1048px) {
    grid-column: span 2;
  }
}
#person-view-chpw {
  display: grid;
  grid: auto / auto auto;
  column-gap: 0.5rem;
  margin-top: 0.75rem;
}
</style>
