<!--
PersonView displays the information about a person, in non-editable form.
-->

<template lang="pug">
#person-view
  #person-view-fullName
    span(v-text="person.fullName")
    span#person-view-callSign(v-if="person.callSign" v-text="person.callSign")
  div(v-for="role in person.roles" v-text="role.name")
  #person-view-emails(v-if="person.emails")
    div(v-for="e in person.emails")
      span.person-view-bad-email(v-if="e.bad" v-text="e.label ? `${e.email} (${e.label})` : e.email")
      template(v-else)
        a(:href="`mailto:${e.email}`" v-text="e.email")
        span.person-view-email-label(v-if="e.label" v-text="` (${e.label})`")
  #person-view-phones(v-if="person.phones")
    div(v-for="p in person.phones")
      a(:href="`tel:${p.phone}`" v-text="p.phone")
      span.person-view-phone-label(v-if="p.label && p.sms" v-text="` (${p.label}, SMS)`")
      span.person-view-phone-label(v-else-if="p.sms" v-text="` (SMS)`")
      span.person-view-phone-label(v-else-if="p.label" v-text="` (${p.label}`")
  template(v-if="person.addresses")
    .person-view-address(v-for="a in person.addresses")
      div(v-if="a.label" v-text="`${a.label}:`")
      div(v-text="a.address")
      div(v-text="`${a.city}, ${a.state} ${a.zip}`")
      .person-view-address-flag(v-if="a.workHours") (during work hours)
      .person-view-address-flag(v-if="a.homeHours") (during non-work hours)
      .person-view-address-flag(v-if="a.mailingAddress") (for postal mail)
  #person-view-attended(v-if="person.attended")
    div Events attended:
    div(v-for="e in person.attended")
      span.person-view-attended-date(v-text="e.date")
      span(v-text="e.name")
</template>

<script>
export default {
  props: {
    person: null,
  },
  computed: {
  },
  methods: {
  },
}
</script>

<style lang="stylus">
#person-view
  margin 1.5rem 0.75rem
#person-view-fullName
  font-weight bold
  font-size 1.25rem
  line-height 1.2
#person-view-callSign
  margin-left 0.5rem
  font-weight normal
#person-view-emails
  margin-top 0.75rem
.person-view-email-label
  color #888
.person-view-bad-email
  text-decoration line-through
#person-view-phones
  margin-top 0.75rem
.person-view-phone-label
  color #888
.person-view-address
  margin-top 0.75rem
.person-view-address-flag
  color #888
  font-size 0.875rem
#person-view-attended
  margin-top 0.75rem
.person-view-attended-date
  margin-right 0.75rem
  font-variant tabular-nums
</style>
