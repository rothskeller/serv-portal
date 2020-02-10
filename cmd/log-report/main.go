package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/config"
)

type sessionData struct {
	person string
	start  string
	end    string
}

var sessionCreateRE = regexp.MustCompile(`^created? session \S+ for person "([^"]+)"`)

func main() {
	var (
		date                string
		filename            string
		file                *os.File
		scanner             *bufio.Scanner
		requestCount        int
		requestElapsedCount int
		requestElapsedSum   time.Duration
		requestElapsedMax   time.Duration
		changes             []map[string]interface{}
		errors              []map[string]interface{}
		out                 bytes.Buffer
		qpw                 *quotedprintable.Writer
		login               loginAuth
		err                 error
		sessions            = map[string]*sessionData{}
	)
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err = os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	if len(os.Args) > 1 {
		date = os.Args[1]
	} else {
		now := time.Now()
		date = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}
	filename = "log/" + date[0:7]
	if file, err = os.Open(filename); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var entry map[string]interface{}
		if err = json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: json: %s\n", err)
			os.Exit(1)
		}
		if time, ok := entry["time"].(string); !ok || !strings.HasPrefix(time, date) {
			continue
		}
		requestCount++
		if elapsed, ok := entry["elapsed"].(float64); ok {
			requestElapsedCount++
			elapsed := time.Duration(elapsed) * time.Millisecond
			requestElapsedSum += elapsed
			if requestElapsedMax < elapsed {
				requestElapsedMax = elapsed
			}
		}
		if token, ok := entry["session"].(string); ok {
			sd := sessions[token]
			if sd == nil {
				sd = new(sessionData)
				sessions[token] = sd
				sd.start = entry["time"].(string)
				if changes, ok := entry["changes"].([]interface{}); ok && len(changes) != 0 {
					if change, ok := changes[0].(string); ok {
						if match := sessionCreateRE.FindStringSubmatch(change); match != nil {
							sd.person = match[1]
						}
					}
				}
			}
		}
		if r := entry["request"].(string); r != "POST /api/login" && r != "POST /api/logout" {
			if _, ok := entry["changes"].([]interface{}); ok {
				changes = append(changes, entry)
			}
		}
		if _, ok := entry["error"].(string); ok {
			errors = append(errors, entry)
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", filename, err)
		os.Exit(1)
	}
	fmt.Fprintf(&out, "From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>\r\nTo: admin@sunnyvaleserv.org\r\nSubject: SunnyvaleSERV.org Usage Report for %s\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n", date)
	qpw = quotedprintable.NewWriter(&out)
	fmt.Fprintf(qpw, `<!DOCTYPE html><html><body><div>%d requests in %d sessions, average %dms, max %dms.</div>`,
		requestCount, len(sessions), int(requestElapsedSum/time.Millisecond)/requestElapsedCount, requestElapsedMax/time.Millisecond)
	if len(errors) != 0 {
		fmt.Fprintf(qpw, `<div style="margin-top:1em;font-weight:bold">Errors</div>`)
		for _, e := range errors {
			var person string
			if token, ok := e["session"].(string); ok {
				if session := sessions[token]; session != nil {
					person = html.EscapeString(session.person)
				}
			}
			fmt.Fprintf(qpw, `<div><span style="font-variant:tabular-nums">%s</span> %s %s:</div>`,
				e["time"].(string)[11:], person, html.EscapeString(e["request"].(string)))
			fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace">%s</div>`, html.EscapeString(e["error"].(string)))
			if stack, ok := e["stack"].(string); ok {
				fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace;white-space:pre">%s</div>`, html.EscapeString(stack))
			}
		}
	}
	if len(changes) != 0 {
		fmt.Fprintf(qpw, `<div style="margin-top:1em;font-weight:bold">Changes</div>`)
		for _, e := range changes {
			var person string
			if token, ok := e["session"].(string); ok {
				if session := sessions[token]; session != nil {
					person = html.EscapeString(session.person)
				}
			}
			fmt.Fprintf(qpw, `<div><span style="font-variant:tabular-nums">%s</span> %s %s:</div>`,
				e["time"].(string)[11:], person, html.EscapeString(e["request"].(string)))
			for _, c := range e["changes"].([]interface{}) {
				fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace">%s</div>`, html.EscapeString(c.(string)))
			}
		}
	}
	if len(sessions) != 0 {
		var slist = make([]*sessionData, 0, len(sessions))
		for _, s := range sessions {
			slist = append(slist, s)
		}
		sort.Slice(slist, func(i, j int) bool { return slist[i].start < slist[j].start })
		fmt.Fprintf(qpw, `<div style="margin-top:1em;font-weight:bold">Sessions</div>`)
		for _, s := range slist {
			var person = s.person
			if person == "" {
				person = "???"
			}
			fmt.Fprintf(qpw, `<div><span style="font-variant:tabular-nums">%s</span> %s</div>`, s.start[11:], person)
		}
	}
	fmt.Fprintf(qpw, `</body></html>`)
	qpw.Close()
	login.username = config.Get("sendGridUsername")
	login.password = config.Get("sendGridPassword")
	if err = smtp.SendMail(config.Get("sendGridServerPort"), &login, "admin@sunnyvaleserv.org", []string{"admin@sunnyvaleserv.org"}, out.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sendmail: %s\n", err)
		os.Exit(1)
	}
}

type loginAuth struct{ username, password string }

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
