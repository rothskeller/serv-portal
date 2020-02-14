<!--
Privileges displays the privilege choices for a role acting on a group.
-->

<template lang="pug">
.privileges
  b-btn.mr-2(
    :variant="privs.member ? 'primary' : 'outline-primary'"
    size="sm"
    @click="toggleMember"
  ) M
  b-btn-group(size="sm")
    b-btn(
      :variant="privs.roster ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleRoster"
    ) R
    b-btn(
      :variant="privs.contact ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleContact"
    ) C
    b-btn(
      :variant="privs.admin ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleAdmin"
    ) A
    b-btn(
      :variant="privs.events ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleEvents"
    ) E
    b-btn(
      :variant="privs.texts ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleTexts"
    ) T
    b-btn(
      :variant="privs.emails ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleEmails"
    ) @
    b-btn(
      :variant="privs.bcc ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleBCC"
    ) B
    b-btn(
      :variant="privs.folders ? 'primary' : 'outline-primary'"
      size="sm"
      @click="toggleFolders"
    ) F
</template>

<script>
export default {
  props: {
    privs: Object,
  },
  model: {
    prop: 'privs',
    event: 'change',
  },
  methods: {
    toggleMember() {
      this.$emit('change', { ...this.privs, member: !this.privs.member })
    },
    toggleRoster() {
      const np = { ...this.privs, roster: !this.privs.roster }
      if (!np.roster)
        np.contact = np.admin = np.events = np.texts = false
      this.$emit('change', np)
    },
    toggleContact() {
      const np = { ...this.privs, contact: !this.privs.contact }
      if (np.contact)
        np.roster = true
      else
        np.texts = false
      this.$emit('change', np)
    },
    toggleAdmin() {
      const np = { ...this.privs, admin: !this.privs.admin }
      if (np.admin)
        np.roster = true
      this.$emit('change', np)
    },
    toggleEvents() {
      const np = { ...this.privs, events: !this.privs.events }
      if (np.events)
        np.roster = true
      this.$emit('change', np)
    },
    toggleTexts() {
      const np = { ...this.privs, texts: !this.privs.texts }
      if (np.texts)
        np.roster = np.contact = true
      this.$emit('change', np)
    },
    toggleEmails() {
      this.$emit('change', { ...this.privs, emails: !this.privs.emails })
    },
    toggleBCC() {
      this.$emit('change', { ...this.privs, bcc: !this.privs.bcc })
    },
    toggleFolders() {
      this.$emit('change', { ...this.privs, bcc: !this.privs.folders })
    },
  },
}
</script>

<style lang="stylus">
.privileges
  margin 0.125rem 0
  white-space nowrap
  .btn-outline-primary:hover
    background-color #fff
    color #007bff
</style>
