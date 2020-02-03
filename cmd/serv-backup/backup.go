package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	var (
		bfn         string
		in          *os.File
		out         *os.File
		err         error
		ifNotExists = flag.Bool("i", false, "back up only if none taken today")
	)
	flag.Parse()
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		os.Chdir("/home/snyserv/sunnyvaleserv.org")
	case "/Users/stever":
		os.Chdir("/Users/stever/src/serv-portal")
	}
	bfn = "data/serv.db." + time.Now().Format("2006-01-02")
	if exists(bfn) {
		if *ifNotExists {
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
}

func exists(fn string) bool {
	_, err := os.Stat(fn)
	return !os.IsNotExist(err)
}
