<!--
PersonEdit displays the editor page for a person.
-->

<template lang="pug">
#person-edit-spinner(v-if='!person')
  SSpinner
SForm(v-else, :submitLabel='submitLabel', @submit='onSubmit')
  template(v-if='canEditDetails')
    .form-item.person-edit-block-head Identification
    SFInput#person-informalName(
      label='Name',
      autofocus,
      trim,
      v-model='person.informalName',
      :errorFn='informalNameError',
      help='What you like to be called, e.g. “Joe Banks”'
    )
    SFInput#person-formalName(
      label='Formal name',
      trim,
      v-model='person.formalName',
      :errorFn='formalNameError',
      help='For formal documents, e.g. “Joseph A. Banks, Jr.”'
    )
    SFInput#person-sortName(
      label='Sort name',
      trim,
      v-model='person.sortName',
      :errorFn='sortNameError',
      help='For appearance in sorted lists, e.g. “Banks, Joe”'
    )
    SFInput#person-username(
      v-if='canEditUsername',
      label='Username',
      trim,
      v-model='person.username',
      :errorFn='usernameError',
      help='Usually the same as primary email address',
      style='text-transform:lowercase'
    )
    SFInput#person-callSign(
      label='Call sign',
      trim,
      v-model='person.callSign',
      :errorFn='callSignError',
      help='FCC amateur radio license',
      style='text-transform:uppercase'
    )
  template(v-if='canEditClearances')
    .form-item.person-edit-block-head Volunteer Status
    SFInput#person-volgistics(
      type='number',
      min='0',
      label='Volgistics ID',
      v-model='volgistics',
      :errorFn='volgisticsError'
    )
    PersonEditDSW(
      v-for='cls in Object.keys(person.dsw).sort()',
      :key='cls',
      :type='cls',
      v-model='person.dsw[cls]'
    )
    SFInput#person-background(
      label='BG Check',
      trim,
      v-model='person.backgroundCheck',
      :errorFn='backgroundError',
      help='Date when background check cleared, or “TRUE” if clearance confirmed but date unknown',
      style='text-transform:uppercase'
    )
  template(v-if='canEditDetails')
    .form-item.person-edit-block-head Change Password
    SFInput#person-oldPassword(
      v-if='!allowBadPassword',
      type='password',
      label='Old Password',
      v-model='oldPassword',
      :errorFn='oldPasswordError'
    )
    SFPassword#person-password(
      label='New Password',
      v-model='password',
      :allowBadPassword='allowBadPassword',
      :passwordHints='myPasswordHints'
    )
    .form-item.person-edit-block-head Contact Information
    SFInput#person-email(
      label='Email',
      trim,
      v-model='person.email',
      :errorFn='emailError',
      style='text-transform:lowercase'
    )
    SFInput#person-email2(
      label='Alt. Email',
      trim,
      v-model='person.email2',
      :errorFn='email2Error',
      style='text-transform:lowercase'
    )
    SFInput#person-cellPhone(
      label='Cell Phone',
      trim,
      v-model='person.cellPhone',
      :errorFn='cellPhoneError'
    )
    SFInput#person-homePhone(
      label='Home Phone',
      trim,
      v-model='person.homePhone',
      :errorFn='homePhoneError'
    )
    SFInput#person-workPhone(
      label='Work Phone',
      trim,
      v-model='person.workPhone',
      :errorFn='workPhoneError'
    )
    PersonEditAddress(type='Home', v-model='person.homeAddress')
    PersonEditAddress(
      type='Work',
      v-model='person.workAddress',
      :hasHome='!!person.homeAddress.address'
    )
    PersonEditAddress(
      type='Mail',
      v-model='person.mailAddress',
      :hasHome='!!person.homeAddress.address'
    )
    .form-item.person-edit-block-head Roles
  SFCheckGroup#person-roles(
    v-if='canEditRoles',
    label='Roles',
    v-model='roles',
    :options='person.roles',
    valueKey='id',
    labelKey='name',
    enabledKey='canAssign',
    :errorFn='rolesError'
  )
  template(v-else)
    .form-item-label Roles
    .form-item-input
      template(v-for='role in person.roles')
        div(v-if='role.held', v-text='role.name')
</template>

