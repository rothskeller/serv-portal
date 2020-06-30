<!--
PersonHours allows editing the hours a person has spent on SERV activities.
-->

<template lang="pug">
div.mt-3.ml-2(v-if="unregistered" style="max-width:600px")
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
        td.person-hours-event
          | {{ eventText(event) }}
          div.person-hours-event-includes(v-if="eventIncludes(event)" v-text="eventIncludes(event)")
        td
          input.person-hours-time(type="number" min="0" step="0.5" :value="eventTime(event)" @change="setEventTime(event, $event)")
      tr
        td.person-hours-total-label TOTAL
        td
          div.person-hours-total-time(v-text="totalHours(month)")
  div.mt-3
    b-btn(type="submit" variant="primary") Save Hours
    b-btn.ml-2(@click="onCancel") Cancel
  table.person-hours-guide
    tr
      td Volunteer Hours
      td Not Volunteer Hours
    tr
      td In general, time you spend helping or preparing to help the community as part of SERV.  For example:
      td In general, time you spend preparing yourself or your household; or time you spend becoming a SERV volunteer.  For example:
    tr
      td Organizing or teaching CERT Basic, LISTOS, PEP, or SNAP events
      td Attending CERT Basic, LISTOS, PEP, or ham cram classes
    tr
      td Preparing and maintaining a CERT or SARES “go kit” for deployment
      td Preparing and maintaining a personal or household evacuation kit
    tr
      td SERV team meetings, radio nets, and drills; CERT continuing education seminars; SARES or county ARES training classes
      td SERV team social gatherings
    tr
      td Responding in an emergency when activated by the city
      td Responding in an emergency when not activated by the city
    tr
      td Travel to and from the above
      td
    tr
      td SERV administration activities
      td
</template>

<script>
export default {
  props: {
    onLoadPerson: Function, // not used
  },
  data: () => ({ months: null, unregistered: false }),
  async created() {
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}/hours`)).data
    if (data)
      this.months = data.months
    else
      this.unregistered = true
  },
  methods: {
    eventIncludes(e) {
      if (!e.placeholder) return null
      switch (e.name) {
        case 'Other CERT Hours': return 'Includes contact tracing'
        case 'Other LISTOS Hours': return 'Includes PEP and Outreach'
        default: return null
      }
    },
    eventText(e) {
      if (e.placeholder) return e.name
      return `${e.date} ${e.name}`
    },
    eventTime(e) {
      if (e.minutes === 0) return ''
      return Math.floor(e.minutes / 30) / 2
    },
    onCancel() { this.$router.go(-1) },
    async onSubmit() {
      const body = new FormData
      this.months.forEach(m => {
        m.events.forEach(e => {
          body.append(`e${e.id}`, e.minutes)
        })
      })
      await this.$axios.post(`/api/people/${this.$route.params.id}/hours`, body)
      this.$router.push(`/people/${this.$route.params.id}`)
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
  line-height 1.2
.person-hours-event-includes
  font-style italic
  font-size 0.75rem
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
.person-hours-guide
  margin-top 1.5rem
  max-width 800px
  tr:nth-child(1)
    background-color #5b9bd5
    td
      color white
      text-align center
      font-weight bold
  tr:nth-child(2)
    font-weight bold
  tr:nth-child(even)
    background-color #deeaf6
  td
    padding 0.25rem
    width 50%
    border 1px solid #eee
    vertical-align top
    line-height 1.2
</style>
