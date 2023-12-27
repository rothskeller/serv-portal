package classes

import (
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleRegister handles /classes/$id/register requests.
func HandleRegister(r *request.Request, cid class.ID) {
	const personFields = person.FID | person.FInformalName | person.FSortName | person.FEmail | person.FEmail2 | person.FCellPhone
	var (
		user        *person.Person
		c           *class.Class
		langfield   class.Fields
		regs        []*classreg.ClassReg
		uregs       []*classreg.Updater
		askReferral bool
	)
	if user = auth.SessionUser(r, personFields, false); user == nil {
		handleRegisterNotLoggedIn(r, cid)
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if r.Language == "es" {
		langfield = class.FEsDesc
	} else {
		langfield = class.FEnDesc
	}
	if c = class.WithID(r, cid, class.FLimit|class.FReferrals|class.FStart|class.FType|langfield); c == nil {
		errpage.NotFound(r, user)
		return
	}
	classreg.AllForClass(r, c.ID(), classreg.UpdaterFields, func(cr *classreg.ClassReg) {
		if cr.RegisteredBy() != user.ID() {
			return
		}
		regs = append(regs, cr.Clone())
		uregs = append(uregs, cr.Updater(r, c, nil, user))
	})
	if len(regs) == 0 {
		uregs = append(uregs, &classreg.Updater{
			Class:        c,
			Person:       user,
			RegisteredBy: user,
			FirstName:    personFirstName(user),
			LastName:     personLastName(user),
			Email:        personEmail(user),
			CellPhone:    user.CellPhone(),
		})
		regs = append(regs, nil)
		askReferral = r.Method == http.MethodGet
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	for i := 0; i < len(regs); i++ {
		emitRow(form, regs[i])
	}
	emitRow(form, nil)
	tmpl := form.E("template id=classregTemplate")
	emitRow(tmpl, nil)
	emitButtons(form)
}

func handleRegisterNotLoggedIn(r *request.Request, cid class.ID) {
	panic("not implemented")
}

func personFirstName(p *person.Person) string {
	parts := strings.SplitN(p.SortName(), ",", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func personLastName(p *person.Person) string {
	parts := strings.Split(p.SortName(), ",")
	return parts[0]
}

func personEmail(p *person.Person) string {
	if p.Email() != "" {
		return p.Email()
	}
	return p.Email2()
}
