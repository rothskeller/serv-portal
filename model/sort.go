package model

type EventSort []*Event

func (es EventSort) Len() int { return len(es) }
func (es EventSort) Less(i, j int) bool {
	switch {
	case es[i].Date < es[j].Date:
		return true
	case es[i].Date > es[j].Date:
		return false
	case es[i].Start < es[j].Start:
		return true
	case es[i].Start > es[j].Start:
		return false
	default:
		return es[i].Name < es[j].Name
	}
}
func (es EventSort) Swap(i, j int) { es[i], es[j] = es[j], es[i] }

type GroupSort []*Group

func (gs GroupSort) Len() int           { return len(gs) }
func (gs GroupSort) Less(i, j int) bool { return gs[i].Name < gs[j].Name }
func (gs GroupSort) Swap(i, j int)      { gs[i], gs[j] = gs[j], gs[i] }

type PersonSort []*Person

func (ps PersonSort) Len() int { return len(ps) }
func (ps PersonSort) Less(i, j int) bool {
	// We need a case-insensitive comparison.  But for our purposes, full
	// unicode support is not needed; we'll just do plain ASCII.
	for x := 0; x < len(ps[i].SortName) && x < len(ps[j].SortName); x++ {
		ic := ps[i].SortName[x]
		jc := ps[j].SortName[x]
		if ic >= 'a' && ic <= 'z' {
			ic -= 32
		}
		if jc >= 'a' && jc <= 'z' {
			jc -= 32
		}
		if ic < jc {
			return true
		}
		if ic > jc {
			return false
		}
	}
	return len(ps[i].SortName) < len(ps[j].SortName)
}
func (ps PersonSort) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

type RoleSort []*Role

func (rs RoleSort) Len() int           { return len(rs) }
func (rs RoleSort) Less(i, j int) bool { return rs[i].Name < rs[j].Name }
func (rs RoleSort) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }

func (vs Venues) Len() int           { return len(vs.Venues) }
func (vs Venues) Less(i, j int) bool { return vs.Venues[i].Name < vs.Venues[j].Name }
func (vs Venues) Swap(i, j int)      { vs.Venues[i], vs.Venues[j] = vs.Venues[j], vs.Venues[i] }