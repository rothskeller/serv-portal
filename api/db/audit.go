package db

import (
	"time"
)

func (tx *Tx) audit(table string, id interface{}, data []byte) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO audit (timestamp, username, request, type, id, data) VALUES (?,?,?,?,?,?)`,
		Time(time.Now()), tx.username, tx.request, table, id, data))
}
