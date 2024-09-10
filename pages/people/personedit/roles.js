; (function () {
  let checks, implies
  function setImplies() {
    Object.values(checks).forEach(cb => {
      if (cb.hasAttribute('disabled')) {
        cb.removeAttribute('disabled')
        cb.removeAttribute('checked')
      }
    })
    let changed = true
    while (changed) {
      changed = false
      Object.values(checks).forEach(cb => {
        if (!cb.hasAttribute('checked')) return
        implies[cb.getAttribute('value')].forEach(impid => {
          const imp = checks[impid]
          if (!imp || imp.hasAttribute('disabled')) return
          imp.setAttribute('disabled', 'disabled')
          imp.setAttribute('checked', 'checked')
          changed = true
        })
      })
    }
  }
  up.compiler('.personeditRoles', elm => {
    checks = {}
    implies = {}
    elm.querySelectorAll('.s-check').forEach(cb => {
      const id = cb.getAttribute('value')
      checks[id] = cb
      implies[id] = cb.dataset.implies ? cb.dataset.implies.split(',') : []
      cb.addEventListener('input', setImplies)
    })
    setImplies()
    return () => { checks = implies = null }
  })
})()
