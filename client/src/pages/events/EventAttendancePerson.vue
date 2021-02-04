<!--
EventAttendancePerson displays the entry for a single person on the attendance
page for an event.
-->

<template lang="pug">
.event-attend-person(@click='onClick')
  span.event-attend-person-hours(:class='hoursClass', v-text='hoursLabel')
  span(v-text='fmtName(person)')
</template>

<script lang="ts">
import { defineComponent, computed, PropType } from 'vue'
import type { GetEventPerson } from './EventAttendance.vue'

export default defineComponent({
  props: {
    person: { type: Object as PropType<GetEventPerson>, required: true },
  },
  emits: ['toggle'],
  setup(props, { emit }) {
    const hoursClass = computed(() =>
      props.person.attended
        ? `event-attend-person-${props.person.attended.type.toLowerCase()}`
        : null
    )
    const hoursLabel = computed(() => {
      if (!props.person.attended) return ''
      if (!props.person.attended.minutes) return props.person.attended.type.substr(0, 1)
      let label = Math.floor(props.person.attended.minutes / 60).toString()
      if (props.person.attended.minutes % 60 !== 0) label += '½'
      return label
    })
    function fmtName(person: GetEventPerson) {
      if (person.callSign) return `${person.sortName} ${person.callSign}`
      return person.sortName
    }
    function onClick() {
      emit('toggle')
    }
    return { fmtName, hoursClass, hoursLabel, onClick }
  },
})
</script>

<style lang="postcss">
.event-attend-person {
  margin-right: 0.75rem;
}
.event-attend-person-hours {
  display: inline-block;
  margin-right: 0.25rem;
  width: 3rem;
  border-radius: 4px;
  color: white;
  text-align: center;
  line-height: 1.2;
}
.event-attend-person-volunteer {
  background-color: #4363d8;
}
.event-attend-person-student {
  background-color: #3cb44b;
}
.event-attend-person-audit {
  background-color: #a9a9a9;
}
.event-attend-person-absent {
  background-color: #800000;
}
</style>
