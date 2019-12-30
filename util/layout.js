window.addEventListener('load', function () {
  if (document.getElementById('layout-menu')) {
    var main = document.getElementById('layout-main');
    document.getElementById('layout-menu-trigger-box').addEventListener('click', function () {
      main.classList.toggle('layout-menu-open');
    });
  }
});
