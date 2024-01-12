package personedit

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/login"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

const passwordPersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | person.FPassword | person.FBadLoginCount | person.FBadLoginTime | person.FPWResetToken | person.FPWResetTime | auth.StrongPasswordPersonFields

// HandlePassword handles requests for /people/$id/edpassword.
func HandlePassword(r *request.Request, idstr string) {
	var (
		user         *person.Person
		p            *person.Person
		f            form.Form
		oldPassword  string
		newPassword1 string
		newPassword2 string
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
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
	f.Attrs = "method=POST up-target=.personviewPassword"
	f.Dialog = true
	f.Title = "Password Change"
	f.Buttons = []*form.Button{
		{Label: "Save", OnClick: func() bool {
			auth.SetPassword(r, p, newPassword1)
			personview.Render(r, user, p, person.ViewFull, "password")
			return true
		}},
	}
	if !user.IsWebmaster() {
		f.Rows = append(f.Rows, &oldPasswordRow{
			PasswordRow: form.PasswordRow{
				InputRow: form.InputRow{
					LabeledRow: form.LabeledRow{
						RowID: "personeditPasswordOld",
						Label: "Old Password",
					},
					Name:   "oldpwd",
					ValueP: &oldPassword,
				},
				Autocomplete: "current-password",
			},
			p: p,
		})
	}
	f.Rows = append(f.Rows, &login.NewPasswordPairRow{
		LabeledRow: form.LabeledRow{
			RowID: "personeditPasswordNew",
			Label: "New Password",
		},
		Name:     "newpwd",
		ValueP1:  &newPassword1,
		ValueP2:  &newPassword2,
		Override: user.IsWebmaster(),
		Person:   p,
	})
	f.Handle(r)
}

type oldPasswordRow struct {
	form.PasswordRow
	p *person.Person
}

func (opr *oldPasswordRow) Read(r *request.Request) bool {
	if !opr.PasswordRow.Read(r) {
		return false
	}
	if *opr.ValueP == "" {
		opr.Error = r.Loc("Please specify your old password.")
		return false
	}
	if !auth.CheckPassword(r, opr.p, *opr.ValueP) {
		opr.Error = r.Loc("This is not the correct old password.")
		return false
	}
	return true
}
