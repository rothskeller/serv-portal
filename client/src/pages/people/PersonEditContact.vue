<!--
PersonEditContact is the dialog box for editing a person's contact information.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit Contact Information',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
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
      template(v-if='person.canEditEmContacts')
        .form-item.person-emContact Emergency Contact 1
        SFInput#person-emContact1Name(label='Name', trim, v-model='person.emContacts[0].name')
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
        .form-item.person-emContact Emergency Contact 2
        SFInput#person-emContact2Name(label='Name', trim, v-model='person.emContacts[1].name')
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
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFInput, SFSelect, SSpinner } from '../../base'
import { GetPersonAddress } from './PersonView.vue'
import PersonEditAddress from './PersonEditAddress.vue'

interface GetPersonContact {
  id: number
  email: string
  email2: string
  cellPhone: string
  homePhone: string
  workPhone: string
  homeAddress: GetPersonAddress
  workAddress: GetPersonAddress
  mailAddress: GetPersonAddress
  canEditEmContacts: boolean
  emContacts: Array<EmContact>
}
interface EmContact {
  name: string
  homePhone: string
  cellPhone: string
  relationship: string
}

const relationshipOptions = [
  'Co-worker', 'Daughter', 'Father', 'Friend', 'Mother', 'Neighbor', 'Other',
  'Relative', 'Son', 'Spouse', 'Supervisor'
]

export default defineComponent({
  components: { Modal, PersonEditAddress, SForm, SFInput, SFSelect, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    function show() {
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const person = ref({} as GetPersonContact)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      person.value = (await axios.get<GetPersonContact>(`/api/people/${props.pid}/contact`)).data
      if (person.value.canEditEmContacts) {
        while (person.value.emContacts.length < 2) {
          person.value.emContacts.push({ name: '', homePhone: '', cellPhone: '', relationship: '' })
        }
      }
      loading.value = false
    }

    // Field validations.
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
    function emContact1HomePhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && person.value.emContacts[0].name && !person.value.emContacts[0].homePhone && !person.value.emContacts[0].cellPhone) {
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
      if (submitted && person.value.emContacts[0].name && !person.value.emContacts[0].relationship) {
        return 'The relationship is required.'
      }
      return ''
    }
    function emContact2HomePhoneError(lostFocus: boolean, submitted: boolean) {
      if (!lostFocus) return ''
      if (submitted && person.value.emContacts[1].name && !person.value.emContacts[1].homePhone && !person.value.emContacts[1].cellPhone) {
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
      if (submitted && person.value.emContacts[1].name && !person.value.emContacts[1].relationship) {
        return 'The relationship is required.'
      }
      return ''
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
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
      if (person.value.canEditEmContacts) {
        body.append('emContact1Name', person.value.emContacts[0].name)
        body.append('emContact1HomePhone', person.value.emContacts[0].homePhone)
        body.append('emContact1CellPhone', person.value.emContacts[0].cellPhone)
        body.append('emContact1Relationship', person.value.emContacts[0].relationship)
        body.append('emContact2Name', person.value.emContacts[1].name)
        body.append('emContact2HomePhone', person.value.emContacts[1].homePhone)
        body.append('emContact2CellPhone', person.value.emContacts[1].cellPhone)
        body.append('emContact2Relationship', person.value.emContacts[1].relationship)
      }
      submitting.value = true
      try {
        await axios.post(`/api/people/${props.pid}/contact`, body)
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        switch (err.response.data) {
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
      loading,
      modal,
      onCancel,
      onSubmit,
      person,
      show,
      submitting,
      emailError,
      email2Error,
      cellPhoneError,
      homePhoneError,
      workPhoneError,
      emContact1HomePhoneError,
      emContact1CellPhoneError,
      emContact1RelationshipError,
      emContact2HomePhoneError,
      emContact2CellPhoneError,
      emContact2RelationshipError,
      relationshipOptions,
    }
  },
})
</script>

<style lang="postcss">
.person-emContact {
  margin: 0 0.75rem;
  border-top: 1px solid #ccc;
  padding-top: 0.25rem;
}
</style>
