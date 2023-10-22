<!--
PersonRegisterVolgistics is the dialog box for registering as a city volunteer.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Register as a City Volunteer',
    submitLabel='Register',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
      #person-volreg-intro.form-item.
        Thank you for your interest in volunteering with the City of Sunnyvale,
        Office of Emergency Services.  Please complete this form to register as
        a City of Sunnyvale Volunteer.  Once we receive your registration
        (which usually takes a few days) we will contact you to schedule an
        appointment for your fingerprinting.
      SFInput#person-informalName(
        label='Name',
        trim,
        autofocus,
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
      SFInput#person-callSign(
        label='Call sign',
        trim,
        v-model='person.callSign',
        :errorFn='callSignError',
        help='FCC amateur radio license',
        style='text-transform: uppercase'
      )
      SFInput#person-birthdate(
        label='Birthdate',
        type='date',
        v-model='person.birthdate',
        :errorFn='birthdateError'
      )
      SFInput#person-email(
        label='Email',
        help='This is the email address you log in with.',
        trim,
        autofocus,
        v-model='person.email',
        :errorFn='emailError',
        style='text-transform: lowercase'
      )
      SFInput#person-email2(
        label='Alt. Email',
        trim,
        v-model='person.email2',
        :errorFn='email2Error',
        style='text-transform: lowercase'
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
      PersonEditAddress(type='Home', required, v-model='person.homeAddress')
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
      SFCheckGroup#person-interests(
        label='Interests',
        :options='interestOptions',
        v-model='interests'
      )
      .form-item.person-heading Emergency Contact 1
      SFInput#person-emContact1Name(
        label='Name',
        trim,
        v-model='person.emContacts[0].name',
        :errorFn='emContact1NameError'
      )
      SFInput#person-emContact1HomePhone(
        label='Home Phone',
        trim,
        v-model='person.emContacts[0].homePhone',
        :errorFn='emContact1HomePhoneError'
      )
      SFInput#person-emContact1CellPhone(
        label='Cell Phone',
        trim,
        v-model='person.emContacts[0].cellPhone',
        :errorFn='emContact1CellPhoneError'
      )
      SFSelect#person-emContact1Relationship(
        label='Relationship',
        :options='relationshipOptions',
        v-model='person.emContacts[0].relationship',
        :errorFn='emContact1RelationshipError'
      )
      .form-item.person-heading Emergency Contact 2
      SFInput#person-emContact2Name(
        label='Name',
        trim,
        v-model='person.emContacts[1].name',
        :errorFn='emContact2NameError'
      )
      SFInput#person-emContact2HomePhone(
        label='Home Phone',
        trim,
        v-model='person.emContacts[1].homePhone',
        :errorFn='emContact2HomePhoneError'
      )
      SFInput#person-emContact2CellPhone(
        label='Cell Phone',
        trim,
        v-model='person.emContacts[1].cellPhone',
        :errorFn='emContact2CellPhoneError'
      )
      SFSelect#person-emContact2Relationship(
        label='Relationship',
        :options='relationshipOptions',
        v-model='person.emContacts[1].relationship',
        :errorFn='emContact2RelationshipError'
      )
      .form-item.person-heading Agreement
      #person-volreg-agreement.form-item.
        By submitting this application, I certify that all statements I have made on this application are true and correct and I
        hereby authorize the City of Sunnyvale to investigate the accuracy of this information.  I am aware that fingerprinting and
        a criminal records search is required for volunteers 18 years of age or older.  I understand that I am working at all times
        on a voluntary basis, without monetary compensation or benefits, and not as a paid employee.  I give the City of Sunnyvale
        permission to use any photographs or videos taken of me during my service without obligation or compensation to me.  I
        understand that the City of Sunnyvale reserves the right to terminate a volunteer's service at any time.  I understand that
        volunteers are covered under the City of Sunnyvale's Worker's Compensation Program for an injury or accident occurring while
        on duty.
      SFCheck#person-agreement(label='I agree', v-model='agreement', :errorFn='agreementError')
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFCheck, SFCheckGroup, SFInput, SFSelect, SSpinner } from '../../base'
import { EmContact, relationshipOptions } from './PersonEditContact.vue'
import { GetPersonAddress } from './PersonView.vue'
import PersonEditAddress from './PersonEditAddress.vue'

