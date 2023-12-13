package roleselect

import (
	"fmt"
	"sort"
	"strings"

	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util/request"
)

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
