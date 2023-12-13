up.on('input', '#textnewMessage', (evt, elm) => {
  let unicode = false
  let runes = 0
  let chars = 0
  for (const ch of elm.value) {
    runes++
    if ('£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !"#¤%&\'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà'.includes(ch))
      chars++
    else if ('\f\n^{}\\[~]|€'.includes(ch)) chars += 2
    else unicode = true
  }
  const oversize = unicode ? (runes > 70) : (chars > 160)
  elm.classList.toggle('oversize', oversize)
})
