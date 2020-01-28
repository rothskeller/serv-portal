package person

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidatePerson validates a Person record, except for its Password field.  It
// enforces all data consistency rules, but does not enforce privileges.
func ValidatePerson(tx *db.Tx, person *model.Person) error {
	var (
		individualHeld map[*model.Role]bool
		people         []*model.Person
		roles          []*model.Role
		roleMap        map[model.RoleID]bool
	)
	if person.InformalName = strings.TrimSpace(person.InformalName); person.InformalName == "" {
		return errors.New("missing informalName")
	}
	if person.FormalName = strings.TrimSpace(person.FormalName); person.FormalName == "" {
		return errors.New("missing formalName")
	}
	if person.SortName = strings.TrimSpace(person.SortName); person.SortName == "" {
		return errors.New("missing sortName")
	}
	person.Username = strings.ToLower(strings.TrimSpace(person.Username))
	person.CallSign = strings.ToUpper(strings.TrimSpace(person.CallSign))
	person.Email = strings.ToLower(strings.TrimSpace(person.Email))
	if person.Email != "" && !emailRE.MatchString(person.Email) {
		return errors.New("invalid email")
	}
	person.Email2 = strings.ToLower(strings.TrimSpace(person.Email2))
	if person.Email2 != "" && !emailRE.MatchString(person.Email2) {
		return errors.New("invalid email2")
	}
	if person.Email2 != "" && (person.Email == "" || person.Email == person.Email2) {
		return errors.New("invalid email2")
	}
	for i := range person.Emails {
		person.Emails[i].Email = strings.ToLower(strings.TrimSpace(person.Emails[i].Email))
		if !emailRE.MatchString(person.Emails[i].Email) {
			return errors.New("invalid email")
		}
		for j := 0; j < i; j++ {
			if person.Emails[i].Email == person.Emails[j].Email {
				return errors.New("duplicate email")
			}
		}
		person.Emails[i].Label = strings.TrimSpace(person.Emails[i].Label)
	}
	switch person.CellPhone = strings.Map(util.KeepDigits, person.CellPhone); len(person.CellPhone) {
	case 0:
		break
	case 10:
		person.CellPhone = person.CellPhone[0:3] + "-" + person.CellPhone[3:6] + "-" + person.CellPhone[6:10]
	default:
		return errors.New("invalid cell phone")
	}
	switch person.HomePhone = strings.Map(util.KeepDigits, person.HomePhone); len(person.HomePhone) {
	case 0:
		break
	case 10:
		person.HomePhone = person.HomePhone[0:3] + "-" + person.HomePhone[3:6] + "-" + person.HomePhone[6:10]
	default:
		return errors.New("invalid home phone")
	}
	switch person.WorkPhone = strings.Map(util.KeepDigits, person.WorkPhone); len(person.WorkPhone) {
	case 0:
		break
	case 10:
		person.WorkPhone = person.WorkPhone[0:3] + "-" + person.WorkPhone[3:6] + "-" + person.WorkPhone[6:10]
	default:
		return errors.New("invalid work phone")
	}
	if person.HomeAddress.Address = strings.TrimSpace(person.HomeAddress.Address); person.HomeAddress.Address != "" {
		if person.HomeAddress.Latitude < -90 || person.HomeAddress.Latitude > 90 {
			return errors.New("invalid latitude")
		}
		if person.HomeAddress.Longitude < -180 || person.HomeAddress.Longitude > 180 {
			return errors.New("invalid longitude")
		}
	} else {
		person.HomeAddress.Latitude = 0
		person.HomeAddress.Longitude = 0
	}
	person.HomeAddress.SameAsHome = false
	if person.MailAddress.Address = strings.TrimSpace(person.MailAddress.Address); person.MailAddress.Address != "" {
		person.MailAddress.SameAsHome = false
	}
	person.MailAddress.Latitude = 0
	person.MailAddress.Longitude = 0
	if person.WorkAddress.Address = strings.TrimSpace(person.WorkAddress.Address); person.WorkAddress.Address != "" {
		if person.WorkAddress.Latitude < -90 || person.WorkAddress.Latitude > 90 {
			return errors.New("invalid latitude")
		}
		if person.WorkAddress.Longitude < -180 || person.WorkAddress.Longitude > 180 {
			return errors.New("invalid longitude")
		}
		person.WorkAddress.SameAsHome = false
	} else {
		person.WorkAddress.Latitude = 0
		person.WorkAddress.Longitude = 0
	}
	people = tx.FetchPeople()
	for _, p := range people {
		if p.ID == person.ID {
			continue
		}
		if strings.EqualFold(p.SortName, person.SortName) {
			return errors.New("duplicate sortName")
		}
		if p.Username != "" && p.Username == person.Username {
			return errors.New("duplicate username")
		}
		if p.CallSign != "" && p.CallSign == person.CallSign {
			return errors.New("duplicate callSign")
		}
		if p.CellPhone != "" && p.CellPhone == person.CellPhone {
			return errors.New("duplicate cellPhone")
		}
	}
	individualHeld = cacheIndividuallyHeldRoles(tx, people, person)
	roles = make([]*model.Role, len(person.Roles))
	roleMap = make(map[model.RoleID]bool)
	for i, rid := range person.Roles {
		if roles[i] = tx.FetchRole(rid); roles[i] == nil {
			return errors.New("invalid role")
		}
		if roleMap[rid] {
			return errors.New("redundant role")
		}
		if individualHeld[roles[i]] {
			return errors.New("individual role already held")
		}
		roleMap[rid] = true
	}
	sort.Sort(model.RoleSort(roles))
	for i := range roles {
		person.Roles[i] = roles[i].ID
	}
	if len(person.Roles) == 0 && person.ID == 0 {
		return errors.New("new user with no roles")
	}
	if person.BadLoginCount < 0 {
		return errors.New("invalid badLoginCount")
	} else if person.BadLoginCount > 0 && person.BadLoginTime.IsZero() {
		return errors.New("invalid badLoginTime")
	}
	if person.PWResetToken != "" && person.PWResetTime.IsZero() {
		return errors.New("invalid pwresetTime")
	}
	for _, a := range person.Archive {
		if !strings.ContainsRune(a, '=') {
			return errors.New("invalid archive string")
		}
	}
	return nil
}

func cacheIndividuallyHeldRoles(tx *db.Tx, people []*model.Person, except *model.Person) (held map[*model.Role]bool) {
	held = make(map[*model.Role]bool)
	for _, p := range people {
		if p.ID == except.ID {
			continue
		}
		for _, role := range tx.FetchRoles() {
			if !role.Individual {
				continue
			}
			if auth.HasRole(p, role) {
				held[role] = true
			}
		}
	}
	return held
}
