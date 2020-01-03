<!--
CERTAttendance displays a CERT attendance report.
-->

<template lang="pug">
Page(title="CERT Attendance" subtitle="CERT Attendance" menuItem="reports")
  div.mt-3(v-if="loading")
    b-spinner(small)
  table#cert-att-table(v-else)
    thead
      tr(v-for="(row, rownum) in header" :key="rownum")
        td(v-for="(cell, colnum) in row" :key="colnum" :colspan="cell.span || 1" :class="headerClass(rownum, colnum)" v-text="cell.text")
    tbody
      tr(v-for="(row, rownum) in body" :key="rownum")
        td(v-for="(cell, colnum) in row" :key="colnum" :class="spanStarts[colnum] ? 'cert-att-span' : null" v-text="cell")
    tfoot
      tr(v-for="(row, rownum) in footer" :key="rownum")
        td(v-for="(cell, colnum) in row" :key="colnum" :class="spanStarts[colnum] ? 'cert-att-span' : null" v-text="cell")
</template>

<script>
export default {
  props: {
    team: String,
    dateFrom: String,
    dateTo: String,
    stats: String,
    detail: String,
  },
  data: () => ({ loading: false, header: null, body: null, footer: null }),
  computed: {
    spanStarts() {
      const spanStarts = {}
      let col = 0
      this.header[0].forEach(cell => {
        if (cell.span) spanStarts[col] = true
        col += cell.span || 1
      })
      return spanStarts
    }
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get('/api/reports/cert-attendance', {
      params: {
        team: this.team,
        dateFrom: this.dateFrom,
        dateTo: this.dateTo,
        stats: this.stats,
        detail: this.detail,
        format: 'JSON',
      }
    })).data
    console.log(data)
    this.header = data.header
    this.body = data.body
    this.footer = data.footer
    this.loading = false
  },
  methods: {
    headerClass(row, col) {
      if ((row === 0 && col !== 0) || (row !== 0 && this.spanStarts[col])) return 'cert-att-span'
      else return null
    },
  },
}
</script>

<style lang="stylus">
#cert-att-table
  margin-top 1.5rem
  border-right 2px solid black
  border-bottom 2px solid black
  thead
    td:first-child
      border-top none
      border-left none
    tr:first-child td:not(:first-child)
      border-top 2px solid black
    tr:nth-child(2) td
      padding 6px 2px
      vertical-align bottom
      text-align center
      writing-mode tb-rl
    tr:first-child td
      padding 2px 6px
  tbody
    border-top 2px solid black
    border-left 2px solid black
    tr:nth-child(even)
      background-color #eee
  td
    min-width 2rem
    border-top 1px solid #ccc
    border-left 1px solid #ccc
    text-align center
    &:first-child
      padding 2px 6px
      text-align left
    &.cert-att-span
      border-left 2px solid black
  tfoot
    border-top 2px solid black
    border-left 2px solid black
</style>
