up.on('change', '.activityYear', (evt, elm) => {
  location.href = location.href.replace(/[^/]*$/, elm.getAttribute('value'))
})
up.on('change', '.activity s-month', (evt, elm) => {
  location.href = location.href.replace(/[^/]*$/, elm.getAttribute('value'))
})
up.on('input', '.activityHours input', () => {
  const elm = document.getElementById('activitySave')
  elm.removeAttribute('disabled')
  elm.classList.toggle('sbtn-disabled', false)
  elm.classList.toggle('sbtn-secondary', false)
  elm.classList.toggle('sbtn-warning', true)
})
