class SYear extends HTMLElement {
  static get observedAttributes() { return ['value'] };
  constructor() {
    super()
    this._left = this.year = this._right = null
  }
  connectedCallback() {
    this.classList.add('s-year')
    this._left = document.createElement('s-icon')
    this._left.className = 's-year-arrow'
    this._left.setAttribute('icon', 'chevron-left-thin')
    this._left.addEventListener('click', this.prevYear.bind(this))
    this.appendChild(this._left)
    this._year = document.createElement('div')
    this._year.className = 's-year-year'
    this.setYearLabel()
    this.appendChild(this._year)
    this._right = document.createElement('s-icon')
    this._right.className = 's-year-arrow'
    this._right.setAttribute('icon', 'chevron-right-thin')
    this._right.addEventListener('click', this.nextYear.bind(this))
    this.appendChild(this._right)
  }
  disconnectedCallback() {
    if (this._left) this.removeChild(this._left)
    if (this._year) this.removeChild(this._year)
    if (this._right) this.removeChild(this._right)
    this._left = this._year = this._right = null
  }
  setYearLabel() {
    const value = this.getAttribute('value')
    if (!value || !this._year) return
    this._year.textContent = value
  }
  prevYear() {
    let value = this.getAttribute('value')
    value = `${parseInt(value) - 1}`
    this.setAttribute('value', value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
    this.setYearLabel()
  }
  nextYear() {
    let value = this.getAttribute('value')
    value = `${parseInt(value) + 1}`
    this.setAttribute('value', value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
    this.setYearLabel()
  }
  attributeChangedCallback() {
    this.setYearLabel()
  }
}
customElements.define('s-year', SYear)
