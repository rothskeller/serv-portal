<!--
PersonView displays the information about a person, in non-editable form.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="!person")
  b-spinner(small)
#person-view(v-else)
  #person-view-name
    #person-view-informalName
      span(v-text="person.informalName")
      span#person-view-callSign(v-if="person.callSign" v-text="person.callSign")
    #person-view-formalName(v-if="person.formalName !== person.informalName" v-text="`(formally ${person.formalName})`")
  #person-view-roles
    div(v-for="role in person.roles" v-text="role.name")
  #person-view-emails(v-if="person.emails")
    div(v-for="e in person.emails")
      a(:href="`mailto:${e.email}`" v-text="e.email")
      span.person-view-email-label(v-if="e.label" v-text="` (${e.label})`")
  .person-view-phone(v-if="person.cellPhone")
      a(:href="`tel:${person.cellPhone}`" v-text="person.cellPhone")
      span.person-view-phone-label (Cell)
  .person-view-phone(v-if="person.homePhone")
      a(:href="`tel:${person.homePhone}`" v-text="person.homePhone")
      span.person-view-phone-label (Home)
  .person-view-phone(v-if="person.workPhone")
      a(:href="`tel:${person.workPhone}`" v-text="person.workPhone")
      span.person-view-phone-label (Work)
  .person-view-address(v-if="person.homeAddress.address")
    div
      span(v-if="person.workAddress.sameAsHome") Home Address (all day):
      span(v-else) Home Address:
      a.person-view-map(target="_blank" :href="`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.homeAddress.address)}`") Map
    div(v-text="person.homeAddress.address.split(',')[0]")
    div(v-text="person.homeAddress.address.replace(/^[^,]*, */, '')")
    div(v-if="person.homeAddress.fireDistrict" v-text="`Sunnyvale Fire District ${person.homeAddress.fireDistrict}`")
  .person-view-address(v-if="person.workAddress.address")
    div
      span Work Address:
      a.person-view-map(target="_blank" :href="`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(person.workAddress.address)}`") Map
    div(v-text="person.workAddress.address.split(',')[0]")
    div(v-text="person.workAddress.address.replace(/^[^,]*, */, '')")
    div(v-if="person.workAddress.fireDistrict" v-text="`Sunnyvale Fire District ${person.workAddress.fireDistrict}`")
  .person-view-address(v-if="person.mailAddress.address")
    div Mailing Address:
    div(v-text="person.mailAddress.address.split(',')[0]")
    div(v-text="person.mailAddress.address.replace(/^[^,]*, */, '')")
  #person-view-attended(v-if="person.attended")
    div Events attended:
    div(v-for="e in person.attended")
      span.person-view-attended-date(v-text="e.date")
      span(v-text="e.name")
</template>

<script>
export default {
  props: {
    onLoadPerson: Function,
  },
  data: () => ({ person: null }),
  async created() {
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}`)).data
    this.person = data.person
    this.onLoadPerson(this.person)
  },
}
</script>

<style lang="stylus">
#person-view
  margin 1.5rem 0.75rem
#person-view-name
  display flex
  flex-direction column
  @media (min-width: 576px)
    flex-direction row
#person-view-informalName
  font-weight bold
  font-size 1.25rem
  line-height 1.2
#person-view-callSign
  margin-left 0.5rem
  font-weight normal
#person-view-formalName
  color #888
  @media (min-width: 576px)
    margin-left 1rem
#person-view-roles
  line-height 1.2
#person-view-emails
  margin-top 0.75rem
.person-view-email-label
  color #888
.person-view-phone
  font-variant tabular-nums
.person-view-phone-label
  margin-left 0.25rem
  color #888
.person-view-address
  margin-top 0.75rem
  line-height 1.2
.person-view-map
  margin-left 1rem
  &::before
    content '['
  &::after
    content ']'
.person-view-address-flag
  color #888
  font-size 0.875rem
#person-view-attended
  margin-top 0.75rem
  line-height 1.2
.person-view-attended-date
  margin-right 0.75rem
  font-variant tabular-nums
</style>
