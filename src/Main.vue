<!--
Main is the main program of the SERV portal.
-->

<template lang="pug">
router-view(v-if="$store.state.me")
Login(v-else-if="$store.state.started")
</template>

<script>
export default {
}
</script>

<style lang="stylus" scoped></style>
