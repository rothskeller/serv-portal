class SMonth extends HTMLElement {
  static get observedAttributes() { return ['value'] };
  constructor() {
    super()
    this._left = this._month = this._right = this._dd = this._year = null
  }
  connectedCallback() {
    this.classList.add('s-month')
    this._left = document.createElement('s-icon')
    this._left.className = 's-month-arrow'
    this._left.setAttribute('icon', 'chevron-left-thin')
    this._left.addEventListener('click', this.prevMonth.bind(this))
    this.appendChild(this._left)
    this._month = document.createElement('div')
    this._month.className = 's-month-month'
    this.setMonthLabel()
    this._month.addEventListener('click', this.toggleDD.bind(this))
    this.appendChild(this._month)
    this._right = document.createElement('s-icon')
    this._right.className = 's-month-arrow'
    this._right.setAttribute('icon', 'chevron-right-thin')
    this._right.addEventListener('click', this.nextMonth.bind(this))
    this.appendChild(this._right)
  }
  disconnectedCallback() {
    if (this._left) this.removeChild(this._left)
    if (this._month) this.removeChild(this._month)
    if (this._right) this.removeChild(this._right)
    if (this._dd) this.removeChild(this._dd)
    this._left = this._month = this._right = this._dd = this._year = null
  }
  setMonthLabel() {
    const value = this.getAttribute('value')
    if (!value || !this._month) return
    const names = document.documentElement.lang === 'es' ? SMonth.esMonthnames : SMonth.monthnames
    this._month.textContent = `${names[value.substr(5, 2)]} ${value.substr(0, 4)}`
  }
  toggleDD() {
    if (this._dd) {
      this._year.textContent = this.getAttribute('value').substr(0, 4)
      this._dd.style.display = this._dd.style.display ? '' : 'none'
      return
    }
    this._dd = document.createElement('div')
    this._dd.className = 's-month-dd'
    const top = document.createElement('div')
    top.className = 's-month-dd-top'
    const left = document.createElement('s-icon')
    left.className = 's-month-arrow'
    left.setAttribute('icon', 'chevron-left')
    left.addEventListener('click', this.prevYear.bind(this))
    top.appendChild(left)
    this._year = document.createElement('div')
    this._year.className = 's-month-dd-year'
    this._year.textContent = this.getAttribute('value').substr(0, 4)
    top.appendChild(this._year)
    const right = document.createElement('s-icon')
    right.className = 's-month-arrow'
    right.setAttribute('icon', 'chevron-right')
    right.addEventListener('click', this.nextYear.bind(this))
    top.appendChild(right)
    this._dd.appendChild(top)
    const abbrs = document.documentElement.lang === 'es' ? SMonth.esMonthabbrs : SMonth.monthabbrs
    Object.keys(abbrs).sort().forEach(month => {
      const btn = document.createElement('button')
      btn.className = 'sbtn sbtn-secondary s-month-dd-month'
      btn.setAttribute('data-month', month)
      btn.textContent = abbrs[month]
      btn.addEventListener('click', this.clickMonth.bind(this))
      this._dd.appendChild(btn)
    })
    this.appendChild(this._dd)
  }
  prevMonth() {
    let value = this.getAttribute('value')
    let month = value.substr(5, 2)
    if (month === '01')
      value = `${(parseInt(value.substr(0, 4)) - 1)}-12`
    else if (month === '12' || month === '11')
      value = `${value.substr(0, 5)}${parseInt(month) - 1}`
    else
      value = `${value.substr(0, 5)}0${parseInt(month) - 1}`
    this.setAttribute('value', value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
    this.setMonthLabel()
  }
  nextMonth() {
    let value = this.getAttribute('value')
    let month = value.substr(5, 2)
    if (month === '12')
      value = `${(parseInt(value.substr(0, 4)) + 1)}-01`
    else if (month === '09' || month === '10' || month === '11')
      value = `${value.substr(0, 5)}${parseInt(month) + 1}`
    else
      value = `${value.substr(0, 5)}0${parseInt(month) + 1}`
    this.setAttribute('value', value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
    this.setMonthLabel()
  }
  prevYear() {
    this._year.textContent = parseInt(this._year.textContent) - 1
  }
  nextYear() {
    this._year.textContent = parseInt(this._year.textContent) + 1
  }
  clickMonth(evt) {
    const month = evt.target.getAttribute('data-month')
    const value = `${this._year.textContent}-${month}`
    this.setAttribute('value', value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
    this.setMonthLabel()
    this._dd.style.display = 'none'
  }
  attributeChangedCallback() {
    this.setMonthLabel()
  }
}
SMonth.monthnames = {
  '01': 'January',
  '02': 'February',
  '03': 'March',
  '04': 'April',
  '05': 'May',
  '06': 'June',
  '07': 'July',
  '08': 'August',
  '09': 'September',
  '10': 'October',
  '11': 'November',
  '12': 'December',
}
SMonth.monthabbrs = {
  '01': 'Jan',
  '02': 'Feb',
  '03': 'Mar',
  '04': 'Apr',
  '05': 'May',
  '06': 'Jun',
  '07': 'Jul',
  '08': 'Aug',
  '09': 'Sep',
  '10': 'Oct',
  '11': 'Nov',
  '12': 'Dec',
}
SMonth.esMonthnames = {
  '01': 'Enero',
  '02': 'Febrero',
  '03': 'Marzo',
  '04': 'Abril',
  '05': 'Mayo',
  '06': 'Junio',
  '07': 'Julio',
  '08': 'Agosto',
  '09': 'Septiembre',
  '10': 'Octubre',
  '11': 'Noviembre',
  '12': 'Diciembre',
}
SMonth.esMonthabbrs = {
  '01': 'Ene',
  '02': 'Feb',
  '03': 'Mar',
  '04': 'Abr',
  '05': 'May',
  '06': 'Jun',
  '07': 'Jul',
  '08': 'Ago',
  '09': 'Sep',
  '10': 'Oct',
  '11': 'Nov',
  '12': 'Dec',
}
customElements.define('s-month', SMonth)
