<!--
PersonEditAddress displays an address entry on the PersonEdit form, and handles
autocomplete and validation of addresses.
-->

<template lang="pug">
b-form-group(
  :label="`${typeFmt} Address`"
  :label-for="`person-address-${type}`"
  label-cols-sm="auto"
  :label-class="{'person-edit-label': true, 'person-edit-sameaddr-label': type !== 'Home'}"
  :state="error ? false : null"
  :invalid-feedback="error"
)
  b-checkbox.person-edit-sameaddr(v-if="type !== 'Home'" v-model="sameAsHome" :disabled="!hasHome") Same as home address
  b-input.person-edit-input(
    v-if="!sameAsHome"
    ref="line1"
    :id="`person-address-${type}`"
    :class="{'mt-1': type !== 'Home'}"
    :state="error ? false : null"
    trim
    v-model="line1"
    @blur="onBlur"
  )
  b-input.person-edit-input.mt-1(
    v-if="!sameAsHome"
    ref="line2"
    :state="error ? false : null"
    trim
    v-model="line2"
    @focus="onFocusLine2"
    @blur="onBlur"
  )
</template>

<script>
import SmartyStreetsSDK from 'smartystreets-javascript-sdk'
const SmartyStreetsCore = SmartyStreetsSDK.core
const Lookup = SmartyStreetsSDK.usStreet.Lookup
const credentials = new SmartyStreetsCore.SharedCredentials('15809213558292353')
const client = SmartyStreetsCore.buildClient.usStreet(credentials)

export default {
  model: {
    prop: 'address',
    event: 'change',
  },
  props: {
    type: String,
    address: Object,
    hasHome: Boolean,
  },
  data: () => ({
    sameAsHome: false,
    line1: '',
    line2: '',
    latitude: 0,
    longitude: 0,
    error: null,
    lastChecked: '',
  }),
  computed: {
    typeFmt() {
      if (this.type === 'Work') return 'Business Hours'
      if (this.type === 'Mail') return 'Mailing'
      return this.type
    },
  },
  mounted() {
    if (this.address && this.address.address) {
      this.lastChecked = this.address.address
      this.line1 = this.address.address.split(',')[0]
      this.line2 = this.address.address.replace(/^[^,]*, */, '')
      this.latitude = this.address.latitude
      this.longitude = this.address.longitude
    }
    if (!this.address) this.sameAsHome = this.type === 'Mail'
    else this.sameAsHome = this.address.sameAsHome || false
  },
  watch: {
    line1() {
      if (this.line1 === '') this.line2 = ''
      else if (this.line2 === '') this.line2 = 'Sunnyvale, CA'
    },
    sameAsHome() {
      if (!this.sameAsHome)
        this.$nextTick(() => { this.$refs.line1.focus() })
    },
  },
  methods: {
    async onBlur(evt) {
      if (evt.relatedTarget === this.$refs.line1.$el || evt.relatedTarget === this.$refs.line2.$el)
        return // focus is still in one of the two inputs
      if (!this.line1 && !this.line2) {
        this.$emit('change', { address: '', latitude: 0, longitude: 0 })
        return
      }
      let check = `${this.line1}, ${this.line2}`
      if (!check.match(/\W[A-Za-z][A-Za-z]\W/)) check += ', CA'
      if (check === this.lastChecked) return
      const lookup = new Lookup()
      lookup.street = check
      this.$emit('change', null) // prevent submit while looking up
      const result = await client.send(lookup).catch(console.error)
      if (result.lookups[0].result.length) {
        const r = result.lookups[0].result[0]
        this.line1 = r.deliveryLine1
        if (r.deliveryLine2) this.line1 += ', ' + r.deliveryLine2
        this.line2 = r.lastLine.replace(/-[0-9]{4}/, '')
        this.lastChecked = this.line1 + ', ' + this.line2
        this.latitude = r.metadata.latitude
        this.longitude = r.metadata.longitude
        this.error = null
        this.$emit('change', { address: this.lastChecked, latitude: this.latitude, longitude: this.longitude })
      } else {
        this.lastChecked = check
        this.error = "We couldn't locate this address."
        this.$emit('change', null)
      }
    },
    onFocusLine2() {
      this.$refs.line2.select()
    },
  },
}
</script>

<style lang="stylus">
.person-edit-sameaddr-label
  padding-top 0
.person-edit-sameaddr
  white-space nowrap
</style>
