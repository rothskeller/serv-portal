; (function () {
  const rowsChanged = new Set()
  let submitted = false
  // The s-validate attribute is similar to up-validate, in that it causes
  // fields marked with it to be validated when they are changed.  However, it's
  // triggered differently, in line with my preferences for this site.
  up.compiler('[s-validate]', (elm) => {
    const row = elm.closest('.formRow')
    const target = elm.getAttribute('s-validate') || up.fragment.toTarget(row)
    rowsChanged.clear()
    elm.addEventListener('input', () => {
      rowsChanged.add(row)
      const error = up.element.get(row, '.formError')
      if (error) error.remove()
    })
    elm.addEventListener('blur', (evt) => {
      // If the field value hasn't changed, we don't need to validate.
      if (!rowsChanged.has(row)) return
      const to = evt.relatedTarget
      // If focus is moving to another control in the same form row, we
      // shouldn't validate.
      if (to && to.closest('.formRow') === row) return
      // If focus is moving to a submit button, we should wait before
      // validating.  We might be actually submitting the form, in which case
      // we shouldn't validate.
      if (to && (to.tagName === 'BUTTON' || (to.tagName === 'INPUT' && to.type === 'submit'))) {
        submitted = false
        setTimeout(() => {
          if (!submitted) up.validate(elm, { target })
        }, 100)
        return
      }
      // We should validate.
      console.log('up.validate(', elm, ', {', target, '})')
      up.validate(elm, { target })
    })
  })
  up.on('submit up:form:submit', () => { submitted = true })
})()
