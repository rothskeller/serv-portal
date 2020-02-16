package store

import "sunnyvaleserv.org/portal/model"

// CreateApproval creates a new Approval in the database.
func (tx *Tx) CreateApproval(a *model.Approval) {
	var (
		alist model.Approvals
	)
	alist = tx.Tx.FetchApprovals()
	alist.Approvals = append(alist.Approvals, a)
	tx.Tx.SaveApprovals(alist)
	if a.SubID != 0 {
		tx.entry.Change("create approval %s [%d.%d] group %q [%d]", model.PrivilegeNames[a.Privilege], a.ID, a.SubID, tx.auth.FetchGroup(a.Group).Name, a.Group)
	} else {
		tx.entry.Change("create approval %s [%d] group %q [%d]", model.PrivilegeNames[a.Privilege], a.ID, tx.auth.FetchGroup(a.Group).Name, a.Group)
	}
}

// DeleteApproval deletes an Approval from the database.
func (tx *Tx) DeleteApproval(a *model.Approval) {
	var (
		alist model.Approvals
		j     int
	)
	alist = tx.Tx.FetchApprovals()
	for _, ai := range alist.Approvals {
		if a.ID != ai.ID || a.SubID != ai.SubID || a.Privilege != ai.Privilege {
			alist.Approvals[j] = ai
			j++
		}
	}
	alist.Approvals = alist.Approvals[:j]
	tx.Tx.SaveApprovals(alist)
	if a.SubID != 0 {
		tx.entry.Change("delete approval %s [%d.%d] group %q [%d]", model.PrivilegeNames[a.Privilege], a.ID, a.SubID, tx.auth.FetchGroup(a.Group).Name, a.Group)
	} else {
		tx.entry.Change("delete approval %s [%d] group %q [%d]", model.PrivilegeNames[a.Privilege], a.ID, tx.auth.FetchGroup(a.Group).Name, a.Group)
	}
}
