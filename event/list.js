window.addEventListener('load', function () {
  var year = document.getElementById('listEvents-year');
  if (year) year.addEventListener('change', function () {
    document.getElementById('listEvents-title').submit();
  });
});
