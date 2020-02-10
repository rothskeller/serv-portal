<!--
Reports displays the index page for generating reports.
-->

<template lang="pug">
div.mt-3(v-if="loading")
  b-spinner(small)
CERTAttendanceForm(v-else-if="certAtt" v-bind="certAtt")
</template>

<script>
import CERTAttendanceForm from './reports/CERTAttendanceForm'

export default {
  components: { CERTAttendanceForm },
  data: () => ({ loading: false, certAtt: null }),
  async created() {
    this.$store.commit('setPage', { title: 'Reports' })
    this.loading = true
    const data = (await this.$axios.get('/api/reports')).data
    this.certAtt = data.certAttendance
    this.loading = false
  },
}
</script>

<style lang="stylus">
.report-title
  font-size 1.5rem
</style>
