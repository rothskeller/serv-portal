package db

import (
	"database/sql"

	"sunnyvaleserv.org/portal/model"
)

// Search executes a search and calls the supplied function for each match.
func (tx *Tx) Search(query string, handler func(string, int, int) bool) error {
	var (
		rows *sql.Rows
		typ  string
		id   int
		id2  int
		err  error
	)
	rows, err = tx.tx.Query(`SELECT type, id, COALESCE(id2, 0) FROM search WHERE search MATCH ? ORDER BY rank`, query)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&typ, &id, &id2))
		if !handler(typ, id, id2) {
			panicOnError(rows.Close())
			break
		}
	}
	// This is where we'll get an error if the search syntax was bad.
	// We'll boldly assume that any error here is the user's fault, and
	// return it.
	return rows.Err()
}

// RebuildSearchIndex rebuilds the entire search index.
func (tx *Tx) RebuildSearchIndex(groups []*model.Group) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM search`))
	for _, tm := range tx.FetchTextMessages() {
		tx.indexTextMessage(tm)
	}
	for _, p := range tx.FetchPeople() {
		tx.indexPerson(p, false)
	}
	tx.IndexGroups(groups, false)
	for _, f := range tx.FetchFolders() {
		tx.indexFolder(f, false)
		for _, d := range f.Documents {
			tx.indexDocument(f, d, false)
		}
	}
	for _, e := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		tx.indexEvent(e, false)
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (search) VALUES ('optimize')`))
}
