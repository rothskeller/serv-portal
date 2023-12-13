package main

import (
	"cmp"
	"slices"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/model"
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/store/taskrole"
	"sunnyvaleserv.org/portal/store/venue"
)

type edata struct {
	events map[event.ID]*hevent
	venues map[venue.ID]*venue.Venue
	people map[person.ID]*person.Person
	roles  map[role.ID]*role.Role
}
type hevent struct {
	event.Updater
	tasks []*htask
}
type htask struct {
	task.Updater
	roles   sets.Set[role.ID]
	people  map[person.ID]taskperson.Flag
	minutes map[person.ID]uint
	shifts  []*hshift
}
type hshift struct {
	shift.Updater
	signups map[person.ID]bool
}

func convertEvents(tx *ostore.Tx, st *nstore.Store) {
	var data = readEvents(tx, st)
	mergeOtherHours(data)
	splitStudents(data)
	addActivations(data)
	mergeEvents(st, data)
	writeEvents(st, data)
}

// readEvents reads the events from the old data store with minimal translation.
// Every event has a single Task, and that task has as many shifts as the old
// event.
func readEvents(tx *ostore.Tx, st *nstore.Store) (data *edata) {
	data = &edata{
		events: make(map[event.ID]*hevent),
		venues: make(map[venue.ID]*venue.Venue),
		people: make(map[person.ID]*person.Person),
		roles:  make(map[role.ID]*role.Role),
	}
	var rolesRequiringBGCheck = role.AllThatImply(st, 87 /* Volunteers */)
	for _, oe := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		var att = tx.FetchAttendanceByEvent(oe)
		var ne hevent
		var nt htask
		ne.ID = event.ID(oe.ID)
		ne.Name = oe.Name
		ne.Start = oe.Date + "T" + oe.Start
		ne.End = oe.Date + "T" + oe.End
		if oe.Venue != 0 {
			if vid := venue.ID(oe.Venue); data.venues[vid] == nil {
				data.venues[vid] = venue.WithID(st, vid, venue.UpdaterFields)
			}
			ne.Venue = data.venues[venue.ID(oe.Venue)]
		}
		ne.Details = oe.Details
		if oe.Type == model.EventHours {
			ne.Flags = event.OtherHours
		}
		ne.tasks = []*htask{&nt}
		data.events[ne.ID] = &ne
		nt.Name = "Tracking"
		nt.Org = enum.Org(oe.Org)
		if oe.CoveredByDSW {
			nt.Flags |= task.CoveredByDSW
		}
		if oe.Type != model.EventSocial {
			nt.Flags |= task.RecordHours
			if slices.ContainsFunc(oe.Roles, func(orid model.RoleID) (ok bool) {
				_, ok = rolesRequiringBGCheck[role.ID(orid)]
				return
			}) {
				nt.Flags |= task.RequiresBGCheck
			}
		}
		nt.Details = oe.SignupText
		nt.roles = make(sets.Set[role.ID])
		for _, rid := range oe.Roles {
			if rid := role.ID(rid); data.roles[rid] == nil {
				data.roles[rid] = role.WithID(st, rid, role.FID|role.FName|role.FPrivLevel)
			}
			nt.roles.Insert(role.ID(rid))
		}
		nt.people = make(map[person.ID]taskperson.Flag)
		nt.minutes = make(map[person.ID]uint)
		for pid, ai := range att {
			if pid := person.ID(pid); data.people[pid] == nil {
				data.people[pid] = person.WithID(st, pid, person.FID|person.FInformalName)
			}
			switch ai.Type {
			case model.AttendAsVolunteer:
				nt.Flags |= task.HasAttended
				nt.people[person.ID(pid)] = taskperson.Attended
				nt.minutes[person.ID(pid)] = uint(ai.Minutes)
			case model.AttendAsAuditor:
				nt.Flags |= task.HasAttended
				nt.people[person.ID(pid)] = taskperson.Attended
			case model.AttendAsStudent:
				nt.Flags |= task.HasAttended | task.HasCredited
				nt.people[person.ID(pid)] = taskperson.Attended | taskperson.Credited
			case model.AttendAsAbsent:
				nt.people[person.ID(pid)] = 0
				nt.minutes[person.ID(pid)] = uint(ai.Minutes)
			}
		}
		for _, os := range oe.Shifts {
			var ns hshift
			nt.shifts = append(nt.shifts, &ns)
			ns.Start = oe.Date + "T" + os.Start
			ns.End = oe.Date + "T" + os.End
			ns.Venue = ne.Venue
			ns.Min = uint(os.Min)
			ns.Max = uint(os.Max)
			if os.Announce {
				nt.Flags |= task.SignupsOpen
			}
			ns.signups = make(map[person.ID]bool)
			for _, pid := range os.SignedUp {
				ns.signups[person.ID(pid)] = true
			}
			for _, pid := range os.Declined {
				ns.signups[person.ID(pid)] = false
			}
		}
	}
	return data
}

