// Clearance displays the clearance report.

import { defineComponent, Fragment, h, ref, watch, watchEffect } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { SButton, SSelect, SSpinner } from '../../base'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import type { GetClearance, GetClearanceParameters, GetClearanceRow } from './api'
import './clearance.css'

const Clearance = defineComponent({
  name: 'Clearance',
  props: {},
  emits: [],
  setup() {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Clearance Report', browserTitle: 'Clearance' })

    // Whenever route query parameters change, request new report.
    const report = ref({} as GetClearance)
    watchEffect(async () => {
      report.value = (
        await axios.get<GetClearance>('/api/reports/clearance', { params: route.query })
      ).data
      report.value.justLoaded = true
    })

    // Watch for changes to the parameters.
    watch(
      report,
      () => {
        if (!report.value) return
        if (report.value.justLoaded) {
          report.value.justLoaded = false
          return
        }
        const query = {} as any
        query.role = report.value.parameters.role
        query.with = report.value.parameters.with
        query.without = report.value.parameters.without
        router.replace({ path: '/reports/clearance', query })
      },
      { deep: true }
    )

    function render() {
      if (!report.value.parameters) return h('div', { id: 'clearrep' }, h(SSpinner))
      const hasBGCDetail =
        report.value.rows.length != 0 && report.value.rows[0].bgCheck === undefined
      const allOrgs = new Set<string>()
      report.value.rows.forEach((r) => {
        Object.keys(r.orgs).forEach((o) => {
          allOrgs.add(o)
        })
      })
      return h('div', { id: 'clearrep' }, [
        renderParameters(report.value.parameters),
        renderTable(report.value.rows, Array.from(allOrgs.keys()).sort(), hasBGCDetail),
        renderCSVButton(report.value),
      ])
    }

    return render
  },
})
export default Clearance

function renderParameters(params: GetClearanceParameters) {
  return h('div', { id: 'clearrep-params' }, [
    h('div', 'Show'),
    h(SSelect, {
      modelValue: params.role,
      options: [{ id: 0, name: 'Everyone' }, ...params.allowedRoles],
      valueKey: 'id',
      labelKey: 'name',
      'onUpdate:modelValue': (v: number) => (params.role = v),
    }),
    h('div', 'with'),
    h(SSelect, {
      modelValue: params.with,
      options: [{ value: '', label: '—' }, ...params.allowedRestrictions],
      'onUpdate:modelValue': (v: string) => (params.with = v),
    }),
    h('div', 'and without'),
    h(SSelect, {
      modelValue: params.without,
      options: [{ value: '', label: '—' }, ...params.allowedRestrictions],
      'onUpdate:modelValue': (v: string) => (params.without = v),
    }),
  ])
}

function renderTable(rows: Array<GetClearanceRow>, allOrgs: Array<string>, hasBGCDetail: boolean) {
  if (!rows.length)
    return h(
      'div',
      { style: 'margin-top:1.5rem;font-weight:bold' },
      'No one matches these report criteria.'
    )
  return h(
    'div',
    {
      id: 'clearrep-tbody',
    },
    [renderTableHeading(hasBGCDetail), ...rows.map((r) => renderTableRow(r, allOrgs, hasBGCDetail))]
  )
}

function renderTableHeading(hasBGCDetail: boolean) {
  return h(Fragment, [
    h('div', { class: 'clearrep-thead' }, 'Orgs'),
    h('div', { class: 'clearrep-thead' }, 'Name'),
    h('div', { class: 'clearrep-thead' }, 'V'),
    h('div', { class: 'clearrep-thead' }, 'DSW'),
    h('div', { class: 'clearrep-thead' }, hasBGCDetail ? 'BG' : 'B'),
    h('div', { class: 'clearrep-thead' }, 'Identification'),
  ])
}

function renderTableRow(r: GetClearanceRow, allOrgs: Array<string>, hasBGCDetail: boolean) {
  return h(Fragment, [
    renderOrgBadgeCells(r, allOrgs),
    renderNameLinkCell(r),
    renderVolgisticsCell(r),
    renderDSWCells(r),
    hasBGCDetail ? renderBGCheckCells(r) : renderBGCheckCell(r),
    renderIdentCells(r),
  ])
}

const orgBadgeLabels: Record<string, string> = {
  admin: 'A',
  'cert-d': 'D',
  'cert-t': 'T',
  listos: 'L',
  sares: 'S',
  snap: 'S',
}

