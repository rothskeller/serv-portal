; (function () {
  (g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
    key: "AIzaSyBPhAKNiwzI4ETL1pw0Nd-I2df1A-Rnp9g",
    v: "weekly",
  });
  let map, markers = []
  up.compiler('#peoplemapCanvas', () => {
    return () => { map = null }
  })
  up.compiler('.peoplemapData', async (elm, people) => {
    await google.maps.importLibrary('maps')
    const { AdvancedMarkerElement } = await google.maps.importLibrary('marker')
    if (!map) {
      const mapOptions = {
        center: new google.maps.LatLng(37.3801648, -122.032706),
        zoom: 13,
        mapId: 'e6a741577eadbdd028704e63',
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
    markers = people.map(p => {
      const box = document.createElement('div')
      box.className = 'peoplemapMarker'
      box.textContent = p.name
      return new AdvancedMarkerElement({ map, content: box, position: p })
    })
  })
})()
up.on('change', '.peoplemapForm *', (evt, elm) => {
  up.submit(elm.form, { target: '.peoplemapData', navigate: false })
})
