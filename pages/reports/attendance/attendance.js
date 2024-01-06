up.on('input', '#attrepParamsDaterange', (evt, elm) => {
  elm.querySelectorAll('option').forEach(opt => {
    if (opt.selected) {
      const dates = document.getElementById('attrepParamsDates')
      dates.textContent = `${opt.dataset.from} to\n${opt.dataset.to}`
    }
  })
})
up.compiler('.attrepForm', form => {
  up.on(form, 'input', () => {
    up.submit(form, { target: '.attrepTable, .attrepPcount', history: true })
  })
})
up.on('click', '#attrepExport', () => {
  const params = new URLSearchParams(new FormData(up.element.get('.attrepForm')))
  params.set('format', 'csv')
  window.location.href = `/reports/attendance?` + params.toString()
})