var otherHoursLabels = map[enum.Org]string{
	enum.OrgAdmin:  "Admin",
	enum.OrgCERTD:  "CERT Deployment",
	enum.OrgCERTT:  "CERT Training",
	enum.OrgListos: "Listos",
	enum.OrgSARES:  "SARES",
	enum.OrgSNAP:   "SNAP",
}

// mergeOtherHours combines the multiple Other Hours events for a month into a
// single event with multiple Tasks.
func mergeOtherHours(data *edata) {
	for _, ne := range data.events {
		if ne.Flags&event.OtherHours == 0 || ne.tasks[0].Org != enum.OrgAdmin {
			continue
		}
		ne.Name = "Other Hours"
		for _, ne2 := range data.events {
			if ne2.Flags&event.OtherHours != 0 && ne2.tasks[0].Org != enum.OrgAdmin && ne2.Start[:7] == ne.Start[:7] {
				ne.tasks = append(ne.tasks, ne2.tasks[0])
				delete(data.events, ne2.ID)
			}
		}
		slices.SortFunc(ne.tasks, func(a, b *htask) int {
			return cmp.Compare(a.Org, b.Org)
		})
		for _, nt := range ne.tasks {
			nt.Name = otherHoursLabels[nt.Org]
		}
	}
}

// splitStudents splits any tracking task that has credits (i.e., had students)
// into Staff and Student tasks.
func splitStudents(data *edata) {
	for _, ne := range data.events {
		if len(ne.tasks) != 1 || ne.tasks[0].Flags&task.HasCredited == 0 {
			continue
		}
		var staff, students htask
		staff.roles, students.roles = make(sets.Set[role.ID]), make(sets.Set[role.ID])
		staff.people, students.people = make(map[person.ID]taskperson.Flag), make(map[person.ID]taskperson.Flag)
		staff.minutes, students.minutes = make(map[person.ID]uint), make(map[person.ID]uint)
		staff.Name, students.Name = "Staff", "Students"
		staff.Org, students.Org = ne.tasks[0].Org, ne.tasks[0].Org
		staff.Flags, students.Flags = ne.tasks[0].Flags, ne.tasks[0].Flags
		students.Flags &^= task.RecordHours | task.RequiresBGCheck
		staff.Details, students.Details = ne.tasks[0].Details, ne.tasks[0].Details
		for rid := range ne.tasks[0].roles {
			if data.roles[rid].PrivLevel() == enum.PrivStudent {
				students.roles.Insert(rid)
			} else {
				staff.roles.Insert(rid)
			}
		}
		for pid, flag := range ne.tasks[0].people {
			if flag&taskperson.Credited != 0 || ne.tasks[0].minutes[pid] == 0 {
				students.people[pid] = flag
			} else {
				staff.people[pid] = flag
				staff.minutes[pid] = ne.tasks[0].minutes[pid]
			}
		}
		if len(staff.people) != 0 {
			ne.tasks = []*htask{&staff, &students}
		} else {
			ne.tasks = []*htask{&students}
		}
	}
}

