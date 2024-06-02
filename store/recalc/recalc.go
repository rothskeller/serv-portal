package recalc

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/listperson"
)

const (
	deleteImpliedRolesSQL         = `DELETE FROM person_role WHERE NOT explicit`
	deletePrivLevelsSQL           = `DELETE FROM person_privlevel`
	deleteListSendersSQL          = `UPDATE list_person SET sender=FALSE`
	deleteUnusedListPersonRowsSQL = `DELETE FROM list_person WHERE NOT sender AND NOT sub AND NOT unsub`
)
const addImpliedRolesSQL = `
WITH RECURSIVE implies AS (
	SELECT ri.implier, ri.implied FROM role_implies ri
	UNION
	SELECT implies.implier, ri.implied FROM role_implies ri, implies WHERE implies.implied=ri.implier
)
INSERT OR IGNORE INTO person_role
SELECT pr.person, implies.implied, 0
FROM   person_role pr, implies
WHERE  pr.role=implies.implier`
const addPrivLevelsSQL = `
INSERT INTO person_privlevel
SELECT   pr.person, r.org, MAX(r.privlevel)
FROM     person_role pr, role r
WHERE    pr.role=r.id AND r.org IS NOT NULL
GROUP BY pr.person, r.org`
const deletePrivLevelsForDisabledSQL = `
DELETE FROM person_privlevel
WHERE  EXISTS (SELECT 1 FROM person_role pr WHERE pr.person=person_privlevel.person AND pr.role=2)`
const addAdminMasterForWebmasterSQL = `
INSERT OR REPLACE INTO person_privlevel
SELECT pr.person, 1, 4 FROM person_role pr WHERE pr.role=1`
const addPrivLevelsForAdminLeaderSQL = `
WITH orgs (org) AS (VALUES (2), (3), (4), (5), (6))
INSERT OR REPLACE INTO person_privlevel
SELECT pp.person, orgs.org, pp.privlevel FROM person_privlevel pp, orgs WHERE pp.org=1 AND pp.privlevel>=3`
const addListSendersSQL = `
INSERT INTO list_person
SELECT lr.list, pr.person, TRUE, FALSE, FALSE
FROM   list_role lr, person_role pr
WHERE  lr.role=pr.role AND lr.sender
ON CONFLICT DO UPDATE SET sender=TRUE`
const deleteDisallowedListReceiversSQL = `
UPDATE list_person SET sub=FALSE
WHERE  NOT EXISTS (
  SELECT 1 FROM list_role lr, person_role pr
  WHERE  pr.person=list_person.person AND lr.list=list_person.list AND lr.role=pr.role AND lr.submodel > 0
)`
const addAutomaticListReceiversSQL = `
INSERT INTO list_person
SELECT lr.list, pr.person, FALSE, TRUE, FALSE
FROM   list_role lr, person_role pr
WHERE  lr.role=pr.role and lr.submodel >= 2
ON CONFLICT DO UPDATE SET sub=TRUE`

// Recalculate recalculates the implicit role assignments, per-organization
// privilege levels, and list privileges and memberships for all people in the
// database.  This should be called whenever role implications are changed,
// roles are deleted, role privileges or memberships on a list are changed, or
// explicit role assignments to people are changed.
func Recalculate(storer phys.Storer) {
	var listdata []byte

	// Some of the cases listed above could be streamlined, acting on only a
	// single person or a single list.  But recalculating everything is fast
	// enough that there's no justification for the added code complexity.
	storer.AsStore().Transaction(func() {
		phys.Exec(storer, deleteImpliedRolesSQL)
		phys.Exec(storer, addImpliedRolesSQL)
		phys.Exec(storer, deletePrivLevelsSQL)
		phys.Exec(storer, addPrivLevelsSQL)
		phys.Exec(storer, deletePrivLevelsForDisabledSQL)
		phys.Exec(storer, addAdminMasterForWebmasterSQL)
		phys.Exec(storer, addPrivLevelsForAdminLeaderSQL)
		phys.Exec(storer, deleteListSendersSQL)
		phys.Exec(storer, addListSendersSQL)
		phys.Exec(storer, deleteDisallowedListReceiversSQL)
		phys.Exec(storer, addAutomaticListReceiversSQL)
		phys.Exec(storer, deleteUnusedListPersonRowsSQL)
		listdata = listperson.ListData(storer)
	})
	phys.UploadEmailListData(storer, listdata)
}
