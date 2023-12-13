up.on('input', '.personeditAddressLine1', (evt, elm) => {
  const csz = up.element.get(elm.closest('.personeditAddress'), '.personeditAddressLine2')
  if (elm.value && !csz.value) csz.value = 'Sunnyvale, CA'
  if (!elm.value && csz.value) csz.value = ''
})
up.on('input', '.personeditAddress input[type=checkbox]', (evt, elm) => {
  if (!elm.checked) setTimeout(() => {
    up.element.get(elm.closest('.personeditAddress'), '.personeditAddressLine1').focus()
  }, 0)
})
