package db

import (
	"database/sql"

	"sunnyvaleserv.org/portal/model"
)

// FetchApprovals retrieves the list of approvals from the database.
func (tx *Tx) FetchApprovals() (a model.Approvals) {
	var (
		data []byte
		err  error
	)
	switch err = tx.tx.QueryRow(`SELECT data FROM approval`).Scan(&data); err {
	case nil:
		panicOnError(a.Unmarshal(data))
		return a
	case sql.ErrNoRows:
		return model.Approvals{}
	default:
		panic(err)
	}
}

// SaveApprovals saves the list of approvals to the database.
func (tx *Tx) SaveApprovals(a model.Approvals) {
	var (
		data []byte
		res  sql.Result
		err  error
	)
	data, err = a.Marshal()
	panicOnError(err)
	res, err = tx.tx.Exec(`UPDATE approval SET data=?`, data)
	panicOnError(err)
	if rc, _ := res.RowsAffected(); rc == 0 {
		panicOnExecError(tx.tx.Exec(`INSERT INTO approval (data) VALUES (?)`, data))
	}
}
