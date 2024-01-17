// A search control is an input field that supports autocomplete from Algolia
// search.  In HTML, it is coded as
//   <input class="s-search" s-filter="...">
// It accepts all the usual <input> attributes.
//
// The control also accepts an "s-filter" attribute, which contains filter(s)
// for the search.  This is technically optional, but in practice it's required.
//
// The control's "s-value" attribute reflects the Algolia object ID of the
// control's current "value".  The "s-value" is empty when the "value" is not a
// valid search result.  When an initial "value" attribute is set for the
// control, the form should also set the "s-value" attribute to the
// corresponding object ID.
//
// When included in a form, the search control submits *two* values for its
// "name" attribute: the control's "value" and its "s-value".
//
// Functionally, the control has two editing states: searching and selected.
// The control starts in the selected state.  In this state, the entire value is
// selected when the input box receives focus, so that typing overwrites it.  If
// the user removes the value from the field, the selection is cleared.  The
// control remains in selected state with nothing selected.  The "value" and
// "s-value" are both empty.
//
// The control enters the searching state when the user makes any change to the
// value *except* removing it.  The control will search for the value, and
// display the first few matches in a drop-down underneath the control, with
// none of them highlighted.  (If there are no matches, it displayes the
// drop-down with a "no match" message in it.)  If the user clicks on a match,
// that result is selected, the drop-down is closed, and the control returns to
// the selected state.  If the user hits the down arrow key, the first match is
// highlighted; subsequent up and down arrows move the highlight in the list of
// matches.
//
// If the input box loses focus while in searching state, the control returns to
// selected state.  If a match was highlighted at the time, that match becomes
// the new selection.  Otherwise the text in the input box is the new "value"
// and the "s-value" is empty.
//
// The control emits a bubbling "s-change" event whenever the selection changes.
// More specifically: whenever the user clicks on a match, or when the control
// loses focus with a different value selected than it had before.
up.compiler('input.s-search', elm => {
  let highlight // which item in the dropdown is highlighted.
  let client    // Algolia search client
  let lastvalue // last s-value reported
  // Set up the dropdown for displaying search results (initially not in DOM).
  const dropdown = document.createElement('div')
  dropdown.className = 's-search-dd'
  // Set up a hidden input for the second value (s-value).
  const svalue = document.createElement('input')
  svalue.type = 'hidden'
  svalue.name = elm.name
  lastvalue = svalue.value = elm.getAttribute('s-value') || ''
  elm.parentElement.insertBefore(svalue, elm.nextSibling)
  // When the input receives focus, select the text.
  elm.addEventListener('focus', () => { elm.select() })
  // When the input loses focus, select any highlighted item in the dropdown and
  // close it.
  elm.addEventListener('blur', () => {
    if (highlight) {
      elm.value = highlight.textContent
      svalue.value = highlight.dataset.value
      if (svalue.value !== lastvalue) {
        elm.setAttribute('s-value', svalue.value)
        elm.dispatchEvent(new Event('s-change', { bubbles: true }))
        lastvalue = svalue.value
      }
    }
    if (dropdown.parentElement) dropdown.parentElement.removeChild(dropdown)
  })
  // When the input is changed, handle that.
  elm.addEventListener('input', async () => {
    // If the input is cleared, clear the selection and close the dropdown.
    if (!elm.value) {
      svalue.value = ''
      if (dropdown.parentElement) dropdown.parentElement.removeChild(dropdown)
      if (svalue.value !== lastvalue) {
        elm.setAttribute('s-value', svalue.value)
        elm.dispatchEvent(new Event('s-change', { bubbles: true }))
        lastvalue = svalue.value
      }
      return
    }
    // Set up the search client.
    if (!client) client = algoliasearch(window.algoliaApplicationID, window.algoliaSearchKey).initIndex(window.algoliaIndex)
    // Run the search.
    const hits = (await client.search(elm.value, {
      filters: elm.getAttribute('s-filter'),
      hitsPerPage: 10,
    })).hits
    // Display the results.
    while (dropdown.firstChild) dropdown.removeChild(dropdown.firstChild)
    highlight = null
    const rect = elm.getBoundingClientRect()
    let bottom = rect.bottom, left = rect.left
    // If the input is in an up-modal (as it usually will be), we'll create the
    // dropdown as a child of the modal viewport, and offset the input rectangle
    // to reflect unpoly's margins on the modal viewport.  Ugh.  Otherwise we'll
    // create it as a child of the <body> with no offset to the input rectangle.
    let parent = elm.closest('up-modal-viewport')
    if (parent) bottom += 25, left += 15
    else parent = document.body
    bottom += parent.scrollTop
    if (!dropdown.parentElement) parent.appendChild(dropdown)
    // Position the dropdown below the input rectangle.  This sort of absolute
    // positioning is ugly, and fragile if someone resizes or changes
    // orientation while the dropdown is open.  But it's the only way to
    // position the dropdown adjacent to the input without a wrapper element.
    dropdown.style.top = `${bottom}px`
    dropdown.style.left = `${left}px`
    // If we didn't find anything, display "no match".
    if (!hits.length) {
      const nm = document.createElement('div')
      nm.className = 's-search-sr nomatch'
      nm.textContent = document.documentElement.lang === 'es' ? '(no hay resultados)' : '(no match)'
      dropdown.appendChild(nm)
      return
    }
    // Create a child of the dropdown for each match.
    hits.forEach(hit => {
      const hd = document.createElement('div')
      hd.className = 's-search-sr'
      hd.dataset.value = hit.objectID
      hd.textContent = hit.label
      hd.addEventListener('click', () => {
        elm.value = hit.label
        svalue.value = hit.objectID
        if (dropdown.parentElement) dropdown.parentElement.removeChild(dropdown)
        if (svalue.value !== lastvalue) {
          elm.setAttribute('s-value', svalue.value)
          elm.dispatchEvent(new Event('s-change', { bubbles: true }))
          lastvalue = svalue.value
        }
      })
      // Prevent mouse down on the hit from causing the input to lose focus.
      // That would hide the dropdown and we'd never get the click event.
      hd.addEventListener('mousedown', evt => { evt.preventDefault() })
      dropdown.appendChild(hd)
    })
  })
  // Handle keystrokes that control the highlight.
  elm.addEventListener('keydown', evt => {
    // Ignore if dropdown is closed.
    if (!up.element.isVisible(dropdown)) return
    // Remove the selection from the highlight in case we're about to change it.
    if (highlight) highlight.classList.remove('selected')
    // Down arrow moves the selection or selects the first item.
    if (evt.key === 'ArrowDown') {
      evt.preventDefault()
      if (highlight) {
        highlight = highlight.nextSibling || highlight
      } else {
        const first = dropdown.children[0]
        if (!first.classList.contains('nomatch')) highlight = first
      }
    }
    // Up arrow moves the selection.
    if (evt.key === 'ArrowUp') {
      evt.preventDefault()
      if (highlight) highlight = highlight.previousSibling || highlight
    }
    // (Re-)add the selection to the highlight.
    if (highlight) highlight.classList.add('selected')
  })
})
