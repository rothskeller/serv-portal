package main

import (
	"serv.rothskeller.net/portal/db"
	"serv.rothskeller.net/portal/model"
)

var roles = []*model.Role{
	{0, "login", "Web Site Users", false, nil},
	{0, "webmaster", "Webmaster", false, nil},
	{0, "", "DPS Staff", true, nil},
	{0, "", "SERV Coordinator", true, nil},
	{0, "", "SERV Volunteers", false, nil},
	{0, "", "PEP Volunteers", true, nil},
	{0, "", "PEP Leads", true, nil},
	{0, "", "SNAP Volunteers", true, nil},
	{0, "", "SNAP Leads", true, nil},
	{0, "", "CERT Volunteers", false, nil},
	{0, "", "CERT Trainers", true, nil},
	{0, "", "CERT Training Leads", true, nil},
	{0, "", "CERT Students", true, nil},
	{0, "", "CERT Graduates", true, nil},
	{0, "", "CERT Responders", false, nil},
	{0, "", "CERT Team Alpha", true, nil},
	{0, "", "CERT Team Alpha Lead", true, nil},
	{0, "", "CERT Team Alpha Operations", true, nil},
	{0, "", "CERT Group A1 Lead", true, nil},
	{0, "", "CERT Group A1 Assistant", true, nil},
	{0, "", "CERT Group A2 Lead", true, nil},
	{0, "", "CERT Group A3 Lead", true, nil},
	{0, "", "CERT Team Bravo", true, nil},
	{0, "", "CERT Team Bravo Lead", true, nil},
	{0, "", "CERT Team Bravo Operations", true, nil},
	{0, "", "CERT Group B1 Lead", true, nil},
	{0, "", "CERT Group B2 Lead", true, nil},
	{0, "", "CERT Group B3 Lead", true, nil},
	{0, "", "CERT Team Leads", false, nil},
	{0, "", "SARES Members", true, nil},
	{0, "", "SARES EC and AECs", true, nil},
	{0, "", "Listos Volunteers", true, nil},
	{0, "", "Listos Leads", true, nil},
	{0, "", "SERV Outreach Volunteers", true, nil},
	{0, "", "SERV Outreach Leads", true, nil},
	{0, "", "SERV Leads", false, nil},
	{0, "", "Web Site Guests", true, nil},
}

var implies = [][]string{
	{"Webmaster", "Web Site Users"},
	{"DPS Staff", "Web Site Users"},
	{"SERV Coordinator", "DPS Staff"},
	{"SERV Coordinator", "SERV Volunteers"},
	{"SERV Volunteers", "Web Site Users"},
	{"PEP Volunteers", "SERV Volunteers"},
	{"PEP Leads", "PEP Volunteers"},
	{"PEP Leads", "SERV Leads"},
	{"SNAP Volunteers", "SERV Volunteers"},
	{"SNAP Leads", "SNAP Volunteers"},
	{"SNAP Leads", "SERV Leads"},
	{"CERT Volunteers", "SERV Volunteers"},
	{"CERT Trainers", "CERT Volunteers"},
	{"CERT Training Leads", "CERT Volunteers"},
	{"CERT Training Leads", "SERV Leads"},
	{"CERT Students", "CERT Volunteers"},
	{"CERT Graduates", "CERT Volunteers"},
	{"CERT Responders", "CERT Volunteers"},
	{"CERT Team Alpha", "CERT Responders"},
	{"CERT Team Alpha Lead", "CERT Team Alpha"},
	{"CERT Team Alpha Lead", "CERT Team Leads"},
	{"CERT Team Alpha Operations", "CERT Team Alpha"},
	{"CERT Team Alpha Operations", "CERT Team Leads"},
	{"CERT Group A1 Lead", "CERT Team Alpha"},
	{"CERT Group A1 Assistant", "CERT Team Alpha"},
	{"CERT Group A2 Lead", "CERT Team Alpha"},
	{"CERT Group A3 Lead", "CERT Team Alpha"},
	{"CERT Team Bravo", "CERT Responders"},
	{"CERT Team Bravo Lead", "CERT Team Bravo"},
	{"CERT Team Bravo Lead", "CERT Team Leads"},
	{"CERT Team Bravo Operations", "CERT Team Bravo"},
	{"CERT Team Bravo Operations", "CERT Team Leads"},
	{"CERT Group B1 Lead", "CERT Team Bravo"},
	{"CERT Group B2 Lead", "CERT Team Bravo"},
	{"CERT Group B3 Lead", "CERT Team Bravo"},
	{"SARES Members", "SERV Volunteers"},
	{"SARES EC and AECs", "SARES Members"},
	{"SARES EC and AECs", "SERV Leads"},
	{"Listos Volunteers", "SERV Volunteers"},
	{"Listos Leads", "Listos Volunteers"},
	{"Listos Leads", "SERV Leads"},
	{"SERV Outreach Volunteers", "SERV Volunteers"},
	{"SERV Outreach Leads", "SERV Outreach Volunteers"},
	{"SERV Outreach Leads", "SERV Leads"},
	{"Web Site Guests", "Web Site Users"},
}

