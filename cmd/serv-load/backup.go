package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func maybeMakeBackup() {
	var (
		bfn string
		in  *os.File
		out *os.File
		err error
	)
	bfn = "serv.db." + time.Now().Format("2006-01-02")
	if _, err = os.Stat(bfn); !os.IsNotExist(err) {
		return
	}
	if in, err = os.Open("serv.db"); err != nil {
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
