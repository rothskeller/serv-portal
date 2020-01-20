package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

var madeBackup bool

func maybeMakeBackup()                           { makeBackup2(false) }
func makeBackup(_ []string, _ map[string]string) { makeBackup2(true) }
func makeBackup2(ifExists bool) {
	var (
		bfn string
		in  *os.File
		out *os.File
		err error
	)
	if madeBackup {
		return
	}
	bfn = "data/serv.db." + time.Now().Format("2006-01-02")
	if exists(bfn) {
		if !ifExists {
			return
		}
		var idx = 2
		for {
			bfn = fmt.Sprintf("%s.%d", bfn[:23], idx)
			if !exists(bfn) {
				break
			}
			idx++
		}
	}
	if in, err = os.Open("data/serv.db"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if out, err = os.Create(bfn); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if _, err = io.Copy(out, in); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR backing up to %s: %s\n", bfn, err)
		os.Remove(bfn)
		os.Exit(1)
	}
	if err = out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR backing up to %s: %s\n", bfn, err)
		os.Remove(bfn)
		os.Exit(1)
	}
	in.Close()
	madeBackup = true
}

func exists(fn string) bool {
	_, err := os.Stat(fn)
	return !os.IsNotExist(err)
}
