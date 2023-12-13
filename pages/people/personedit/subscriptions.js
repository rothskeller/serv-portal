up.on('input', '.personeditSubscriptions input', (evt, elm) => {
  const warnings = document.getElementById('personeditSubscriptionsWarnings')
  warnings.textContent = ''
  up.element.all('.personeditSubscriptions .s-check').forEach(check => {
    if (check.hasAttribute('checked')) return
    if (!check.dataset.warnroles) return
    const p = document.createElement('p')
    p.textContent = `Messages sent to ${check.getAttribute('label')} are considered required for the ${check.dataset.warnroles}.  Unsubscribing from it may cause you to lose that role.`
    warnings.appendChild(p)
  })
})
