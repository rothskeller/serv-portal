package roleselect

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// NewRoleSelectRow creates and returns a form row with a role selection control
// in it.
func NewRoleSelectRow(
	r *request.Request, fields role.Fields, filter func(*role.Role) bool, label, name string, valueP *[]role.ID, rollup bool,
) form.Row {
	var rsr = roleSelectRow{
		LabeledRow: form.LabeledRow{Label: label},
		name:       name, valueP: valueP, rollup: rollup,
	}
	rsr.treedesc = MakeRoleTree(r, fields, filter)
	return &rsr
}

type roleSelectRow struct {
	form.LabeledRow
	name     string
	valueP   *[]role.ID
	rollup   bool
	treedesc string
}

func (rsr *roleSelectRow) Read(r *request.Request) bool {
	*rsr.valueP = (*rsr.valueP)[:0]
	for _, ridstr := range strings.Fields(r.FormValue(rsr.name)) {
		*rsr.valueP = append(*rsr.valueP, role.ID(util.ParseID(ridstr)))
		// We could check for existence, duplicates, etc., but it's not
		// worth bothering.
	}
	return true
}

func (rsr *roleSelectRow) ShouldEmit(vl request.ValidationList) bool {
	return !vl.Enabled()
}

func (rsr *roleSelectRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	row := rsr.EmitPrefix(r, parent, "")
	var vstrs = make([]string, len(*rsr.valueP))
	for i, rid := range *rsr.valueP {
		vstrs[i] = strconv.Itoa(int(rid))
	}
	row.E("s-seltree name=%s class=formInput value=%s", rsr.name, strings.Join(vstrs, " "),
		rsr.rollup, "rollup").R(rsr.treedesc)
	rsr.EmitSuffix(r, row)
}

type node struct {
	role      *role.Role
	impliedBy []*node
	disabled  bool
}

// MakeRoleTree builds a "tree" of roles based on implications.  Note that the
// same role may appear in multiple places, since the implications are a
// directed acyclic graph rather than a strict tree.  If filter is non-nil, only
// roles accepted by the filter, or roles implied by such roles, are included in
// the tree.  The return value is a textual tree representation of the form
// accepted by the s-seltree control.
func MakeRoleTree(r *request.Request, fields role.Fields, filter func(*role.Role) bool) (treedesc string) {
	var (
		tree  []*node
		nodes []*node
		sb    strings.Builder
		idmap = make(map[role.ID]*node)
	)
	fields |= role.FID | role.FName | role.FImplies | role.FPriority
	// Fetch the list of roles.
	role.All(r, fields, func(rl *role.Role) {
		node := &node{role: rl.Clone()}
		nodes = append(nodes, node)
		idmap[rl.ID()] = node
	})
	// Now build the tree based on the implies.
	for _, n := range nodes {
		var implied bool

		for _, impid := range n.role.Implies() {
			imp := idmap[impid]
			imp.impliedBy = append(imp.impliedBy, n)
			implied = true
		}
		if !implied {
			tree = append(tree, n)
		}
	}
	// Remove filtered leaves.
	if filter != nil {
		tree = removeFiltered(tree, filter)
	}
	// Sort each level of the tree by priority.
	sortTree(tree)
	// Render the tree in string format.
	renderTree(&sb, tree, 0)
	return sb.String()
}
func removeFiltered(tree []*node, filter func(*role.Role) bool) []*node {
	j := 0
	for _, n := range tree {
		n.impliedBy = removeFiltered(n.impliedBy, filter)
		if filter(n.role) {
			tree[j] = n
			j++
		} else if len(n.impliedBy) != 0 {
			n.disabled = true
			tree[j] = n
			j++
		}
	}
	return tree[:j]
}
func sortTree(tree []*node) {
	sort.Slice(tree, func(i, j int) bool { return tree[i].role.Priority() < tree[j].role.Priority() })
	for _, child := range tree {
		sortTree(child.impliedBy)
	}
}
func renderTree(sb *strings.Builder, tree []*node, indent int) {
	for _, n := range tree {
		sb.WriteString("                    "[:indent])
		if n.disabled {
			sb.WriteByte('-')
		}
		fmt.Fprintf(sb, "%d %s\n", n.role.ID(), n.role.Name())
		renderTree(sb, n.impliedBy, indent+1)
	}
}
