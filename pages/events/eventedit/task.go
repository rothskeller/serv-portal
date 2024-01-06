package eventedit

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/roleselect"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/store/taskrole"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type taskEditor struct {
	user       *person.Person
	e          *event.Event
	t          *task.Task
	ut         *task.Updater
	roles      []*role.Role
	copyShifts task.ID
	canDelete  bool
	hasShifts  bool
	nameError  string
	orgError   string
	hasError   bool
	op         string
	validate   request.ValidationList
}

// HandleTask handles requests for /events/edtask/$id.  $id may be "NEW", in
// which case an eid=$eid parameter must be included in the request.
func HandleTask(r *request.Request, tidstr string) {
	var te *taskEditor

	if te = getTaskEditor(r, tidstr); te == nil {
		return
	}
	if te.op == "delete" {
		te.handleDelete(r)
		return
	}
	if te.op != "get" {
		te.copyShifts = task.ID(util.ParseID(r.FormValue("copyShifts")))
		te.nameError = readTaskName(r, te.ut)
		te.orgError = readOrg(r, te.user, te.ut)
		te.roles = readRoles(r, te.user, te.roles)
		readTaskFlags(r, te.ut)
		readTaskDetails(r, te.ut)
		te.hasError = te.nameError != "" || te.orgError != ""
	}
	if !te.hasError && (te.op == "save" || (te.op == "copy" && te.t == nil)) {
		if te.t == nil {
			te.create(r)
			if te.op == "save" {
				te.e = event.WithID(r, te.e.ID(), eventview.EventFields)
				eventview.Render(r, te.user, te.e, "")
				return
			}
		} else {
			te.update(r)
			if te.op == "save" {
				eventview.Render(r, te.user, te.e, fmt.Sprintf("task%d", te.ut.ID))
				return
			}
		}
	}
	if te.op == "copy" && te.t != nil && te.ut.Name == te.t.Name() {
		// We're copying a task whose name has not been changed.  Make a
		// unique name for it.
		te.makeUniqueName(r)
	}
	if te.op == "copy" {
		te.copyShifts = te.t.ID()
		te.t = nil
		te.ut.ID = 0
	}
	r.HTMLNoCache()
	if te.op != "get" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	te.writeForm(r, html)
}

func getTaskEditor(r *request.Request, tidstr string) (te *taskEditor) {
	const eventFields = event.FID | event.FName | event.FStart | event.FEnd | event.FVenue | event.FFlags
	var eid event.ID

	// Get a valid user.
	te = new(taskEditor)
	if te.user = auth.SessionUser(r, 0, true); te.user == nil {
		return nil
	}
	if !auth.CheckCSRF(r, te.user) {
		return nil
	}
	// Get and check the event and task.
	if tidstr == "NEW" {
		if !te.user.HasPrivLevel(0, enum.PrivLeader) {
			errpage.Forbidden(r, te.user)
			return nil
		}
		eid = event.ID(util.ParseID(r.FormValue("eid")))
	} else {
		if te.t = task.WithID(r, task.ID(util.ParseID(tidstr)), task.UpdaterFields); te.t == nil {
			errpage.NotFound(r, te.user)
			return nil
		}
		if !te.user.HasPrivLevel(te.t.Org(), enum.PrivLeader) {
			errpage.Forbidden(r, te.user)
			return nil
		}
		eid = te.t.Event()
	}
	if te.e = event.WithID(r, eid, eventFields); te.e == nil || te.e.Flags()&event.OtherHours != 0 {
		errpage.NotFound(r, te.user)
		return nil
	}
	// Get an updater.
	if te.t == nil {
		te.ut = &task.Updater{Event: te.e}
	} else {
		te.ut = te.t.Updater(r, te.e)
	}
	te.canDelete = te.t != nil && !shiftperson.TaskHasSignups(r, te.t.ID()) && !taskperson.ExistsForTask(r, te.t.ID()) && task.CountForEvent(r, te.e.ID()) > 1
	te.hasShifts = shift.ExistsForTask(r, te.ut.ID)
	taskrole.Get(r, te.ut.ID, role.FID|role.FName|role.FOrg, func(rl *role.Role) {
		te.roles = append(te.roles, rl.Clone())
	})
	sort.Slice(te.roles, func(i, j int) bool { return te.roles[i].Name() < te.roles[j].Name() })
	// Get the operation.
	if r.Method != http.MethodPost {
		te.op = "get"
	} else if te.validate = r.ValidationList(); te.validate.Enabled() {
		te.op = "validate"
	} else if r.FormValue("delete") != "" && te.canDelete {
		te.op = "delete"
	} else if r.FormValue("copy") != "" {
		te.op = "copy"
	} else {
		te.op = "save"
	}
	return te
}

