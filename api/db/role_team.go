package db

import (
	"database/sql"

	"rothskeller.net/serv/model"
)

func (tx *Tx) cacheTeams() {
	var (
		rows    *sql.Rows
		err     error
		parents = map[model.TeamID]model.TeamID{}
	)
	tx.teams = make(map[model.TeamID]*model.Team)
	rows, err = tx.tx.Query(`SELECT id, parent, tag, type, name, email FROM team`)
	panicOnError(err)
	for rows.Next() {
		var (
			team   model.Team
			tag    sql.NullString
			parent model.TeamID
		)
		panicOnError(rows.Scan(&team.ID, &parent, &tag, &team.Type, &team.Name, &team.Email))
		parents[team.ID] = parent
		team.Tag = model.TeamTag(tag.String)
		team.PrivMap = make(model.PrivilegeMap)
		tx.teams[team.ID] = &team
	}
	panicOnError(rows.Err())
	for c, p := range parents {
		if c == p {
			continue
		}
		tx.teams[c].Parent = tx.teams[p]
		tx.teams[p].Children = append(tx.teams[p].Children, tx.teams[c])
	}

	rows, err = tx.tx.Query(`SELECT actor, target, privileges FROM team_privilege`)
	panicOnError(err)
	for rows.Next() {
		var (
			actor, target model.TeamID
			priv          model.Privilege
		)
		panicOnError(rows.Scan(&actor, &target, &priv))
		tx.teams[actor].PrivMap.Set(tx.teams[target], priv)
	}
	panicOnError(rows.Err())

	tx.roles = make(map[model.RoleID]*model.Role)
	rows, err = tx.tx.Query(`SELECT id, team, name FROM role ORDER BY name`)
	panicOnError(err)
	for rows.Next() {
		var (
			role model.Role
			team model.TeamID
		)
		panicOnError(rows.Scan(&role.ID, &team, &role.Name))
		role.Team = tx.teams[team]
		role.PrivMap = make(model.PrivilegeMap)
		tx.teams[team].Roles = append(tx.teams[team].Roles, &role)
		tx.roles[role.ID] = &role
	}
	panicOnError(rows.Err())

	rows, err = tx.tx.Query(`SELECT role, team, privileges FROM role_privilege`)
	panicOnError(err)
	for rows.Next() {
		var (
			rid  model.RoleID
			tid  model.TeamID
			priv model.Privilege
		)
		panicOnError(rows.Scan(&rid, &tid, &priv))
		tx.roles[rid].PrivMap.Set(tx.teams[tid], priv)
	}
	panicOnError(rows.Err())
}

// FetchRole retrieves a single role from the database.  It returns nil if the
// specified role doesn't exist.
func (tx *Tx) FetchRole(id model.RoleID) *model.Role {
	return tx.roles[id]
}

// FetchTeam retrieves a single team from the database.  It returns nil if the
// specified team doesn't exist.
func (tx *Tx) FetchTeam(id model.TeamID) *model.Team {
	return tx.teams[id]
}

// FetchTeamByTag retrieves the team with the specified tag from the database.
// It returns nil if no such team exists.
func (tx *Tx) FetchTeamByTag(tag model.TeamTag) *model.Team {
	for _, team := range tx.teams {
		if team.Tag == tag {
			return team
		}
	}
	return nil
}

// FetchTeams retrieves all of the teams from the database.
func (tx *Tx) FetchTeams() (teams []*model.Team) {
	teams = make([]*model.Team, 0, len(tx.teams))
	for _, team := range tx.teams {
		teams = append(teams, team)
	}
	return teams
}

// SaveRole saves a role definition to the database.  If its supplied ID is
// zero, it creates a new role in the database, and puts its ID in the supplied
// role structure.
func (tx *Tx) SaveRole(role *model.Role) {
	var err error

	if role.ID == 0 {
		var result sql.Result
		result, err = tx.tx.Exec(`INSERT INTO role (team, name) VALUES (?,?)`, role.Team.ID, role.Name)
		panicOnError(err)
		role.ID = model.RoleID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE role SET team=?, name=? WHERE id=?`, role.Team.ID, role.Name, role.ID))
		panicOnExecError(tx.tx.Exec(`DELETE FROM role_privilege WHERE role=?`, role.ID))
	}
	tx.roles[role.ID] = role
	for _, t := range tx.teams {
		panicOnExecError(tx.tx.Exec(`INSERT INTO role_privilege (role, team, privileges) VALUES (?,?,?)`, role.ID, t.ID, role.PrivMap.Get(t)))
	}
	tx.audit(model.AuditRecord{Role: role})
}

// SaveTeam saves a team definition to the database.  If its supplied ID is
// zero, it creates a new team in the database, and puts its ID in the supplied
// team structure.
func (tx Tx) SaveTeam(team *model.Team) {
	var (
		parent model.TeamID
		stmt   *sql.Stmt
		err    error
		tag    = sql.NullString{Valid: team.Tag != "", String: string(team.Tag)}
	)
	if team.Parent != nil {
		parent = team.Parent.ID
	} else {
		parent = team.ID
	}
	if team.ID == 0 {
		var result sql.Result
		result, err = tx.tx.Exec(`INSERT INTO team (parent, tag, type, name, email) VALUES (?,?,?,?,?)`, parent, tag, team.Type, team.Name, team.Email)
		panicOnError(err)
		team.ID = model.TeamID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE team SET parent=?, tag=?, type=?, name=?, email=? WHERE id=?`, parent, tag, team.Type, team.Name, team.Email, team.ID))
	}
	tx.teams[team.ID] = team
	// Changing a team usually involves changes to other teams' privileges
	// on the target team.  So, we rewrite the entire team_privilege table.
	panicOnExecError(tx.tx.Exec(`DELETE FROM team_privilege`))
	stmt, err = tx.tx.Prepare(`INSERT INTO team_privilege (actor, target, privileges) VALUES (?,?,?)`)
	panicOnError(err)
	for _, a := range tx.teams {
		for _, t := range tx.teams {
			panicOnExecError(stmt.Exec(a.ID, t.ID, a.PrivMap.Get(t)))
		}
	}
	panicOnError(stmt.Close())
	tx.audit(model.AuditRecord{Team: team})
}

// DeleteRole deletes a role definition from the database.
func (tx *Tx) DeleteRole(role *model.Role) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM role WHERE id=?`, role.ID))
	delete(tx.roles, role.ID)
	j := 0
	for _, r := range role.Team.Roles {
		if r != role {
			role.Team.Roles[j] = r
			j++
		}
	}
	role.Team.Roles = role.Team.Roles[:j]
	tx.audit(model.AuditRecord{Role: &model.Role{ID: role.ID}})
}

// DeleteTeam deletes a team definition from the database.
func (tx *Tx) DeleteTeam(team *model.Team) {
	for _, c := range team.Children {
		tx.DeleteTeam(c)
	}
	panicOnNoRows(tx.tx.Exec(`DELETE FROM team WHERE id=?`, team.ID))
	delete(tx.teams, team.ID)
	for _, role := range team.Roles {
		delete(tx.roles, role.ID)
		role.Team = nil
	}
	if team.Parent != nil {
		j := 0
		for _, c := range team.Parent.Children {
			if c != team {
				team.Parent.Children[j] = c
				j++
			}
		}
		team.Parent.Children = team.Parent.Children[:j]
	}
	tx.audit(model.AuditRecord{Team: &model.Team{ID: team.ID}})
}
