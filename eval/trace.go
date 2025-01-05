package eval

import (
	"fmt"
)

// nodeAddr returns the pointer address of node, short version.
func ptrAddr(v any) string {
	p := fmt.Sprintf("%p", v)
	return p[:2] + p[9:] // unique bits
}

func (n *node) String() string {
	s := n.kind.String()
	if n.ident != "" {
		s += " " + n.ident
	}
	s += " " + ptrAddr(n)
	if n.sym != nil {
		s += " sym:" + n.sym.String()
	} else if n.typ != nil {
		s += " typ:" + n.typ.String()
	}
	if n.findex >= 0 {
		s += fmt.Sprintf(" fidx: %d lev: %d", n.findex, n.level)
	}
	if n.start != nil && n.start != n {
		s += fmt.Sprintf(" ->start: %s %s", n.start.kind.String(), ptrAddr(n.start))
	}
	if n.tnext != nil {
		s += fmt.Sprintf(" ->tnext: %s %s", n.tnext.kind.String(), ptrAddr(n.tnext))
	}
	if n.fnext != nil {
		s += fmt.Sprintf(" ->fnext: %s %s", n.fnext.kind.String(), ptrAddr(n.fnext))
	}
	return s
}

func (n *node) Depth() int {
	if n.anc != nil {
		return n.anc.Depth() + 1
	}
	return 0
}

func (sy *symbol) String() string {
	s := sy.kind.String()
	if sy.typ != nil {
		s += " (" + sy.typ.String() + ")"
	}
	if sy.rval.IsValid() {
		s += " = " + sy.rval.String()
	}
	if sy.index >= 0 {
		s += fmt.Sprintf(" idx: %d", sy.index)
	}
	if sy.node != nil {
		s += " " + sy.node.String()
	}
	return s
}

func (t *itype) String() string {
	if t.str != "" {
		return t.str
	}
	s := t.cat.String()
	if t.name != "" {
		s += " (" + t.name + ")"
	}
	return s
}
