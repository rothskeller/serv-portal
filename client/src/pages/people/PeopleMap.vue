<!--
PeopleMap displays people on a map.
-->

<template lang="pug">
#people-map
  #people-map-title
    select#people-map-group(v-if='groups && groups.length > 1', v-model='group')
      option(v-for='t in groups', :key='t.id', :value='t.id', v-text='t.name')
    SCheck#people-map-home.people-map-option(v-model='home', label='Home')
    SCheck#people-map-work.people-map-option(v-model='work', label='Business Hours')
  #people-map-container
    #people-map-map(ref='map')
</template>

<script lang="ts">
import { defineComponent, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import Cookies from 'js-cookie'
import axios from '../../plugins/axios'
import { SCheck } from '../../base'
import * as districts from './districts'
import type { GetPeople, GetPeoplePerson, GetPeopleViewableGroup } from './PeopleList.vue'

const mapScriptPromise = new Promise((resolve, reject) => {
  const script = document.createElement('script')
  script.src =
    'https://maps.googleapis.com/maps/api/js?key=AIzaSyDYiDjdYhCKZnM4qbK68KZRjKZqJiQ1dZw&callback=initMap'
  script.defer = true
  ;(window as any).initMap = resolve
  document.head.appendChild(script)
})

export default defineComponent({
  components: { SCheck },
  setup() {
    const route = useRoute()

    // Create the map.
    const map = ref(null as null | HTMLElement)
    let gmap: google.maps.Map<HTMLElement>
    const markers = [] as Array<google.maps.Marker>
    onMounted(() => {
      mapScriptPromise.then(() => {
        gmap = new google.maps.Map(map.value!, {
          center: { lat: 37.3801648, lng: -122.032706 },
          zoom: 13,
        })
        new google.maps.Polygon({
          paths: districts.district1,
          strokeWeight: 0,
          fillColor: '#9900CC',
          fillOpacity: 0.36,
        }).setMap(gmap)
        new google.maps.Polygon({
          paths: districts.district2,
          strokeWeight: 0,
          fillColor: '#00CC66',
          fillOpacity: 0.36,
        }).setMap(gmap)
        new google.maps.Polygon({
          paths: districts.district3,
          strokeWeight: 0,
          fillColor: '#FF9966',
          fillOpacity: 0.36,
        }).setMap(gmap)
        new google.maps.Polygon({
          paths: districts.district4,
          strokeWeight: 0,
          fillColor: '#00CCCC',
          fillOpacity: 0.36,
        }).setMap(gmap)
        new google.maps.Polygon({
          paths: districts.district5,
          strokeWeight: 0,
          fillColor: '#336633',
          fillOpacity: 0.36,
        }).setMap(gmap)
        new google.maps.Polygon({
          paths: districts.district6,
          strokeWeight: 0,
          fillColor: '#CC99CC',
          fillOpacity: 0.36,
        }).setMap(gmap)
        resetMarkers()
      })
    })

    // Create the set of markers.
    function resetMarkers() {
      if (!gmap) return
      markers.forEach((m) => {
        m.setMap(null)
      })
      markers.length = 0
      people.value.forEach((p) => {
        if (home.value && p.homeAddress && p.homeAddress.latitude && p.homeAddress.longitude)
          markers.push(
            new google.maps.Marker({
              position: { lat: p.homeAddress.latitude, lng: p.homeAddress.longitude },
              title: p.informalName,
            })
          )
        if (work.value && p.workAddress && p.workAddress.latitude && p.workAddress.longitude)
          markers.push(
            new google.maps.Marker({
              position: { lat: p.workAddress.latitude, lng: p.workAddress.longitude },
              title: p.informalName,
            })
          )
        if (
          work.value &&
          !home.value &&
          p.workAddress &&
          p.workAddress.sameAsHome &&
          p.homeAddress &&
          p.homeAddress.latitude &&
          p.homeAddress.longitude
        )
          markers.push(
            new google.maps.Marker({
              position: { lat: p.homeAddress.latitude, lng: p.homeAddress.longitude },
              title: p.informalName,
            })
          )
      })
      markers.forEach((m) => {
        m.setMap(gmap)
      })
    }

    // The group being viewed.
    const group = ref(
      parseInt((route.query.group as string) || Cookies.get('serv-people-group') || '0')
    )
    const groups = ref([] as Array<GetPeopleViewableGroup>)
    const people = ref([] as Array<GetPeoplePerson>)
    watch(
      group,
      async () => {
        Cookies.set('serv-people-group', group.value.toString(), { expires: 3650 })
        const data = (await axios.get<GetPeople>('/api/people', { params: { group: group.value } }))
          .data
        people.value = data.people
        if (data.viewableGroups.length > 1) {
          data.viewableGroups.unshift({ id: 0, name: '(all)' })
          groups.value = data.viewableGroups
        }
        resetMarkers()
      },
      { immediate: true }
    )

    // The address types being viewed.
    const home = ref(true)
    const work = ref(false)
    watch([home, work], resetMarkers)

    return { group, groups, home, map, people, work }
  },
})
</script>

<style lang="postcss">
#people-map {
  display: flex;
  flex: auto;
  flex-direction: column;
  padding: 0.75rem 0 0 0;
  height: 100%;
}
#people-map-title {
  display: flex;
  flex: none;
  flex-direction: row;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 0.75rem;
}
.people-map-option {
  margin-left: 0.75rem;
  font-size: 1rem;
}
#people-map-group {
  font-size: 1rem;
  @media (min-width: 576px) {
    margin-left: 1rem;
  }
}
#people-map-container {
  position: relative;
  flex: auto;
}
#people-map-map {
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  position: absolute;
}
</style>
