<!--
Modal displays its contents (default slot) as a modal dialog box when requested
by a call to its show() method.  The dialog is closed when the modal's close()
method is called.  The show() method returns a Promise which is fulfilled with
the argument passed to the close() method.

Components can reach the show() and close() methods through a reference to the
Modal component.  The close() method is also provided to the default slot as a
scoped slot property.
-->

<template lang="pug">
teleport(v-if='showing', to='#modal-port')
  #modal-base
    #modal-dialog
      #modal-content
        slot(:close='close')
  #modal-backdrop
</template>

<script lang="ts">
import { defineComponent, onBeforeUnmount, ref } from 'vue'

export default defineComponent({
  setup() {
    const showing = ref(false)
    let resolve: undefined | ((value?: unknown) => void)
    let reject: undefined | ((reason?: any) => void)
    function close(value: any) {
      showing.value = false
      if (resolve) resolve(value)
      resolve = reject = undefined
    }
    onBeforeUnmount(() => {
      showing.value = false
      if (reject) reject(new Error('modal unmounted without calling close()'))
      resolve = reject = undefined
    })
    function show() {
      showing.value = true
      return new Promise((res, rej) => {
        resolve = res
        reject = rej
      })
    }
    return { close, show, showing }
  },
})
</script>

<style lang="postcss">
#modal-base {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1050;
  width: 100vw;
  height: 100vh;
  outline: 0;
  overflow-x: hidden;
  overflow-y: auto;
}
#modal-dialog {
  max-width: 500px;
  margin: 1.75rem auto;
  position: relative;
  width: auto;
  pointer-events: none;
}
#modal-content {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  pointer-events: auto;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid rgba(0 0 0 0.2);
  border-radius: 0.3rem;
  outline: 0;
}
#modal-backdrop {
  opacity: 0.5;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1040;
  width: 100vw;
  height: 100vh;
  background-color: #000;
}
</style>
