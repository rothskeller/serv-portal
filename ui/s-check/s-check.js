// For a SERV-styled checkbox, use
//   <input type=checkbox class=s-check label="Label">
// Note: setting the input state to "indeterminate" displays a loading indicator
// in place of the checkbox.
; (function () {
  let idseq = 0
  up.compiler('input[type=checkbox].s-check', input => {
    let id = input.id
    if (!id) {
      idseq++
      id = `s-check-${idseq}`
      input.id = id
    }
    const label = document.createElement('label')
    label.className = 's-check-lb'
    label.htmlFor = id
    label.textContent = input.getAttribute('label')
    if (input.hasAttribute('title')) label.setAttribute('title', input.getAttribute('title'))
    input.parentElement.insertBefore(label, input.nextSibling)
  })
}())
