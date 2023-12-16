package errpage

import (
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// NotFound sends a 404 Not Found error page.
func NotFound(r *request.Request, user *person.Person) {
	ui.Page(r, user, ui.PageOpts{
		StatusCode: http.StatusNotFound,
	}, func(main *htmlb.Element) {
		main.A("class=errpage")
		main.E("h1").R(r.LangString("No Such Page", "No existe esa página"))
		main.E("p").R(r.LangString(`Sorry, the page you asked for
doesn’t exist.  But we have plenty of other good ones!  You can
<a href="javascript:history.back()">go&nbsp;back</a>
to where you were, or return to <a href="/">the&nbsp;home&nbsp;page</a>.
Look around; you’re sure to find a page you like.`,
			`Lo sentimos, la página que solicitó no existe.
¡Pero tenemos muchas otras buenas! Puede
<a href="javascript:history.back()">volver</a>
a donde estaba o regrese a la <a href="/">la página de inicio</a>.
Mire alrededor; Seguro que encontrará una página que le gusta.`))
	})
}

// Forbidden sends a 403 Forbidden error page.
func Forbidden(r *request.Request, user *person.Person) {
	r.Problems().Add("insufficient privilege")
	if !strings.Contains(r.Request.Header.Get("Accept"), "text/html") {
		http.Error(r, "403 Forbidden", http.StatusForbidden)
		return
	}
	ui.Page(r, user, ui.PageOpts{StatusCode: http.StatusForbidden}, func(main *htmlb.Element) {
		main.A("class=errpage")
		main.E("h1").R(r.LangString("Forbidden", "Prohibido"))
		main.E("p").R(r.LangString(`Sorry, but your account doesn’t have
permissions for the operation you requested.  If you think you should have
permissions, contact
<a href="mailto:admin@sunnyvaleserv.org">admin@SunnyvaleSERV.org</a> for
assistance.`,
			`Lo sentimos, pero su cuenta no tiene permisos para la
operación que usted solicitó.  Si cree que debería tener permisos, póngase en
contacto con
<a href="mailto:admin@sunnyvaleserv.org">admin@SunnyvaleSERV.org</a> para
asistencia.`))
	})
}

// ServerError sends a 500 Internal Server Error error page.
func ServerError(r *request.Request, user *person.Person) {
	ui.Page(r, user, ui.PageOpts{
		StatusCode: http.StatusInternalServerError,
	}, func(main *htmlb.Element) {
		main.A("class=errpage")
		main.E("h1").R(r.LangString("Web Site Error", "Error del sitio web"))
		main.E("p").R(r.LangString(`We’re sorry, but this web site isn’t
working correctly right now.  This problem has been reported to the site
administrator.  We’ll get it fixed as soon as possible.`,
			`Lo sentimos, pero este sitio web no funciona
correctamente en este momento.  Este problema ha sido informado al administrador
del sitio.  Lo solucionaremos lo antes posible.`))
	})
}
