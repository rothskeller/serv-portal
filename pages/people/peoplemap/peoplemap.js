; (function () {
  let loadp, mapp, mapresolve, map, markers = []
  function loadMapScript() {
    if (!loadp) loadp = new Promise(resolve => {
      const tag = document.createElement('script')
      tag.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyCi9J9RDZh5ouo3zk23yDmtY5Pp-NNBsBo'
      document.head.appendChild(tag)
      tag.addEventListener('load', () => { resolve() })
    })
    return loadp
  }
  up.compiler('#peoplemapCanvas', () => {
    return () => { map = null }
  })
  up.compiler('.peoplemapData', async (elm, people) => {
    await loadMapScript()
    if (!map) {
      const mapOptions = {
        center: new google.maps.LatLng(37.3801648, -122.032706),
        zoom: 13,
      }
      map = new google.maps.Map(up.element.get('#peoplemapCanvas'), mapOptions)
      const districts = up.element.jsonAttr(up.element.get('.peoplemapDistricts'), 'up-data')
      districts.forEach(dist => {
        new google.maps.Polygon({
          map,
          paths: dist.points.map(p => ({ lng: p[0], lat: p[1] })),
          strokeWeight: 0,
          fillColor: dist.color,
          fillOpacity: 0.36,
        })
      })
    }
    markers.forEach(m => { m.setMap(null) })
    markers = people.map(p => new google.maps.Marker({ map, title: p.name, position: p }))
  })
})()
up.on('change', '.peoplemapForm *', (evt, elm) => {
  up.submit(elm.form, { target: '.peoplemapData', navigate: false })
})
