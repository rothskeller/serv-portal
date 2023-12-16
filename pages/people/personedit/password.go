package personedit

import (
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const passwordPersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | person.FPassword | person.FBadLoginCount | person.FBadLoginTime | person.FPWResetToken | person.FPWResetTime | auth.StrongPasswordPersonFields

// HandlePassword handles requests for /people/$id/edpassword.
func HandlePassword(r *request.Request, idstr string) {
	var (
		user             *person.Person
		p                *person.Person
		oldPassword      string
		oldPasswordError string
		newPassword      string
		newPasswordError string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), passwordPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() != p.ID() && !user.IsWebmaster() {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		if !user.IsWebmaster() {
			oldPassword, oldPasswordError = readOldPassword(r, p)
		}
		newPassword, newPasswordError = readNewPassword(r, user, p)
		if oldPasswordError == "" && newPasswordError == "" {
			auth.SetPassword(r, p, newPassword)
			personview.Render(r, user, p, person.ViewFull, "password")
			return
		}
	}
	r.HTMLNoCache()
	if oldPasswordError != "" || newPasswordError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col personeditPassword' method=POST up-main up-layer=parent up-target=.personviewPassword")
	form.E("div class='formTitle formTitle-primary'").R(r.LangString("Change Password", "Cambiar de contraseña"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if !user.IsWebmaster() {
		emitOldPassword(r, form, oldPassword, oldPasswordError, oldPasswordError != "" || newPasswordError == "")
	}
	emitNewPassword(r, form, user, p, newPassword)
	emitButtons(r, form)
}

func readOldPassword(r *request.Request, p *person.Person) (oldPassword, oldPasswordError string) {
	if oldPassword = r.FormValue("oldpwd"); oldPassword == "" {
		oldPasswordError = r.LangString("Please specify your old password.", "Por favor ingrese su contraseña anterior.")
	} else if !auth.CheckPassword(r, p, oldPassword) {
		oldPasswordError = r.LangString("This is not the correct old password.", "Esta no es la contraseña anterior correcta.")
	}
	return
}

func emitOldPassword(r *request.Request, form *htmlb.Element, oldPassword, oldPasswordError string, focus bool) {
	row := form.E("div class=formRow")
	row.E("label for=personeditPasswordOld").R(r.LangString("Old Password", "Contraseña anterior"))
	row.E("input type=password id=personeditPasswordOld name=oldpwd autocomplete=current-password value=%s", oldPassword,
		focus, "autofocus")
	if oldPasswordError != "" {
		row.E("div class=formError>%s", oldPasswordError)
	}
}

func readNewPassword(r *request.Request, user, p *person.Person) (newPassword, newPasswordError string) {
	if newPassword = r.FormValue("newpwd"); newPassword == "" {
		newPasswordError = r.LangString("Please specify a valid new password.", "Por favor ingrese una nueva contraseña válida.")
	} else if !user.IsWebmaster() && !auth.StrongPassword(p, newPassword) {
		newPasswordError = r.LangString("The new password is too weak.", "La nueva contraseña es demasiado débil.")
	}
	return
}

func emitNewPassword(r *request.Request, form *htmlb.Element, user, p *person.Person, newPassword string) {
	row := form.E("div class=formRow")
	row.E("label for=personeditPasswordNew").R(r.LangString("New Password", "Contraseña nueva"))
	row.E("div class=formInput-2col").E("s-password id=personeditPasswordNew name=newpwd hints=%s value=%s",
		strings.Join(auth.StrongPasswordHints(p), ","), newPassword,
		user.IsWebmaster(), "override")
}
