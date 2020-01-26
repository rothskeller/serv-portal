<!--
Logout handles a logout request.
-->

<template lang="pug">
div.mt-3.ml-2 Logging out...
</template>

<script>
export default {
  async created() {
    await this.$axios.post('/api/logout')
    this.$store.commit('login', null)
    this.$router.replace('/login')
  },
}
</script>
