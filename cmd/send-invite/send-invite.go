// send-invite generates a username and password for a person and sends them an
// invitation email.
//
// usage: send-invite «personID|username|email»
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"sunnyvaleserv.org/portal/api/authn"
	"sunnyvaleserv.org/portal/api/email"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/sendmail"
)

func main() {
	var (
		entry    *log.Entry
		tx       *store.Tx
		person   *model.Person
		password string
		buf      bytes.Buffer
		body     io.Writer
		err      error
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
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: send-invite personID-or-username-or-email\n")
		os.Exit(2)
	}
	store.Open("serv.db")
	entry = log.New("", "send-invite")
	defer entry.Log()
	tx = store.Begin(entry)
	if pid := model.PersonID(util.ParseID(os.Args[1])); pid > 0 {
		if person = tx.FetchPerson(pid); person == nil {
			fmt.Fprintf(os.Stderr, "ERROR: no existing person with the ID %s\n", os.Args[1])
			os.Exit(1)
		}
	} else if person = tx.FetchPersonByUsername(os.Args[1]); person == nil {
		if person = tx.FetchPersonByEmail(os.Args[1]); person == nil {
			fmt.Fprintf(os.Stderr, "ERROR: no existing person with the username or email %q\n", os.Args[1])
			os.Exit(1)
		}
	}
	tx.WillUpdatePerson(person)
	if person.Username != "" && len(person.Password) != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: person already has a username and password\n")
		os.Exit(1)
	}
	if tx.Authorizer().MemberPG(person.ID, tx.Authorizer().FetchGroupByTag(model.GroupDisabled).ID) {
		fmt.Fprintf(os.Stderr, "ERROR: person is disabled\n")
		os.Exit(1)
	}
	if person.Username == "" {
		if person.Email == "" {
			fmt.Fprintf(os.Stderr, "ERROR: person has no email address\n")
			os.Exit(1)
		}
		if tx.FetchPersonByUsername(person.Email) != nil {
			fmt.Fprintf(os.Stderr, "ERROR: person's email %q is already in use as a username\n", person.Email)
			os.Exit(1)
		}
		person.Username = person.Email
	}
	password = authn.RandomPassword()
	authn.SetPassword(&util.Request{Tx: tx}, person, password)
	tx.UpdatePerson(person)
	tx.Commit()
	body = email.NewCRLFWriter(&buf)
	fmt.Fprintf(body, `From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>
To: %s <%s>
Subject: Welcome to SunnyvaleSERV.org!
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0

--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

Hello, %s,

An account has been created for you on the SunnyvaleSERV.org web site.  Thi=
s web site, run by the Sunnyvale Department of Public Safety, provides info=
rmation to volunteers and students in the Sunnyvale Emergency Response Volu=
nteers (SERV) programs, including:
    Community Emergency Response Team (CERT)
    LISTOS (disaster preparedness in Spanish)
    Personal Emergency Preparedness (PEP)
    Sunnyvale Amateur Radio Emergency Service (SARES)
    Sunnyvale Neighborhoods Actively Prepare (SNAP)
On this web site, you will find a calendar of SERV events, course materials=
 for SERV classes, and a wealth of other information.

SunnyvaleSERV.org is our replacement for Samariteam, which is no longer ser=
ving our needs well.  We will be transitioning to it gradually over the nex=
t few months.

To log in, visit https://SunnyvaleSERV.org and provide the following creden=
tials:
        Email:    %s
        Password: %s
You may wish to change this password to something you can more easily remem=
ber.  You can do so by clicking on the "Profile" button after you log in.

Your account has also been subscribed to the email lists for your volunteer=
 programs and classes.  You can adjust your email preferences on your Profi=
le page, or by clicking the "Unsubscribe" link at the bottom of any email.

If you have any questions about the web site, simply reply to this email.

The Sunnyvale Department of Public Safety thanks you for you time and inter=
est in emergency response.


--BOUNDARY
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

<html>
<head>
<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=utf-8">
</head>
<body>
<div>Hello, %[3]s,</div>
<div><br></div>
<div>An account has been created for you on the SunnyvaleSERV.org web site.=
 This web site, run by the Sunnyvale Department of Public Safety, provides =
information to volunteers and students in the Sunnyvale Emergency Response =
Volunteers (SERV) programs, including:</div>
<blockquote style=3D"margin-top:0px; margin-bottom:0px">
<div>Community Emergency Response Team (CERT)</div>
<div>LISTOS (disaster preparedness in Spanish)</div>
<div>Personal Emergency Preparedness (PEP)</div>
<div>Sunnyvale Amateur Radio Emergency Service (SARES)</div>
<div>Sunnyvale Neighborhoods Actively Prepare (SNAP)</div>
</blockquote>
<div>On this web site, you will find a calendar of SERV events, course mate=
rials for SERV classes, and a wealth of other information.</div>
<div><br></div>
<div>SunnyvaleSERV.org is our replacement for Samariteam, which is no longe=
r serving our needs well.  We will be transitioning to it gradually over th=
e next few months.</div>
<div><br></div>
<div>To log in, visit <a href=3D"https://SunnyvaleSERV.org">https://Sunnyva=
leSERV.org</a> and provide the following credentials:</div>
<blockquote style=3D"margin-top:0px; margin-bottom:0px">
<div>Email:<span style=3D"font-family:monospace">  %[4]s</span></div>
<div>Password:<span style=3D"font-family:monospace">  %[5]s</span></div>
</blockquote>
<div>You may wish to change this password to something you can more easily =
remember. You can do so by clicking on the "Profile" button after you log i=
n.</div>
<div><br></div>
<div>Your account has also been subscribed to the email lists for your volu=
nteer programs and classes. You can adjust your email preferences on your P=
rofile page, or by clicking the &quot;Unsubscribe&quot; link at the bottom =
of any email.</div>
<div><br></div>
<div>If you have any questions about the web site, simply reply to this ema=
il.</div>
<div><br></div>
<div>The Sunnyvale Department of Public Safety thanks you for you time and =
interest in emergency response.</div>
</body>
</html>

--BOUNDARY--
`, email.QuoteIfNeeded(person.InformalName), person.Email, person.InformalName, person.Username, password)
	if err = sendmail.SendMessage("admin@sunnyvaleserv.org", []string{person.Email}, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sending email: %s\n", err)
		os.Exit(1)
	}
}
