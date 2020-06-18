<!--
VolunteerHours allows editing the hours a person has spent on SERV activities.
It does not require a login but does require a token as part of the URL.  The
tokens are different for each person and change each month.  This is nearly
identical to person/PersonHours, except for PublicPage framing and handling the
saved path differently.
-->

<template lang="pug">
PublicPage(:title="title")
  div.mt-3.ml-2(v-if="saved")
    div Your volunteer hours have been saved.  Thank you for volunteering!
    div.mt-3
      b-btn(to="/" variant="primary") Go to Login Page
  div.mt-3.ml-2(v-else-if="notFound")
    div The link you used is not valid or has expired.
    div.mt-3
      b-btn(to="/" variant="primary") Go to Login Page
  div.mt-3.ml-2(v-else-if="unregistered" style="max-width:600px")
    | You are not currently registered as a City of Sunnyvale volunteer.  We
    | appreciate your volunteer efforts, but we cannot record your hours until you
    | are registered.  To do so, please fill out
    |
    a(href="https://www.volgistics.com/ex/portal.dll/ap?AP=929478828" target="_blank") this form
    | .  In the “City employee status or referral” box, please enter
    pre.ml-4.mt-3
      | Rebecca Elizondo
      | Department of Public Safety
    | and the names of the organizations you're volunteering for (CERT, LISTOS,
    | SNAP, and/or SARES).  Come back a week or so later and we should have your
    | registration on file.  If you have any difficulties with this, please
    | contact Rebecca at RElizondo@sunnyvale.ca.gov.
  div.mt-3.ml-2(v-else-if="!months")
    b-spinner(small)
  form#person-hours(v-else @submit.prevent="onSubmit")
    .person-hours(v-for="month in months")
      .person-hours-heading(v-text="`Volunteer Hours for ${month.month}`")
      table.person-hours-table
        tr(v-for="event in month.events")
          td.person-hours-event(v-text="eventText(event)")
          td
            input.person-hours-time(type="number" min="0" step="0.5" :value="eventTime(event)" @change="setEventTime(event, $event)")
        tr
          td.person-hours-total-label TOTAL
          td
            div.person-hours-total-time(v-text="totalHours(month)")
    div.mt-3
      b-btn(type="submit" variant="primary") Save Hours
</template>

<script>
import PublicPage from '@/base/PublicPage'

export default {
  components: { PublicPage },
  data: () => ({ months: null, notFound: false, saved: false, title: 'Volunteer Hours', unregistered: false }),
  async created() {
    const response = await this.$axios.get(`/api/people/${this.$route.params.id}/hours`).catch(e => {
      if (e.response && e.response.status === 404)
        return null
      else
        throw (e)
    })
    if (!response)
      this.notFound = true
    else if (!response.data)
      this.unregistered = true
    else {
      this.months = response.data.months
      this.title = response.data.name
    }
  },
  methods: {
    eventText(e) {
      if (e.placeholder) return e.name
      return `${e.date} ${e.name}`
    },
    eventTime(e) {
      if (e.minutes === 0) return ''
      return Math.floor(e.minutes / 30) / 2
    },
    async onSubmit() {
      const body = new FormData
      this.months.forEach(m => {
        m.events.forEach(e => {
          body.append(`e${e.id}`, e.minutes)
        })
      })
      await this.$axios.post(`/api/people/${this.$route.params.id}/hours`, body)
      this.saved = true
    },
    setEventTime(e, evt) {
      e.minutes = evt.target.value * 60
    },
    totalHours(month) {
      const minutes = month.events.reduce((sum, e) => sum + e.minutes, 0)
      return Math.floor(minutes / 30) / 2
    },
  },
}
</script>

<style lang="stylus">
#person-hours
  margin 1.5rem 0.75rem
.person-hours-heading
  font-weight bold
  font-size 1.5rem
.person-hours-event
  padding-right 1rem
  padding-left 0.5rem
  text-indent -0.5rem
.person-hours-time
  width 4rem
  text-align right
.person-hours-total-label
  padding-right 1rem
  text-align right
  font-weight bold
.person-hours-total-time
  width 3rem
  text-align right
  font-weight bold
</style>
