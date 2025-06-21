package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"k8s.io/apimachinery/pkg/util/sets"
)

type mailMetadata struct {
	lists      sets.Set[string]
	verdicts   []string
	approved   sets.Set[string]
	rejects    sets.Set[string]
	moderating sets.Set[string]
	sentToList sets.Set[string]
	sent       sets.Set[string]
}

func readTracking(tf *os.File) (mail *mailMetadata, err error) {
	var scan *bufio.Scanner

	if err = syscall.Flock(int(tf.Fd()), syscall.LOCK_EX); err != nil {
		return nil, fmt.Errorf("lock: %w", err)
	}
	scan = bufio.NewScanner(tf)
	mail = &mailMetadata{
		lists:      sets.New[string](),
		approved:   sets.New[string](),
		rejects:    sets.New[string](),
		moderating: sets.New[string](),
		sentToList: sets.New[string](),
		sent:       sets.New[string](),
	}
	for scan.Scan() {
		line := scan.Text()
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			goto ERROR
		}
		switch fields[0] {
		case "R": // (received) timestamp recipient...
			if len(fields) < 3 {
				goto ERROR
			}
			mail.lists.Insert(fields[2:]...)
		case "V": // (verdicts) dkim dmarc spf spam virus
			if len(fields) != 6 {
				goto ERROR
			}
			mail.verdicts = fields[1:]
		case "A": // (moderation approval) timestamp list from
			if len(fields) != 4 {
				goto ERROR
			}
			mail.approved.Insert(fields[2])
		case "X": // (unknown recipient) timestamp recipient
			if len(fields) != 3 {
				goto ERROR
			}
			mail.rejects.Insert(fields[2])
		case "M": // (approval requested) timestamp list
			if len(fields) != 3 {
				goto ERROR
			}
			mail.moderating.Insert(fields[2])
		case "L": // (sent to list) timestamp list
			if len(fields) != 3 {
				goto ERROR
			}
			mail.sentToList.Insert(fields[2])
		case "S": // (sent to recipient) timestamp destination
			if len(fields) != 3 {
				goto ERROR
			}
			mail.sent.Insert(fields[2])
		case "E": // (error) message
			continue
		default:
			return nil, fmt.Errorf("unknown line code %s", fields[0])
		}
	}
	if err = scan.Err(); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return mail, nil
ERROR:
	return nil, fmt.Errorf("read: invalid line")
}
