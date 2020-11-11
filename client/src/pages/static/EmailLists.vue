<!--
EmailLists displays a page with information about our email lists.
-->

<template lang="pug">
#help-email-lists
  h2 SERV Email Lists
  p.
    The SunnyvaleSERV.org site offers a number of email distribution lists.  We
    have one for each volunteer program, that we give out to the general public
    who might want more information about the program.  Email sent to these
    lists is delivered to designated public contact people for each program:
  pre.
    cert@sunnyvaleserv.org
    listos@sunnyvaleserv.org
    sares@sunnyvaleserv.org
    snap@sunnyvaleserv.org
  p.
    There are also lists for the volunteers on each of our teams:
  pre.
    cert-alpha@sunnyvaleserv.org
    cert-bravo@sunnyvaleserv.org
    cert-committee@sunnyvaleserv.org
    listos-team@sunnyvaleserv.org
    outreach-team@sunnyvaleserv.org
    sares-active@sunnyvaleserv.org
    sares-leads@sunnyvaleserv.org
    snap-team@sunnyvaleserv.org
  p.
    and for the students in each CERT class:
  pre.
    cert-54@sunnyvaleserv.org
    cert-55@sunnyvaleserv.org
    cert-56@sunnyvaleserv.org
  p.
    Finally, there are some broader lists for special purposes:
  pre.
    serv-all@sunnyvaleserv.org
    volunteer-hours@sunnyvaleserv.org
  p.
    All of these email lists have restricted access.  For the team lists, only
    members of the team can send mail to them; for the class lists, only the
    instructors can send mail to them; and for the broader lists, only DPS staff
    can send mail to them.  Any mail sent to any of our lists from someone else
    is held for approval before being routed to the list.  Messages on topics
    unrelated to SERV will generally be rejected.
  p.
    If you are receiving email from one of these lists that you do not want,
    there is an “unsubscribe” link at the bottom of every email.  If you are
    receiving email at the wrong address, you can change your email address in
    the “Profile” section of this web site.
  div(style='margin: 1.5rem 0')
    SButton(@click='$router.go(-1)', variant='primary') Back
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import setPage from '../../plugins/page'
import SButton from '../../base/SButton.vue'

export default defineComponent({
  components: { SButton },
  setup() {
    setPage({ title: 'SERV Email Lists' })
  },
})
</script>

<style lang="postcss">
#help-email-lists {
  margin: 1.5rem 0.75rem;
  max-width: 40rem;
}
</style>
