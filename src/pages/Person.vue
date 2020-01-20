<!--
Person displays the person viewing/editing page.
-->

<template lang="pug">
Page(:title="title" menuItem="people" noPadding)
  div.mt-3(v-if="loading")
    b-spinner(small)
  //-b-card#person-card(v-else-if="tabs" no-body)
    b-tabs(card)
      b-tab.person-tab-pane(v-if="!newp" title="Details" no-body)
        PersonView(:person="person")
      b-tab.person-tab-pane(v-if="canEdit" title="Edit" no-body)
        PersonEdit(:person="person")
  //-PersonEdit(v-else-if="canEdit" :person="person")
  PersonView(v-else :person="person")
</template>

<script>
export default {
  data: () => ({
    loading: false,
    title: 'Person',
    canEdit: false,
    person: null,
  }),
  computed: {
    newp() { return this.$route.params.id === 'NEW' },
    tabs() {
      return (this.newp ? 0 : 1) + (this.canEdit ? 1 : 0) > 1
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}`)).data
    this.canEdit = data.canEdit
    this.person = data.person
    this.title = data.person.id ? data.fullName : 'New Person'
    this.loading = false
  },
}
</script>

<style lang="stylus">
#person-card
  height calc(100vh - 40px)
  border none
.person-tab-pane
  overflow-y auto
  height calc(100vh - 3.25rem - 42px)
</style>
