<!--
Reports displays the index page for generating reports.
-->

<template lang="pug">
Page(title="Reports" menuItem="reports")
  div.mt-3(v-if="loading")
    b-spinner(small)
  template(v-else)
    CERTAttendanceForm(v-bind="certAtt")
</template>

<script>
export default {
  data: () => ({ loading: false, certAtt: null }),
  async created() {
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
