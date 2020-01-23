<!--
EventAttendancePerson displays the entry for a single person on the attendance
page for an event.
-->

<template lang="pug">
.event-attend-person(@click="$emit('toggle', person)")
  span(:class="classes" v-text="label")
  span(v-text="person.sortName")
</template>

<script>
export default {
  props: {
    person: Object,
  },
  computed: {
    classes() {
      if (!this.person.attended) return 'event-attend-person-absent'
      else return 'event-attend-person-' + this.person.attended.type.toLowerCase()
    },
    label() {
      if (!this.person.attended) return ''
      if (!this.person.attended.minutes)
        return this.person.attended.type.substr(0, 1)
      let label = Math.floor(this.person.attended.minutes / 60)
      if (this.person.attended.minutes % 60 !== 0)
        label += '½'
      return label
    },
  }
}
</script>

<style lang="stylus">
.event-attend-person
  margin-right 0.75rem
.event-attend-person-absent
  display inline-block
  margin-right 0.25rem
  width 3rem
  border-radius 4px
  color white
  text-align center
  line-height 1.2
.event-attend-person-volunteer
  @extend .event-attend-person-absent
  background-color #4363d8
.event-attend-person-student
  @extend .event-attend-person-absent
  background-color #3cb44b
.event-attend-person-audit
  @extend .event-attend-person-absent
  background-color #a9a9a9
</style>
