//go:build darwin

package osdep

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var (
	// HomeDir is the user's home directory (or the root, if for some reason
	// the user's home directory isn't discernable).
	HomeDir string
)

func init() {
	if HomeDir = os.Getenv("HOME"); HomeDir == "" {
		fmt.Fprintln(os.Stderr, "ERROR: can't locate data files: $HOME not set")
		os.Exit(1)
	}
}

// ReadLock locks the file for reading.
func ReadLock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_SH)
}

// WriteLock locks the file for writing.
func WriteLock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_EX)
}

// Unlock unlocks the file.
func Unlock(fh *os.File) (err error) {
	return syscall.Flock(int(fh.Fd()), syscall.LOCK_UN)
}

// DetachChild is the SysProcAttr to use in a *exec.Cmd when creating a process
// that should run independent of its parent (i.e., the server).
var DetachChild = &syscall.SysProcAttr{
	// Setsid:  true,
	// Setpgid: true,
	// Noctty: true,
}

// OpenURLCommand is the command to cause the system default browser to open a
// URL.
func OpenURLCommand(url string) *exec.Cmd {
	return exec.Command("open", url)
}

// OpenFileCommand is the command to open a file using the system default
// application for its file type.
func OpenFileCommand(file string) *exec.Cmd {
	return exec.Command("open", file)
}

// IsAdmin returns whether the user is an Administrator.
func IsAdmin() bool { return false } // only valid on Windows
