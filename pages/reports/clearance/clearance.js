up.compiler('.clearrepForm', form => {
  up.on(form, 'input', () => {
    up.submit(form, { target: '.clearrepTable, .clearrepCount', history: true })
  })
})
up.on('click', '#clearrepExport', () => {
  const params = new URLSearchParams(new FormData(up.element.get('.clearrepForm')))
  params.set('format', 'csv')
  window.location.href = `/reports/clearance?` + params.toString()
})
