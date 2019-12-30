package db

import (
	"os"
	"syscall"
	"time"

	"github.com/mailru/easyjson"

	"serv.rothskeller.net/portal/model"
)

const auditFileName = "audit.log"

func (tx *Tx) audit(ar model.AuditRecord) {
	var (
		auditfile *os.File
		err       error
	)
	ar.Timestamp = time.Now()
	ar.Username = tx.username
	ar.Request = tx.request
	if auditfile, err = os.OpenFile(auditFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		panic("can't open audit file: " + err.Error())
	}
	defer auditfile.Close()
	if err = syscall.Flock(int(auditfile.Fd()), syscall.LOCK_EX); err != nil {
		panic("can't lock audit file: " + err.Error())
	}
	easyjson.MarshalToWriter(&ar, auditfile)
}
