package orgdot

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// OrgDot emits a dot of the appropriate color for the specified Org.
func OrgDot(r *request.Request, elm *htmlb.Element, org enum.Org) {
	// colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/
	switch org {
	case enum.OrgAdmin:
		elm.E("span class=orgdot style=background-color:#a9a9a9 title=Admin")
	case enum.OrgCERTD:
		elm.E("span class=orgdot style='border:2px solid #3cb44b;background-color:white' title=%s", r.LangString("CERT Deployment", "Displiegue de CERT"))
	case enum.OrgCERTT:
		elm.E("span class=orgdot style=background-color:#3cb44b title=%s", r.LangString("CERT Training", "Capacitación CERT"))
	case enum.OrgListos:
		elm.E("span class=orgdot style=background-color:#f58231 title=Listos")
	case enum.OrgSARES:
		elm.E("span class=orgdot style=background-color:#ffe119 title=SARES")
	case enum.OrgSNAP:
		elm.E("span class=orgdot style=background-color:#4363d8 title=SNAP")
	}
}