function renderOrgBadgeCells(r: GetClearanceRow, orgs: Array<string>) {
  return h(
    'div',
    { class: 'clearrep-boxes' },
    orgs.map((o) =>
      r.orgs[o]
        ? h(
            'div',
            {
              class: `clearrep-org-${o}-${r.orgs[o].privLevel}`,
              title: r.orgs[o].title,
            },
            orgBadgeLabels[o]
          )
        : h('div', { class: 'clearrep-org-placeholder' })
    )
  )
}

function renderNameLinkCell(r: GetClearanceRow) {
  return h(RouterLink, { to: `/people/${r.id}` }, () => r.sortName)
}

function renderVolgisticsCell(r: GetClearanceRow) {
  return h(
    'div',
    { class: 'clearrep-volgistics', title: r.volgistics ? 'City Volunteer' : null },
    r.volgistics ? 'V' : ''
  )
}

function renderDSWCells(r: GetClearanceRow) {
  return h('div', { class: 'clearrep-boxes' }, [
    h(
      'div',
      { class: 'clearrep-dswCERT', title: r.dswCERT ? 'DSW for CERT' : null },
      r.dswCERT ? 'C' : ''
    ),
    h(
      'div',
      { class: 'clearrep-dswComm', title: r.dswComm ? 'DSW for Communications' : null },
      r.dswComm ? 'S' : ''
    ),
  ])
}

function renderBGCheckCell(r: GetClearanceRow) {
  return h(
    'div',
    { class: 'clearrep-bgCheck', title: r.bgCheck ? 'Background Check' : null },
    r.bgCheck ? 'B' : ''
  )
}

function renderBGCheckCells(r: GetClearanceRow) {
  return h('div', { class: 'clearrep-boxes' }, [
    h(
      'div',
      {
        class: `clearrep-bgCheckDOJ-${r.bgCheckDOJ}`,
        title:
          r.bgCheckDOJ == 'recorded'
            ? 'LiveScan/DOJ'
            : r.bgCheckDOJ === 'assumed'
            ? 'LiveScan/DOJ (assumed)'
            : null,
      },
      r.bgCheckDOJ ? 'D' : ''
    ),
    h(
      'div',
      {
        class: `clearrep-bgCheckFBI-${r.bgCheckFBI}`,
        title:
          r.bgCheckFBI == 'recorded'
            ? 'LiveScan/FBI'
            : r.bgCheckFBI === 'assumed'
            ? 'LiveScan/FBI (assumed)'
            : null,
      },
      r.bgCheckFBI ? 'F' : ''
    ),
    h(
      'div',
      {
        class: `clearrep-bgCheckPHS-${r.bgCheckPHS}`,
        title:
          r.bgCheckPHS == 'recorded'
            ? 'Personal History'
            : r.bgCheckPHS === 'assumed'
            ? 'Personal History (assumed)'
            : null,
      },
      r.bgCheckPHS ? 'P' : ''
    ),
  ])
}

function renderIdentCells(r: GetClearanceRow) {
  return h('div', { class: 'clearrep-boxes' }, [
    h(
      'div',
      { class: 'clearrep-idPhoto', title: r.idPhoto ? 'Photo ID' : null },
      r.idPhoto ? 'P' : ''
    ),
    h(
      'div',
      { class: 'clearrep-cardKey', title: r.cardKey ? 'Card Key' : null },
      r.cardKey ? 'C' : ''
    ),
    h(
      'div',
      { class: 'clearrep-certShirtLS', title: r.certShirtLS ? 'Green CERT Shirt (LS)' : null },
      r.certShirtLS ? 'S' : ''
    ),
    h(
      'div',
      { class: 'clearrep-certShirtSS', title: r.certShirtSS ? 'Green CERT Shirt (SS)' : null },
      r.certShirtSS ? 'S' : ''
    ),
    h(
      'div',
      { class: 'clearrep-servShirt', title: r.servShirt ? 'Tan SERV Shirt' : null },
      r.servShirt ? 'S' : ''
    ),
  ])
}

function renderCSVButton(report: GetClearance) {
  if (!report.rows.length) return null
  function exportCSV() {
    window.location.href = `/api/reports/clearance?role=${report.parameters.role}&with=${report.parameters.with}&without=${report.parameters.without}&format=csv`
  }
  return h(
    'div',
    { id: 'clearrep-buttons' },
    h(SButton, { variant: 'primary', onClick: exportCSV }, () => 'Export')
  )
}
