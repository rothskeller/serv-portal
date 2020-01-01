// +build mage

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var Default = All

func All() {
	mg.Deps(Styles, Scripts, Binary)
}

func Styles() error {
	var stylesheets = []string{
		"./util/layout.styl",
		"./auth/login.styl",
		"./event/attendance.styl",
		"./event/edit.styl",
		"./event/list.styl",
		"./person/edit.styl",
		"./person/list.styl",
		"./report/cert-attendance.styl",
		"./report/index.styl",
		"./role/edit.styl",
		"./team/edit.styl",
		"./team/list.styl",
	}
	changed, err := target.Dir("./portal.css", stylesheets...)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	out, err := os.Create("./portal.css")
	if err != nil {
		return err
	}
	for _, infn := range stylesheets {
		css, err := sh.Output("stylus", "--compress", "--print", infn)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, css)
	}
	out.Close()
	return nil
}

func Scripts() error {
	var scriptFiles = []string{
		"./util/layout.js",
		"./auth/login.js",
		"./event/list.js",
		"./person/edit.js",
		"./person/list.js",
		"./role/edit.js",
		"./team/edit.js",
	}
	changed, err := target.Dir("./portal.js", scriptFiles...)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	out, err := os.Create("./portal.js")
	if err != nil {
		return err
	}
	for _, infn := range scriptFiles {
		in, err := os.Open(infn)
		if err != nil {
			return err
		}
		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		in.Close()
	}
	out.Close()
	return nil
}

func Binary() error {
	return sh.Run(mg.GoCmd(), "build", ".")
}
