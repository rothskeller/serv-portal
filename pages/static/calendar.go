package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SubscribeCalendarPage handles GET /subscribe-calendar requests.
func SubscribeCalendarPage(r *request.Request) {
	var user *person.Person

	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	ui.Page(r, user, ui.PageOpts{Title: r.LangString("SERV Calendar Subscription", "Suscripción al calendario de SERV")}, func(main *htmlb.Element) {
		main.A("class=static").R(r.LangString(`<p>You can subscribe to the SERV calendar so that SERV events will automatically appear in the calendar app on your
  phone, or in your desktop calendar software. Please see the instructions for your phone or software below.
<h1>iPhone or iPad Calendar App</h1>
<ol>
  <li>Open the Settings app.
  <li>Go to “Calendar”.
  <li>Go to “Accounts”.
  <li>Go to “Add Account”.
  <li>Tap on “Other”.
  <li>Tap on “Add Subscribed Calendar”.
  <li>In the “Server” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Tap “Next”.
  <li>Optional: change the “Description” field to a name that’s meaningful to you, such as “SERV Calendar”.
  <li>Tap “Save”.
</ol>
<h1>Google Calendar (including Android Phones)</h1>
<ol>
  <li>In a web browser, go to Google Calendar (<code>https://calendar.google.com</code>). Log in if necessary.
  <li>In the left sidebar, click the large “+” sign next to “Other Calendars”.
  <li>Click “From URL”.
  <li>In the “URL of calendar” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Add calendar”.
</ol>
<h1>Microsoft Outlook</h1>
<ol>
  <li>Open Microsoft Outlook.
  <li>Go to the calendar page.
  <li>In the Home ribbon, click on “Open Calendar”, then “From Internet”.
  <li>Enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Yes”.
  <li>In the left sidebar, under “Other Calendars”, right-click on “Untitled” and choose “Rename Calendar”.
  <li>Give the calendar a name meaningful to you, such as “SERV Calendar”.
</ol>
<h1>Mac Calendar App</h1>
<ol>
  <li>Open the Calendar app.
  <li>From the menu, choose File → New Calendar Subscription.
  <li>Enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Subscribe”.
  <li>Set the options to suit your preferences and click “OK”.
</ol>
<h1>Other Software</h1>
<p>Most calendar software has the ability to subscribe to Internet calendars. Consult the documentation for your
  software to find out how. The address of the SERV calendar is <code>https://sunnyvaleserv.org/calendar.ics</code>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`,

			`<p>Puede suscribirse al calendario de SERV para que los eventos de SERV
  aparezcan automáticamente en la aplicación de calendario de su teléfono o en
  el software de calendario de su computadora.  Consulte las instrucciones para
  su teléfono o software a continuación.
<h1>Aplicación de calendario de iPhone o iPad</h1>
<ol>
  <li>Abra la aplicación Configuración.
  <li>Vaya a “Calendario”.
  <li>Vaya a “Cuentas”.
  <li>Vaya a “Agregar cuenta”.
  <li>Toque “Otro”.
  <li>Toque “Agregar calendario suscrito”.
  <li>En el campo "Servidor", ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Toque “Siguiente”.
  <li>Opcional: cambie el campo “Descripción” por un nombre que sea significativo para usted, como “Calendario de SERV”.
  <li>Toque “Guardar”.
</ol>
<h1>Google Calendar (incluidos los teléfonos Android)</h1>
<ol>
  <li>En un navegador web, vaya a Google Calendar (<code>https://calendar.google.com</code>). Inicie sesión si es necesario.
  <li>En la barra lateral izquierda, haga clic en el signo grande “+” junto a “Otros calendarios”.
  <li>Haga clic en “Desde URL”.
  <li>En el campo “URL del calendario”, ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Haga clic en “Agregar calendario”.
</ol>
<h1>Microsoft Outlook</h1>
<ol>
  <li>Abra Microsoft Outlook.
  <li>Vaya a la página del calendario.
  <li>En la cinta Inicio, haga clic en “Abrir calendario” y luego en “Desde Internet”.
  <li>Ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Haga clic en "Sí".
  <li>En la barra lateral izquierda, en "Otros calendarios", haga clic derecho en "Untitled" y seleccione "Cambiar nombre de calendario".
  <li>Asigne al calendario un nombre significativo para usted, como "Calendario de SERV".
</ol>
<h1>Aplicación de calendario para Mac</h1>
<ol>
  <li>Abra la aplicación Calendario.
  <li>En el menú, elija Archivo → Nueva suscripción a calendario.
  <li>Ingrese <código>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Haga clic en "Suscribir".
  <li>Ajuste las opciones a sus preferencias y haga clic en "Aceptar".
</ol>
<h1>Otro software</h1>
<p>La mayoría del software de calendario tiene la capacidad de suscribirse a
  calendarios de Internet. Consulta la documentación de tu software para
  descubrir cómo. La dirección del calendario de SERV es
  <code>https://sunnyvaleserv.org/calendar.ics</code>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regrese</button></div>`))
	})
}
