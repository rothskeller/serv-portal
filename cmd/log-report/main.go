package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"mime/quotedprintable"
	"os"
	"regexp"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/util/sendmail"
)

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
		requestElapsedAvg   time.Duration
		requestElapsedMax   time.Duration
		changes             []map[string]interface{}
		errors              []map[string]interface{}
		out                 bytes.Buffer
		qpw                 *quotedprintable.Writer
		err                 error
		sessions            = map[string]string{}
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
		if token, ok := entry["session"].(string); ok {
			if sessions[token] == "" {
				if changes, ok := entry["changes"].([]interface{}); ok && len(changes) != 0 {
					if change, ok := changes[0].(string); ok {
						if match := sessionCreateRE.FindStringSubmatch(change); match != nil {
							sessions[token] = match[1]
						}
					}
				}
			}
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
		if _, ok := entry["changes"].([]interface{}); ok {
			changes = append(changes, entry)
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
	if requestElapsedCount != 0 {
		requestElapsedAvg = requestElapsedSum / time.Duration(requestElapsedCount)
	}
	fmt.Fprintf(qpw, `<!DOCTYPE html><html><body><div>%d requests in %d sessions, average %dms, max %dms.</div>`,
		requestCount, len(sessions), requestElapsedAvg/time.Millisecond, requestElapsedMax/time.Millisecond)
	if len(errors) != 0 {
		fmt.Fprintf(qpw, `<div style="margin-top:1em;font-weight:bold">Errors</div>`)
		for _, e := range errors {
			var person string
			if token, ok := e["session"].(string); ok {
				person = html.EscapeString(sessions[token])
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
				person = html.EscapeString(sessions[token])
			}
			fmt.Fprintf(qpw, `<div><span style="font-variant:tabular-nums">%s</span> %s %s:</div>`,
				e["time"].(string)[11:], person, html.EscapeString(e["request"].(string)))
			for _, c := range e["changes"].([]interface{}) {
				fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace">%s</div>`, html.EscapeString(c.(string)))
			}
		}
	}
	fmt.Fprintf(qpw, `</body></html>`)
	qpw.Close()
	if err = sendmail.SendMessage("admin@sunnyvaleserv.org", []string{"admin@sunnyvaleserv.org"}, out.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sendmail: %s\n", err)
		os.Exit(1)
	}
}
