package peoplemap

import (
	"encoding/json"
	"math"
	"net/http"
	"slices"
	"sort"

	"github.com/paulmach/orb"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

type personData struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}
type districtData struct {
	Points orb.Ring `json:"points"`
	Color  string   `json:"color"`
}

var districtList = []districtData{
	{Points: person.District1, Color: "#9900CC"},
	{Points: person.District2, Color: "#00CC66"},
	{Points: person.District3, Color: "#FF9966"},
	{Points: person.District4, Color: "#00CCCC"},
	{Points: person.District5, Color: "#336633"},
	{Points: person.District6, Color: "#CC99CC"},
}

// Handle handles GET and POST /people/map requests.
func Handle(r *request.Request) {
	const personFields = person.FInformalName | person.FAddresses | person.FPrivLevels | person.CanViewTargetFields
	var (
		user   *person.Person
		opts   ui.PageOpts
		focus  *role.Role
		home   bool
		work   bool
		title  string
		people []*personData
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	// Figure out what role to focus on, if any.
	if ridstr := r.FormValue("role"); ridstr != "" {
		focus = role.WithID(r, role.ID(util.ParseID(ridstr)), role.FOrg|role.FName)
	} else if rid := state.GetFocusRole(r); rid != 0 {
		focus = role.WithID(r, rid, role.FOrg|role.FName)
	}
	if focus != nil && !user.HasPrivLevel(focus.Org(), enum.PrivMember) {
		focus = nil
	}
	if focus != nil {
		title = focus.Name()
		state.SetFocusRole(r, focus.ID())
	} else {
		title = r.Loc("People")
		state.SetFocusRole(r, 0)
	}
	// Showing home or business addresses or both?
	if r.FormValue("home") != "" || r.FormValue("work") != "" {
		home = r.FormValue("home") != ""
		work = r.FormValue("work") != ""
	} else {
		home, work = true, false
	}
	// Fetch the list of people and narrow it down to those (a) whom the
	// caller can view; (b) who have GPS coordinates; and (c), if there is a
	// focus role, those who hold the focus role.
	person.All(r, personFields, func(p *person.Person) {
		if user.CanView(p) != person.ViewFull {
			return
		}
		if focus != nil {
			if held, _ := personrole.PersonHasRole(r, p.ID(), focus.ID()); !held {
				return
			}
		}
		if h := p.Addresses().Home; h != nil && h.Latitude != 0 &&
			(home || (work && p.Addresses().Work != nil && p.Addresses().Work.SameAsHome)) {
			people = append(people, &personData{
				Name: p.InformalName(),
				Lat:  h.Latitude,
				Lng:  h.Longitude,
			})
		}
		if w := p.Addresses().Work; w != nil && w.Latitude != 0 && work &&
			(!home || p.Addresses().Home == nil || p.Addresses().Home.Latitude != w.Latitude || p.Addresses().Home.Longitude != w.Longitude) {
			people = append(people, &personData{
				Name: p.InformalName() + " (W)",
				Lat:  w.Latitude,
				Lng:  w.Longitude,
			})
		}
	})
	// Where two entries have the exact same latitude and longitude, merge
	// the names.
	for i := 0; i < len(people); i++ {
		for j := i + 1; j < len(people); {
			if sameLocation(people[i], people[j]) {
				people[i].Name += "\n" + people[j].Name
				people = slices.Delete(people, j, j+1)
			} else {
				j++
			}
		}
	}
	// Show the page.
	opts = ui.PageOpts{
		Title: title,
		// ExternalScript: "https://maps.googleapis.com/maps/api/js?key=AIzaSyCi9J9RDZh5ouo3zk23yDmtY5Pp-NNBsBo",
		MenuItem: "people",
		Tabs: []ui.PageTab{
			{Name: r.Loc("List"), URL: "/people", Target: "main"},
			{Name: r.Loc("Map"), URL: "/people/map", Target: "main", Active: true},
		},
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		main.A("class=peoplemap")
		mapControls(r, user, main, focus, home, work)
		main.E("div id=peoplemapCanvas")
		if r.Method == http.MethodGet {
			dd, _ := json.Marshal(districtList)
			main.E("div class=peoplemapDistricts up-data=%s", string(dd))
		}
		pd, _ := json.Marshal(people)
		main.E("div class=peoplemapData up-data=%s", string(pd))
	})
}

func mapControls(r *request.Request, user *person.Person, main *htmlb.Element, focus *role.Role, home, work bool) {
	var roleOptions []*role.Role

	form := main.E("form class=peoplemapForm method=POST")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	// Get the list of roles the caller is allowed to focus on.
	role.All(r, role.FID|role.FName|role.FOrg|role.FFlags, func(rl *role.Role) {
		if user.HasPrivLevel(rl.Org(), enum.PrivMember) && rl.Flags()&role.Filter != 0 {
			clone := *rl
			roleOptions = append(roleOptions, &clone)
		}
	})
	// If they have more than one choice, give them a select box; otherwise
	// just add the choice as a hidden element.
	if len(roleOptions) > 1 {
		sort.Slice(roleOptions, func(i, j int) bool { return roleOptions[i].Name() < roleOptions[j].Name() })
		sel := form.E("select id=peoplemapRole name=role")
		sel.E("option value=0", focus == nil, "selected").R(r.Loc("(all)"))
		for _, role := range roleOptions {
			sel.E("option value=%d", role.ID(), focus != nil && focus.ID() == role.ID(), "selected").T(role.Name())
		}
	} else if focus != nil {
		form.E("input type=hidden name=role value=%d", focus.ID())
	}
	form.E("input type=checkbox class=s-check id=peoplemapHome name=home label=%s", r.Loc("Home[ADDR]"), home, "checked")
	form.E("input type=checkbox class=s-check id=peoplemapWork name=work label=%s", r.Loc("Business"), work, "checked")
}

func sameLocation(a, b *personData) bool {
	// To be considered the same location, they need to be within 30 feet
	// in each direction.  (Yes, that's totally arbitrary.)  That means
	// within .00008 degree latitude and .0001 degree longitude.
	return math.Abs(a.Lat-b.Lat) < 0.00008 && math.Abs(a.Lng-b.Lng) < 0.0001
}
