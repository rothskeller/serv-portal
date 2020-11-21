package list

import (
	"errors"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

var listEmailRE = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

// ValidateList verifies the parameters of a list.
func ValidateList(tx *store.Tx, list *model.List) error {
	if _, ok := model.ListTypeNames[list.Type]; !ok {
		return errors.New("invalid type")
	}
	if list.Name = strings.TrimSpace(list.Name); list.Name == "" {
		return errors.New("missing name")
	}
	if list.Type == model.ListEmail && !listEmailRE.MatchString(list.Name) {
		return errors.New("invalid name")
	}
	for pid := range list.People {
		if tx.FetchPerson(pid) == nil {
			return errors.New("nonexistent person in people")
		}
	}
	for _, l := range tx.FetchLists() {
		if l == list {
			continue
		}
		if l.Name == list.Name {
			return errors.New("duplicate name")
		}
	}
	return nil
}
