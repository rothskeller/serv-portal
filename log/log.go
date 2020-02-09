package log

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/mailru/easyjson/jwriter"
)

// An Entry encapsulates all of the information that might be included in a log
// entry.  The only strictly required fields are Timestamp and Request.
type Entry struct {
	Timestamp time.Time
	Session   string
	Request   string
	Params    map[string][]string
	Status    int
	Error     string
	Stack     []byte
	Changes   []string
	Elapsed   time.Duration
}

// New creates a new Entry and populates it with the current time and the
// specified request method and path.
func New(method, path string) (e *Entry) {
	e = new(Entry)
	e.Timestamp = time.Now()
	if method != "" && path != "" {
		e.Request = method + " " + path
	} else {
		e.Request = method + path
	}
	return e
}

// Change records a change in the log entry, using printf-style arguments.
func (e *Entry) Change(s string, a ...interface{}) {
	e.Changes = append(e.Changes, fmt.Sprintf(s, a...))
}

// Log saves the log entry to the log file, atomically.  If it is unable to do
// so, it emits it to stderr instead.
func (e *Entry) Log() {
	var (
		out      jwriter.Writer
		filename string
		logfile  *os.File
		err      error
	)
	e.ToJSON(&out)
	filename = fmt.Sprintf("log/%04d-%02d", e.Timestamp.Year(), e.Timestamp.Month())
	if logfile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		goto FAIL
	}
	defer logfile.Close()
	if err = syscall.Flock(int(logfile.Fd()), syscall.LOCK_EX); err != nil {
		goto FAIL
	}
	if _, err = out.DumpTo(logfile); err != nil {
		// There won't be anything left in out if this fails, so we need
		// to re-render the entry before sending it to stderr.
		e.ToJSON(&out)
		goto FAIL
	}
	if err = logfile.Close(); err != nil {
		e.ToJSON(&out)
		return
	}
	return
FAIL:
	fmt.Fprintf(os.Stderr, "ERROR: unable to log to %s: %s\n", filename, err)
	out.DumpTo(os.Stderr)
}

// ToJSON renders the log entry in JSON form, leaving it in the provided writer.
func (e *Entry) ToJSON(out *jwriter.Writer) {
	var timebuf [21]byte

	if e.Elapsed == 0 {
		e.Elapsed = time.Since(e.Timestamp)
	}
	out.RawString(`{"time":`)
	e.Timestamp.In(time.Local).AppendFormat(timebuf[:], `"2006-01-02 15:04:05"`)
	out.Raw(timebuf[:], nil)
	if e.Session != "" {
		out.RawString(`,"session":`)
		out.String(string(e.Session))
	}
	if e.Request != "" {
		out.RawString(`,"request":`)
		out.String(e.Request)
	}
	if len(e.Params) != 0 {
		out.RawString(`,"params":{`)
		first := true
		for k, va := range e.Params {
			if len(va) == 0 || k == "auth" || k == "password" || k == "oldPassword" {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(k)
			if len(va) == 1 {
				out.RawByte(':')
				out.String(va[0])
			} else {
				out.RawString(`:[`)
				for i, v := range va {
					if i != 0 {
						out.RawByte(',')
					}
					out.String(v)
				}
				out.RawByte(']')
			}
		}
		out.RawByte('}')
	}
	if e.Status != 0 {
		out.RawString(`,"status":`)
		out.Int(e.Status)
	}
	if e.Error != "" {
		out.RawString(`,"error":`)
		out.String(e.Error)
	}
	if len(e.Stack) != 0 {
		out.RawString(`,"stack":`)
		out.String(string(e.Stack))
	}
	if len(e.Changes) != 0 {
		out.RawString(`,"changes":[`)
		for i, c := range e.Changes {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(c)
		}
		out.RawByte(']')
	}
	out.RawString(`,"elapsed":`)
	out.Int(int(e.Elapsed / time.Millisecond))
	out.RawString("}\n")
}