interface GetPersonVolgisticsRegister {
  id: number
  informalName: string
  formalName: string
  sortName: string
  callSign: string
  birthdate: string
  email: string
  email2: string
  cellPhone: string
  homePhone: string
  workPhone: string
  homeAddress: GetPersonAddress
  workAddress: GetPersonAddress
  mailAddress: GetPersonAddress
  emContacts: Array<EmContact>
}

const interestOptions = [
  { value: 'CERT-D', label: 'CERT Deployment Team' },
  { value: 'Outreach', label: 'Community Outreach' },
  { value: 'SARES', label: 'Amateur Radio (SARES)' },
  { value: 'SNAP', label: 'Neighborhood Preparedness Facilitator' },
  { value: 'Listos', label: 'Preparedness Instructor' },
  { value: 'CERT-T', label: 'CERT Instructor' },
]

export default defineComponent({
  components: { Modal, PersonEditAddress, SForm, SFCheck, SFCheckGroup, SFInput, SFSelect, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    const interests = ref(new Set())
    const agreement = ref(false)
    function show() {
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const person = ref({} as GetPersonVolgisticsRegister)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      person.value = (await axios.get<GetPersonVolgisticsRegister>(`/api/people/${props.pid}/volreg`)).data
      while (person.value.emContacts.length < 2) {
        person.value.emContacts.push({ name: '', homePhone: '', cellPhone: '', relationship: '' })
      }
      loading.value = false
    }

    // Field validation.
    const duplicateCallSign = ref('')
    const duplicateSortName = ref('')
    function informalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.informalName) return 'A name is required.'
      return ''
    }
    function formalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.formalName) return 'A name is required.'
      return ''
    }
    function sortNameError(lostFocus: boolean) {
      if (!person.value.sortName) return lostFocus ? 'A name is required.' : ''
      if (duplicateSortName.value === person.value?.sortName)
        return 'A different person has this name.'
      return ''
    }
    function callSignError(lostFocus: boolean) {
      if (!person.value.callSign) return ''
      if (lostFocus && !person.value.callSign.match(/^[AKNW][A-Z]?[0-9][A-Z]{1,3}$/i))
        return 'This is not a valid call sign.'
      if (duplicateCallSign.value === person.value?.callSign)
        return 'A different person has this call sign.'
      return ''
    }
    function birthdateError(lostFocus: boolean) {
      if (!person.value.birthdate) return lostFocus ? 'Your birthdate is required.' : ''
      if (lostFocus && !person.value.birthdate.match(/^(?:19|20)\d\d-\d\d-\d\d$/)) {
        return 'This is not a valid YYYY-MM-DD date.'
      }
      return ''
    }
    const duplicateEmail = ref('')
    function emailError(lostFocus: boolean) {
      if (!lostFocus || !person.value?.email) return ''
      if (
        !person.value.email.match(
          /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
        )
      )
        return 'This is not a valid email address.'
      if (duplicateEmail.value === person.value.email) return 'Another user has this email address.'
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
    const duplicateCellPhone = ref('')
    function cellPhoneError(lostFocus: boolean, submitted: boolean) {
      if (submitted && !person.value.cellPhone && !person.value.homePhone) {
        return 'A cell phone or home phone number is required.'
      }
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
    function emContact1NameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.emContacts[0].name) {
        return 'The emergency contact name is required.'
      }
      return ''
    }
    function emContact1HomePhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && !person.value.emContacts[0].homePhone && !person.value.emContacts[0].cellPhone) {
        return 'The emergency contact must have at least one phone number.'
      }
      if (!person.value?.emContacts || !person.value?.emContacts[0].homePhone) return ''
      if (person.value.emContacts[0].homePhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function emContact1CellPhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus || !person.value?.emContacts[0].cellPhone) return ''
      if (person.value.emContacts[0].cellPhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function emContact1RelationshipError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && !person.value.emContacts[0].relationship) {
        return 'The relationship is required.'
      }
      return ''
    }
    function emContact2NameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.emContacts[1].name) {
        return 'The emergency contact name is required.'
      }
      return ''
    }
    function emContact2HomePhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && !person.value.emContacts[1].homePhone && !person.value.emContacts[1].cellPhone) {
        return 'The emergency contact must have at least one phone number.'
      }
      if (!person.value?.emContacts || !person.value?.emContacts[1].homePhone) return ''
      if (person.value.emContacts[1].homePhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function emContact2CellPhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus || !person.value?.emContacts[1].cellPhone) return ''
      if (person.value.emContacts[1].cellPhone.replace(/[^0-9]/g, '').length !== 10)
        return 'A valid phone number must have 10 digits.'
      return ''
    }
    function emContact2RelationshipError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && !person.value.emContacts[1].relationship) {
        return 'The relationship is required.'
      }
      return ''
    }
    function agreementError(lostFocus: boolean, submitted: boolean) {
      if (submitted && !agreement.value) {
        return 'You must agree to this statement in order to register.'
      }
      return ''
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      body.append('informalName', person.value.informalName)
      body.append('formalName', person.value.formalName)
      body.append('sortName', person.value.sortName)
      body.append('callSign', person.value.callSign)
      body.append('birthdate', person.value.birthdate)
      body.append('email', person.value.email)
      body.append('email2', person.value.email2)
      body.append('cellPhone', person.value.cellPhone)
      body.append('homePhone', person.value.homePhone)
      body.append('workPhone', person.value.workPhone)
      if (person.value.homeAddress.address) {
        body.append('homeAddress', person.value.homeAddress.address)
        body.append('homeAddressLatitude', (person.value.homeAddress.latitude || 0).toString())
        body.append('homeAddressLongitude', (person.value.homeAddress.longitude || 0).toString())
      }
      if (person.value.workAddress.address) {
        body.append('workAddress', person.value.workAddress.address)
        body.append('workAddressLatitude', (person.value.workAddress.latitude || 0).toString())
        body.append('workAddressLongitude', (person.value.workAddress.longitude || 0).toString())
      } else {
        body.append('workAddressSameAsHome', person.value.workAddress.sameAsHome.toString())
      }
      if (person.value.mailAddress.address)
        body.append('mailAddress', person.value.mailAddress.address)
      else body.append('mailAddressSameAsHome', person.value.mailAddress.sameAsHome.toString())
      body.append('emContact1Name', person.value.emContacts[0].name)
      body.append('emContact1HomePhone', person.value.emContacts[0].homePhone)
      body.append('emContact1CellPhone', person.value.emContacts[0].cellPhone)
      body.append('emContact1Relationship', person.value.emContacts[0].relationship)
      body.append('emContact2Name', person.value.emContacts[1].name)
      body.append('emContact2HomePhone', person.value.emContacts[1].homePhone)
      body.append('emContact2CellPhone', person.value.emContacts[1].cellPhone)
      body.append('emContact2Relationship', person.value.emContacts[1].relationship)
      body.append('interests', Array.from(interests.value.keys()).join(','))
      body.append('agreement', agreement.value.toString())
      submitting.value = true
      try {
        await axios.post(`/api/people/${props.pid}/volreg`, body)
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        switch (err.response.data) {
          case 'sortName':
            duplicateSortName.value = person.value.sortName
            break
          case 'callSign':
            duplicateCallSign.value = person.value.callSign
            break
          case 'email':
            duplicateEmail.value = person.value.email
            break
          case 'cellPhone':
            duplicateCellPhone.value = person.value.cellPhone
            break
        }
      } finally {
        submitting.value = false
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      agreement,
      agreementError,
      birthdateError,
      callSignError,
      cellPhoneError,
      emailError,
      email2Error,
      emContact1NameError,
      emContact1CellPhoneError,
      emContact1HomePhoneError,
      emContact1RelationshipError,
      emContact2NameError,
      emContact2CellPhoneError,
      emContact2HomePhoneError,
      emContact2RelationshipError,
      formalNameError,
      homePhoneError,
      informalNameError,
      interestOptions,
      interests,
      loading,
      modal,
      onCancel,
      onSubmit,
      person,
      relationshipOptions,
      show,
      sortNameError,
      submitting,
      workPhoneError,
    }
  },
})
</script>

<style lang="postcss">
.person-heading {
  margin: 0 0.75rem;
  border-top: 1px solid #ccc;
  padding-top: 0.25rem;
}
#person-volreg-intro {
  margin: 0 0.75rem 1.5rem;
}
#person-volreg-agreement {
  margin: 1rem 0.75rem;
}
</style>
