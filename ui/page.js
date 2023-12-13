up.fragment.config.mainTargets = ['.pageCanvas', ...up.fragment.config.mainTargets]
up.on('click', '#pageMenuTrigger', elm => {
  document.body.classList.toggle('page-menuOpen')
})
up.on('up:link:follow', () => {
  document.body.classList.remove('page-menuOpen')
})
