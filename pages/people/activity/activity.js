up.on('change', '.activityYear', (evt, elm) => {
  location.href = location.href.replace(/[^/]*$/, elm.getAttribute('value'))
})
up.on('change', '.activity s-month', (evt, elm) => {
  location.href = location.href.replace(/[^/]*$/, elm.getAttribute('value'))
})
up.on('input', '.activityHours input', () => {
  up.element.toggle(document.querySelector('.activityButtons'), true)
})
