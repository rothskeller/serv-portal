<!--
Logout handles a logout request.
-->

<template lang="pug">
Page(title="Logout" subtitle="Logging out..." menuItem="logout")
</template>

<script>
export default {
  async created() {
    await this.$axios.post('/api/logout')
    this.$router.replace('/')
    this.$store.commit('login', null)
  },
}
</script>
