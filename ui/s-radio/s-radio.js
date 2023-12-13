class SRadio extends HTMLElement {
  static get observedAttributes() { return ['checked', 'disabled', 'form', 'label', 'name', 'value'] };
  constructor() {
    super()
    this._rb = this._lb = null
  }
  connectedCallback() {
    if (!this.id) this.id = `s-radio-${SRadio.nextID++}`
    this.classList.add('s-radio')
    this._rb = document.createElement('input')
    this._rb.type = 'radio'
    this._rb.id = `${this.id}-rb`
    this._rb.className = 's-radio-rb'
    this._rb.checked = this.hasAttribute('checked')
    this._rb.disabled = this.hasAttribute('disabled')
    if (this.hasAttribute('form')) this._rb.setAttribute('form', this.getAttribute('form'))
    if (this.hasAttribute('name')) this._rb.setAttribute('name', this.getAttribute('name'))
    if (this.hasAttribute('value')) this._rb.setAttribute('value', this.getAttribute('value'))
    this.appendChild(this._rb)
    this._lb = document.createElement('label')
    this._lb.id = `${this.id}-lb`
    this._lb.className = 's-radio-lb'
    this._lb.htmlFor = this._rb.id
    if (this.hasAttribute('label')) this._lb.textContent = this.getAttribute('label')
    this.appendChild(this._lb)
  }
  disconnectedCallback() {
    this.removeChild(this._lb)
    this.removeChild(this._rb)
  }
  attributeChangedCallback(name, oldValue, newValue) {
    if (name == 'label') {
      if (this._lb) this._lb.textContent = newValue
    } else if (name == 'checked' || name == 'disabled') {
      if (newValue === null && this._rb) this._rb.removeAttribute(name)
      else if (this._rb) this._rb.setAttribute(name, name)
    } else {
      if (this._rb) this._rb.setAttribute(name, newValue)
    }
  }
  get checked() { return this._rb.checked }
  set checked(t) { this._rb.checked = t }
}
SRadio.nextID = 1
customElements.define('s-radio', SRadio)
