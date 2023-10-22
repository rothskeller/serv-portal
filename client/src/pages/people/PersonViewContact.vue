<!--
PersonViewContact displays a person's contact information.
-->

<template lang="pug">
PersonViewSection(
  v-if='hasContactInfo || person.canEdit',
  title='Contact Information',
  :editable='person.canEdit',
  @edit='onEditContact'
)
  #person-view-emails(v-if='person.contact.email || person.contact.email2')
    div(v-if='person.contact.email')
      a(:href='`mailto:${person.contact.email}`', v-text='person.contact.email')
    div(v-if='person.contact.email2')
      a(:href='`mailto:${person.contact.email2}`', v-text='person.contact.email2')
  .person-view-phone(v-if='person.contact.cellPhone')
    a(:href='`tel:${person.contact.cellPhone}`', v-text='person.contact.cellPhone')
    span.person-view-phone-label (Cell)
  .person-view-phone(v-if='person.contact.homePhone')
    a(:href='`tel:${person.contact.homePhone}`', v-text='person.contact.homePhone')
    span.person-view-phone-label (Home)
  .person-view-phone(v-if='person.contact.workPhone')
    a(:href='`tel:${person.contact.workPhone}`', v-text='person.contact.workPhone')
    span.person-view-phone-label (Work)
  .person-view-address(v-if='person.contact.homeAddress && person.contact.homeAddress.address')
    div
      span(v-if='person.contact.workAddress && person.contact.workAddress.sameAsHome') Home Address (all day):
      span(v-else) Home Address:
      a.person-view-map(
        target='_blank',
        :href='`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.contact.homeAddress.address)}`'
      ) Map
    div(v-text='person.contact.homeAddress.address.split(",")[0]')
    div(v-text='person.contact.homeAddress.address.replace(/^[^,]*, */, "")')
    div(
      v-if='person.contact.homeAddress.fireDistrict',
      v-text='`Sunnyvale Fire District ${person.contact.homeAddress.fireDistrict}`'
    )
  .person-view-address(v-if='person.contact.workAddress && person.contact.workAddress.address')
    div
      span Work Address:
      a.person-view-map(
        target='_blank',
        :href='`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.contact.workAddress.address)}`'
      ) Map
    div(v-text='person.contact.workAddress.address.split(",")[0]')
    div(v-text='person.contact.workAddress.address.replace(/^[^,]*, */, "")')
    div(
      v-if='person.contact.workAddress.fireDistrict',
      v-text='`Sunnyvale Fire District ${person.contact.workAddress.fireDistrict}`'
    )
  .person-view-address(v-if='person.contact.mailAddress && person.contact.mailAddress.address')
    div Mailing Address:
    div(v-text='person.contact.mailAddress.address.split(",")[0]')
    div(v-text='person.contact.mailAddress.address.replace(/^[^,]*, */, "")')
  #person-view-emContacts(
    v-if='person.emContacts',
    v-text='`${person.emContacts} emergency contact${person.emContacts > 1 ? "s" : ""} on file`'
  )
  PersonEditContact(v-if='person.canEdit', ref='editContactModal', :pid='person.id')
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import { GetPerson } from './PersonView.vue'
import PersonEditContact from './PersonEditContact.vue'
import PersonViewSection from './PersonViewSection.vue'

export default defineComponent({
  components: { PersonEditContact, PersonViewSection },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const hasContactInfo = computed(() => {
      if (!props.person.contact) return false
      return (
        props.person.contact.email ||
        props.person.contact.email2 ||
        props.person.contact.cellPhone ||
        props.person.contact.homePhone ||
        props.person.contact.workPhone ||
        props.person.contact.homeAddress ||
        props.person.contact.workAddress ||
        props.person.contact.mailAddress
      )
    })
    const editContactModal = ref(null as any)
    async function onEditContact() {
      if (!(await editContactModal.value.show())) return
      emit('reload')
    }
    return {
      editContactModal,
      hasContactInfo,
      onEditContact,
    }
  },
})
</script>

<style lang="postcss">
#person-view-emails,
#person-view-emContacts {
  margin-top: 0.75rem;
}
.person-view-email-label {
  color: #888;
}
.person-view-phone {
  font-variant: tabular-nums;
}
.person-view-phone-label {
  margin-left: 0.25rem;
  color: #888;
}
.person-view-address,
.person-view-clearances {
  margin-top: 0.75rem;
  line-height: 1.2;
}
.person-view-map {
  margin-left: 1rem;
  &::before {
    content: '[';
  }
  &::after {
    content: ']';
  }
}
.person-view-address-flag {
  color: #888;
  font-size: 0.875rem;
}
</style>
