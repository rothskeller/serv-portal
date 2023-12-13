up.on('click', '.signupShiftHave>a', (evt, elm) => {
  evt.preventDefault()
  elm = elm.parentElement
  while (!elm.classList.contains('signupShiftList'))
    elm = elm.nextElementSibling
  up.element.toggle(elm)
})
up.on('click', '.signupShiftCheck[title] + .s-check-lb', (evt, elm) => {
  while (!elm.classList.contains('signupShiftDisabled'))
    elm = elm.nextElementSibling
  up.element.toggle(elm, true)
})
up.on('change', '.signupShiftCheck', (evt, elm) => {
  evt.preventDefault()
  const signedup = elm.checked ? 'true' : 'false'
  elm.checked = false
  elm.indeterminate = true
  const form = elm.closest('form')
  form.elements['signedup'].value = signedup
  form.elements['shift'].value = elm.dataset['shift']
  up.submit(form, { navigate: false })
})
up.on('click', '.signupShiftRemove', (evt, elm) => {
  evt.preventDefault()
  const form = elm.closest('form')
  form.elements['signedup'].value = 'false'
  form.elements['shift'].value = elm.dataset['shift']
  form.elements['person'].value = elm.dataset['person']
  up.submit(form, { navigate: false })
})
