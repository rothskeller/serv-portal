up.on('input', '.personeditSubscriptions input', (evt, elm) => {
  const warnings = document.getElementById('personeditSubscriptionsWarnings')
  warnings.textContent = ''
  document.querySelectorAll('.personeditSubscriptions .s-check').forEach(check => {
    if (check.hasAttribute('checked')) return
    if (!check.dataset.warnroles) return
    const p = document.createElement('p')
    p.textContent = check.dataset.warnroles
    warnings.appendChild(p)
  })
})
