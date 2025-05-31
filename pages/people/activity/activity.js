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
up.on('change', '.activityHours', () => {
  let hours = 0, halves = 0
  document.querySelectorAll('s-hours').forEach(sh => {
    let v = sh.getAttribute('value')
    if (v.endsWith('½')) {
      halves++
      v = v.substring(0, v.length-1)
    }
    if (v) hours += parseInt(v)
  })
  if (halves >= 2) {
    wholes = Math.floor(halves/2)
    hours += wholes
    halves -= wholes*2
  }
  let total = hours.toString()
  if (halves) total += '½'
  document.getElementById('activityTotal').textContent = total
})
