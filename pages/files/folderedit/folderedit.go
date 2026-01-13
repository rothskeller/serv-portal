package folderedit

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/files"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /folderedit/${id} requests.  ${id} may be a folder ID, or it
// may be the word "NEW".  In the latter case, a URL parameter "parent" must
// give the folder ID of the parent of the folder to be created.
func Handle(r *request.Request, idstr string) {
	var (
		user      *person.Person
		f         *folder.Folder
		uf        *folder.Updater
		nameError string
		validate  []string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if idstr == "NEW" {
		uf = new(folder.Updater)
		if uf.Parent = folder.WithID(r, folder.ID(util.ParseID(r.FormValue("parent"))), files.FolderFields|folder.FParent); uf.Parent == nil {
			errpage.NotFound(r, user)
			return
		}
		if !user.HasPrivLevel(uf.Parent.Editor()) {
			errpage.Forbidden(r, user)
			return
		}
		uf.ViewOrg, uf.ViewPriv = uf.Parent.Viewer()
		uf.EditOrg, uf.EditPriv = uf.Parent.Editor()
	} else {
		if f = folder.WithID(r, folder.ID(util.ParseID(idstr)), files.FolderFields|folder.FParent); f == nil {
			errpage.NotFound(r, user)
			return
		}
		if !user.HasPrivLevel(f.Editor()) {
			errpage.Forbidden(r, user)
			return
		}
		if r.FormValue("delete") != "" {
			if document.ExistInFolder(r, f.ID()) || folder.ExistsWithParent(r, f.ID()) {
				errpage.Forbidden(r, user)
				return
			}
			r.Transaction(func() {
				f.Delete(r)
			})
			f = folder.WithID(r, f.Parent(), files.FolderFields|folder.FParent)
			files.GetFolder(r, user, f.FolderPath(r, files.FolderFields), 0, nil)
			return
		}
		uf = f.Updater(r, nil)
	}
	validate = strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		nameError = readName(r, uf)
		readViewer(r, user, uf)
		readEditor(r, user, uf)
		if len(validate) == 0 && nameError == "" {
			r.Transaction(func() {
				if uf.ID == 0 {
					f = folder.Create(r, uf)
				} else {
					f.Update(r, uf)
				}
			})
			files.GetFolder(r, user, f.FolderPath(r, files.FolderFields), f.ID(), nil)
			return
		}
	}
	r.HTMLNoCache()
	if nameError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	if uf.ID == 0 {
		form.E("div class='formTitle formTitle-primary'>Add Folder")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit Folder")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "name") {
		emitName(form, uf, nameError)
	}
	if len(validate) == 0 {
		emitViewer(form, user, uf)
		emitEditor(form, user, uf)
		emitButtons(form, uf.ID != 0 && !folder.ExistsWithParent(r, uf.ID) && !document.ExistInFolder(r, uf.ID))
	}
}

func readName(r *request.Request, uf *folder.Updater) string {
	if uf.Name = strings.TrimSpace(r.FormValue("name")); uf.Name == "" {
		return "The folder name is required."
	}
	if uf.Name[0] == '.' || strings.ContainsAny(uf.Name, "/:") {
		return fmt.Sprintf("%q is not a valid name.  Names may not start with a period, and may not contain slashes or colons.", uf.Name)
	}
	files.MakeURLName(uf)
	if uf.URLName == "" {
		return fmt.Sprintf("%q is not a valid name.  Names must contain at least one alphanumeric character.", uf.Name)
	}
	if uf.DuplicateURLName(r) {
		return fmt.Sprintf("The name %q conflicts with that of another folder under the same parent.", uf.Name)
	}
	return ""
}

