package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"mime/quotedprintable"
	"os"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/sendmail"
)

func main() {
	var (
		date                string
		filename            string
		file                *os.File
		decoder             *json.Decoder
		requestCount        int
		requestElapsedCount int
		requestElapsedSum   time.Duration
		requestElapsedAvg   time.Duration
		requestElapsedMax   time.Duration
		changes             [][]string
		authn               [][]string
		errors              []map[string]interface{}
		out                 bytes.Buffer
		qpw                 *quotedprintable.Writer
		err                 error
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
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
	decoder = json.NewDecoder(file)
	for {
		var entry map[string]interface{}
		if err = decoder.Decode(&entry); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ERROR: json: %s\n", err)
			os.Exit(1)
		}
		if err == io.EOF {
			break
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
		// This next bit adds an entry to either the "changes" or
		// "authn" slice.  In either case, the first string in the entry
		// is the time, the second contains the username and request,
		// and the remainder are the "::"-separated parts of the change.
		if cs, ok := entry["changes"].([]interface{}); ok {
			for _, c := range cs {
				var components []string
				components = append(components, entry["time"].(string)[11:])
				if person, _ := entry["user"].(string); person != "" {
					components = append(components, person+" "+entry["request"].(string))
				} else {
					components = append(components, entry["request"].(string))
				}
				parts := strings.Split(c.(string), "::")
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				if parts[0] == "AuthN" {
					components = append(components, parts[1:]...)
					authn = append(authn, components)
				} else {
					components = append(components, parts...)
					changes = append(changes, components)
				}
			}
		}
		if _, ok := entry["error"].(string); ok {
			errors = append(errors, entry)
		} else if _, ok := entry["errors"].([]interface{}); ok {
			errors = append(errors, entry)
		}
	}
	reorder(changes)
	reorder(authn)
	fmt.Fprintf(&out, "From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>\r\nTo: admin@sunnyvaleserv.org\r\nSubject: SunnyvaleSERV.org Usage Report for %s\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n", date)
	qpw = quotedprintable.NewWriter(&out)
	if requestElapsedCount != 0 {
		requestElapsedAvg = requestElapsedSum / time.Duration(requestElapsedCount)
	}
	fmt.Fprintf(qpw, `<!DOCTYPE html><html><body><div>%d requests, average %dms, max %dms.</div>`,
		requestCount, requestElapsedAvg/time.Millisecond, requestElapsedMax/time.Millisecond)
	if len(errors) != 0 {
		fmt.Fprintf(qpw, `<div style="margin-top:1em;font-weight:bold">Errors</div>`)
		for _, e := range errors {
			var person string
			person, _ = e["user"].(string)
			fmt.Fprintf(qpw, `<div><span style="font-variant:tabular-nums">%s</span> %s %s:</div>`,
				e["time"].(string)[11:], person, html.EscapeString(e["request"].(string)))
			if estr, ok := e["error"].(string); ok {
				fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace">%s</div>`, html.EscapeString(estr))
			}
			if estrs, ok := e["errors"].([]interface{}); ok {
				for _, estr := range estrs {
					fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace">%s</div>`, html.EscapeString(estr.(string)))
				}
			}
			if stack, ok := e["stack"].(string); ok {
				fmt.Fprintf(qpw, `<div style="margin-left:2em;font-family:monospace;white-space:pre">%s</div>`, html.EscapeString(stack))
			}
		}
	}
	showlist(qpw, changes, "Changes")
	showlist(qpw, authn, "Authentication Updates")
	fmt.Fprintf(qpw, `</body></html>`)
	qpw.Close()
	if err = sendmail.SendMessage(context.Background(), config.Get("fromAddr"), []string{config.Get("adminEmail")}, out.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sendmail: %s\n", err)
		os.Exit(1)
	}
}

// reorder rearranges the list so that all items with common prefixes are
// adjacent.  It retains the original order (i.e., chronological) otherwise.
func reorder(list [][]string) {
	if len(list) == 0 {
		return
	}
	var newlist = make([][]string, 0, len(list))
	for i := range list {
		if list[i] == nil {
			continue
		}
		newlist = append(newlist, list[i])
		for prefixlen := len(list[i]); prefixlen > 0; prefixlen-- {
			for j := i + 1; j < len(list); j++ {
				if list[j] == nil {
					continue
				}
				match := true
				for k := 0; k < prefixlen; k++ {
					if len(list[j]) <= k || list[i][k] != list[j][k] {
						match = false
						break
					}
				}
				if match {
					newlist = append(newlist, list[j])
					list[j] = nil
				}
			}
		}
	}
	for i := range list {
		list[i] = newlist[i]
	}
}

// Each entry in the list is a slice of strings.  The first string is the time;
// the second string is the user and request, and the subsequent strings are
// details of the change.  We show this as a hierarchical list where each
// string is indented under the one before it.
func showlist(w io.Writer, list [][]string, label string) {
	var stack []string

	if len(list) == 0 {
		return
	}
	fmt.Fprintf(w, `<div style="margin-top:1em;font-weight:bold">%s</div>`, label)
	for i, item := range list {
		if len(stack) < 2 || item[1] != stack[1] {
			var multipleTimes bool
			for j := i + 1; j < len(list); j++ {
				if list[j][1] == item[1] && list[j][0] != item[0] {
					multipleTimes = true
					break
				}
			}
			if multipleTimes {
				fmt.Fprintf(w, `<div><span style="font-variant:tabular-nums">%sff</span> %s:</div>`, item[0], item[1])
			} else {
				fmt.Fprintf(w, `<div><span style="font-variant:tabular-nums">%s</span> %s:</div>`, item[0], item[1])
			}
			stack = append(stack[:0], item[0], item[1])
		}
		for i := 2; i < len(item); i++ {
			part := item[i]
			if i < len(stack) && part == stack[i] {
				continue
			}
			if i == len(item)-1 {
				fmt.Fprintf(w, `<div style="margin-left:%dem;font-family:monospace">%s</div>`, 2*i+2, part)
			} else {
				stack = append(stack[:i], part)
				fmt.Fprintf(w, `<div style="margin-left:%dem;font-family:monospace">%s::</div>`, 2*i+2, part)
			}
		}
	}
}
