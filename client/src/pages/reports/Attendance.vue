<!--
Attendance displays the attendance report.
-->

<template lang="pug">
#attrep(v-if='params.cells')
  table#attrep-params
    tr
      td Date range
      td
        SSelect(
          :options='options.dateRanges',
          valueKey='tag',
          labelKey='label',
          v-model='params.dateRange'
        )
        |
        | {{ params.dateFrom }} to {{ params.dateTo }}
    tr
      td Rows
      td
        SRadioGroup#attrep-rows(inline, :options='rowOptions', v-model='params.rows')
        SCheck#attrep-includeZerosY(label='Include Zeros', inline, v-model='params.includeZerosY')
    tr
      td Columns
      td
        SRadioGroup#attrep-columns(inline, :options='columnOptions', v-model='params.columns')
        SCheck#attrep-includeZerosX(label='Include Zeros', inline, v-model='params.includeZerosX')
    tr
      td Cells
      td
        SRadioGroup#attrep-cells(inline, :options='cellOptions', v-model='params.cells')
    tr
      td Organizations
      td
        SCheckGroup#attrep-orgs(inline, :options='options.orgs', valueKey='id', v-model='orgs')
    tr
      td Events
      td
        SCheckGroup#attrep-eventTypes(
          inline,
          :options='options.eventTypes',
          valueKey='id',
          v-model='eventTypes'
        )
    tr
      td Attendees
      td
        SCheckGroup#attrep-attendanceTypes(
          inline,
          :options='options.attendanceTypes',
          valueKey='id',
          v-model='attendanceTypes'
        )
  table#attrep-table
    tbody
      tr(v-for='(r, ri) in rows', :class='`attrep-row-${r}`')
        template(v-if='ri === 0')
          td(v-for='(c, ci) in cells[ri]', :class='`attrep-col-${columns[ci]}`')
            .attrep-vertical(v-if='c', v-text='c')
        template(v-else)
          td(v-for='(c, ci) in cells[ri]', :class='`attrep-col-${columns[ci]}`', v-text='c')
  #attrep-count(v-if='count', v-text='count > 1 ? `${count} people listed` : "1 person listed"')
</template>

<script lang="ts">
import { defineComponent, ref, watch, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SCheck, SCheckGroup, SRadioGroup, SSelect } from '../../base'

interface GetReportsAttendanceDateRange {
  tag: string
  label: string
  dateFrom: string
  dateTo: string
}
interface GetReportsAttendanceIDLabel {
  id: number
  label: string
}
interface GetReportsAttendanceOrg {
  id: number
  label: string
  fmtOrg: string // added locally
}
interface AttendanceReportParams {
  dateRange: string
  dateFrom: string
  dateTo: string
  rows: string
  columns: string
  cells: string
  includeZerosX: boolean
  includeZerosY: boolean
  orgs: Array<number>
  eventTypes: Array<number>
  attendanceTypes: Array<number>
}
interface AttendanceReportOptions {
  dateRanges: Array<GetReportsAttendanceDateRange>
  orgs: Array<GetReportsAttendanceOrg>
  eventTypes: Array<GetReportsAttendanceIDLabel>
  attendanceTypes: Array<GetReportsAttendanceIDLabel>
}
interface AttendanceReport {
  parameters: AttendanceReportParams
  options: AttendanceReportOptions
  rows: Array<string>
  columns: Array<string>
  cells: Array<Array<string>>
}

const rowOptions = [
  { value: 'p', label: 'Person' },
  { value: 'o', label: 'Org' },
  { value: 'po', label: 'Person, Org' },
  { value: 'op', label: 'Org, Person' },
]
const columnOptions = [
  { value: 'e', label: 'Events' },
  { value: 'm', label: 'Months' },
]
const cellOptions = [
  { value: 'h', label: 'Cumulative Hours' },
  { value: 'c', label: 'Attendance Counts' },
]