var people = []struct {
	lastname  string
	firstname string
	phone     string
	email     string
}{
	{"Allard", "Jeff", "(408) 835-9943", "jj1841@hotmail.com"},
	{"Cebron", "Ellie", "(669) 246-1437", "ellie.cebron@gmail.com"},
	{"Cheng", "Ka Yun", "(650) 520-3095", "fujiphotog@yahoo.com"},
	{"Chia", "Yip Fong", "(408) 733-6063", "yfchia@gmail.com"},
	{"Chien", "Jane", "(703) 587-2027", "enaj99@hotmail.com"},
	{"Chopin", "Scott", "(408) 605-5946", "scottchopin@yahoo.com"},
	{"Cohen", "Fran", "", "fran.cohen10@gmail.com"},
	{"Cohen", "Michel", "(408) 832-8582", "michelcohen@sbcglobal.net"},
	{"Davalos", "John", "(408) 718-6474", "jrdav01@comcast.net"},
	{"Duque", "Emma", "(408) 306-5831", "emma_duque@yahoo.com"},
	{"Erikson", "Kurt", "(408) 739-4634", "kurt_erikson@yahoo.com"},
	{"Flack", "Patricia", "(408) 747-7865", "pattyflack@comcast.net"},
	{"Freund", "Alice", "(408) 732-0421", "affreund@yahoo.com"},
	{"Guazelli", "Annick", "(408) 739-8497", "annickguazelli@hotmail.com"},
	{"Gupta", "Akshay", "(669) 251-9060", "aks05gupta@yahoo.co.in"},
	{"Hales", "Wendy", "(415) 690-6741", "wendyhales@gmail.com"},
	{"Hartford", "Willy", "(408) 420-7154", "willy.hartford@gmail.com"},
	{"He", "Mel", "(650) 387-1523", "miao.he@gmail.com"},
	{"Hood", "Elizabeth", "(408) 220-3641", "lizzym.hood@yahoo.com"},
	{"Howey", "Andrew", "(650) 279-7449", "ajhowey@gmail.com"},
	{"Hsin", "Heart", "(408) 802-6589", "hearthsin@gmail.com"},
	{"Icasiano", "Alexander", "(408) 973-8956", "aicasiano@gmail.com"},
	{"Knoefels", "Andreas", "(408) 910-0397", "andreas@knoefels.org"},
	{"Letourneau", "Emilie", "(702) 610-8460", "emilieariel@gmail.com"},
	{"Moynihan", "Andrew", "(205) 792-5437", "arm0002@gmail.com"},
	{"Mullins", "Brian", "(408) 480-1261", "habu0313@yahoo.com"},
	{"Oseguera", "Jazmin", "(408) 661-7480", "jazmin_toraloseguera@yahoo.com"},
	{"Pease", "Barbara", "(408) 832-8399", "bvpease051@gmail.com"},
	{"Pease", "Roger", "", "rogermpease@yahoo.com"},
	{"Quait", "Sara", "(408) 691-3109", "sbquait@gmail.com"},
	{"Respicio", "Loralyn", "(408) 515-1764", "loralyn830@hotmail.com"},
	{"Roth", "Steve", "(408) 234-5674", "sroth@sunnyvale.ca.gov"},
	{"Sehgal", "Sumit", "(408) 839-3761", "sumitsehgal23@gmail.com"},
	{"Steichen", "Greg", "(408) 205-1214", "gregsteichen@gmail.com"},
	{"Sysmans", "Jan", "(408) 582-3475", "jsysmans@gmail.com"},
	{"Tawil", "Afifa", "(623) 340-8830", "afifa.tawil@gmail.com"},
	{"Veilande", "Una", "(408) 813-5104", "una.veilande@gmail.com"},
	{"Vilcans", "Didzis", "(410) 343-9471", "events@vilcans.eu"},
	{"Wang", "Tom", "(408) 218-3514", "tcwang76@hotmail.com"},
	{"Wheeler", "Linda", "(408) 730-4014", "lwheeler88@yahoo.com"},
	{"Woods", "Emerick", "(410) 343-9471", "ewoods@emerickwoods.com"},
	{"Wright", "Amy", "(650) 799-3112", "amytwright@sbcglobal.net"},
	{"Yamada", "Kelly", "(408) 209-1914", "dracos_ky@yahoo.com"},
	{"DiGiovanna", "Liz", "", "edigiovanna@sunnyvale.ca.gov"},
	{"Werges", "Kent", "", "kwerges@sunnyvale.ca.gov"},
}

