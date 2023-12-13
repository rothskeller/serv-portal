// An <s-searchcombo> element is a search/select input control driven by data
// from Algolia search.  Visually, it appears to be a text input element.  It
// has two editing states:  searching and selected.
//
// The control starts in the "selected" state.  In this state, the label for the
// currently selected item (if any) is displayed in the text input element, and
// its corresponding Algolia ID is the value of the control on form submission.
// The entire label is selected when the input box receives focus, so that
// typing overwrites it.  If the user removes the text from the field, the
// selection is cleared; the control remains in "selected" state with nothing
// selected, and its value on form submission is an empty string.
//
// The control enters the "searching" state when the user makes any change to
// the input element text *except* emptying it.  The control will search for the
// text in the input element, and display the first few matches in a drop-down
// underneath the text input, with none of them highlighted.  (If there are no
// matches, it displays the drop-down with a "no match" message in it.)  If the
// user clicks on a match, that result is selected, the drop-down is closed, and
// the control returns to "selected" state.  If the user hits the down arrow
// key, the first match is highlighted; subsequent down arrows move the
// highlight down the list of matches.  The up arrow moves the highlight up in
// the list of matches.
//
// If the input box loses focus while in "searching" state, the control returns
// to "selected" state.  If a match was highlighted at the time, or if there was
// only one match in the drop-down, that match becomes the new selection.
// Otherwise, the selection is unchanged (and the text input returns to showing
// its label).
//
// Attributes:
//
// autofocus - if set, puts autofocus on the input field
// facet - facet filter string to restrict search results; technically optional
//     but essential in practice
// form - identifies the form to which to submit the control's value (optional;
//     defaults to the containing <form> element)
// name - parameter name under which the control's value is submitted with the
//     form (required if form submission is expected, otherwise optional)
// optionalFilters - optional filters for the search; these affect ranking
//     without changing overall resullts
// placeholder - text to display in the input element when its value is empty
//     (optional)
// value - initial value of the control (optional)
// valuelabel - text to display in input field while the initial value remains
//     selected; should be non-empty iff value is non-empty
//
// Events dispatched:
//
// change - dispatched when the value of the control changes
// edit - dispatched when the edit button (if any) is pressed
class SSearchCombo extends HTMLElement {
  static get observedAttributes() { return ['value', 'valuelabel'] };
  connectedCallback() {
    this.classList.add('s-searchcombo')
    this._hi = document.createElement('input')
    this._hi.type = 'hidden'
    if (this.hasAttribute('form')) this._hi.setAttribute('form', this.getAttribute('form'))
    if (this.hasAttribute('name')) this._hi.setAttribute('name', this.getAttribute('name'))
    if (this.hasAttribute('value')) this._hi.setAttribute('value', this.getAttribute('value'))
    this.appendChild(this._hi)
    this._in = document.createElement('input')
    this._in.className = 's-searchcombo-in formInput'
    if (this.hasAttribute('placeholder')) this._in.setAttribute('placeholder', this.getAttribute('placeholder'))
    if (this.hasAttribute('autofocus')) this._in.setAttribute('autofocus', this.getAttribute('autofocus'))
    this._saveID = this.id || ''
    if (this.id) {
      this.id += '_top'
      this._in.id = this._saveID
    }
    this.appendChild(this._in)
    if (this.hasAttribute('autofocus')) { this._in.focus() }
    this._dd = document.createElement('div')
    this._dd.className = 's-searchcombo-dd'
    this._dd.style.display = 'none'
    this.appendChild(this._dd)
    this._in.addEventListener('focus', this.onFocus.bind(this))
    this._in.addEventListener('blur', this.onBlur.bind(this))
    this._in.addEventListener('input', this.onInput.bind(this))
    this._in.addEventListener('keydown', this.onKeyDown.bind(this))
    this.attributeChangedCallback('value')
  }
  disconnectedCallback() {
    this.removeChild(this._hi)
    this.removeChild(this._in)
    this.removeChild(this._dd)
    this._hi = this._in = this._dd = this._hl = null
    this.id = this._saveID
  }
  attributeChangedCallback(name, oldValue, newValue) {
    if (!this._hi) return
    if (name === 'value' || name === 'valuelabel') {
      if (this.hasAttribute('value')) {
        const old = this._hi.value
        this._hi.value = this.getAttribute('value')
        if (old !== this._hi.value) {
          this.dispatchEvent(new Event('change'))
        }
      }
      if (this.hasAttribute('valuelabel')) this._in.value = this.getAttribute('valuelabel')
      else if (this.hasAttribute('value')) this._in.value = this.getAttribute('value')
    }
  }
  onFocus() {
    this._in.select()
  }
  onBlur() {
    if (this._dd.style.display) return
    if (this._hl) {
      this.setAttribute('valuelabel', this._hl.textContent)
      this.setAttribute('value', this._hl.getAttribute('data-value'))
    } else if (this._dd.childElementCount == 1) {
      const one = this._dd.children[0]
      if (one.classList.contains('nomatch')) {
        if (this.hasAttribute('valuelabel')) this._in.value = this.getAttribute('valuelabel')
        else if (this.hasAttribute('value')) this._in.value = this.getAttribute('value')
      } else {
        this.setAttribute('valuelabel', one.textContent)
        this.setAttribute('value', one.getAttribute('data-value'))
      }
    }
    this._dd.style.display = 'none'
  }
  async onInput(evt) {
    if (!this._in.value) {
      this._dd.style.display = 'none'
      this.setAttribute('valuelabel', '')
      this.setAttribute('value', '')
    } else {
      if (!this._sc) this._sc = algoliasearch(window.algoliaApplicationID, window.algoliaSearchKey).initIndex(window.algoliaIndex)
      const hits = (await this._sc.search(this._in.value, {
        filters: this.hasAttribute('facet') ? this.getAttribute('facet') : null,
        optionalFilters: this.hasAttribute('optionalFilters') ? this.getAttribute('optionalFilters') : null,
        hitsPerPage: 10,
      })).hits
      while (this._dd.firstChild) this._dd.removeChild(this._dd.firstChild)
      this._dd.style.display = null
      this._hl = null
      if (!hits.length) {
        const nm = document.createElement('div')
        nm.className = 's-searchcombo-sr nomatch'
        nm.textContent = '(no match)'
        this._dd.appendChild(nm)
        return
      }
      hits.forEach(hit => {
        const hd = document.createElement('div')
        hd.className = 's-searchcombo-sr'
        hd.setAttribute('data-value', hit.objectID)
        let title = hit.label
        if (hit.context) title += ` (${hit.context})`
        hd.setAttribute('title', title)
        hd.textContent = hit.label
        hd.addEventListener('mousedown', evt => { evt.preventDefault() })
        // This prevents mouse down on the hit from causing the input to lose
        // focus, which would hide the dropdown, which would mean we'd never get
        // the click event.
        hd.addEventListener('click', this.onResultClick.bind(this))
        this._dd.appendChild(hd)
      })
    }
  }
  onKeyDown(evt) {
    if (this._dd.style.display) return
    switch (evt.key) {
      case 'ArrowDown':
        evt.preventDefault()
        if (this._hl) {
          const ns = this._hl.nextElementSibling
          if (ns) {
            this._hl.classList.remove('selected')
            this._hl = ns
            this._hl.classList.add('selected')
          }
        } else if (this._dd.childElementCount) {
          const one = this._dd.children[0]
          if (!one.classList.contains('nomatch')) {
            this._hl = one
            this._hl.classList.add('selected')
          }
        }
        break
      case 'ArrowUp':
        evt.preventDefault()
        if (this._hl) {
          const ps = this._hl.previousElementSibling
          if (ps) {
            this._hl.classList.remove('selected')
            this._hl = ps
            this._hl.classList.add('selected')
          }
        }
        break
    }
  }
  onResultClick(evt) {
    this.setAttribute('valuelabel', evt.currentTarget.textContent)
    this.setAttribute('value', evt.currentTarget.getAttribute('data-value'))
    this._dd.style.display = 'none'
  }
  get value() { return this.getAttribute('value') }
}
customElements.define('s-searchcombo', SSearchCombo)
up.compiler('s-searchcombo,input.s-search', function () {
  const ascript = document.getElementById('algoliaScript')
  if (!ascript) {
    const ascript = document.createElement('script')
    ascript.id = 'algoliaScript'
    ascript.src = 'https://cdn.jsdelivr.net/npm/algoliasearch@4.20.0/dist/algoliasearch-lite.umd.js'
    ascript.crossOrigin = 'anonymous'
    ascript.setAttribute('integrity', 'sha256-DABVk+hYj0mdUzo+7ViJC6cwLahQIejFvC+my2M/wfM=')
    document.head.appendChild(ascript)
  }
})
