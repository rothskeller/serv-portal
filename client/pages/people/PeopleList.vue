<!--
PeopleList displays the list of people.
-->

<template lang="pug">
#people-list
  select#people-group(v-if="groups && groups.length > 1" v-model="group")
    option(v-for="t in groups" :key="t.id" :value="t.id" v-text="t.name")
  div.mt-3(v-if="loading")
    b-spinner(small)
  template(v-else)
    #people-table
      .people-person.people-heading Person
      .people-contact.people-heading Contact Info
      .people-roles.people-heading Roles
      template(v-for="p in people")
        .people-person
          router-link(:to="`/people/${p.id}`" v-text="p.callSign ? `${p.sortName} (${p.callSign})` : p.sortName")
          div.people-badge(v-for="b in p.badges" :class="b.startsWith('No') ? 'people-badge-no' : 'people-badge-yes'" v-text="b")
        .people-contact
          div(v-if="p.email")
            a(:href="`mailto:${p.email}`" v-text="p.email")
          div(v-if="p.email2")
            a(:href="`mailto:${p.email2}`" v-text="p.email2")
          .people-phone(v-if="p.cellPhone")
            a(:href="`tel:${p.cellPhone}`" v-text="p.cellPhone")
            |
            | (Cell)
          .people-phone(v-if="p.homePhone")
            a(:href="`tel:${p.homePhone}`" v-text="p.homePhone")
            |
            | (Home)
          .people-phone(v-if="p.workPhone")
            a(:href="`tel:${p.workPhone}`" v-text="p.workPhone")
            |
            | (Work)
        .people-roles
          div(v-for="(r, i) in p.roles" :key="i" v-text="r")
          div(v-if="!p.roles.length") &mdash;
    div(v-if="people.length === 1") 1 person listed.
    div(v-else v-text="`${people.length} people listed.`")
</template>

<script>
import Cookies from 'js-cookie'

export default {
  data: () => ({ group: 0, groups: null, people: null, canAdd: false, loading: false }),
  created() {
    this.group = this.$route.query.group || Cookies.get('serv-people-group') || 0
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
      this.loading = true
      const data = (await this.$axios.get('/api/people', { params: { group: this.group } })).data
      this.people = data.people
      this.canAdd = data.canAdd
      if (data.viewableGroups.length > 1) {
        data.viewableGroups.unshift({ id: 0, name: '(all)' })
        this.groups = data.viewableGroups
      }
      this.loading = false
    },
  },
}
</script>

<style lang="stylus">
#people-list
  padding 1.5rem 0.75rem 0.75rem
#people-table
  display flex
  flex-direction column
  margin-top 0.75rem
  @media (min-width: 576px)
    display grid
    grid auto / 10rem 1fr
  @media (min-width: 800px)
    grid auto / 10rem 1fr 1fr
  @media (min-width: 960px)
    grid auto / 10rem 21rem 1fr
.people-heading
  display none
  @media (min-width: 576px)
    display block
    font-weight bold
.people-person
  overflow hidden
  margin-top 0.25rem
  padding-left 1.5rem
  text-indent -1.5rem
  text-overflow ellipsis
  white-space nowrap
  line-height 1.2
  @media (min-width: 576px)
    margin-top 0.75rem
    text-overflow clip
    white-space normal
.people-badge
  display inline-block
  margin-left 1.5rem
  padding 0 0.25rem
  border-radius 4px
  color white
  text-indent 0
  font-size 0.75rem
.people-badge-no
  background-color red
.people-badge-yes
  background-color green
.people-contact
  margin-left 6rem
  div
    overflow hidden
    text-overflow ellipsis
    white-space nowrap
  @media (min-width: 576px)
    margin-top 0.75rem
    margin-left 0.25rem
.people-phone
  font-variant tabular-nums
.people-roles
  display none
  @media (min-width: 800px)
    display block
    margin-top 0.75rem
    margin-left 0.25rem
    div
      overflow hidden
      text-overflow ellipsis
      white-space nowrap
</style>
