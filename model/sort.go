package model

type BGCheckSort []*BackgroundCheck

func (bs BGCheckSort) Len() int { return len(bs) }
func (bs BGCheckSort) Less(i, j int) bool {
	if bs[i].Date != bs[j].Date {
		return bs[i].Date < bs[j].Date
	}
	return bs[i].Type < bs[j].Type
}
func (bs BGCheckSort) Swap(i, j int) { bs[i], bs[j] = bs[j], bs[i] }

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

type DocumentSort []*Document

func (ds DocumentSort) Len() int           { return len(ds) }
func (ds DocumentSort) Less(i, j int) bool { return ds[i].Name < ds[j].Name }
func (ds DocumentSort) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }

type FolderSort []*Folder

func (fs FolderSort) Len() int           { return len(fs) }
func (fs FolderSort) Less(i, j int) bool { return fs[i].Name < fs[j].Name }
func (fs FolderSort) Swap(i, j int)      { fs[i], fs[j] = fs[j], fs[i] }

func (ls Lists) Len() int { return len(ls.Lists) }
func (ls Lists) Less(i, j int) bool {
	if ls.Lists[i].Type != ls.Lists[j].Type {
		return ls.Lists[i].Type < ls.Lists[j].Type
	}
	return ls.Lists[i].Name < ls.Lists[j].Name
}
func (ls Lists) Swap(i, j int) { ls.Lists[i], ls.Lists[j] = ls.Lists[j], ls.Lists[i] }

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

type NoteSort []*PersonNote

func (ns NoteSort) Len() int { return len(ns) }
func (ns NoteSort) Less(i, j int) bool {
	if ns[i].Date != ns[j].Date {
		return ns[j].Date < ns[i].Date
	}
	return ns[i].Note < ns[j].Note
}
func (ns NoteSort) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }

func (rs Roles) Len() int           { return len(rs.Roles) }
func (rs Roles) Less(i, j int) bool { return rs.Roles[i].Priority < rs.Roles[j].Priority }
func (rs Roles) Swap(i, j int)      { rs.Roles[i], rs.Roles[j] = rs.Roles[j], rs.Roles[i] }

func (vs Venues) Len() int           { return len(vs.Venues) }
func (vs Venues) Less(i, j int) bool { return vs.Venues[i].Name < vs.Venues[j].Name }
func (vs Venues) Swap(i, j int)      { vs.Venues[i], vs.Venues[j] = vs.Venues[j], vs.Venues[i] }
