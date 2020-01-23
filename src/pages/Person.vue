<!--
Person displays the person viewing/editing page.
-->

<template lang="pug">
Page(:title="title" menuItem="people" noPadding)
  div.mt-3(v-if="loading")
    b-spinner(small)
  b-card#person-card(v-else-if="tabs" no-body)
    b-tabs(card)
      b-tab.person-tab-pane(v-if="!newp" title="Details" no-body)
        PersonView(:person="person")
      b-tab.person-tab-pane(v-if="canEdit" title="Edit" no-body)
        PersonEdit(:person="person" :allowBadPassword="allowBadPassword" :canEditDetails="canEditDetails" :canEditRoles="canEditRoles" :canEditUsername="canEditUsername" :passwordHints="passwordHints")
      b-tab.person-tab-pane(title="Map" no-body)
        GmapMap(
          style="height:100%"
          :center="{lat:37.3801648,lng:-122.032706}"
          :zoom="13"
        )
          GmapPolygon(:path="districts.district1" :options="{strokeWeight: 0, fillColor: '#9900CC', fillOpacity: 0.36}")
          GmapPolygon(:path="districts.district2" :options="{strokeWeight: 0, fillColor: '#00CC66', fillOpacity: 0.36}")
          GmapPolygon(:path="districts.district3" :options="{strokeWeight: 0, fillColor: '#FF9966', fillOpacity: 0.36}")
          GmapPolygon(:path="districts.district4" :options="{strokeWeight: 0, fillColor: '#00CCCC', fillOpacity: 0.36}")
          GmapPolygon(:path="districts.district5" :options="{strokeWeight: 0, fillColor: '#336633', fillOpacity: 0.36}")
          GmapPolygon(:path="districts.district6" :options="{strokeWeight: 0, fillColor: '#CC99CC', fillOpacity: 0.36}")
  PersonEdit(v-else-if="canEdit" :person="person" :allowBadPassword="allowBadPassword" :canEditDetails="canEditDetails" :canEditRoles="canEditRoles" :canEditUsername="canEditUsername" :passwordHints="passwordHints")
  PersonView(v-else :person="person")
</template>

<script>
import * as districts from '../districts'
export default {
  data: () => ({
    loading: false,
    title: 'Person',
    canEditDetails: false,
    canEditRoles: false,
    canEditUsername: false,
    allowBadPassword: false,
    passwordHints: null,
    person: null,
    districts,
  }),
  computed: {
    canEdit() { return this.canEditDetails || this.canEditRoles },
    newp() { return this.$route.params.id === 'NEW' },
    tabs() {
      return (this.newp ? 0 : 1) + (this.canEdit ? 1 : 0) > 1
    },
  },
  async created() {
    this.loading = true
    const data = (await this.$axios.get(`/api/people/${this.$route.params.id}`)).data
    this.allowBadPassword = data.allowBadPassword
    this.canEditDetails = data.canEditDetails
    this.canEditRoles = data.canEditRoles
    this.canEditUsername = data.canEditUsername
    this.passwordHints = data.passwordHints
    this.person = data.person
    this.title = data.person.id ? data.person.informalName : 'New Person'
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
