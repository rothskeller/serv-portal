import { ref } from "vue";

export const touch = ref(false)

try {
  const mql = window.matchMedia('(pointer:coarse),(hover:none),(-moz-touch-enabled:1)')
  touch.value = mql.matches
  mql.addListener(evt => { touch.value = evt.matches })
} catch (err) {
  console.error('Touch detection:', err)
}
