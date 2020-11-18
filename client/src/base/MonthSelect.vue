<!--
MonthSelect is a control allowing convenient selection and browsing of months.
-->

<template lang="pug">
.monsel
  .monsel-arrow(@click='decMonth')
    SIcon.monsel-icon(icon='chevron-left')
  SSelect#monsel-month(:options='monthOptions', v-model='month')
  .monsel-arrow(@click='incMonth')
    SIcon.monsel-icon(icon='chevron-right')
  .monsel-space
  .monsel-arrow(@click='decYear')
    SIcon.monsel-icon(icon='chevron-left')
  SSelect#monsel-year(:options='yearOptions', v-model='year')
  .monsel-arrow(@click='incYear')
    SIcon.monsel-icon(icon='chevron-right')
</template>

<script lang="ts">
import { defineComponent, ref, watch, watchEffect } from 'vue'
import moment from 'moment-mini'
import SIcon from './SIcon.vue'
import SSelect from './SSelect.vue'

const monthOptions = [
  { value: 1, label: 'January' },
  { value: 2, label: 'February' },
  { value: 3, label: 'March' },
  { value: 4, label: 'April' },
  { value: 5, label: 'May' },
  { value: 6, label: 'June' },
  { value: 7, label: 'July' },
  { value: 8, label: 'August' },
  { value: 9, label: 'September' },
  { value: 10, label: 'October' },
  { value: 11, label: 'November' },
  { value: 12, label: 'December' },
]

export default defineComponent({
  components: { SIcon, SSelect },
  props: {
    modelValue: { type: String, required: true },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    // List years from 2013 through "next" year.
    const yearOptions: Array<{ value: number; label: string }> = []
    const now = moment().year()
    for (let i = 2013; i <= now + 1; i++) {
      yearOptions.push({ value: i, label: i.toString() })
    }

    // If modelValue is empty, start with current month; otherwise start with
    // modelValue.
    const month = ref(moment().month())
    const year = ref(now)
    watchEffect(() => {
      if (props.modelValue) {
        month.value = parseInt(props.modelValue.substr(5, 2), 10)
        year.value = parseInt(props.modelValue.substr(0, 4))
      }
    })

    // Update model value with components change.
    watch([month, year], () => {
      const ym = month.value > 9 ? `${year.value}-${month.value}` : `${year.value}-0${month.value}`
      emit('update:modelValue', ym)
    })

    // Increments and decrements.
    function incMonth() {
      if (month.value !== 12) month.value++
      else if (year.value !== now + 1) {
        month.value = 1
        year.value++
      }
    }
    function decMonth() {
      if (month.value !== 1) month.value--
      else if (year.value !== 2013) {
        month.value = 12
        year.value--
      }
    }
    function incYear() {
      if (year.value < now + 1) year.value++
    }
    function decYear() {
      if (year.value > 2013) year.value--
    }

    return { decMonth, decYear, incMonth, incYear, month, monthOptions, year, yearOptions }
  },
})
</script>

<style lang="postcss">
.monsel {
  flex: auto;
  display: flex;
  justify-content: center;
  & .sselect {
    width: calc(50% - 80px);
    height: 40px;
    @media (min-width: 576px) {
      width: auto;
      font-weight: bold;
      font-size: 1.25rem;
    }
  }
  @media print {
    display: none;
  }
}
.monsel-arrow {
  flex: none;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  cursor: pointer;
  user-select: none;
  &:hover {
    background-color: #efefef;
  }
}
.monsel-icon {
  width: 0.5rem;
  @media (min-width: 576px) {
    width: 1rem;
  }
}
.monsel-space {
  @media (min-width: 576px) {
    width: 40px;
  }
}
</style>
