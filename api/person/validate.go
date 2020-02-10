package person

import (
	"errors"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/authz"
	"sunnyvaleserv.org/portal/util"
)

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var dateRE = regexp.MustCompile(`^20\d\d-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12]\d|3[01])$`)

// ValidatePerson validates a Person record, except for its Password field.  It
// enforces all data consistency rules, but does not enforce privileges.
func ValidatePerson(tx *store.Tx, person *model.Person, roles []model.RoleID) error {
	var (
		individualHeld map[model.RoleID]model.PersonID
		people         []*model.Person
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
		person.HomeAddress.FireDistrict = FireDistrict(&person.HomeAddress)
	} else {
		person.HomeAddress.Latitude = 0
		person.HomeAddress.Longitude = 0
		person.HomeAddress.FireDistrict = 0
	}
	person.HomeAddress.SameAsHome = false
	if person.MailAddress.Address = strings.TrimSpace(person.MailAddress.Address); person.MailAddress.Address != "" {
		person.MailAddress.SameAsHome = false
	}
	person.MailAddress.Latitude = 0
	person.MailAddress.Longitude = 0
	person.MailAddress.FireDistrict = 0
	if person.WorkAddress.Address = strings.TrimSpace(person.WorkAddress.Address); person.WorkAddress.Address != "" {
		if person.WorkAddress.Latitude < -90 || person.WorkAddress.Latitude > 90 {
			return errors.New("invalid latitude")
		}
		if person.WorkAddress.Longitude < -180 || person.WorkAddress.Longitude > 180 {
			return errors.New("invalid longitude")
		}
		person.WorkAddress.SameAsHome = false
		person.WorkAddress.FireDistrict = FireDistrict(&person.WorkAddress)
	} else {
		person.WorkAddress.Latitude = 0
		person.WorkAddress.Longitude = 0
		person.WorkAddress.FireDistrict = 0
	}
	if person.UnsubscribeToken == "" {
		person.UnsubscribeToken = util.RandomToken()
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
		if p.UnsubscribeToken == person.UnsubscribeToken {
			return errors.New("duplicate unsubscribeToken")
		}
	}
	individualHeld = cacheIndividuallyHeldRoles(tx.Authorizer(), person.ID)
	roleMap = make(map[model.RoleID]bool)
	for _, rid := range roles {
		if tx.Authorizer().FetchRole(rid) == nil {
			return errors.New("invalid role")
		}
		if roleMap[rid] {
			return errors.New("redundant role")
		}
		if individualHeld[rid] != 0 {
			return errors.New("individual role already held")
		}
		roleMap[rid] = true
	}
	if roles != nil && len(roles) == 0 && person.ID == 0 {
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
	for _, n := range person.Notes {
		if n.Note = strings.TrimSpace(n.Note); n.Note == "" {
			return errors.New("missing note text")
		}
		if !dateRE.MatchString(n.Date) {
			return errors.New("invalid note date")
		}
		if n.Privilege != 0 {
			found := false
			for _, p := range model.AllPrivileges {
				if p == n.Privilege {
					found = true
					break
				}
			}
			if !found {
				return errors.New("invalid note privilege")
			}
		}
	}
	return nil
}

func cacheIndividuallyHeldRoles(auth *authz.Authorizer, except model.PersonID) (held map[model.RoleID]model.PersonID) {
	held = auth.RolesIndividuallyHeld()
	for rid, pid := range held {
		if pid == except {
			delete(held, rid)
		}
	}
	return held
}