func addActivations(data *edata) {
	data.events[10].Activation = "CERT-2020-01T"
	data.events[19].Activation = "CERT-2020-05T"
	data.events[20].Activation = "CERT-2020-07T"
	data.events[444].Activation = "CERT-2020-02"
	data.events[445].Activation = "CERT-2020-02"
	data.events[454].Activation = "CERT-2020-03"
	data.events[455].Activation = "CERT-2020-03"
	data.events[456].Activation = "CERT-2020-03"
	data.events[460].Activation = "CERT-2020-04"
	data.events[461].Activation = "CERT-2020-04"
	data.events[475].Activation = "CERT-2020-06"
	data.events[476].Activation = "CERT-2020-06"
	data.events[477].Activation = "CERT-2020-06"
	data.events[478].Activation = "CERT-2020-06"
	data.events[496].Activation = "CERT-2020-08"
	data.events[497].Activation = "CERT-2020-08"
	data.events[504].Activation = "CERT-2020-10"
	data.events[505].Activation = "CERT-2020-10"
	data.events[506].Activation = "CERT-2020-09"
	data.events[513].Activation = "CERT-2020-12T"
	data.events[514].Activation = "CERT-2020-12T"
	data.events[540].Activation = "CERT-2020-11"
	data.events[541].Activation = "CERT-2020-11"
	data.events[542].Activation = "CERT-2020-11"
	data.events[543].Activation = "CERT-2020-11"
	data.events[613].Activation = "CERT-2021-04T"
	data.events[614].Activation = "CERT-2021-05T"
	data.events[615].Activation = "CERT-2021-06T"
	data.events[619].Activation = "SNY-2021-10T"
	data.events[639].Activation = "SNY-2021-07T"
	data.events[675].Activation = "CERT-2021-01"
	data.events[707].Activation = "CERT-2021-03T"
	data.events[708].Activation = "CERT-2021-03T"
	data.events[721].Activation = "CERT-2021-02T"
	data.events[746].Activation = "SNY-2021-08T"
	data.events[747].Activation = "SNY-2021-08T"
	data.events[748].Activation = "SNY-2021-08T"
	data.events[749].Activation = "SNY-2021-08T"
	data.events[752].Activation = "SNY-2021-08T"
	data.events[803].Activation = "SNY-2022-05a"
	data.events[804].Activation = "SNY-2022-05a"
	data.events[805].Activation = "SNY-2022-05a"
	data.events[806].Activation = "SNY-2022-05a"
	data.events[811].Activation = "SNY-2021-09"
	data.events[812].Activation = "SNY-2021-09"
	data.events[821].Activation = "SNY-2022-09"
	data.events[822].Activation = "SNY-2022-09"
	data.events[837].Activation = "SNY-2022-05a"
	data.events[981].Activation = "SNY-2022-03T"
	data.events[983].Activation = "SNY-2022-07T"
	data.events[985].Activation = "SNY-2022-08T"
	data.events[986].Activation = "SNY-2022-09a"
	data.events[987].Activation = "SNY-2022-10T"
	data.events[988].Activation = "SNY-2022-13T"
	data.events[989].Activation = "SNY-2022-15T"
	data.events[990].Activation = "SNY-2022-16T"
	data.events[995].Activation = "SNY-2022-02T"
	data.events[1025].Activation = "SNY-2021-11"
	data.events[1050].Activation = "SNY-2022-01"
	data.events[1051].Activation = "SNY-2022-01"
	data.events[1052].Activation = "SNY-2022-01"
	data.events[1053].Activation = "SNY-2022-01"
	data.events[1054].Activation = "SNY-2022-01"
	data.events[1066].Activation = "SNY-2022-04"
	data.events[1074].Activation = "SNY-2022-05"
	data.events[1090].Activation = "SNY-2022-06"
	data.events[1101].Activation = "SNY-2022-12T"
	data.events[1102].Activation = "SNY-2022-12T"
	data.events[1103].Activation = "SNY-2022-12T"
	data.events[1104].Activation = "SNY-2022-12T"
	data.events[1105].Activation = "SNY-2022-12T"
	data.events[1136].Activation = "SNY-2022-11T"
	data.events[1147].Activation = "SNY-2022-14"
	data.events[1150].Activation = "SNY-2023-05"
	data.events[1151].Activation = "SNY-2023-05"
	data.events[1152].Activation = "SNY-2023-05"
	data.events[1154].Activation = "SNY-2023-05"
	data.events[1172].Activation = "SNY-2022-17"
	data.events[1173].Activation = "SNY-2023-05"
	data.events[1174].Activation = "SNY-2023-05"
	data.events[1175].Activation = "SNY-2023-05"
	data.events[1176].Activation = "SNY-2023-05"
	data.events[1177].Activation = "SNY-2023-05"
	data.events[1199].Activation = "SNY-2023-08"
	data.events[1207].Activation = "SNY-2023-03"
	data.events[1209].Activation = "SNY-2023-06"
	data.events[1212].Activation = "SNY-2023-15"
	data.events[1357].Activation = "SNY-2023-01"
	data.events[1359].Activation = "SNY-2023-02"
	data.events[1360].Activation = "SNY-2023-04"
	data.events[1365].Activation = "SNY-2023-07"
	data.events[1366].Activation = "SNY-2023-10"
	data.events[1367].Activation = "SNY-2023-11"
	data.events[1368].Activation = "SNY-2023-14"
	data.events[1407].Activation = "SNY-2023-09"
	data.events[1408].Activation = "SNY-2023-09"
	data.events[1409].Activation = "SNY-2023-09"
	data.events[1410].Activation = "SNY-2023-09"
	data.events[1411].Activation = "SNY-2023-09"
	data.events[1412].Activation = "SNY-2023-09"
	data.events[1413].Activation = "SNY-2023-09"
	data.events[1414].Activation = "SNY-2023-09"
	data.events[1463].Activation = "SNY-2023-12"
	data.events[1482].Activation = "SNY-2023-13"
	data.events[1492].Activation = "SNY-2023-16"
}

