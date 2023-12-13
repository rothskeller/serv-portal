up.on('click', '.textviewGridPerson', (evt, elm) => {
  if (!navigator.userAgent.match(/Android/i) && !navigator.userAgent.match(/iPhone/i)) return
  const ndiv = up.element.get(elm, '.textviewGridNumber')
  if (!ndiv || !ndiv.textContent) return
  window.location.href = `sms:${ndiv.textContent}`
})
