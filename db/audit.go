package db

import (
	"database/sql"
	"fmt"
	"time"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) audit(table string, id interface{}, data []byte) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO audit (timestamp, username, request, type, id, data) VALUES (?,?,?,?,?,?)`,
		Time(time.Now()), tx.username, tx.request, table, id, data))
}

// FetchAudit calls the supplied handler function for every audit record in the
// database.
func (tx *Tx) FetchAudit(handler func(time.Time, string, string, string, interface{}, interface{})) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT timestamp, username, request, type, id, data FROM audit ORDER BY timestamp`)
	panicOnError(err)
	for rows.Next() {
		var (
			timestamp time.Time
			username  string
			request   string
			otype     string
			id        interface{}
			data      []byte
			object    interface{}
		)
		panicOnError(rows.Scan((*Time)(&timestamp), &username, &request, &otype, &id, &data))
		switch otype {
		case "authz":
			var authz model.AuthzData
			panicOnError(authz.Unmarshal(data))
			object = authz
		case "event":
			var event model.Event
			panicOnError(event.Unmarshal(data))
			object = event
		case "person":
			var person model.Person
			panicOnError(person.Unmarshal(data))
			object = person
		case "session":
			var session model.Session
			var username string
			var expires string
			session.Token = model.SessionToken(id.(string))
			fmt.Sscanf(string(data), "person:%s expires:%s", &username, &expires)
			session.Person = tx.FetchPersonByUsername(username)
			session.Expires, _ = time.ParseInLocation("2006-01-02 15:04:05", expires, time.Local)
			object = session
		case "text_message":
			var text model.TextMessage
			panicOnError(text.Unmarshal(data))
			object = text
		case "venues":
			var venues model.Venues
			panicOnError(venues.Unmarshal(data))
			object = venues
		default:
			panic("unknown object type in audit record: " + otype)
		}
		handler(timestamp, username, request, otype, id, object)
	}
	panicOnError(rows.Err())
}
