package list

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

const withIDSQL = `SELECT type, name, moderators FROM list WHERE id=?`

// WithID returns the list with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID) (l *List) {
	phys.SQL(storer, withIDSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			l = new(List)
			l.ID = id
			l.Type = Type(stmt.ColumnInt())
			l.Name = stmt.ColumnText()
			l.Moderators = unpackModerators(stmt.ColumnText())
		}
	})
	return l
}

const allSQL = `SELECT id, type, name, moderators FROM list ORDER BY name`

// All reads each list from the database, in order by name.
func All(storer phys.Storer, fn func(*List)) {
	phys.SQL(storer, allSQL, func(stmt *phys.Stmt) {
		var l List
		for stmt.Step() {
			l.ID = ID(stmt.ColumnInt())
			l.Type = Type(stmt.ColumnInt())
			l.Name = stmt.ColumnText()
			l.Moderators = unpackModerators(stmt.ColumnText())
			fn(&l)
		}
	})
}

func unpackModerators(s string) (mods sets.Set[string]) {
	if s == "" {
		return nil
	}
	return sets.New[string](strings.Split(s, ",")...)
}
