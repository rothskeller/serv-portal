package classpeople

import (
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /classes/$cid/people requests.
func Handle(r *request.Request, cidstr string) {
	const classFields = class.FType
	var (
		user *person.Person
		c    *class.Class
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), classFields); c == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.IsWebmaster() {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		handlePost(r, c)
	} else {
		handleGet(r, c)
	}
}

func handleGet(r *request.Request, c *class.Class) {
	const classFields = classreg.FCellPhone | classreg.FEmail | classreg.FFirstName | classreg.FID | classreg.FLastName | classreg.FPerson | classreg.FRegisteredBy | classreg.FWaitlist
	const personFields = person.FCellPhone | person.FEmail | person.FEmail2 | person.FHomePhone | person.FID | person.FInformalName | person.FSortName | person.FWorkPhone
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form' method=POST up-main")
	form.E("div class='formTitle formTitle-primary'>Class Person Assignments")
	main := form.E("div class=formRow-3col")
	used := sets.New[person.ID]()
	classreg.AllForClass(r, c.ID(), classreg.FPerson, func(cr *classreg.ClassReg) {
		if cr.Person() != 0 {
			used.Insert(cr.Person())
		}
	})
	classreg.AllForClass(r, c.ID(), classFields, func(cr *classreg.ClassReg) {
		var registrarFound bool
		if cr.Person() != 0 {
			return
		}
		registrant := main.E("div class=classpersonRegistrant>%s %s %s %s", cr.FirstName(), cr.LastName(), cr.Email(), cr.CellPhone())
		if cr.Waitlist() {
			registrant.R(" (waitlist)")
		}
		person.All(r, personFields, func(p *person.Person) {
			if used.Has(p.ID()) {
				return
			}
			lastname, _, _ := strings.Cut(p.SortName(), ",")
			switch {
			case cr.Email() != "" && (strings.EqualFold(cr.Email(), p.Email()) || strings.EqualFold(cr.Email(), p.Email2())):
			case cr.CellPhone() != "" && (cr.CellPhone() == p.CellPhone() || cr.CellPhone() == p.HomePhone() || cr.CellPhone() == p.WorkPhone()):
			case strings.EqualFold(cr.LastName(), lastname):
			default:
				return
			}
			candidate := main.E("div class=classpersonCandidate")
			button := candidate.E("input type=radio name=r%d value=%d", cr.ID(), p.ID())
			if p.ID() == cr.RegisteredBy() {
				button.A("checked")
				registrarFound = true
			}
			candidate.TF(" #%d %s %s %s %s %s %s", p.ID(), p.InformalName(), p.Email(), p.Email2(), p.CellPhone(), p.HomePhone(), p.WorkPhone())
			if p.ID() == cr.RegisteredBy() {
				candidate.R(" (registrar)")
			}
		})
		if !registrarFound {
			main.E("div class=classpersonCandidate").E("input type=radio name=r%d value=NEW").P().T(" Create new person")
		}
		main.E("div class=classpersonCandidate").E("input type=radio name=r%d value=0").P().T(" Skip")
	})
	form.E("div class=formButtons").E("button type=submit class='sbtn sbtn-primary'>Save")
}
func handlePost(r *request.Request, c *class.Class) {}