<script lang="ts">
import { computed, defineComponent, inject, Ref, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import { SFCheckGroup, SFInput, SFPassword, SForm, SSpinner } from '../../base'
import PersonEditAddress from './PersonEditAddress.vue'
import PersonEditDSW from './PersonEditDSW.vue'
import { GetPersonPersonBase, GetPersonRole } from './PersonView.vue'
import { LoginData } from '../../plugins/login'

interface GetPersonEditRole extends GetPersonRole {
  canAssign: boolean
  held: boolean
}
interface GetPersonEditPerson extends GetPersonPersonBase {
  roles: Array<GetPersonEditRole>
  volgistics?: number
  dsw?: Record<string, string>
  backgroundCheck?: string
}
interface GetPersonEditBase {
  person: GetPersonEditPerson
  canEditRoles: boolean
  canEditDetails: boolean
  canEditClearances: boolean
  allowBadPassword: boolean
  canEditUsername: boolean
}
interface GetPersonEditNED extends GetPersonEditBase {
  canEditDetails: false
}
interface GetPersonEditED extends GetPersonEditBase {
  canEditDetails: true
  passwordHints?: Array<string>
}
type GetPersonEdit = GetPersonEditNED | GetPersonEditED

type PostPersonEdit = {
  duplicateCallSign?: boolean
  duplicateCellPhone?: boolean
  duplicateSortName?: boolean
  duplicateUsername?: boolean
  wrongOldPassword?: boolean
}

export default defineComponent({
  components: {
    PersonEditAddress,
    PersonEditDSW,
    SFCheckGroup,
    SFInput,
    SFPassword,
    SForm,
    SSpinner,
  },
  props: {
    onLoadPerson: { type: Function, required: true },
  },
  setup(props) {
    // Get editing capabilities, and details of the person to be edited.
    const route = useRoute()
    const router = useRouter()
    const newPerson = route.params.id === 'NEW'
    const allowBadPassword = ref(false)
    const canEditClearances = ref(false)
    const canEditDetails = ref(false)
    const canEditRoles = ref(false)
    const canEditUsername = ref(false)
    let passwordHints: Array<string>
    const person = ref(null as null | GetPersonEditPerson)
    const volgistics = ref('')
    const roles = ref(new Set() as Set<number>)
    axios.get<GetPersonEdit>(`/api/people/${route.params.id}?edit=1`).then((resp) => {
      allowBadPassword.value = resp.data.allowBadPassword
      canEditClearances.value = resp.data.canEditClearances
      canEditDetails.value = resp.data.canEditDetails
      canEditRoles.value = resp.data.canEditRoles
      canEditUsername.value = resp.data.canEditUsername
      if (resp.data.canEditDetails) passwordHints = resp.data.passwordHints || []
      if (resp.data.person.backgroundCheck === 'true') resp.data.person.backgroundCheck = 'TRUE'
      person.value = resp.data.person
      volgistics.value = resp.data.person.volgistics ? resp.data.person.volgistics.toString() : ''
      roles.value = new Set(resp.data.person.roles.filter((r) => r.held).map((r) => r.id))
      props.onLoadPerson(person.value)
    })

    // When checking whether a password is good, we want to disallow inclusion
    // of personal data (their name, email, etc.).  We use the current values
    // from the form to do that.
    const myPasswordHints = computed(() => {
      const hints = [...passwordHints]
      if (!person.value) return hints
      if (person.value.informalName) hints.push(person.value.informalName)
      if (person.value.formalName) hints.push(person.value.formalName)
      if (person.value.callSign) hints.push(person.value.callSign)
      if (person.value.username) hints.push(person.value.username)
      if (person.value.email) hints.push(person.value.email)
      if (person.value.email2) hints.push(person.value.email2)
      if (person.value.homeAddress && person.value.homeAddress.address)
        hints.push(person.value.homeAddress.address)
      if (person.value.mailAddress && person.value.mailAddress.address)
        hints.push(person.value.mailAddress.address)
      if (person.value.workAddress && person.value.workAddress.address)
        hints.push(person.value.workAddress.address)
      if (person.value.cellPhone) hints.push(person.value.cellPhone)
      if (person.value.homePhone) hints.push(person.value.homePhone)
      if (person.value.workPhone) hints.push(person.value.workPhone)
      return hints
    })

    // These hold references to errors that come back from a submission
    // attempt.  See onSubmit() below for where they are set.
    const duplicateCallSign = ref('')
    const duplicateCellPhone = ref('')
    const duplicateSortName = ref('')
    const duplicateUsername = ref('')
    const wrongOldPassword = ref('')

    // Validation functions for the various input fields (except password and
    // address fields; those specialized input controls have their validation
    // functions built in).
    function informalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value?.informalName) return 'A name is required.'
      return ''
    }
    function formalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value?.formalName) return 'A name is required.'
      return ''
    }
    function sortNameError(lostFocus: boolean) {
      if (!person.value?.sortName) return lostFocus ? 'A name is required.' : ''
      if (duplicateSortName.value === person.value?.sortName)
        return 'A different person has this name.'
      return ''
    }
    function usernameError() {
      if (!person.value?.username) return ''
      if (duplicateUsername.value === person.value.username)
        return 'A different person has this username.'
      return ''
    }
    function callSignError(lostFocus: boolean) {
      if (!person.value?.callSign) return ''
      if (lostFocus && !person.value.callSign.match(/^[AKNW][A-Z]?[0-9][A-Z]{1,3}$/i))
        return 'This is not a valid call sign.'
      if (duplicateCallSign.value === person.value?.callSign)
        return 'A different person has this call sign.'
      return ''
    }
    function volgisticsError(lostFocus: boolean) {
      if (!lostFocus || !volgistics.value) return ''
      if (parseInt(volgistics.value) < 1) return 'This is not a valid Volgistics ID number.'
      return ''
    }
    function backgroundError(lostFocus: boolean) {
      if (
        !lostFocus ||
        !person.value?.backgroundCheck ||
        person.value?.backgroundCheck.toUpperCase() === 'TRUE'
      )
        return ''
      if (!person.value.backgroundCheck.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }
    const oldPassword = ref('')
    const password = ref('')
    function oldPasswordError(lostFocus: boolean, submitted: boolean) {
      if (!submitted) return ''
      if (password.value && !oldPassword.value && !allowBadPassword.value)
        return 'You must supply your old password in order to change your password.'
      if (oldPassword.value && oldPassword.value === wrongOldPassword.value)
        return 'This is not the correct old password.'
      return ''
    }
    function emailError(lostFocus: boolean) {
      if (!lostFocus || !person.value?.email) return ''
      if (
        !person.value.email.match(
          /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
        )
      )
        return 'This is not a valid email address.'
      return ''
    }
    function email2Error(lostFocus: boolean) {
      if (!lostFocus || !person.value?.email2) return ''
      if (
        !person.value.email2.match(
          /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
        )
      )
        return 'This is not a valid email address.'
      if (person.value.email === person.value.email2)
        return 'The two email addresses should not be the same.  (Leave this field empty if you only have one.)'
      return ''
    }
    function cellPhoneError(lostFocus: boolean) {
      if (!person.value?.cellPhone) return ''
      if (lostFocus && person.value.cellPhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      if (duplicateCellPhone.value === person.value.cellPhone)
        return 'A different person has this cell phone number.'
      return ''
    }
    function homePhoneError(lostFocus: boolean) {
      if (!lostFocus || !person.value?.homePhone) return ''
      if (person.value.homePhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function workPhoneError(lostFocus: boolean) {
      if (!lostFocus || !person.value?.workPhone) return ''
      if (person.value.workPhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function rolesError(lostFocus: boolean, submitted: boolean) {
      if (!submitted || !newPerson) return ''
      if (roles.value.size === 0) return 'At least one role must be selected.'
      return ''
    }

    // When the informal name is changed, we may update the formal and/or
    // sort name.
    function informalToSort(n: string | undefined): string | undefined {
      if (!n) return n
      const parts = n.split(/\s+/, 2)
      return parts.length > 1 ? `${parts[1]}, ${parts[0]}` : n
    }
    watch(
      () => person.value?.informalName,
      (n, o) => {
        if (!person.value) return
        if (person.value.formalName === o) person.value.formalName = n!
        if (person.value.sortName === informalToSort(o!)) person.value.sortName = informalToSort(n)!
      }
    )

    // When the email address is changed, we may update the username.
    watch(
      () => person.value?.email,
      (n, o) => {
        if (person.value && person.value.username === o) person.value.username = n!
      }
    )

    // Label for the Submit button.
    const me = inject<Ref<LoginData>>('me')!
    const submitLabel =
      (route.params.id as string) === me.value.id.toString()
        ? 'Save Changes'
        : newPerson
        ? 'Create Person'
        : 'Save Person'

    // Handle submission of the form.
    async function onSubmit() {
      if (!person.value) return // squelch TypeScript warnings
      const body = new FormData()
      if (canEditDetails.value) {
        body.append('informalName', person.value.informalName)
        body.append('formalName', person.value.formalName)
        body.append('sortName', person.value.sortName)
        body.append('username', person.value.username.toLowerCase())
        body.append('callSign', person.value.callSign.toUpperCase())
        body.append('email', (person.value.email || person.value.email2 || '').toLowerCase())
        body.append('email2', person.value.email ? (person.value.email2 || '').toLowerCase() : '')
        body.append('cellPhone', person.value.cellPhone || '')
        body.append('homePhone', person.value.homePhone || '')
        body.append('workPhone', person.value.workPhone || '')
        if (oldPassword.value) body.append('oldPassword', oldPassword.value)
        if (password.value) body.append('password', password.value)
        if (person.value.homeAddress?.address) {
          body.append('homeAddress', person.value.homeAddress.address)
          body.append('homeAddressLatitude', (person.value.homeAddress.latitude || 0).toString())
          body.append('homeAddressLongitude', (person.value.homeAddress.longitude || 0).toString())
        }
        if (person.value.workAddress?.address) {
          body.append('workAddress', person.value.workAddress.address)
          body.append('workAddressLatitude', (person.value.workAddress.latitude || 0).toString())
          body.append('workAddressLongitude', (person.value.workAddress.longitude || 0).toString())
        } else {
          body.append(
            'workAddressSameAsHome',
            (person.value.workAddress?.sameAsHome || false).toString()
          )
        }
        if (person.value.mailAddress?.address) {
          body.append('mailAddress', person.value.mailAddress.address)
        } else {
          body.append(
            'mailAddressSameAsHome',
            (person.value.mailAddress?.sameAsHome || false).toString()
          )
        }
      }
      if (canEditClearances.value) {
        body.append('volgistics', volgistics.value)
        body.append('backgroundCheck', (person.value.backgroundCheck || '').toLowerCase())
        Object.keys(person.value.dsw!).forEach((k) => {
          body.append(`dsw-${k}`, person.value!.dsw![k])
        })
      }
      if (canEditRoles.value) {
        person.value.roles
          .filter((role) => role.canAssign && roles.value.has(role.id))
          .forEach((role) => {
            body.append('role', role.id.toString())
          })
      }
      const resp = (await axios.post<PostPersonEdit>(`/api/people/${route.params.id}`, body)).data
      if (resp) {
        if (resp.duplicateSortName) duplicateSortName.value = person.value.sortName
        if (resp.duplicateUsername) duplicateUsername.value = person.value.username
        if (resp.duplicateCallSign) duplicateCallSign.value = person.value.callSign
        if (resp.duplicateCellPhone) duplicateCellPhone.value = person.value.cellPhone!
        if (resp.wrongOldPassword) wrongOldPassword.value = oldPassword.value
        // disregarding resp.weakPassword since we catch that locally
      } else {
        router.push('/people')
      }
    }

    return {
      allowBadPassword,
      backgroundError,
      callSignError,
      canEditClearances,
      canEditDetails,
      canEditRoles,
      canEditUsername,
      cellPhoneError,
      emailError,
      email2Error,
      formalNameError,
      homePhoneError,
      informalNameError,
      myPasswordHints,
      oldPassword,
      oldPasswordError,
      password,
      person,
      onSubmit,
      roles,
      rolesError,
      sortNameError,
      submitLabel,
      usernameError,
      volgistics,
      volgisticsError,
      workPhoneError,
    }
  },
  data: () => ({
    person: null,
    allowBadPassword: false,
    canEditClearances: false,
    canEditDetails: false,
    canEditRoles: false,
    canEditUsername: false,
    passwordHints: [],
    informalNameError: null,
    formalNameError: null,
    sortNameError: null,
    duplicateSortName: null,
    usernameError: null,
    duplicateUsername: null,
    callSignError: null,
    volgisticsError: null,
    dswError: {},
    backgroundError: null,
    duplicateCallSign: null,
    emailError: null,
    email2Error: null,
    cellPhoneError: null,
    duplicateCellPhone: null,
    homePhoneError: null,
    workPhoneError: null,
    oldPassword: null,
    oldPasswordError: null,
    wrongOldPassword: null,
    password: '',
    rolesError: null,
    submitted: false,
    suggestions: null,
  }),
})
</script>

<style lang="postcss">
#person-edit-spinner {
  margin: 1.5rem 0.75rem;
}
.person-edit-block-head {
  overflow: hidden;
  margin-top: 2rem;
  margin-bottom: 1rem;
  padding-left: 2rem;
  max-width: 28rem;
  color: #888;
  &:first-child {
    margin-top: 0;
  }
  &::before {
    display: inline-block;
    margin-right: 0.5rem;
    margin-left: -100%;
    width: 100%;
    border-top: 1px solid #888;
    content: '';
    vertical-align: middle;
  }
  &::after {
    display: inline-block;
    margin-right: -100%;
    margin-left: 0.5rem;
    width: 100%;
    border-top: 1px solid #888;
    content: '';
    vertical-align: middle;
  }
}
.person-edit-label {
  width: 8rem;
}
.person-edit-input {
  min-width: 14rem;
  max-width: 20rem;
}
.person-edit-label-input {
  margin-top: 0.25rem;
  min-width: 14rem;
  max-width: 20rem;
  @media (min-width: 41.75rem) {
    display: inline;
    margin-top: 0;
    margin-left: 0.25rem;
    min-width: 6rem;
    width: 6rem;
  }
}
</style>
