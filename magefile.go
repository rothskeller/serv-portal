//go:build mage
// +build mage

package main

import (
	"compress/gzip"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

// Default target.
var Default = Run

// Run runs the main program.
func Run() {
	mg.Deps(Assets)
	println("Running...")
	sh.Run(mg.GoCmd(), "run", ".")
}

// Build runs all build steps, resulting in a compiled portal executable.
func Build() {
	mg.Deps(Assets)
	sh.Run(mg.GoCmd(), "build", ".")
}

// Install builds the server (as an FCGI executable) and all associated commands
// and installs them.
func Install() error {
	mg.Deps(Assets)
	if err := sh.Run(mg.GoCmd(), "install", "./cmd/gen-ical"); err != nil {
		return err
	}
	if err := sh.Run(mg.GoCmd(), "install", "./cmd/log-report"); err != nil {
		return err
	}
	if err := sh.Run(mg.GoCmd(), "install", "./cmd/rebuild-search-index"); err != nil {
		return err
	}
	if err := sh.Run(mg.GoCmd(), "build", "-o", "received-text-hook", "./cmd/received-text-hook"); err != nil {
		return err
	}
	if err := os.Rename("received-text-hook", "/home/snyserv/sunnyvaleserv.org/received-text-hook"); err != nil {
		os.Remove("received-text-hook")
		return err
	}
	if err := sh.Run(mg.GoCmd(), "build", "-o", "text-status-hook", "./cmd/text-status-hook"); err != nil {
		return err
	}
	if err := os.Rename("text-status-hook", "/home/snyserv/sunnyvaleserv.org/text-status-hook"); err != nil {
		os.Remove("text-status-hook")
		return err
	}
	if err := sh.Run(mg.GoCmd(), "install", "./cmd/volunteer-hours"); err != nil {
		return err
	}
	if err := sh.Run(mg.GoCmd(), "build", "-o", "index.fcgi", "./cmd/portal.fcgi"); err != nil {
		os.Remove("index.fcgi")
		return err
	}
	if err := os.Rename("index.fcgi", "/home/snyserv/sunnyvaleserv.org/index.fcgi"); err != nil {
		os.Remove("index.fcgi")
		return err
	}
	sh.Run("killall", "-USR1", "-q", "index.fcgi")
	sh.Run("/home/snyserv/sunnyvaleserv.org/index.fcgi", "-writeassets")
	return nil
}

// Style sheets are listed explicitly because their order matters.
var stylesheets = []string{
	// Normalize must come first.
	"ui/normalize.css",
	// Third party library styles.
	"ui/unpoly.css",
	// Globals.
	"ui/globals.css",
	// Individual controls and web components.
	"ui/orgdot/orgdot.css",
	"ui/s-check/s-check.css",
	"ui/s-hours/s-hours.css",
	"ui/s-icon/s-icon.css",
	"ui/s-month/s-month.css",
	"ui/s-radio/s-radio.css",
	"ui/s-searchcombo/s-search.css",
	"ui/s-searchcombo/s-searchcombo.css",
	"ui/s-seltree/s-seltree.css",
	"ui/s-year/s-year.css",
	"ui/controls.css",
	// Major shared layouts.
	"ui/form.css",
	"ui/page.css",
	"ui/dialog.css",
	// Individual pages.
	"pages/admin/classedit/classedit.css",
	"pages/admin/classlist/classlist.css",
	"pages/admin/listedit/listedit.css",
	"pages/admin/listlist/listlist.css",
	"pages/admin/listpeople/listpeople.css",
	"pages/admin/roleedit/roleedit.css",
	"pages/admin/rolelist/rolelist.css",
	"pages/admin/venuelist/venuelist.css",
	"pages/classes/all.css",
	"pages/classes/cert.css",
	"pages/classes/common.css",
	"pages/classes/myn.css",
	"pages/classes/pep.css",
	"pages/classes/register.css",
	"pages/classes/reglist.css",
	"pages/errpage/errpage.css",
	"pages/events/eventattend/attendance.css",
	"pages/events/eventcopy/eventcopy.css",
	"pages/events/eventedit/details.css",
	"pages/events/eventedit/shift.css",
	"pages/events/eventscal/eventscal.css",
	"pages/events/eventslist/eventslist.css",
	"pages/events/eventview/details.css",
	"pages/events/eventview/eventview.css",
	"pages/events/eventview/ident.css",
	"pages/events/eventview/task.css",
	"pages/events/signups/shared.css",
	"pages/events/signups/signups.css",
	"pages/files/files.css",
	"pages/homepage/homepage.css",
	"pages/login/login.css",
	"pages/login/newpwd.css",
	"pages/people/activity/activity.css",
	"pages/people/peoplelist/peoplelist.css",
	"pages/people/peoplemap/peoplemap.css",
	"pages/people/personedit/contact.css",
	"pages/people/personedit/roles.css",
	"pages/people/personedit/status.css",
	"pages/people/personedit/subscriptions.css",
	"pages/people/personedit/vregister.css",
	"pages/people/personview/contact.css",
	"pages/people/personview/names.css",
	"pages/people/personview/notes.css",
	"pages/people/personview/password.css",
	"pages/people/personview/personview.css",
	"pages/people/personview/roles.css",
	"pages/people/personview/status.css",
	"pages/people/personview/subscriptions.css",
	"pages/reports/attendance/attendance.css",
	"pages/reports/clearance/clearance.css",
	"pages/search/search.css",
	"pages/static/static.css",
	"pages/texts/textlist/textlist.css",
	"pages/texts/textnew/textnew.css",
	"pages/texts/textview/textview.css",
}

// Styles generates the minified, merged style sheet.
func Styles() error {
	var (
		readers []io.Reader
		mr      io.Reader
		comp    *os.File
		gz      *gzip.Writer
		sum     hash.Hash32
		mw      io.Writer
		err     error
	)
	if outdated, err := target.Path("ui/assets/styles.css.gz", stylesheets...); err != nil || !outdated {
		return err
	}
	fmt.Println(" * build ui/assets/styles.css.gz")
	readers = make([]io.Reader, len(stylesheets))
	for i, sheet := range stylesheets {
		var fh *os.File
		if fh, err = os.Open(sheet); err != nil {
			return err
		}
		defer fh.Close()
		readers[i] = fh
	}
	mr = io.MultiReader(readers...)
	if comp, err = os.Create("ui/assets/styles.css.gz"); err != nil {
		return err
	}
	gz = gzip.NewWriter(comp)
	sum = crc32.NewIEEE()
	mw = io.MultiWriter(sum, gz)
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	if err = minifier.Minify("text/css", mw, mr); err != nil {
		return err
	}
	if err = gz.Close(); err != nil {
		return err
	}
	if err = comp.Close(); err != nil {
		return err
	}
	return nil
}

var scripts = []string{
	// Error reporting.
	"ui/jserror.js",
	// Third party libraries.
	"ui/unpoly.js",
	// Web components.
	"ui/s-check/s-check.js",
	"ui/s-hours/s-hours.js",
	"ui/s-icon/s-icon.js",
	"ui/s-month/s-month.js",
	"ui/s-radio/s-radio.js",
	"ui/s-searchcombo/s-search.js",
	"ui/s-searchcombo/s-searchcombo.js",
	"ui/s-seltree/s-seltree.js",
	"ui/s-year/s-year.js",
	// Major shared layouts.
	"ui/page.js",
	"ui/form.js",
	// Individual pages.
	"pages/admin/listedit/listedit.js",
	"pages/admin/listrole/listrole.js",
	"pages/admin/roleedit/roleedit.js",
	"pages/classes/all.js",
	"pages/classes/register.js",
	"pages/events/eventattend/attendance.js",
	"pages/events/eventedit/details.js",
	"pages/events/eventscal/eventscal.js",
	"pages/events/eventslist/eventslist.js",
	"pages/events/eventview/task.js",
	"pages/events/proxysignup/proxy.js",
	"pages/events/signups/shared.js",
	"pages/files/files.js",
	"pages/people/activity/activity.js",
	"pages/people/peoplelist/peoplelist.js",
	"pages/people/peoplemap/peoplemap.js",
	"pages/people/personedit/contact.js",
	"pages/people/personedit/password.js",
	"pages/people/personedit/roles.js",
	"pages/people/personedit/status.js",
	"pages/people/personedit/subscriptions.js",
	"pages/reports/attendance/attendance.js",
	"pages/reports/clearance/clearance.js",
	"pages/texts/textnew/textnew.js",
	"pages/texts/textview/textview.js",
}

// Scripts generates the minified, merged Javascript file.
func Scripts() error {
	var (
		readers []io.Reader
		mr      io.Reader
		comp    *os.File
		gz      *gzip.Writer
		sum     hash.Hash32
		mw      io.Writer
		err     error
	)
	if outdated, err := target.Path("ui/assets/script.js.gz", scripts...); err != nil || !outdated {
		return err
	}
	fmt.Println(" * build ui/assets/script.js.gz")
	readers = make([]io.Reader, len(scripts))
	for i, script := range scripts {
		var fh *os.File
		if fh, err = os.Open(script); err != nil {
			return err
		}
		defer fh.Close()
		readers[i] = fh
	}
	mr = io.MultiReader(readers...)
	if comp, err = os.Create("ui/assets/script.js.gz"); err != nil {
		return err
	}
	gz = gzip.NewWriter(comp)
	sum = crc32.NewIEEE()
	mw = io.MultiWriter(sum, gz)
	minifier := minify.New()
	minifier.AddFunc("text/javascript", js.Minify)
	if err = minifier.Minify("text/javascript", mw, mr); err != nil {
		return fmt.Errorf("minify JS: %s", err)
	}
	if err = gz.Close(); err != nil {
		return err
	}
	if err = comp.Close(); err != nil {
		return err
	}
	return nil
}

// Assets generates the asset registration file.
func Assets() error {
	var (
		assets []string
		out    *os.File
		err    error
	)
	mg.Deps(Styles, Scripts)
	assets, _ = filepath.Glob("ui/assets/*")
	if outdated, err := target.Path("ui/ui.assets.go", assets...); err != nil || !outdated {
		return err
	}
	fmt.Println(" * build ui/ui.assets.go")
	if out, err = os.Create("ui/ui.assets.go"); err != nil {
		return err
	}
	fmt.Fprint(out, `package ui

import _ "embed"

`)
	for i, asset := range assets {
		path := asset[3:] // skip "ui/"
		fmt.Fprintf(out, "//go:embed %q\nvar asset%d []byte\n\n", path, i)
	}
	fmt.Fprint(out, "func init() {\n")
	for i, asset := range assets {
		name := filepath.Base(asset)
		by, err := os.ReadFile(asset)
		if err != nil {
			os.Remove("ui/ui.assets.go")
			return fmt.Errorf("%s: %s", asset, err)
		}
		ext := filepath.Ext(asset)
		compressed := ext == ".gz"
		if compressed {
			name = name[:len(name)-3]
			ext = filepath.Ext(name)
		}
		mediatype := ""
		switch ext {
		case ".js":
			mediatype = "text/javascript; charset=utf-8"
		case ".css":
			mediatype = "text/css; charset=utf-8"
		case ".png":
			mediatype = "image/png"
		default:
			panic("unknown asset extension")
		}
		fmt.Fprintf(out, "\tRegisterAsset(%q, %q, asset%d, 0x%08x, %v)\n",
			name, mediatype, i, crc32.ChecksumIEEE(by), compressed)
	}
	fmt.Fprint(out, "}\n")
	out.Close()
	return nil
}

// Clean removes all transient build files and build products.
func Clean() {
	os.Remove("portal")
	os.Remove("ui/assets/styles.css.gz")
	os.Remove("ui/assets/script.js.gz")
	os.Remove("ui/ui.assets.go")
}