export default defineComponent({
  components: { SCheck, SCheckGroup, SRadioGroup, SSelect },
  setup() {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Attendance Report', browserTitle: 'Attendance' })

    // Whenever route query parameters change, request new report.
    const params = ref({} as AttendanceReportParams)
    const options = ref({} as AttendanceReportOptions)
    const columns = ref([] as Array<string>)
    const rows = ref([] as Array<string>)
    const cells = ref([] as Array<Array<string>>)
    const orgs = ref(new Set<number>())
    const eventTypes = ref(new Set<number>())
    const attendanceTypes = ref(new Set<number>())
    const count = ref(0)
    watchEffect(async () => {
      const report = (
        await axios.get<AttendanceReport>('/api/reports/attendance', { params: route.query })
      ).data
      params.value = report.parameters
      options.value = report.options
      columns.value = report.columns
      rows.value = report.rows
      cells.value = report.cells
      orgs.value = new Set(report.parameters.orgs)
      eventTypes.value = new Set(report.parameters.eventTypes)
      attendanceTypes.value = new Set(report.parameters.attendanceTypes)
      switch (report.parameters.rows) {
        case 'p':
          count.value = report.cells.length
          break
        case 'po':
          count.value = report.cells.filter(row => row[0]).length
          break
        default: count.value = 0
      }
    })

    // Watch for changes to the parameters.
    watch(
      [params, attendanceTypes, eventTypes, orgs],
      () => {
        const query = {} as any
        query.dateRange = params.value.dateRange
        query.rows = params.value.rows
        query.columns = params.value.columns
        query.cells = params.value.cells
        query.orgs = Array.from(orgs.value.keys(), (v) => v.toString())
          .sort()
          .join(',')
        query.eventTypes = Array.from(eventTypes.value.keys(), (v) => v.toString())
          .sort()
          .join(',')
        query.attendanceTypes = Array.from(attendanceTypes.value.keys(), (v) => v.toString())
          .sort()
          .join(',')
        if (params.value.includeZerosX) query.includeZerosX = 'true'
        if (params.value.includeZerosY) query.includeZerosY = 'true'
        router.replace({ path: '/reports/attendance', query })
      },
      { deep: true }
    )

    return {
      attendanceTypes,
      cells,
      cellOptions,
      columns,
      columnOptions,
      count,
      eventTypes,
      options,
      orgs,
      params,
      rows,
      rowOptions,
    }
  },
})
</script>

<style lang="postcss">
#attrep-params {
  margin-bottom: 1.5rem;
  & td:first-child {
    padding-right: 1rem;
  }
  @media print {
    display: none;
  }
}
#attrep {
  margin: 1.5rem 0.75rem;
  overflow-x: auto;
}
.attrep-col-h {
  border-left: 2px solid #888;
}
.attrep-col-h2 {
  border-left: 1px solid #ccc;
}
.attrep-col-s {
  border-left: 2px solid #888;
}
.attrep-col-c {
  border-left: 1px solid #ccc;
}
.attrep-col-t {
  border-left: 1px solid #ccc;
  border-right: 2px solid #888;
}
.attrep-col-1 {
  border-left: 2px solid #888;
  border-right: 2px solid #888;
}
.attrep-row-h,
.attrep-row-s,
.attrep-row-t {
  border-top: 2px solid #888;
}
.attrep-row-h2 {
  border-top: 1px solid #ccc;
}
.attrep-row-t {
  border-bottom: 2px solid #888;
}
.attrep-row-c2,
.attrep-row-tc2 {
  background-color: #eee;
}
.attrep-row-h > .attrep-col-h,
.attrep-row-h > .attrep-col-h2,
.attrep-row-h2 > .attrep-col-h,
.attrep-row-h2 > .attrep-col-h2 {
  border-left: hidden;
  border-top: hidden;
}
.attrep-row-h > .attrep-col-s,
.attrep-row-h > .attrep-col-1,
.attrep-row-h > .attrep-col-c,
.attrep-row-h > .attrep-col-t {
  text-align: right;
  vertical-align: bottom;
  padding-right: 0.5rem;
}
.attrep-vertical {
  writing-mode: vertical-rl;
  font-variant-numeric: tabular-nums;
  line-height: 1;
  width: 100%;
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}
.attrep-col-s,
.attrep-col-1,
.attrep-col-c,
.attrep-col-t {
  padding-right: 0.5rem;
  min-width: 4rem;
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.attrep-col-h,
.attrep-col-h2 {
  padding-left: 0.5rem;
  padding-right: 0.5rem;
  white-space: nowrap;
}
.attrep-col-t,
.attrep-row-t {
  font-weight: bold;
}
</style>
