package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SNAPPage handles GET /sares requests.
func SNAPPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: "SNAP"}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<p><b>Sunnyvale Neighborhoods Actively Prepare (SNAP)</b>
<p>The SNAP program (“Sunnyvale Neighborhoods Actively Prepare”) is our
neighborhood disaster preparedness program.  While our
<a href="/listos" up-target=main>Listos</a> program teaches preparedness for
individuals and families, SNAP prepares neighbors to support each other in a
disaster.  When someone is willing to host a preparedness event for their
neighborhood (ideally 15–20 homes), we can help facilitate the event and raise
the preparedness level of the whole group.
<p>To do this, we make use of the “Map Your Neighborhood” program provided by
the Washington State Emergency Management Division.  We help the neighbors
build a map and a common understanding of the resources and skills available in
their neighborhood in a disaster, and any special challenges or people with
particular needs.  In the process, the neighbors get to know each other better
and a more prepared to face a disaster together.
<p>For more information about SNAP, or to arrange an event for your
neighborhood, write to
<a href="mailto:snap@sunnyvale.ca.gov" target="_blank">snap@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
		} else {
			main.R(`<p><b>Vecindarios de Sunnyvale Se Preparan Activamente</b>
<p>El programa SNAP (“Vecindarios de Sunnyvale se preparan activamente”, por sus
siglas en inglés) es nuestro programa vecinal de preparación ante desastres.
Mientras que nuestro programa
<a href="/listos" up-target=main>Listos</a> enseña preparación a individuos y
familias, SNAP prepara a los vecinos para que se apoyen mutuamente en caso de
desastre.  Cuando alguien está dispuesto a organizar un evento de preparación
para su vecindario (idealmente 15–20 hogares), podemos ayudar a facilitar el
evento y elevar el nivel de preparación de todo el grupo.
<p>Para ello, utilizamos el programa “Map Your Neighborhood” (“Mapear su
vecindario”) de la División de Gestión de Emergencias del Estado de Washington.
Ayudamos a los vecinos a crear un mapa y un conocimiento común de los recursos y
capacidades disponibles en su barrio en caso de desastre, así como de los
problemas especiales o las personas con necesidades particulares.  De este modo,
los vecinos se conocen mejor y están mejor preparados para afrontar juntos un
desastre.
<p>Para más información sobre SNAP, o para organizar un evento para su
vecindario, escriba a
<a href="mailto:snap@sunnyvale.ca.gov" target="_blank">snap@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regresar</button></div>`)
		}
	})
}
