<!--
CERTAttendanceForm displays the form for generating a CERT Attendance report.
-->

<template lang="pug">
form#cert-att-form(@submit.prevent="onSubmit")
  .report-title CERT Attendance Report
  b-form-group(label="Report on team" label-cols-sm="auto" label-class="cert-att-label")
    b-form-radio-group.cert-att-radio-group(:options="teamList" v-model="team")
  b-form-group(label="Date range" label-cols-sm="auto" label-class="cert-att-label" :state="dateError ? false : null" :invalid-feedback="dateError")
    b-form-input#cert-att-date-from(type="date" v-model="dateFromI" :state="dateError ? false : null")
    |
    | through
    |
    b-form-input#cert-att-date-to(type="date" v-model="dateToI" :state="dateError ? false : null")
  b-form-group(label="Statistics by" label-cols-sm="auto" label-class="cert-att-label")
    b-form-radio-group.cert-att-radio-group(:options="statsList" v-model="stats")
  b-form-group(label="Show detail" label-cols-sm="auto" label-class="cert-att-label")
    b-form-radio-group.cert-att-radio-group(:options="detailList" v-model="detail")
  div.mt-3
    b-btn(type="submit" variant="primary" :disabled="!!dateError") Generate Report
</template>

<script>
const teamList = ['Alpha', 'Bravo', 'Both']
const statsList = [{ value: 'count', text: 'Number of Events' }, { value: 'hours', text: 'Cumulative Hours' }]
const detailList = [
  { value: 'date', text: 'Show by date' },
  { value: 'month', text: 'Show by month' },
  { value: 'total', text: 'Show totals only' },
]
const dateRE = /^20\d\d-(?:0[1-9]|1[012])-(?:0[1-9]|[12][0-9]|3[01])$/

export default {
  props: {
    dateFrom: String,
    dateTo: String,
  },
  data: () => ({
    teamList, statsList, detailList,
    team: 'Both',
    dateFromI: null, dateToI: null,
    stats: 'count',
    detail: 'month',
    dateError: null,
  }),
  created() {
    this.dateFromI = this.dateFrom
    this.dateToI = this.dateTo
  },
  watch: {
    dateFromI: 'checkDates',
    dateToI: 'checkDates',
  },
  methods: {
    checkDates() {
      if (!this.dateFromI || !this.dateToI)
        this.dateError = 'Starting and ending dates are required.'
      else if (!this.dateFromI.match(dateRE) || !this.dateToI.match(dateRE))
        this.dateError = 'Valid dates have the form YYYY-MM-DD.'
      else if (this.dateFromI > this.dateToI)
        this.dateError = 'The starting date must be before the ending date.'
      else
        this.dateError = null
    },
    onSubmit() {
      if (this.dateError) return
      this.$router.push(`/reports/cert-attendance?team=${this.team}&dateFrom=${this.dateFrom}&dateTo=${this.dateTo}&stats=${this.stats}&detail=${this.detail}`)
    },
  },
}
</script>

<style lang="stylus">
.cert-att-label
  @media (min-width: 576px)
    width 8rem
.cert-att-radio-group
  @media (min-width: 576px)
    padding-top 0.375rem
#cert-att-date-from, #cert-att-date-to
  display inline
  max-width 12rem
</style>