var proles = []struct {
	email string
	role  string
}{
	{"jj1841@hotmail.com", "CERT Group A3 Lead"},
	{"ellie.cebron@gmail.com", "CERT Group B3 Lead"},
	{"fujiphotog@yahoo.com", "CERT Team Bravo"},
	{"yfchia@gmail.com", "CERT Team Alpha"},
	{"enaj99@hotmail.com", "CERT Team Alpha"},
	{"scottchopin@yahoo.com", "CERT Team Alpha"},
	{"fran.cohen10@gmail.com", "CERT Team Alpha"},
	{"michelcohen@sbcglobal.net", "CERT Team Bravo"},
	{"jrdav01@comcast.net", "CERT Team Bravo"},
	{"emma_duque@yahoo.com", "CERT Team Alpha Lead"},
	{"kurt_erikson@yahoo.com", "CERT Team Bravo"},
	{"pattyflack@comcast.net", "CERT Team Alpha"},
	{"affreund@yahoo.com", "CERT Team Bravo"},
	{"annickguazelli@hotmail.com", "CERT Team Alpha"},
	{"aks05gupta@yahoo.co.in", "CERT Team Bravo"},
	{"wendyhales@gmail.com", "CERT Team Bravo Lead"},
	{"willy.hartford@gmail.com", "CERT Team Bravo"},
	{"miao.he@gmail.com", "CERT Team Alpha"},
	{"lizzym.hood@yahoo.com", "CERT Team Alpha"},
	{"ajhowey@gmail.com", "CERT Team Bravo"},
	{"hearthsin@gmail.com", "CERT Group B2 Lead"},
	{"aicasiano@gmail.com", "CERT Group A1 Lead"},
	{"aicasiano@gmail.com", "CERT Training Leads"},
	{"andreas@knoefels.org", "CERT Team Alpha"},
	{"emilieariel@gmail.com", "CERT Group A1 Assistant"},
	{"arm0002@gmail.com", "CERT Team Bravo"},
	{"habu0313@yahoo.com", "CERT Team Alpha"},
	{"jazmin_toraloseguera@yahoo.com", "CERT Team Alpha"},
	{"jazmin_toraloseguera@yahoo.com", "Listos Leads"},
	{"bvpease051@gmail.com", "CERT Group B1 Lead"},
	{"rogermpease@yahoo.com", "CERT Team Alpha"},
	{"sbquait@gmail.com", "CERT Team Bravo"},
	{"sbquait@gmail.com", "PEP Leads"},
	{"loralyn830@hotmail.com", "CERT Team Bravo"},
	{"sroth@sunnyvale.ca.gov", "CERT Team Alpha"},
	{"sroth@sunnyvale.ca.gov", "SERV Coordinator"},
	{"sroth@sunnyvale.ca.gov", "SARES Members"},
	{"sumitsehgal23@gmail.com", "CERT Team Alpha"},
	{"gregsteichen@gmail.com", "CERT Team Bravo Operations"},
	{"jsysmans@gmail.com", "CERT Team Bravo"},
	{"afifa.tawil@gmail.com", "CERT Team Bravo"},
	{"una.veilande@gmail.com", "CERT Team Bravo"},
	{"events@vilcans.eu", "CERT Group A2 Lead"},
	{"tcwang76@hotmail.com", "CERT Team Bravo"},
	{"lwheeler88@yahoo.com", "CERT Team Bravo"},
	{"ewoods@emerickwoods.com", "CERT Team Alpha"},
	{"amytwright@sbcglobal.net", "CERT Team Bravo"},
	{"dracos_ky@yahoo.com", "CERT Team Alpha"},
	{"edigiovanna@sunnyvale.ca.gov", "DPS Staff"},
	{"kwerges@sunnyvale.ca.gov", "DPS Staff"},
}

func main() {
	db.Open("serv.db")
	tx := db.Begin()
	tx.SetRequest("populate")
	tx.SetUsername("populate")
	for _, role := range roles {
		tx.SaveRole(role)
	}
	for _, ipair := range implies {
		rname, impname := ipair[0], ipair[1]
		var role, implies *model.Role
		for _, r := range roles {
			if r.Name == rname {
				role = r
			}
			if r.Name == impname {
				implies = r
			}
		}
		role.Implies = append(role.Implies, implies)
	}
	for _, role := range roles {
		tx.SaveRole(role)
	}
	var plist []*model.Person
	for _, p := range people {
		plist = append(plist, &model.Person{FirstName: p.firstname, LastName: p.lastname, Phone: p.phone, Email: p.email})
	}
	for _, pr := range proles {
		var person *model.Person
		for _, p := range plist {
			if p.Email == pr.email {
				person = p
			}
		}
		var role *model.Role
		for _, r := range roles {
			if r.Name == pr.role {
				role = r
			}
		}
		person.ExplicitRoles = append(person.ExplicitRoles, role)
	}
	for _, p := range plist {
		tx.SavePerson(p)
	}
	tx.Commit()
}