func emitName(form *htmlb.Element, uf *folder.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=foldereditName>Folder Name")
	row.E("input id=foldereditName name=name s-validate autofocus value=%s", uf.Name)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

var possibleViewers = []struct {
	org   enum.Org
	priv  enum.PrivLevel
	label string
}{
	{0, 0, "General Public (no login required)"},
	{enum.OrgCERTT, enum.PrivStudent, "CERT Basic Training Students"},
	{enum.OrgCERTD, enum.PrivMember, "CERT Deployment Team"},
	{enum.OrgCERTT, enum.PrivMember, "CERT Training Committee / Instructors"},
	{enum.OrgListos, enum.PrivMember, "Listos Team"},
	{enum.OrgSARES, enum.PrivMember, "SARES Members"},
	{enum.OrgSNAP, enum.PrivMember, "SNAP Team"},
	{0, enum.PrivMember, "SERV Volunteers (any organization)"},
	{enum.OrgAdmin, enum.PrivMember, "SERV Leads"},
	{enum.OrgAdmin, enum.PrivLeader, "OES Staff"},
	{0, enum.PrivMaster, "Webmaster"},
}

func readViewer(r *request.Request, user *person.Person, uf *folder.Updater) {
	vstr := r.FormValue("viewer")
	for _, pv := range possibleViewers {
		if vstr == fmt.Sprintf("%d.%d", pv.org, pv.priv) {
			if !user.HasPrivLevel(pv.org, pv.priv) {
				return
			}
			if pv.priv == 0 && !user.IsAdminLeader() {
				return
			}
			uf.ViewOrg, uf.ViewPriv = pv.org, pv.priv
			return
		}
	}
}

func emitViewer(form *htmlb.Element, user *person.Person, uf *folder.Updater) {
	row := form.E("div class=formRow")
	row.E("label>Folder Viewers")
	input := row.E("div class='formInput-2col foldereditPrivs'")
	for _, pv := range possibleViewers {
		if !user.HasPrivLevel(pv.org, pv.priv) {
			continue
		}
		if pv.priv == 0 && !user.IsAdminLeader() {
			continue
		}
		input.E("s-radio name=viewer value=%d.%d label=%s", pv.org, pv.priv, pv.label,
			uf.ViewOrg == pv.org && uf.ViewPriv == pv.priv, "checked")
	}
}

var possibleEditors = []struct {
	org   enum.Org
	priv  enum.PrivLevel
	label string
}{
	{enum.OrgCERTD, enum.PrivMember, "CERT Deployment Team"},
	{enum.OrgCERTD, enum.PrivLeader, "CERT Deployment Team Leads"},
	{enum.OrgCERTT, enum.PrivMember, "CERT Training Committee / Instructors"},
	{enum.OrgCERTT, enum.PrivLeader, "CERT Training Leads"},
	{enum.OrgListos, enum.PrivMember, "Listos Team"},
	{enum.OrgListos, enum.PrivLeader, "Listos Team Leads"},
	{enum.OrgSARES, enum.PrivMember, "SARES Members"},
	{enum.OrgSARES, enum.PrivLeader, "SARES EC / AECs"},
	{enum.OrgSNAP, enum.PrivMember, "SNAP Team"},
	{enum.OrgSNAP, enum.PrivLeader, "SNAP Team Leads"},
	{enum.OrgAdmin, enum.PrivMember, "SERV Leads"},
	{enum.OrgAdmin, enum.PrivLeader, "OES Staff"},
	{0, enum.PrivMaster, "Webmaster"},
}

func readEditor(r *request.Request, user *person.Person, uf *folder.Updater) {
	vstr := r.FormValue("editor")
	for _, pe := range possibleEditors {
		if vstr == fmt.Sprintf("%d.%d", pe.org, pe.priv) {
			if !user.HasPrivLevel(pe.org, pe.priv) {
				return
			}
			uf.EditOrg, uf.EditPriv = pe.org, pe.priv
			return
		}
	}
}

func emitEditor(form *htmlb.Element, user *person.Person, uf *folder.Updater) {
	row := form.E("div class=formRow")
	row.E("label>Folder Editors")
	input := row.E("div class='formInput-2col foldereditPrivs'")
	for _, pe := range possibleEditors {
		if !user.HasPrivLevel(pe.org, pe.priv) {
			continue
		}
		input.E("s-radio name=editor value=%d.%d label=%s", pe.org, pe.priv, pe.label,
			uf.EditOrg == pe.org && uf.EditPriv == pe.priv, "checked")
	}
}

func emitButtons(form *htmlb.Element, canDelete bool) {
	buttons := form.E("div class=formButtons")
	if canDelete {
		buttons.E("div class=formButtonSpace")
	}
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	if canDelete {
		// This button comes last in the tree order so that it is not
		// the default.  But it comes first in the visual order because
		// of the formButton-beforeAll class.
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
}
