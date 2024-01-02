package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SARESPage handles GET /sares requests.
func SARESPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Sunnyvale ARES")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<p><b>Sunnyvale Amateur Radio Emergency Service</b>
<p>The Sunnyvale Amateur Radio Emergency Service (SARES) is the local chapter of
the nationwide Amateur Radio Emergency Service operated by the Amateur Radio
Relay League (ARRL).  During times of emergency, it also operates as a local
branch of the federal Radio Amateur Civil Emergency Service (RACES).  SARES
provides emergency communications services, usually but not always using amateur
radio, when regular communications methods are unavailable or saturated.
<p>In a disaster, telephones and the Internet will likely be down.  Or, if they
are working, they will be unable to keep up with demand.  Radio communications
serve as an effective backup because they do not rely on massive, fragile
infrastructure.  SARES operators can provide essential emergency communications
when no other methods are working.  Outside of emergencies, SARES operators also
provide communications services at public events to keep themselves in practice.
<p>Membership in SARES requires a current FCC amateur radio license.  If you are
interested in emergency communications but do not have a license, you should
pursue getting your license first before trying to join SARES.  However, the
members of SARES will be happy to connect you with resources to help you get
your license.
<p>For more information about SARES or amateur radio, write to
<a href="mailto:sares@sunnyvale.ca.gov" target="_blank">sares@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
		} else {
			main.R(`<p><b>Servicio de Emergencias de Radioaficionados de Sunnyvale</b>
<p>El Servicio de Emergencia de Radioaficionados de Sunnyvale (SARES, por siglas
en inglés) es el capítulo local del Servicio de Emergencia de Radioaficionados
(ARES) a nivel nacional operado por la Liga de Radioaficionados (ARRL).  En
situaciones de emergencia, también funciona como una rama local del Servicio de
Emergencia Civil de Radioaficionados (RACES).  El SARES proporciona servicios de
comunicaciones en emergencias, normalmente pero no siempre utilizando
radioafición, cuando los métodos de comunicación habituales no están disponibles
o están saturados.
<p>En caso de desastre, es probable que los teléfonos e Internet no funcionen.
O, si funcionan, serán incapaces de satisfacer la demanda.  Las comunicaciones
por radio son un medio de reserva eficaz porque no dependen de infraestructuras
masivas y frágiles.  Los operadores de SARES pueden proporcionar comunicaciones
de emergencia esenciales cuando no funcionan otros métodos.  Fuera de las
emergencias, los operadores de SARES también prestan servicios de comunicaciones
en eventos públicos para mantenerse en activo.
<p>Para ser miembro de SARES se requiere una licencia de radioaficionado de la
FCC en vigor.  Si se interesan las comunicaciones de emergencia pero no tiene
licencia, primero debe obtenerla antes de intentar unirte a SARES.  No obstante,
los miembros de SARES estarán encantados de ponerle en contacto con recursos que
le ayuden a obtener su licencia.
<p>Para más información sobre SARES or radioafición, escriba a
<a href="mailto:sares@sunnyvale.ca.gov" target="_blank">sares@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regresar</button></div>`)
		}
	})
}
