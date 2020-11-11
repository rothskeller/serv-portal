<!--
SSpinner displays a spinner for a process in progress.
-->

<template lang="pug">
span.sspinner(aria-hidden)
</template>

<style lang="postcss">
.sspinner {
  display: inline-block;
  color: #212529;
  width: 1rem;
  height: 1rem;
  vertical-align: text-bottom;
  border: 0.2em solid;
  border-right: 0.2em solid transparent;
  border-radius: 50%;
  animation: sspinner 0.75s linear infinite;
}
@keyframes sspinner {
  100% {
    transform: rotate(1turn);
  }
}
</style>