func (te *taskEditor) writeForm(r *request.Request, html *htmlb.Element) {
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent")
	if te.t == nil {
		form.Attr("up-target=main action=/events/edtask/NEW")
		form.E("input type=hidden name=eid value=%d", te.e.ID())
		if te.copyShifts != 0 {
			form.E("input type=hidden name=copyShifts value=%d", te.copyShifts)
		}
		form.E("div class='formTitle formTitle-primary'>New Task")
	} else {
		form.Attr("up-target=#eventviewTask%d up-fallback=.eventview action=/events/edtask/%d", te.t.ID(), te.t.ID())
		form.E("div class='formTitle formTitle-primary'>Edit Task")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if te.validate.Validating("name") {
		emitTaskName(form, te.ut, te.nameError != "" || !te.hasError, te.nameError)
	}
	if te.validate.Validating("org") {
		emitOrg(form, te.user, te.ut, te.orgError != "", te.orgError)
	}
	if !te.validate.Enabled() {
		emitRoles(r, form, te.user, te.roles)
		emitTaskFlags(form, te.ut, "task")
		emitTaskDetails(form, te.ut)
		emitTaskButtons(form, te.canDelete)
	}
}

func readTaskName(r *request.Request, ut *task.Updater) string {
	if ut.Name = strings.TrimSpace(r.FormValue("name")); ut.Name == "" {
		return "The task name is required."
	}
	if ut.DuplicateName(r) {
		return fmt.Sprintf("Another task on this event has the name %q.", ut.Name)
	}
	return ""
}
func emitTaskName(form *htmlb.Element, ut *task.Updater, focus bool, err string) {
	row := form.E("div id=eventeditTaskNameRow class=formRow")
	row.E("label for=eventeditTaskName>Name")
	row.E("input id=eventeditTaskName name=name s-validate value=%s", ut.Name, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

var orgNames = map[enum.Org]string{
	enum.OrgAdmin:  "SERV Admin",
	enum.OrgCERTD:  "CERT Deployment",
	enum.OrgCERTT:  "CERT Training",
	enum.OrgListos: "Listos",
	enum.OrgSARES:  "SARES",
	enum.OrgSNAP:   "SNAP",
}

func readOrg(r *request.Request, user *person.Person, ut *task.Updater) string {
	var allowed []enum.Org

	for _, org := range enum.AllOrgs {
		if user.HasPrivLevel(org, enum.PrivLeader) {
			allowed = append(allowed, org)
		}
	}
	if len(allowed) == 1 {
		ut.Org = allowed[0]
		return ""
	}
	orgstr := r.FormValue("org")
	if orgstr == "" {
		return "The task organization is required."
	}
	ut.Org = enum.Org(util.ParseID(orgstr))
	if !ut.Org.Valid() {
		return "The task organization is not valid."
	}
	if !slices.Contains(allowed, ut.Org) {
		return "You do not have privilege to schedule tasks for this organization."
	}
	return ""
}
func emitOrg(form *htmlb.Element, user *person.Person, ut *task.Updater, focus bool, err string) {
	var allowed []enum.Org

	for _, org := range enum.AllOrgs {
		if user.HasPrivLevel(org, enum.PrivLeader) {
			allowed = append(allowed, org)
		}
	}
	row := form.E("div id=eventeditOrgRow class=formRow")
	row.E("label for=eventeditOrg>Organization")
	sel := row.E("select id=eventeditOrg name=org s-validate", len(allowed) == 1, "disabled", focus, "autofocus")
	if len(allowed) == 1 {
		sel.Attr("disabled")
		sel.E("option value=%d>%s", allowed[0], orgNames[allowed[0]])
		return
	}
	if ut.Org == 0 {
		sel.E("option value=0 selected>(select organization)")
	}
	for _, org := range allowed {
		sel.E("option value=%d", org, org == ut.Org, "selected").T(orgNames[org])
	}
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readRoles(r *request.Request, user *person.Person, roles []*role.Role) (np []*role.Role) {
	var ids = make(map[role.ID]*role.Role)
	for _, idstr := range strings.Fields(r.FormValue("roles")) {
		if rl := role.WithID(r, role.ID(util.ParseID(idstr)), role.FID|role.FName|role.FOrg); rl != nil {
			if user.HasPrivLevel(rl.Org(), enum.PrivLeader) {
				ids[rl.ID()] = rl
			}
		}
	}
	// Preserve roles that this user can't change.
	for _, rl := range roles {
		if !user.HasPrivLevel(rl.Org(), enum.PrivLeader) {
			ids[rl.ID()] = rl
		}
	}
	for _, rl := range ids {
		np = append(np, rl)
	}
	return np
}
func emitRoles(r *request.Request, form *htmlb.Element, user *person.Person, roles []*role.Role) {
	row := form.E("div id=eventeditRolesRow class=formRow")
	row.E("label for=eventeditRoles0 class=checkLabel>Roles")
	tree := roleselect.MakeRoleTree(r, role.FID|role.FOrg, func(rl *role.Role) bool {
		return user.HasPrivLevel(rl.Org(), enum.PrivLeader)
	})
	var selids []string
	for _, rl := range roles {
		selids = append(selids, strconv.Itoa(int(rl.ID())))
	}
	row.E("div class=formInput").E("s-seltree name=roles value=%s>%s", strings.Join(selids, " "), tree)
}

func readTaskFlags(r *request.Request, ut *task.Updater) {
	ut.Flags &^= task.RecordHours | task.CoveredByDSW | task.SignupsOpen
	if r.FormValue("recordHours") != "" {
		ut.Flags |= task.RecordHours
	}
	if r.FormValue("coveredByDSW") != "" {
		ut.Flags |= task.CoveredByDSW
	}
	if r.FormValue("requiresBGCheck") != "" {
		ut.Flags |= task.RequiresBGCheck
	}
	if r.FormValue("signupsOpen") != "" {
		ut.Flags |= task.SignupsOpen
	}
}
func emitTaskFlags(form *htmlb.Element, ut *task.Updater, item string) {
	row := form.E("div class=formRow")
	row.E("label for=eventeditRecordHours class=checkLabel>Task Flags")
	in := row.E("div class=formInput")
	in.E("input type=checkbox class=s-check id=eventeditRecordHours name=recordHours label=%s", "Record volunteer hours for this "+item,
		ut.Flags&task.RecordHours != 0, "checked")
	in.E("input type=checkbox class=s-check id=eventeditCoveredByDSW name=coveredByDSW label=%s", strings.ToUpper(item[:1])+item[1:]+" is covered by Sunnyvale DSW",
		ut.Flags&task.CoveredByDSW != 0, "checked")
	in.E("input type=checkbox class=s-check id=eventeditRequiresBGCheck name=requiresBGCheck label=%s", strings.ToUpper(item[:1])+item[1:]+" requires background check",
		ut.Flags&task.RequiresBGCheck != 0, "checked")
	in.E("input type=checkbox class=s-check id=eventeditSignupsOpen name=signupsOpen label=%s", "People can sign up for this "+item,
		ut.Flags&task.SignupsOpen != 0, "checked")
}

func readTaskDetails(r *request.Request, ut *task.Updater) {
	ut.Details = htmlSanitizer.Sanitize(strings.TrimSpace(r.FormValue("details")))
}
func emitTaskDetails(form *htmlb.Element, ut *task.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=eventeditTaskDetails>Details")
	row.E("textarea id=eventeditTaskDetails name=details wrap=soft rows=3").T(ut.Details)
	row.E("div class=formHelp>This may contain HTML &lt;a&gt; tags for links, but no other tags.")
}

func emitTaskButtons(form *htmlb.Element, canDelete bool) {
	buttons := form.E("div class=formButtons")
	buttons.E("div class=formButtonSpace")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	// These buttons must appear lexically after Save, even though they
	// appear visually before it, so that Save is the default button when
	// the user presses Enter.  The formButton-beforeAll class implements
	// that.
	if canDelete {
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
	buttons.E("input type=submit name=copy class='sbtn sbtn-secondary formButton-beforeAll' value=Copy")
}

func (te *taskEditor) create(r *request.Request) {
	r.Transaction(func() {
		te.t = task.Create(r, te.ut)
		taskrole.Set(r, te.e, te.t, te.roles, nil)
		if te.copyShifts == 0 {
			return
		}
		if ct := task.WithID(r, te.copyShifts, task.FEvent|task.FOrg); ct == nil || ct.Event() != te.e.ID() || !te.user.HasPrivLevel(ct.Org(), enum.PrivLeader) {
			return
		}
		shift.AllForTask(r, te.copyShifts, shift.FStart|shift.FEnd|shift.FMin|shift.FMax, venue.FID|venue.FName, func(s *shift.Shift, v *venue.Venue) {
			shift.Create(r, &shift.Updater{
				Event: te.e,
				Task:  te.t,
				Start: s.Start(),
				End:   s.End(),
				Venue: v,
				Min:   s.Min(),
				Max:   s.Max(),
			})
		})
	})
}

func (te *taskEditor) update(r *request.Request) {
	r.Transaction(func() {
		te.t.Update(r, te.ut)
		taskrole.Set(r, te.e, te.t, te.roles, nil)
		if te.ut.Flags&task.SignupsOpen != 0 && !te.hasShifts {
			shift.Create(r, &shift.Updater{
				Event: te.e,
				Task:  te.t,
				Start: te.e.Start(),
				End:   te.e.End(),
				Venue: venue.WithID(r, te.e.Venue(), venue.FID|venue.FName),
			})
		}
	})
}

var numsufRE = regexp.MustCompile(` (\d+)$`)

func (te *taskEditor) makeUniqueName(r *request.Request) {
	var seq = 1
	var base = te.ut.Name

	if match := numsufRE.FindStringSubmatch(te.ut.Name); match != nil {
		seq, _ = strconv.Atoi(match[1])
		base = te.ut.Name[:len(te.ut.Name)-len(match[1])]
	} else if n, err := strconv.Atoi(base); err == nil {
		seq = n
		base = ""
	} else {
		base += " "
	}
	for {
		seq++
		te.ut.Name = fmt.Sprintf("%s%d", base, seq)
		if !te.ut.DuplicateName(r) {
			return
		}
	}
}

func (te *taskEditor) handleDelete(r *request.Request) {
	r.Transaction(func() {
		te.t.Delete(r, te.e)
	})
	te.e = event.WithID(r, te.e.ID(), eventview.EventFields)
	eventview.Render(r, te.user, te.e, "")
}
