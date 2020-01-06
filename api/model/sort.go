package model

type RoleSort []*Role

func (rs RoleSort) Len() int           { return len(rs) }
func (rs RoleSort) Less(i, j int) bool { return rs[i].Name < rs[j].Name }
func (rs RoleSort) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
