<!--
PeopleMap displays people on a map.
-->

<template lang="pug">
#people-map
  #people-map-title
    select#people-map-group(v-if="groups && groups.length > 1" v-model="group")
      option(v-for="t in groups" :key="t.id" :value="t.id" v-text="t.name")
    b-checkbox.people-map-option(v-model="home") Home
    b-checkbox.people-map-option(v-model="work") Business Hours
  GmapMap#people-map-map(
    :center="{lat:37.3801648,lng:-122.032706}"
    :zoom="13"
  )
    GmapPolygon(:path="districts.district1" :options="{strokeWeight: 0, fillColor: '#9900CC', fillOpacity: 0.36}")
    GmapPolygon(:path="districts.district2" :options="{strokeWeight: 0, fillColor: '#00CC66', fillOpacity: 0.36}")
    GmapPolygon(:path="districts.district3" :options="{strokeWeight: 0, fillColor: '#FF9966', fillOpacity: 0.36}")
    GmapPolygon(:path="districts.district4" :options="{strokeWeight: 0, fillColor: '#00CCCC', fillOpacity: 0.36}")
    GmapPolygon(:path="districts.district5" :options="{strokeWeight: 0, fillColor: '#336633', fillOpacity: 0.36}")
    GmapPolygon(:path="districts.district6" :options="{strokeWeight: 0, fillColor: '#CC99CC', fillOpacity: 0.36}")
    template(v-for="person in people")
      GmapMarker(
        v-if="home && person.homeAddress && person.homeAddress.latitude && person.homeAddress.longitude"
        :key="person.id"
        :position="{lat: person.homeAddress.latitude, lng: person.homeAddress.longitude}"
        :options="{title: person.informalName}"
        :title="person.informalName"
      )
      GmapMarker(
        v-if="work && person.workAddress && person.workAddress.latitude && person.workAddress.longitude"
        :key="`w${person.id}`"
        :position="{lat: person.workAddress.latitude, lng: person.workAddress.longitude}"
        :options="{title: person.informalName}"
        :title="person.informalName"
      )
      GmapMarker(
        v-else-if="work && !home && person.workAddress.sameAsHome && person.homeAddress && person.homeAddress.latitude && person.homeAddress.longitude"
        :key="`w${person.id}`"
        :position="{lat: person.homeAddress.latitude, lng: person.homeAddress.longitude}"
        :options="{title: person.informalName}"
        :title="person.informalName"
      )
</template>

<script>
import Cookies from 'js-cookie'
import * as districts from '@/districts'

export default {
  data: () => ({ group: 0, groups: null, people: [], canAdd: false, home: true, work: false, districts }),
  created() {
    this.group = Cookies.get('serv-people-group') || 0
    this.load()
  },
  watch: {
    group() {
      Cookies.set('serv-people-group', this.group, { expires: 3650 })
      this.load()
    },
  },
  methods: {
    async load() {
      const data = (await this.$axios.get('/api/people', { params: { group: this.group } })).data
      this.people = data.people
      this.canAdd = data.canAdd
      if (data.viewableGroups.length > 1) {
        data.viewableGroups.unshift({ id: 0, name: '(all)' })
        this.groups = data.viewableGroups
      }
    },
  },
}
</script>

<style lang="stylus">
#people-map
  display flex
  flex auto
  flex-direction column
  padding 0.75rem 0 0 0
  height 100%
#people-map-title
  display flex
  flex none
  flex-direction row
  flex-wrap wrap
  align-items center
  margin-bottom 0.75rem
.people-map-option
  margin-left 0.75rem
  font-size 1rem
#people-map-group
  font-size 1rem
  @media (min-width: 576px)
    margin-left 1rem
#people-map-map
  flex auto
</style>