// mergeEvents merges events together.
func mergeEvents(st *nstore.Store, data *edata) {
	for i := venue.ID(25); i <= 42; i++ {
		data.venues[i] = venue.WithID(st, i, venue.UpdaterFields)
	}
	data.events[721].tasks = append(data.events[721].tasks, data.events[722].tasks...)
	delete(data.events, 722)
	data.events[721].Name = "Surprise Drill"
	data.events[721].tasks[0].Name = "CERT"
	data.events[721].tasks[1].Name = "SARES"

	data.events[811].tasks = append(data.events[811].tasks, data.events[819].tasks...)
	delete(data.events, 819)
	data.events[811].Name = "Sunnyvale Art & Wine Festival"
	data.events[811].tasks[0].Name = "CERT"
	data.events[811].tasks[1].Name = "Outreach"
	data.events[811].tasks[1].shifts = []*hshift{
		data.events[811].tasks[0].shifts[0],
		data.events[811].tasks[0].shifts[2],
		data.events[811].tasks[0].shifts[5],
	}
	data.events[811].tasks[0].shifts = []*hshift{
		data.events[811].tasks[0].shifts[1],
		data.events[811].tasks[0].shifts[3],
		data.events[811].tasks[0].shifts[4],
	}

	data.events[812].tasks = append(data.events[812].tasks, data.events[820].tasks...)
	delete(data.events, 820)
	data.events[812].Name = "Sunnyvale Art & Wine Festival"
	data.events[812].tasks[0].Name = "CERT"
	data.events[812].tasks[1].Name = "Outreach"
	data.events[812].tasks[1].shifts = []*hshift{
		data.events[812].tasks[0].shifts[0],
		data.events[812].tasks[0].shifts[2],
		data.events[812].tasks[0].shifts[4],
		data.events[812].tasks[0].shifts[6],
	}
	data.events[812].tasks[0].shifts = []*hshift{
		data.events[812].tasks[0].shifts[1],
		data.events[812].tasks[0].shifts[3],
		data.events[812].tasks[0].shifts[5],
	}

	data.events[821].tasks = append(data.events[821].tasks, data.events[1099].tasks...)
	delete(data.events, 1099)
	data.events[821].Name = "Sunnyvale Art & Wine Festival"
	data.events[821].tasks[0].Name = "Outreach"
	data.events[821].tasks[1].Name = "Patrol"

	data.events[822].tasks = append(data.events[822].tasks, data.events[1100].tasks...)
	delete(data.events, 1100)
	data.events[822].Name = "Sunnyvale Art & Wine Festival"
	data.events[822].tasks[0].Name = "Outreach"
	data.events[822].tasks[1].Name = "Patrol"

	data.events[1178].tasks = append(data.events[1178].tasks, data.events[1179].tasks...)
	delete(data.events, 1179)
	data.events[1178].Name = "Los Altos Festival of Lights Parade"
	data.events[1178].tasks[0].Name = "CERT"
	data.events[1178].tasks[1].Name = "SARES"

	data.events[1126].tasks[0].shifts[0].Venue = data.venues[32]
	data.events[1126].tasks[0].shifts[1].Venue = data.venues[25]
	data.events[1126].tasks[0].shifts[2].Venue = data.venues[33]
	data.events[1126].tasks[0].shifts[3].Venue = data.venues[34]
	data.events[1126].tasks[0].shifts[4].Venue = data.venues[35]
	data.events[1126].tasks[0].shifts[5].Venue = data.venues[36]
	data.events[1126].tasks[0].shifts[6].Venue = data.venues[37]
	data.events[1126].tasks[0].shifts[7].Venue = data.venues[38]
	data.events[1126].tasks[0].shifts[8].Venue = data.venues[39]
	data.events[1126].tasks[0].shifts[9].Venue = data.venues[40]

	data.events[1415].tasks = append(data.events[1415].tasks, data.events[1416].tasks...)
	delete(data.events, 1416)
	data.events[1415].tasks = append(data.events[1415].tasks, data.events[1419].tasks...)
	delete(data.events, 1419)
	data.events[1415].tasks = append(data.events[1415].tasks, data.events[1430].tasks...)
	delete(data.events, 1430)
	data.events[1415].Name = "Sunnyvale Art & Wine Festival"
	data.events[1415].tasks[0].Name = "English Outreach"
	data.events[1415].tasks[1].Name = "Patrol"
	data.events[1415].tasks[2].Name = "Spanish Outreach"
	data.events[1415].tasks[3].Name = "SARES"

	data.events[1417].tasks = append(data.events[1417].tasks, data.events[1418].tasks...)
	delete(data.events, 1418)
	data.events[1417].tasks = append(data.events[1417].tasks, data.events[1420].tasks...)
	delete(data.events, 1420)
	data.events[1417].tasks = append(data.events[1417].tasks, data.events[1431].tasks...)
	delete(data.events, 1431)
	data.events[1417].Name = "Sunnyvale Art & Wine Festival"
	data.events[1417].tasks[0].Name = "English Outreach"
	data.events[1417].tasks[1].Name = "Patrol"
	data.events[1417].tasks[2].Name = "SARES"
	data.events[1417].tasks[3].Name = "Spanish Outreach"

	data.events[1439].tasks[0].shifts[0].Venue = data.venues[25]
	data.events[1439].tasks[0].shifts[1].Venue = data.venues[41]
	data.events[1439].tasks[0].shifts[2].Venue = data.venues[33]
	data.events[1439].tasks[0].shifts[3].Venue = data.venues[38]
	data.events[1439].tasks[0].shifts[4].Venue = data.venues[34]
	data.events[1439].tasks[0].shifts[5].Venue = data.venues[37]
	data.events[1439].tasks[0].shifts[6].Venue = data.venues[42]
	data.events[1439].tasks[0].shifts[7].Venue = data.venues[36]
	data.events[1439].tasks[0].shifts[8].Venue = data.venues[40]

	data.events[1463].tasks = []*htask{
		{
			Updater: task.Updater{
				Name:  "Fire Engine Rides",
				Org:   enum.OrgCERTD,
				Flags: task.HasAttended | task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
			},
			roles: sets.New[role.ID](66, 102),
			people: map[person.ID]taskperson.Flag{
				858: data.events[1463].tasks[0].people[858],
				990: data.events[1463].tasks[0].people[990],
				999: data.events[1463].tasks[0].people[999],
			},
			minutes: map[person.ID]uint{
				858: data.events[1463].tasks[0].minutes[858],
				990: data.events[1463].tasks[0].minutes[990],
				999: data.events[1463].tasks[0].minutes[999],
			},
			shifts: []*hshift{
				data.events[1463].tasks[0].shifts[1],
			},
		},
		{
			Updater: task.Updater{
				Name:  "Fire Tower Tours",
				Org:   enum.OrgCERTD,
				Flags: task.HasAttended | task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
			},
			roles: sets.New[role.ID](66, 102),
			people: map[person.ID]taskperson.Flag{
				1014: data.events[1463].tasks[0].people[1014],
				1027: data.events[1463].tasks[0].people[1027],
			},
			minutes: map[person.ID]uint{
				1014: data.events[1463].tasks[0].minutes[1014],
				1027: data.events[1463].tasks[0].minutes[1027],
			},
			shifts: []*hshift{
				data.events[1463].tasks[0].shifts[2],
			},
		},
		{
			Updater: task.Updater{
				Name:    "MEOC Tours",
				Org:     enum.OrgSARES,
				Flags:   task.HasAttended | task.RecordHours | task.RequiresBGCheck,
				Details: data.events[1463].tasks[0].Details[6:],
			},
			roles: sets.New[role.ID](75),
		},
		{
			Updater: task.Updater{
				Name:  "Outreach Booth",
				Org:   enum.OrgListos,
				Flags: task.HasAttended | task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
			},
			roles: sets.New[role.ID](66, 71, 72, 78, 102),
			people: map[person.ID]taskperson.Flag{
				834:  data.events[1463].tasks[0].people[834],
				1029: data.events[1463].tasks[0].people[1029],
			},
			minutes: map[person.ID]uint{
				834:  data.events[1463].tasks[0].minutes[834],
				1029: data.events[1463].tasks[0].minutes[1029],
			},
			shifts: []*hshift{
				data.events[1463].tasks[0].shifts[0],
			},
		},
	}

	data.events[1492].tasks = []*htask{
		{
			Updater: task.Updater{
				Name:    "Fire Engine Rides",
				Org:     enum.OrgCERTD,
				Flags:   task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
				Details: "CERT volunteers will be managing the fire engine rides.",
			},
			roles: sets.New[role.ID](66, 102),
			people: map[person.ID]taskperson.Flag{
				756:  data.events[1492].tasks[0].people[756],
				805:  data.events[1492].tasks[0].people[805],
				913:  data.events[1492].tasks[0].people[913],
				1034: data.events[1492].tasks[0].people[1034],
				1055: data.events[1492].tasks[0].people[1055],
			},
			minutes: map[person.ID]uint{
				756:  data.events[1492].tasks[0].minutes[756],
				805:  data.events[1492].tasks[0].minutes[805],
				913:  data.events[1492].tasks[0].minutes[913],
				1034: data.events[1492].tasks[0].minutes[1034],
				1055: data.events[1492].tasks[0].minutes[1055],
			},
			shifts: []*hshift{
				data.events[1492].tasks[0].shifts[1],
			},
		},
		{
			Updater: task.Updater{
				Name:    "English Outreach",
				Org:     enum.OrgListos,
				Flags:   task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
				Details: "SERV Outreach will be hosting our usual disaster preparedness (preparedness & response preparedness) information for this Rides for Toys event, in both English and Spanish.",
			},
			roles: sets.New[role.ID](66, 71, 72, 78, 102),
			people: map[person.ID]taskperson.Flag{
				289: data.events[1492].tasks[0].people[289],
				557: data.events[1492].tasks[0].people[557],
			},
			minutes: map[person.ID]uint{
				289: data.events[1492].tasks[0].minutes[289],
				557: data.events[1492].tasks[0].minutes[557],
			},
			shifts: []*hshift{
				data.events[1492].tasks[0].shifts[0],
			},
		},
		{
			Updater: task.Updater{
				Name:    "Spanish Outreach",
				Org:     enum.OrgListos,
				Flags:   task.RecordHours | task.RequiresBGCheck | task.SignupsOpen,
				Details: "SERV Outreach will be hosting our usual disaster preparedness (preparedness & response preparedness) information for this Rides for Toys event, in both English and Spanish.",
			},
			roles: sets.New[role.ID](66, 71, 72, 78, 102),
			people: map[person.ID]taskperson.Flag{
				143: data.events[1492].tasks[0].people[143],
				540: data.events[1492].tasks[0].people[540],
				947: data.events[1492].tasks[0].people[947],
			},
			minutes: map[person.ID]uint{
				143: data.events[1492].tasks[0].minutes[143],
				540: data.events[1492].tasks[0].minutes[540],
				947: data.events[1492].tasks[0].minutes[947],
			},
			shifts: []*hshift{
				data.events[1492].tasks[0].shifts[2],
			},
		},
	}

	data.events[1495].tasks = append(data.events[1495].tasks, data.events[1496].tasks...)
	delete(data.events, 1496)
	data.events[1495].Name = "Los Altos Festival of Lights Parade"
	data.events[1495].tasks[0].Name = "CERT"
	data.events[1495].tasks[1].Name = "SARES"
}

// writeEvents writes the new events out to the new database.
func writeEvents(st *nstore.Store, data *edata) {
	for _, he := range data.events {
		ne := event.Create(st, &he.Updater)
		for _, ht := range he.tasks {
			ht.Event = ne
			nt := task.Create(st, &ht.Updater)
			var roles []*role.Role
			for _, rid := range ht.roles.UnsortedList() {
				roles = append(roles, data.roles[rid])
			}
			taskrole.Set(st, ne, nt, roles, []*role.Role{})
			for pid, flag := range ht.people {
				taskperson.Set(st, ne, nt, data.people[person.ID(pid)], ht.minutes[pid], flag)
			}
			for _, hs := range ht.shifts {
				hs.Event = ne
				hs.Task = nt
				ns := shift.Create(st, &hs.Updater)
				for pid, flag := range hs.signups {
					if flag {
						shiftperson.SignUp(st, ne, nt, ns, data.people[person.ID(pid)])
					} else {
						shiftperson.Decline(st, ne, nt, ns, data.people[person.ID(pid)])
					}
				}
			}
		}
	}
}
