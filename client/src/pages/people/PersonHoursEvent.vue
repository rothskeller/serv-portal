<!--
PersonHoursEvent displays one event on the PersonHours page.
-->

<template lang="pug">
tr
  td.person-hours-event
    | {{ eventText }}
    .person-hours-event-includes(v-if='eventIncludes', v-text='eventIncludes')
  td
    input.person-hours-time(
      type='number',
      min='0',
      step='0.5',
      :value='eventTime',
      @change='setEventTime'
    )
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import type { GetPersonHoursEvent } from './PersonHours.vue'

export default defineComponent({
  props: {
    event: { type: Object as PropType<GetPersonHoursEvent>, required: true },
  },
  setup(props) {
    const eventIncludes = computed(() => {
      if (!props.event.placeholder) return null
      switch (props.event.name) {
        case 'Other CERT Hours':
          return 'Includes contact tracing'
        case 'Other LISTOS Hours':
          return 'Includes PEP and Outreach'
        default:
          return null
      }
    })
    const eventText = computed(() => {
      if (props.event.placeholder) return props.event.name
      return `${props.event.date} ${props.event.name}`
    })
    const eventTime = computed(() => {
      if (props.event.minutes === 0) return ''
      return Math.floor(props.event.minutes / 30) / 2
    })
    function setEventTime(evt: Event) {
      props.event.minutes = parseFloat((evt.target as HTMLInputElement).value) * 60
    }
    return {
      eventText,
      eventIncludes,
      eventTime,
      setEventTime,
    }
  },
})
</script>

<style lang="postcss">
.person-hours-event {
  padding-right: 1rem;
  padding-left: 0.5rem;
  text-indent: -0.5rem;
  line-height: 1.2;
}
.person-hours-event-includes {
  font-style: italic;
  font-size: 0.75rem;
}
.person-hours-time {
  width: 4rem;
  text-align: right;
}
</style>
